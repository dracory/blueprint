package auth

import (
	"net/http"
	"project/app/links"
	"project/internal/helpers"

	"github.com/gouniverse/auth"
)

type logoutController struct{}

func NewLogoutController() *logoutController {
	return &logoutController{}
}

func (controller *logoutController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	auth.AuthCookieRemove(w, r)

	return helpers.ToFlashSuccess(w, r, "You have been logged out successfully", links.Website().Home(), 5)
}
