package admin

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/req"
)

// == CONTROLLER ==============================================================

// userImpersonateController represents a controller for handling user impersonation.
type userImpersonateController struct {
	app types.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewUserImpersonateController(app types.RegistryInterface) *userImpersonateController {
	return &userImpersonateController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (c *userImpersonateController) Handler(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "User not found", links.Admin().Users(), 15)
	}

	if !authUser.IsAdministrator() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Not authorized", links.Admin().Users(), 15)
	}

	userID := req.GetStringTrimmed(r, "user_id")

	if userID == "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "User ID not found", links.Admin().Users(), 15)
	}

	err := Impersonate(c.app.GetSessionStore(), w, r, userID)

	if err != nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, err.Error(), links.Admin().Users(), 15)
	}

	return helpers.ToFlashSuccess(c.app.GetCacheStore(), w, r, "Impersonation is successful", links.User().Home(), 15)
}
