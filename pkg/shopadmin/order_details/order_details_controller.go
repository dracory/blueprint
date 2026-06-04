package order_details

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

type orderDetailsController struct {
	app app.AppInterface
}

func NewOrderDetailsController(app app.AppInterface) *orderDetailsController {
	return &orderDetailsController{app: app}
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
