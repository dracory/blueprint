package categories

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/app"

	"github.com/dracory/hb"
)

type categoryManagerController struct {
	app app.AppInterface
}

func NewCategoryManagerController(app app.AppInterface) *categoryManagerController {
	return &categoryManagerController{app: app}
}

func (controller *categoryManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Categories | Shop",
		Content: hb.Div().HTML("Category Manager - TODO: Implement full migration"),
	}).ToHTML()
}

func NewCategoryCreateController(app app.AppInterface) *categoryManagerController {
	return &categoryManagerController{app: app}
}

func NewCategoryUpdateController(app app.AppInterface) *categoryManagerController {
	return &categoryManagerController{app: app}
}
