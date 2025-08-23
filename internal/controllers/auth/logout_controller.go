package auth

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/gouniverse/auth"
)

type logoutController struct {
	app types.AppInterface
}

func NewLogoutController(app types.AppInterface) *logoutController {
	return &logoutController{app: app}
}

func (controller *logoutController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	auth.AuthCookieRemove(w, r)

	return helpers.ToFlashSuccess(controller.app.GetCacheStore(), w, r, "You have been logged out successfully", links.Website().Home(), 5)
}
