package discounts

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/hb"
)

type discountController struct {
	registry registry.RegistryInterface
}

func NewDiscountController(registry registry.RegistryInterface) *discountController {
	return &discountController{registry: registry}
}

func (controller *discountController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Discounts | Shop",
		Content: hb.Div().HTML("Discount Manager - TODO: Implement full migration"),
	}).ToHTML()
}
