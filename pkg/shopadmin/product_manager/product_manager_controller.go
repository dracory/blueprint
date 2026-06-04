package product_manager

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

const (
	actionLoadProducts          = "load-products"
	actionProductDelete         = "delete-product"
	actionProductDeleteSelected = "delete-selected"
	actionCreateProduct         = "create-product-ajax"
)

// == CONTROLLER ==============================================================

type productManagerController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewProductManagerController(app app.AppInterface) *productManagerController {
	return &productManagerController{app: app}
}

func (controller *productManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadProducts:
		return controller.handleLoadProducts(w, r)
	case actionProductDelete:
		return controller.handleProductDelete(w, r)
	case actionProductDeleteSelected:
		return controller.handleProductDeleteSelected(w, r)
	case actionCreateProduct:
		return controller.handleProductCreateAjax(w, r)
	default:
		return controller.renderPage(w, r)
	}
}
