package admin

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/base/req"
)

// == CONTROLLER ==============================================================

// userImpersonateController represents a controller for handling user impersonation.
type userImpersonateController struct {
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewUserImpersonateController(app types.AppInterface) *userImpersonateController {
	return &userImpersonateController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (c *userImpersonateController) Handler(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "User not found", links.NewAdminLinks().Users(map[string]string{}), 15)
	}

	if !authUser.IsAdministrator() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Not authorized", links.NewAdminLinks().Users(map[string]string{}), 15)
	}

	userID := req.Value(r, "user_id")

	if userID == "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "User ID not found", links.NewAdminLinks().Users(map[string]string{}), 15)
	}

	err := Impersonate(c.app.GetSessionStore(), w, r, userID)

	if err != nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, err.Error(), links.NewAdminLinks().Users(map[string]string{}), 15)
	}

	return helpers.ToFlashSuccess(c.app.GetCacheStore(), w, r, "Impersonation is successful", links.NewUserLinks().Home(map[string]string{}), 15)
}
