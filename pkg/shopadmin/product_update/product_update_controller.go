package product_update

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

const (
	actionLoadProduct   = "load-product"
	actionUpdateProduct = "update-product"
)

// == CONTROLLER ==============================================================

type productUpdateController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewProductUpdateController(app app.AppInterface) *productUpdateController {
	return &productUpdateController{app: app}
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
