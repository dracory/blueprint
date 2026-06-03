package order_manager

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

const (
	actionLoadOrdersAjax = "load-orders-ajax"
)

// == CONTROLLER ==============================================================

type orderManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewOrderManagerController(registry registry.RegistryInterface) *orderManagerController {
	return &orderManagerController{registry: registry}
}

func (controller *orderManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadOrdersAjax:
		return controller.handleOrdersLoadAjax(w, r)
	default:
		return controller.renderPage(r)
	}
}
