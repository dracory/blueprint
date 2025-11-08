package blog_settings

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

type blogSettingsController struct {
	app types.AppInterface
}

type blogSettingsData struct {
	blogTopic                string
	formErrorMessage         string
	formSuccessMessage       string
	formInfoMessage          string
	formRedirect             string
	formRedirectDelaySeconds int
	isEnvOverride            bool
}

func NewBlogSettingsController(app types.AppInterface) *blogSettingsController {
	return &blogSettingsController{app: app}
}

func (c *blogSettingsController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errMessage := c.prepareData(r)
	if errMessage != "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, errMessage, shared.NewLinks().Home(), 10)
	}

	var formComponent liveflux.ComponentInterface

	if r.Method == http.MethodPost {
		// Create the component for handling POST
		formComponent = NewFormBlogSettings(c.app)
		if formComponent != nil {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				data.formErrorMessage = "Failed to parse form data"
			} else {
				// Handle the action using the component
				action := req.GetStringTrimmed(r, "action")
				if action == "" {
					action = "apply" // default action
				}

				params := map[string]string{
					"return_url": shared.NewLinks().PostManager(),
				}

				// Mount the component first
				ctx := r.Context()
				if err := formComponent.Mount(ctx, params); err != nil {
					data.formErrorMessage = "Failed to initialize form"
				} else {
					// Handle the action
					if err := formComponent.Handle(ctx, action, r.Form); err != nil {
						data.formErrorMessage = "Failed to process form"
					} else {
						// Update data from component state
						if comp, ok := formComponent.(*formBlogSettings); ok {
							data.blogTopic = comp.BlogTopic
							data.formErrorMessage = comp.FormErrorMessage
							data.formSuccessMessage = comp.FormSuccessMessage
							data.formInfoMessage = comp.FormInfoMessage
							data.formRedirect = comp.FormRedirect
							data.formRedirectDelaySeconds = comp.FormRedirectDelaySeconds
							data.isEnvOverride = comp.IsEnvOverride
						}
					}
				}
			}
		} else {
			data.formErrorMessage = "Failed to initialize form component"
		}
	}

	// Use the existing component if available, otherwise create a new one
	if formComponent == nil {
		formComponent = NewFormBlogSettings(c.app)
	}

	rendered := liveflux.SSR(formComponent, map[string]string{
		"return_url": shared.NewLinks().PostManager(),
	})

	if rendered == nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Error rendering blog settings form", shared.NewLinks().Home(), 10)
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "Settings | Blog",
		Content: c.page(data, rendered),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
		},
		Scripts: []string{
			ext.HxHideIndicatorCSS(),
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}

func (c *blogSettingsController) prepareData(r *http.Request) (blogSettingsData, string) {
	data := blogSettingsData{}

	if helpers.GetAuthUser(r) == nil {
		return data, "You are not logged in. Please login to continue."
	}

	store := c.app.GetSettingStore()
	if store == nil {
		c.app.GetLogger().Error("Blog settings controller: setting store is not configured")
		return data, "Blog settings are unavailable. Please contact an administrator."
	}

	value, err := store.Get(r.Context(), SettingKeyBlogTopic, "")
	if err != nil {
		c.app.GetLogger().Error("Blog settings controller: failed to load blog topic", slog.String("error", err.Error()))
		return data, "Failed to load blog settings. Please try again later."
	}

	data.blogTopic = value

	envTopic := strings.TrimSpace(os.Getenv("BLOG_TOPIC"))
	if envTopic != "" {
		data.blogTopic = envTopic
		data.isEnvOverride = true
		data.formInfoMessage = "The BLOG_TOPIC environment variable is set, so updates are disabled here."
	}

	return data, ""
}

func (c *blogSettingsController) page(data blogSettingsData, component hb.TagInterface) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "Settings",
			URL:  shared.NewLinks().BlogSettings(),
		},
	})

	heading := hb.Heading1().
		HTML("Blog Settings")

	buttonBack := hb.Hyperlink().
		Class("btn btn-secondary ms-3").
		HTML("Back to Blog").
		Href(shared.NewLinks().Home())

	cardBody := hb.Div().
		Class("card-body").
		Child(component)

	card := hb.Div().
		Class("card shadow-sm").
		Child(hb.Div().
			Class("card-header d-flex justify-content-between align-items-center").
			Child(hb.Heading4().Class("mb-0").HTML("General Settings"))).
		Child(cardBody)

	return hb.Div().
		Class("container py-4 min-vh-100").
		Child(breadcrumbs).
		Child(hb.Div().
			Class("d-flex align-items-center mb-4").
			Child(heading).
			Child(buttonBack)).
		Child(card)
}
