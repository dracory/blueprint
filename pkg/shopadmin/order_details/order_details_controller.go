package order_details

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

type orderDetailsController struct {
	registry registry.RegistryInterface
}

func NewOrderDetailsController(registry registry.RegistryInterface) *orderDetailsController {
	return &orderDetailsController{registry: registry}
}

func (controller *orderDetailsController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadOrderDetailsAjax:
		return controller.handleOrderDetailsLoadAjax(w, r)
	default:
		return controller.renderPage(r)
	}
}
