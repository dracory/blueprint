package blog_settings

import (
	"project/internal/controllers/admin/blog/shared"

	"github.com/dracory/hb"
)

type blogSettingsFormOptions struct {
	Data blogSettingsData
}

func blogSettingsForm(options blogSettingsFormOptions) hb.TagInterface {
	textarea := hb.TextArea().
		Class("form-control").
		ID("blog_topic").
		Name("blog_topic").
		Text(options.Data.blogTopic).
		Attr("placeholder", "e.g. Contract review insights")

	if options.Data.isEnvOverride {
		textarea = textarea.
			Attr("readonly", "true").
			Attr("disabled", "true")
	}

	formGroup := hb.Div().
		Class("mb-3").
		Child(hb.Label().
			Class("form-label fw-semibold").
			Attr("for", "blog_topic").
			HTML("Blog Topic")).
		Child(textarea)

	applyIndicator := hb.Span().
		ID("BlogSettingsApplyIndicator").
		Class("htmx-indicator spinner-border spinner-border-sm ms-2").
		Attr("role", "status").
		Attr("aria-hidden", "true")

	saveIndicator := hb.Span().
		ID("BlogSettingsSaveIndicator").
		Class("htmx-indicator spinner-border spinner-border-sm ms-2").
		Attr("role", "status").
		Attr("aria-hidden", "true")

	buttonApply := hb.Button().
		Type("button").
		Class("btn btn-primary d-inline-flex align-items-center")

	buttonSaveClose := hb.Button().
		Type("button").
		Class("btn btn-success d-inline-flex align-items-center")

	if options.Data.isEnvOverride {
		buttonApply = buttonApply.
			Attr("disabled", "true")
		buttonSaveClose = buttonSaveClose.
			Attr("disabled", "true")
	} else {
		buttonApply = buttonApply.
			Attr("hx-post", shared.NewLinks().BlogSettings()).
			Attr("hx-target", "#BlogSettingsFormWrapper").
			Attr("hx-swap", "outerHTML").
			Attr("hx-include", "#BlogSettingsFormWrapper").
			Attr("hx-vals", "{\"action\":\"apply\"}").
			Attr("hx-indicator", "#BlogSettingsApplyIndicator")

		buttonSaveClose = buttonSaveClose.
			Attr("hx-post", shared.NewLinks().BlogSettings()).
			Attr("hx-target", "#BlogSettingsFormWrapper").
			Attr("hx-swap", "outerHTML").
			Attr("hx-include", "#BlogSettingsFormWrapper").
			Attr("hx-vals", "{\"action\":\"save_close\"}").
			Attr("hx-indicator", "#BlogSettingsSaveIndicator")
	}

	buttonApply = buttonApply.
		Child(hb.I().Class("bi bi-check2 me-2")).
		Child(hb.Span().Text("Apply")).
		Child(applyIndicator)

	buttonSaveClose = buttonSaveClose.
		Child(hb.I().Class("bi bi-check2-all me-2")).
		Child(hb.Span().Text("Save & Close")).
		Child(saveIndicator)

	submitRow := hb.Div().
		Class("d-flex justify-content-between align-items-center flex-wrap gap-2").
		Child(hb.Hyperlink().
			Class("btn btn-secondary d-inline-flex align-items-center").
			Href(shared.NewLinks().PostManager()).
			Child(hb.I().Class("bi bi-chevron-left me-2")).
			Child(hb.Span().Text("Cancel"))).
		Child(hb.Div().
			Class("d-flex gap-2").
			Child(buttonApply).
			Child(buttonSaveClose))

	flashMessages := hb.Div()

	errorMessage := hb.SwalError(hb.SwalOptions{
		Title:            "Error",
		Text:             options.Data.formErrorMessage,
		TimerProgressBar: true,
		Timer:            5000,
		Position:         "top-end",
	})

	successMessage := hb.SwalSuccess(hb.SwalOptions{
		Title:            "Success",
		Text:             options.Data.formSuccessMessage,
		TimerProgressBar: true,
		Timer:            5000,
		Position:         "top-end",
		RedirectURL:      options.Data.formRedirect,
		RedirectSeconds:  options.Data.formRedirectDelaySeconds,
	})

	infoMessage := hb.Div().
		Class("alert alert-info d-flex align-items-center gap-2 mb-3").
		Attr("role", "alert").
		Child(hb.I().Class("bi bi-info-circle-fill")).
		Child(hb.Span().Text(options.Data.formInfoMessage))

	flashMessages = flashMessages.ChildIf(options.Data.formErrorMessage != "", errorMessage)
	flashMessages = flashMessages.ChildIf(options.Data.formSuccessMessage != "", successMessage)
	flashMessages = flashMessages.ChildIf(options.Data.formInfoMessage != "", infoMessage)

	return hb.Div().
		ID("BlogSettingsFormWrapper").
		Child(flashMessages).
		Child(formGroup).
		Child(submitRow)
}
