package aipostcontentupdate

import (
	"embed"
	"net/http"
	"project/internal/app"
)

//go:embed *.html
//go:embed *.js
var editorFiles embed.FS

const (
	actionFetchData       = "fetch-data"
	actionRegenerateBlock = "regenerate-block"
	actionSave            = "save"
)

type Controller struct {
	app app.AppInterface
}

func NewController(app app.AppInterface) *Controller {
	return &Controller{app: app}
}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")
	if action == "" {
		action = r.PostFormValue("action")
	}

	switch action {
	case actionFetchData:
		return c.handleFetchData(r)
	case actionRegenerateBlock:
		return c.handleRegenerateBlock(r)
	case actionSave:
		return c.handleSave(w, r)
	default:
		return c.renderPage(w, r)
	}
}
