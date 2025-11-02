package post_create

import (
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

type postCreateController struct {
	app types.AppInterface
}

type postCreateControllerData struct {
	title          string
	successMessage string
}

func NewPostCreateController(app types.AppInterface) *postCreateController {
	return &postCreateController{app: app}
}

func (controller *postCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
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

	return modalPostCreate(data).ToHTML()
}

func (controller *postCreateController) prepareDataAndValidate(r *http.Request) (data postCreateControllerData, errorMessage string) {
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return data, "You are not logged in. Please login to continue."
	}

	data.title = req.GetStringTrimmed(r, "post_title")

	if r.Method != "POST" {
		return data, ""
	}

	if data.title == "" {
		return data, "post title is required"
	}

	post := blogstore.NewPost()
	post.SetTitle(data.title)

	err := controller.app.GetBlogStore().PostCreate(post)

	if err != nil {
		controller.app.GetLogger().Error("At postCreateController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Creating post failed. Please contact an administrator."
	}

	data.successMessage = "post created successfully."

	return data, ""

}
