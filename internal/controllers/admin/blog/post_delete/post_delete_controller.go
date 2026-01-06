package post_delete

import (
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

type postDeleteController struct {
	app types.RegistryInterface
}

type postDeleteControllerData struct {
	postID         string
	post           *blogstore.Post
	successMessage string
}

func NewPostDeleteController(app types.RegistryInterface) *postDeleteController {
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

	return modalPostDelete(data).ToHTML()
}

func (controller *postDeleteController) prepareDataAndValidate(r *http.Request) (data postDeleteControllerData, errorMessage string) {
	authUser := helpers.GetAuthUser(r)
	data.postID = req.GetStringTrimmed(r, "post_id")

	if authUser == nil {
		return data, "You are not logged in. Please login to continue."
	}

	if data.postID == "" {
		return data, "post id is required"
	}

	post, err := controller.app.GetBlogStore().PostFindByID(r.Context(), data.postID)

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

	err = controller.app.GetBlogStore().PostTrash(r.Context(), post)

	if err != nil {
		controller.app.GetLogger().Error("At postDeleteController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Deleting post failed. Please contact an administrator."
	}

	data.successMessage = "post deleted successfully."

	return data, ""

}
