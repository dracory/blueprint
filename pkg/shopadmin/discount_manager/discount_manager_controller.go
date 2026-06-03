package discount_manager

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

const (
	actionLoadDiscounts          = "load-discounts"
	actionDiscountDelete         = "delete-discount"
	actionDiscountDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type discountManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewDiscountManagerController(registry registry.RegistryInterface) *discountManagerController {
	return &discountManagerController{registry: registry}
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
