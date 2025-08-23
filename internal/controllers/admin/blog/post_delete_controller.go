package admin

import (
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/base/req"
	"github.com/gouniverse/blogstore"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
)

type postDeleteController struct {
	app types.AppInterface
}

type postDeleteControllerData struct {
	postID         string
	post           *blogstore.Post
	successMessage string
}

func NewPostDeleteController(app types.AppInterface) *postDeleteController {
	return &postDeleteController{app: app}
}

func (controller *postDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return hb.Swal(hb.SwalOptions{
			Icon: "error",
			Text: errorMessage,
		}).ToHTML()
	}

	if data.successMessage != "" {
		return hb.Wrap().
			Child(hb.Swal(hb.SwalOptions{
				Icon: "success",
				Text: data.successMessage,
			})).
			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)")).
			ToHTML()
	}

	return controller.
		modal(data).
		ToHTML()
}

func (controller *postDeleteController) modal(data postDeleteControllerData) hb.TagInterface {
	submitUrl := links.NewAdminLinks().BlogPostDelete(map[string]string{
		"post_id": data.postID,
	})

	modalID := "ModalPostDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupPostId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("post_id").
		Value(data.postID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalPostDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Post").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalPostDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(hb.Paragraph().Text("Are you sure you want to delete this post?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupPostId)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(
						hb.Button().HTML("Close").
							Class("btn btn-secondary float-start").
							Data("bs-dismiss", "modal").
							OnClick(modalCloseScript)).
					Child(buttonDelete)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().
		Children([]hb.TagInterface{
			modal,
			backdrop,
		})
}

func (controller *postDeleteController) prepareDataAndValidate(r *http.Request) (data postDeleteControllerData, errorMessage string) {
	authUser := helpers.GetAuthUser(r)
	data.postID = req.Value(r, "post_id")

	if authUser == nil {
		return data, "You are not logged in. Please login to continue."
	}

	if data.postID == "" {
		return data, "post id is required"
	}

	post, err := controller.app.GetBlogStore().PostFindByID(data.postID)

	if err != nil {
		controller.app.GetLogger().Error("At postDeleteController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Post not found"
	}

	if post == nil {
		return data, "Post not found"
	}

	data.post = post

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.app.GetBlogStore().PostTrash(post)

	if err != nil {
		controller.app.GetLogger().Error("At postDeleteController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Deleting post failed. Please contact an administrator."
	}

	data.successMessage = "post deleted successfully."

	return data, ""

}
