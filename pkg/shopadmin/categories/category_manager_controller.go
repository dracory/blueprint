package categories

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/hb"
)

type categoryManagerController struct {
	registry registry.RegistryInterface
}

func NewCategoryManagerController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}

func (controller *categoryManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Categories | Shop",
		Content: hb.Div().HTML("Category Manager - TODO: Implement full migration"),
	}).ToHTML()
}

func NewCategoryCreateController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}

func NewCategoryUpdateController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}
