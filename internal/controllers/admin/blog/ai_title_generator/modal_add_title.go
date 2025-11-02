package aititlegenerator

import (
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

const (
	modalAddTitleID            = "ModalAiAddTitle"
	modalAddTitleBackdropClass = "ModalBackdrop"
)

func (c *AiTitleGeneratorController) onAddTitleModal(r *http.Request) string {
	customTitle := req.GetStringTrimmed(r, "custom_title")

	formGroupTitle := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Custom Title").Class("form-label fw-semibold")).
		Child(
			bs.FormInput().
				Name("custom_title").
				Value(customTitle).
				Placeholder("Enter a custom blog title").
				Attr("required", "required"),
		).
		Child(hb.Small().Class("text-muted").Text("Approved titles move to the AI post generator."))

	modalCloseScript := fmt.Sprintf("closeModal%s();", modalAddTitleID)
	jsCloseFn := fmt.Sprintf(
		"function closeModal%s(){var modal=document.getElementById('%s');if(modal){modal.remove();}var backdrops=document.getElementsByClassName('%s');while(backdrops.length){backdrops[0].remove();}}",
		modalAddTitleID,
		modalAddTitleID,
		modalAddTitleBackdropClass,
	)

	buttonSubmit := hb.Button().
		Attr("id", "ButtonSubmitCustomTitle").
		Class("btn btn-primary float-end").
		HTML(`Save Title <span class="htmx-indicator spinner-border spinner-border-sm" role="status"></span>`).
		HxInclude("#"+modalAddTitleID).
		HxPost(shared.NewLinks().AiTitleGenerator(map[string]string{"action": ACTION_ADD_TITLE})).
		HxTarget("body").
		HxSwap("beforeend").
		Attr("hx-indicator", "this")

	buttonCancel := hb.Button().
		Class("btn btn-secondary float-start").
		HTML("Close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	modal := bs.Modal().
		ID(modalAddTitleID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(bs.ModalHeader().
					Child(hb.Heading5().Class("mb-0").HTML("Add Custom Title")).
					Child(hb.Button().Type("button").Class("btn-close").Data("bs-dismiss", "modal").OnClick(modalCloseScript))).
				Child(bs.ModalBody().Child(formGroupTitle)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSubmit))))

	backdrop := hb.Div().
		Class(modalAddTitleBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1050;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	}).ToHTML()
}
