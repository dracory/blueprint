package user_update

import (
	"net/http"

	"project/internal/app"

	"github.com/dracory/req"
)

type userUpdateController struct {
	app app.AppInterface
}

const (
	actionUserFetch    = "user-fetch-ajax"
	actionGetTimezones = "get-timezones-ajax"
	actionUserUpdate   = "user-update-ajax"
)

func NewUserUpdateController(app app.AppInterface) *userUpdateController {
	return &userUpdateController{app: app}
}

func (controller *userUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionUserFetch:
		controller.handleUserFetchAjax(w, r)
		return ""
	case actionGetTimezones:
		controller.handleTimezonesFetchAjax(w, r)
		return ""
	case actionUserUpdate:
		controller.handleUserUpdateAjax(w, r)
		return ""
	default:
		return controller.renderPage(w, r)
	}
}
