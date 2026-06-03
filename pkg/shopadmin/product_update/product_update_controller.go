package product_update

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

const (
	actionLoadProduct   = "load-product"
	actionUpdateProduct = "update-product"
)

// == CONTROLLER ==============================================================

type productUpdateController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewProductUpdateController(registry registry.RegistryInterface) *productUpdateController {
	return &productUpdateController{registry: registry}
}

func (controller *productUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadProduct:
		return controller.handleLoadProduct(w, r)
	case actionUpdateProduct:
		return controller.handleUpdateProduct(w, r)
	default:
		return controller.renderPage(w, r)
	}
}
