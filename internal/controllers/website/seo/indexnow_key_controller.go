package seo

import (
	"net/http"
	"project/internal/app"
)

type indexNowKeyController struct {
	app app.AppInterface
}

func NewIndexNowKeyController(app app.AppInterface) *indexNowKeyController {
	return &indexNowKeyController{
		app: app,
	}
}

func (c indexNowKeyController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/plain")
	return c.app.GetConfig().GetIndexNowKey()
}
