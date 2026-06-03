package category_manager

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

const (
	actionLoadCategories         = "load-categories"
	actionCategoryDelete         = "delete-category"
	actionCategoryDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type categoryManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewCategoryManagerController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}

func (controller *categoryManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadCategories:
		return controller.handleLoadCategories(w, r)
	case actionCategoryDelete:
		return controller.handleCategoryDelete(w, r)
	case actionCategoryDeleteSelected:
		return controller.handleCategoryDeleteSelected(w, r)
	default:
		return controller.renderPage(w, r)
	}
}
