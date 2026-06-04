package discount_manager

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

const (
	actionLoadDiscounts          = "load-discounts"
	actionDiscountDelete         = "delete-discount"
	actionDiscountDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type discountManagerController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewDiscountManagerController(app app.AppInterface) *discountManagerController {
	return &discountManagerController{app: app}
}

func (controller *discountManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadDiscounts:
		return controller.handleLoadDiscounts(w, r)
	case actionDiscountDelete:
		return controller.handleDiscountDelete(w, r)
	case actionDiscountDeleteSelected:
		return controller.handleDiscountDeleteSelected(w, r)
	default:
		return controller.renderPage(w, r)
	}
}
