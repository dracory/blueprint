package user_manager

import (
	"net/http"

	"project/internal/registry"

	"github.com/dracory/req"
)

type userManagerController struct{ registry registry.RegistryInterface }

const (
	actionLoadUsers  = "load-users-ajax"
	actionDeleteUser = "delete-user-ajax"
	actionCreateUser = "create-user-ajax"
)

func NewUserManagerController(registry registry.RegistryInterface) *userManagerController {
	return &userManagerController{registry: registry}
}

func (controller *userManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadUsers:
		return controller.handleUsersFetchAjax(w, r)
	case actionDeleteUser:
		return controller.handleUserDeleteAjax(w, r)
	case actionCreateUser:
		return controller.handleUserCreateAjax(w, r)
	default:
		return controller.renderPage(w, r)
	}
}
