package blog_settings

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/types"

	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

type formBlogSettings struct {
	liveflux.Base
	App                      types.RegistryInterface
	BlogTopic                string
	FormErrorMessage         string
	FormSuccessMessage       string
	FormInfoMessage          string
	FormRedirect             string
	FormRedirectDelaySeconds int
	IsEnvOverride            bool
	ReturnURL                string
}

func NewFormBlogSettings(app types.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formBlogSettings{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if c, ok := inst.(*formBlogSettings); ok {
		c.App = app
	}
	return inst
}

func (c *formBlogSettings) GetKind() string {
	return "admin_blog_settings_form_component"
}

func (c *formBlogSettings) Mount(ctx context.Context, params map[string]string) error {
	if c.App == nil {
		c.FormErrorMessage = "Application not initialized"
		return nil
	}

	c.ReturnURL = strings.TrimSpace(params["return_url"])
	if c.ReturnURL == "" {
		c.ReturnURL = shared.NewLinks().PostManager()
	}

	store := c.App.GetSettingStore()
	if store == nil {
		c.FormErrorMessage = "Setting store is not configured"
		return nil
	}

	value, err := store.Get(ctx, SettingKeyBlogTopic, "")
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Blog settings form: failed to load blog topic", "error", err.Error())
		}
		c.FormErrorMessage = "Failed to load blog settings"
		return nil
	}

	c.BlogTopic = value

	envTopic := strings.TrimSpace(os.Getenv("BLOG_TOPIC"))
	if envTopic != "" {
		c.BlogTopic = envTopic
		c.IsEnvOverride = true
		c.FormInfoMessage = "The BLOG_TOPIC environment variable is set, so updates are disabled here."
	}

	return nil
}

func (c *formBlogSettings) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "apply", "save_close":
		return c.handleUpdate(ctx, action, data)
	default:
		return nil
	}
}

func (c *formBlogSettings) handleUpdate(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}

	topic := strings.TrimSpace(data.Get("blog_topic"))
	if topic == "" {
		c.FormErrorMessage = "Blog topic is required"
		c.FormSuccessMessage = ""
		return nil
	}

	if c.IsEnvOverride {
		c.FormErrorMessage = "Blog topic is managed via environment and cannot be changed here."
		c.FormSuccessMessage = ""
		return nil
	}

	store := c.App.GetSettingStore()
	if store == nil {
		c.FormErrorMessage = "Setting store is not configured"
		c.FormSuccessMessage = ""
		return nil
	}

	if err := store.Set(ctx, SettingKeyBlogTopic, topic); err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Blog settings form: failed to save blog topic", "error", err.Error())
		}
		c.FormErrorMessage = "Failed to save blog topic. Please try again later."
		c.FormSuccessMessage = ""
		return nil
	}

	c.BlogTopic = topic
	c.FormErrorMessage = ""
	c.FormSuccessMessage = "Blog settings saved successfully"

	switch action {
	case "apply":
		c.FormRedirect = ""
		c.FormRedirectDelaySeconds = 0
	case "save_close":
		c.FormRedirect = c.ReturnURL
		c.FormRedirectDelaySeconds = 2
	default:
		c.FormRedirect = shared.NewLinks().BlogSettings()
		c.FormRedirectDelaySeconds = 2
	}

	return nil
}

func (c *formBlogSettings) Render(ctx context.Context) hb.TagInterface {
	alerts := hb.Div()
	if c.FormErrorMessage != "" {
		alerts = alerts.Child(hb.SwalError(hb.SwalOptions{
			Text:             c.FormErrorMessage,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		}))
	}
	if c.FormSuccessMessage != "" {
		if c.FormRedirect != "" {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.FormSuccessMessage,
				RedirectURL:      c.FormRedirect,
				RedirectSeconds:  c.FormRedirectDelaySeconds,
				Timer:            5000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		} else {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.FormSuccessMessage,
				Timer:            5000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		}
	}
	if c.FormInfoMessage != "" {
		alerts = alerts.Child(hb.Div().
			Class("alert alert-info d-flex align-items-center gap-2 mb-3").
			Attr("role", "alert").
			Child(hb.I().Class("bi bi-info-circle-fill")).
			Child(hb.Span().Text(c.FormInfoMessage)))
	}

	textareaBlogTopic := hb.TextArea().
		Class("form-control").
		ID("blog_topic").
		Name("blog_topic").
		Text(c.BlogTopic).
		Placeholder("e.g. Contract review insights").
		Attr("rows", "5")

	if c.IsEnvOverride {
		textareaBlogTopic = textareaBlogTopic.
			Attr("readonly", "true").
			Attr("disabled", "true")
	}

	formGroup := hb.Div().
		Class("mb-3").
		Child(hb.Label().
			Class("form-label fw-semibold").
			For("blog_topic").
			HTML("Blog Topic")).
		Child(textareaBlogTopic).
		Child(hb.Small().
			Class("form-text text-muted").
			Text("This is the topic of your blog. It will be used by the AI to generate titles and content for your blog."))

	applyIndicator := hb.Span().
		ID("BlogSettingsApplyIndicator").
		Class("spinner-border spinner-border-sm ms-2").
		Style(`display: none;`).
		Role("status").
		Aria("hidden", "true")

	saveIndicator := hb.Span().
		ID("BlogSettingsSaveIndicator").
		Class("spinner-border spinner-border-sm ms-2").
		Style(`display: none;`).
		Role("status").
		Aria("hidden", "true")

	buttonApply := hb.Button().
		Type("submit").
		Class("btn btn-primary").
		Attr(liveflux.DataFluxAction, "apply").
		Attr(liveflux.DataFluxIndicator, "#BlogSettingsApplyIndicator").
		Child(hb.I().Class("bi bi-check2 me-2")).
		Child(hb.Span().Text("Apply")).
		Child(applyIndicator)

	buttonSaveClose := hb.Button().
		Type("submit").
		Class("btn btn-success").
		Attr(liveflux.DataFluxAction, "save_close").
		Attr(liveflux.DataFluxIndicator, "#BlogSettingsSaveIndicator").
		Child(hb.I().Class("bi bi-check2-all me-2")).
		Child(hb.Span().Text("Save & Close")).
		Child(saveIndicator)

	if c.IsEnvOverride {
		buttonApply = buttonApply.Attr("disabled", "true")
		buttonSaveClose = buttonSaveClose.Attr("disabled", "true")
	}

	submitRow := hb.Div().
		Class("d-flex justify-content-between align-items-center flex-wrap gap-2").
		Child(hb.A().
			Href(c.ReturnURL).
			Class("btn btn-secondary").
			Child(hb.I().Class("bi bi-chevron-left me-2")).
			Text("Cancel")).
		Child(hb.Div().
			Class("d-flex gap-2").
			Child(buttonApply).
			Child(buttonSaveClose))

	form := hb.Div().
		ID("BlogSettingsFormWrapper").
		Child(alerts).
		Child(formGroup).
		Child(submitRow)

	return c.Root(form)
}

func init() {
	if err := liveflux.Register(&formBlogSettings{}); err != nil {
		log.Printf("Failed to register formBlogSettings component: %v", err)
	}
}
