package order_manager

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

const (
	actionLoadOrdersAjax = "load-orders-ajax"
)

// == CONTROLLER ==============================================================

type orderManagerController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewOrderManagerController(app app.AppInterface) *orderManagerController {
	return &orderManagerController{app: app}
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
