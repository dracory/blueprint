package logout

import (
	"log/slog"
	"net/http"
	"project/internal/app"
	"project/internal/config"
	"project/internal/controllers/auth/shared"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/auth"
)

type logoutController struct {
	app app.AppInterface
}

func NewLogoutController(app app.AppInterface) *logoutController {
	return &logoutController{app: app}
}

func (controller *logoutController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	// Delete the server-side auth session — removing the cookie alone would
	// leave the session valid for anyone still holding the key.
	controller.deleteSessionByCookie(w, r, auth.CookieName, false)

	// Delete the remember session and expire its cookie so a remembered
	// device is fully logged out.
	controller.deleteSessionByCookie(w, r, config.COOKIE_NAME_REMEMBER_TOKEN, true)
	shared.RemoveRememberCookie(controller.app, w, r)

	auth.AuthCookieRemove(w, r)

	return helpers.ToFlashSuccess(controller.app.GetCacheStore(), w, r, "You have been logged out successfully", links.Website().Home(), 5)
}

// deleteSessionByCookie looks up the session named by cookieName and deletes
// it from the session store. Failures are logged, not surfaced — logout must
// always proceed to clear the client-side cookies. When rememberOnly is set,
// only sessions marked as remember sessions are deleted.
func (controller *logoutController) deleteSessionByCookie(w http.ResponseWriter, r *http.Request, cookieName string, rememberOnly bool) {
	if controller.app.IsDisabledSessionStore() {
		return
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	session, err := controller.app.GetSessionStore().SessionFindByKey(r.Context(), cookie.Value)
	if err != nil || session == nil {
		return
	}

	if rememberOnly && session.GetValue() != shared.RememberSessionValue {
		return
	}

	if err := controller.app.GetSessionStore().SessionDelete(r.Context(), session); err != nil {
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At Logout > failed to delete session", slog.String("error", err.Error()))
		}
	}
}
