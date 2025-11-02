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

	if r.Method == http.MethodPost {
		if data.isEnvOverride {
			data.formErrorMessage = "BLOG_TOPIC is managed via environment and cannot be changed here."
		} else {
			var postError string
			data, postError = c.processForm(r, data)
			if postError != "" {
				return helpers.ToFlashError(c.app.GetCacheStore(), w, r, postError, shared.NewLinks().Home(), 10)
			}
		}
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "Settings | Blog",
		Content: c.page(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
		},
		Scripts: []string{
			ext.HxHideIndicatorCSS(),
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

func (c *blogSettingsController) processForm(r *http.Request, data blogSettingsData) (blogSettingsData, string) {
	action := req.GetStringTrimmed(r, "action")
	topic := req.GetStringTrimmed(r, "blog_topic")

	if topic == "" {
		data.formErrorMessage = "Blog topic is required"
		return data, ""
	}

	store := c.app.GetSettingStore()
	if store == nil {
		c.app.GetLogger().Error("Blog settings controller: setting store is not configured for POST")
		return data, "Blog settings are unavailable. Please contact an administrator."
	}

	if err := store.Set(r.Context(), SettingKeyBlogTopic, topic); err != nil {
		c.app.GetLogger().Error("Blog settings controller: failed to save blog topic", slog.String("error", err.Error()))
		data.formErrorMessage = "Failed to save blog topic. Please try again later."
		return data, ""
	}

	data.blogTopic = topic
	data.formSuccessMessage = "Blog settings saved successfully"

	switch action {
	case "apply":
		data.formRedirect = ""
		data.formRedirectDelaySeconds = 0
	case "save_close":
		data.formRedirect = shared.NewLinks().PostManager()
		data.formRedirectDelaySeconds = 2
	default:
		data.formRedirect = shared.NewLinks().BlogSettings()
		data.formRedirectDelaySeconds = 2
	}
	return data, ""
}

func (c *blogSettingsController) page(data blogSettingsData) hb.TagInterface {
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
		Child(blogSettingsForm(blogSettingsFormOptions{Data: data}))

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
