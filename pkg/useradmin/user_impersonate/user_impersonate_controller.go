package user_impersonate

import (
	"net/http"

	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/req"
)

// == CONTROLLER ==============================================================

// userImpersonateController represents a controller for handling user impersonation.
type userImpersonateController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewUserImpersonateController(app app.AppInterface) *userImpersonateController {
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

	// In development (HTTP), use insecure cookie so the browser sends it
	// back over plain HTTP. In production (HTTPS), use Secure cookie.
	secure := true
	if c.app.GetConfig() != nil && c.app.GetConfig().IsEnvDevelopment() {
		secure = false
	}

	err := Impersonate(c.app.GetSessionStore(), w, r, userID, secure)

	if err != nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, err.Error(), links.Admin().Users(), 15)
	}

	return helpers.ToFlashSuccess(c.app.GetCacheStore(), w, r, "Impersonation is successful", links.User().Home(), 15)
}
