package users

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/useradmin/shared"
	"project/pkg/useradmin/user_create"
	"project/pkg/useradmin/user_delete"
	"project/pkg/useradmin/user_impersonate"
	"project/pkg/useradmin/user_manager"
	"project/pkg/useradmin/user_update"

	"github.com/dracory/req"
)

// usersAdminController handles all user admin requests
type usersAdminController struct {
	registry registry.RegistryInterface
}

// NewUsersAdminController creates a new users admin controller
func NewUsersAdminController(registry registry.RegistryInterface) *usersAdminController {
	return &usersAdminController{registry: registry}
}

// Handler processes users admin requests
func (controller *usersAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetAuthUser(r)
	if user == nil {
		http.Redirect(w, r, links.Admin().Home(), http.StatusSeeOther)
		return
	}

	controllerParam := req.GetStringTrimmed(r, "controller")

	var html string
	switch controllerParam {
	case shared.CONTROLLER_USER_MANAGER:
		html = user_manager.NewUserManagerController(controller.registry).Handler(w, r)
	case shared.CONTROLLER_USER_CREATE:
		html = user_create.NewUserCreateController(controller.registry).Handler(w, r)
	case shared.CONTROLLER_USER_DELETE:
		html = user_delete.NewUserDeleteController(controller.registry).Handler(w, r)
	case shared.CONTROLLER_USER_UPDATE:
		html = user_update.NewUserUpdateController(controller.registry).Handler(w, r)
	case shared.CONTROLLER_USER_IMPERSONATE:
		html = user_impersonate.NewUserImpersonateController(controller.registry).Handler(w, r)
	default:
		html = user_manager.NewUserManagerController(controller.registry).Handler(w, r)
	}

	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(html)); err != nil {
			if logger := controller.registry.GetLogger(); logger != nil {
				logger.Error("At usersAdminController > Handler", "write_error", err.Error())
			}
		}
	}
}
