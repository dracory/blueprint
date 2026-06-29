package blog_settings

import (
	"embed"
	"net/http"
	"project/internal/app"
)

//go:embed *.html
//go:embed *.js
var settingsFiles embed.FS

const (
	actionFetchData = "fetch-data"
	actionSubmit    = "submit"
)

type blogSettingsController struct {
	app app.AppInterface
}

func NewBlogSettingsController(app app.AppInterface) *blogSettingsController {
	return &blogSettingsController{app: app}
}

func (controller *blogSettingsController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")
	if action == "" {
		action = r.PostFormValue("action")
	}

	switch action {
	case actionFetchData:
		return controller.handleFetchData(r)
	case actionSubmit:
		return controller.handleSubmit(r)
	default:
		return controller.renderPage(w, r)
	}
}
