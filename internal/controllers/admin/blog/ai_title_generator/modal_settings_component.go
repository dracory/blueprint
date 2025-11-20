package aititlegenerator

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/types"

	livefluxctl "project/internal/controllers/liveflux"

	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

type titleGeneratorSettingsModal struct {
	liveflux.Base

	App                      types.AppInterface
	FormBlogTopic            string
	FormErrorMessage         string
	FormSuccessMessage       string
	FormInfoMessage          string
	FormRedirect             string
	FormRedirectDelaySeconds int
	ReturnURL                string
	IsOpen                   bool
}

func NewTitleGeneratorSettingsModal(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&titleGeneratorSettingsModal{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*titleGeneratorSettingsModal); ok {
		c.App = app
	}

	return inst
}

func (c *titleGeneratorSettingsModal) GetKind() string {
	return "admin_ai_title_generator_settings_modal"
}

func (c *titleGeneratorSettingsModal) Mount(ctx context.Context, params map[string]string) error {
	if c.App == nil {
		if app, ok := ctx.Value(livefluxctl.AppContextKey).(types.AppInterface); ok {
			c.App = app
		}
	}

	c.ReturnURL = strings.TrimSpace(params["return_url"])
	if c.ReturnURL == "" {
		c.ReturnURL = shared.NewLinks().AiTitleGenerator()
	}

	store := c.App.GetSettingStore()
	if store == nil {
		c.FormErrorMessage = "Setting store is not configured"
		return nil
	}

	value, err := store.Get(ctx, SETTING_KEY_BLOG_TOPIC, "")
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("AI title generator settings modal: failed to load blog title", "error", err.Error())
		}
		c.FormErrorMessage = "Failed to load title generator settings"
		return nil
	}

	c.FormBlogTopic = strings.TrimSpace(value)

	c.IsOpen = false

	return nil
}

func (c *titleGeneratorSettingsModal) Handle(ctx context.Context, action string, data url.Values) error {
	// get app from context
	// c.App = ctx.Value("app").(types.AppInterface)

	switch action {
	case "open":
		c.IsOpen = true
		c.FormErrorMessage = ""
		c.FormSuccessMessage = ""
		return nil
	case "close":
		c.IsOpen = false
		return nil
	case "apply", "save_close":
		return c.onSave(ctx, action, data)
	default:
		return nil
	}
}

func (c *titleGeneratorSettingsModal) onSave(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}

	blogTopic := strings.TrimSpace(data.Get("blog_topic"))
	if blogTopic == "" {
		c.FormErrorMessage = "Blog topic is required"
		c.FormSuccessMessage = ""
		return nil
	}

	store := c.App.GetSettingStore()
	if store == nil {
		c.FormErrorMessage = "Setting store is not configured"
		c.FormSuccessMessage = ""
		return nil
	}

	if err := store.Set(ctx, SETTING_KEY_BLOG_TOPIC, blogTopic); err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("AI title generator settings modal: failed to save blog topic", "error", err.Error())
		}
		c.FormErrorMessage = "Failed to save blog topic. Please try again later."
		c.FormSuccessMessage = ""
		return nil
	}

	c.FormBlogTopic = blogTopic
	c.FormErrorMessage = ""
	c.FormSuccessMessage = "Settings saved successfully"

	switch action {
	case "apply":
		c.FormRedirect = ""
		c.FormRedirectDelaySeconds = 0
	case "save_close":
		c.FormRedirect = c.ReturnURL
		c.FormRedirectDelaySeconds = 2
	default:
		c.FormRedirect = shared.NewLinks().AiTitleGenerator()
		c.FormRedirectDelaySeconds = 2
	}

	return nil
}

func (c *titleGeneratorSettingsModal) Render(ctx context.Context) hb.TagInterface {
	if !c.IsOpen {
		return c.Root(hb.Div().ID("AiTitleGeneratorSettingsModalContainer"))
	}

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

	textareaSystemPrompt := hb.TextArea().
		Class("form-control").
		ID("blog_topic").
		Name("blog_topic").
		Text(c.FormBlogTopic).
		Attr("rows", "8").
		Placeholder("Set the main blog topic used to generate AI titles.")

	formGroup := hb.Div().
		Class("mb-3").
		Child(hb.Label().
			Class("form-label fw-semibold").
			For("blog_topic").
			HTML("Blog Topic")).
		Child(textareaSystemPrompt).
		Child(hb.Small().
			Class("form-text text-muted").
			Text("This blog topic guides how the AI generates suggestions. Make it specific to your niche."))

	applyIndicator := hb.Span().
		ID("AiTitleSettingsApplyIndicator").
		Class("spinner-border spinner-border-sm ms-2").
		Style(`display: none;`).
		Role("status").
		Aria("hidden", "true")

	saveIndicator := hb.Span().
		ID("AiTitleSettingsSaveIndicator").
		Class("spinner-border spinner-border-sm ms-2").
		Style(`display: none;`).
		Role("status").
		Aria("hidden", "true")

	buttonApply := hb.Button().
		Type("submit").
		Class("btn btn-primary").
		Attr(liveflux.DataFluxAction, "apply").
		Attr(liveflux.DataFluxIndicator, "#AiTitleSettingsApplyIndicator").
		Child(hb.I().Class("bi bi-check2 me-2")).
		Child(hb.Span().Text("Apply")).
		Child(applyIndicator)

	buttonSaveClose := hb.Button().
		Type("submit").
		Class("btn btn-success").
		Attr(liveflux.DataFluxAction, "save_close").
		Attr(liveflux.DataFluxIndicator, "#AiTitleSettingsSaveIndicator").
		Child(hb.I().Class("bi bi-check2-all me-2")).
		Child(hb.Span().Text("Save & Close")).
		Child(saveIndicator)

	buttonClose := hb.Button().
		Type("button").
		Class("btn btn-outline-secondary float-start").
		Attr(liveflux.DataFluxAction, "close").
		Attr(liveflux.DataFluxIndicator, "this").
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		Child(hb.Span().Text("Close"))

	footerButtons := hb.Div().
		Class("d-flex w-100 justify-content-between align-items-center flex-wrap gap-2").
		Child(buttonClose).
		Child(hb.Div().
			Class("d-flex gap-2").
			Child(buttonApply).
			Child(buttonSaveClose))

	modalBody := hb.Div().
		Class("modal-body").
		Child(alerts).
		Child(formGroup)

	modalContent := hb.Div().
		Class("modal-content").
		Child(hb.Div().
			Class("modal-header").
			Child(hb.Heading5().Class("modal-title mb-0").HTML("Title Generator Settings")).
			Child(hb.Button().
				Type("button").
				Class("btn-close").
				Attr(liveflux.DataFluxAction, "close")),
		).
		Child(modalBody).
		Child(hb.Div().Class("modal-footer").Child(footerButtons))

	modalDialog := hb.Div().
		Class("modal-dialog modal-lg modal-dialog-centered").
		Child(modalContent)

	modal := hb.Div().
		ID("AiTitleGeneratorSettingsModal").
		Class("modal fade show").
		Attr("tabindex", "-1").
		Style("display:block; background: rgba(0,0,0,0.5);").
		Child(hb.Div().
			Class("modal-dialog-wrapper d-flex align-items-center justify-content-center min-vh-100").
			Child(modalDialog))

	return c.Root(modal)
}

func init() {
	if err := liveflux.Register(&titleGeneratorSettingsModal{}); err != nil {
		log.Printf("Failed to register titleGeneratorSettingsModal component: %v", err)
	}
}
