package middlewares

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// NewUserMiddleware checks if the user is authenticated and active
// before allowing access to the protected route.
//
// Business logic:
//  1. user must be authenticated
//  2. user must be active
//  3. user must be registered
func NewUserMiddleware(app app.AppInterface) rtr.MiddlewareInterface {
	return rtrMiddleware.UserMiddleware(rtrMiddleware.UserMiddlewareConfig{
		GetUser: func(r *http.Request) rtrMiddleware.UserMiddlewareUser {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return nil
			}
			return user
		},
		RegistrationEnabled: app.GetConfig().GetRegistrationEnabled(),
		RegistrationPaths:   []string{links.USER_PROFILE, links.AUTH_REGISTER},
		OnNotAuthenticated: func(w http.ResponseWriter, r *http.Request) {
			returnURL := links.URL(r.URL.Path, map[string]string{})
			loginURL := links.Auth().Login(returnURL)
			helpers.ToFlashError(app.GetCacheStore(), w, r, "Only authenticated users can access this page", loginURL, 15)
		},
		OnNotActive: func(w http.ResponseWriter, r *http.Request) {
			homeURL := links.Website().Home()
			helpers.ToFlashError(app.GetCacheStore(), w, r, "User account not active", homeURL, 15)
		},
		OnRegistrationIncomplete: func(w http.ResponseWriter, r *http.Request) {
			registerURL := links.Auth().Register()
			helpers.ToFlashInfo(app.GetCacheStore(), w, r, "Please complete your registration to continue", registerURL, 15)
		},
	})
}
