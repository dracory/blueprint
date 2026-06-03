package user_update

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

type userUpdateController struct {
	registry registry.RegistryInterface
}

const (
	actionUserFetch    = "user-fetch-ajax"
	actionGetTimezones = "get-timezones-ajax"
	actionUserUpdate   = "user-update-ajax"
)

func NewUserUpdateController(registry registry.RegistryInterface) *userUpdateController {
	return &userUpdateController{registry: registry}
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
