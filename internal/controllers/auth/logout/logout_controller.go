package logout

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/auth"
)

type logoutController struct {
	app app.AppInterface
}

func NewLogoutController(app app.AppInterface) *logoutController {
	return &logoutController{app: app}
}

func (controller *logoutController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	auth.AuthCookieRemove(w, r)

	return helpers.ToFlashSuccess(controller.app.GetCacheStore(), w, r, "You have been logged out successfully", links.Website().Home(), 5)
}
