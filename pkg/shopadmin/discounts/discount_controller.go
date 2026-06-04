package discounts

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/app"

	"github.com/dracory/hb"
)

type discountController struct {
	app app.AppInterface
}

func NewDiscountController(app app.AppInterface) *discountController {
	return &discountController{app: app}
}

func (controller *discountController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Discounts | Shop",
		Content: hb.Div().HTML("Discount Manager - TODO: Implement full migration"),
	}).ToHTML()
}
