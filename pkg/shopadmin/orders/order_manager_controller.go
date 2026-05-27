package orders

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/hb"
)

type orderManagerController struct {
	registry registry.RegistryInterface
}

func NewOrderManagerController(registry registry.RegistryInterface) *orderManagerController {
	return &orderManagerController{registry: registry}
}

func (controller *orderManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Orders | Shop",
		Content: hb.Div().HTML("Order Manager - TODO: Implement full migration"),
	}).ToHTML()
}
