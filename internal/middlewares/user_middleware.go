package middlewares

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"strings"

	"github.com/dracory/rtr"
)

// NewUserMiddleware checks if the user is authenticated and active
// before allowing access to the protected route.
//
// Business logic:
//  1. user must be authenticated
//  2. user must be active
//  3. user must be registered
func NewUserMiddleware(registry registry.RegistryInterface) rtr.MiddlewareInterface {
	m := rtr.NewMiddleware().
		SetName("User Middleware").
		SetHandler(userMiddlewareHandler(registry))

	return m
}

func userMiddlewareHandler(registry registry.RegistryInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			returnURL := links.URL(r.URL.Path, map[string]string{})
			loginURL := links.Auth().Login(returnURL)
			homeURL := links.Website().Home()
			registrationEnabled := registry.GetConfig().GetRegistrationEnabled()
			registerURL := links.Auth().Register()

			// User validation logic here. Change with your own

			authUser := helpers.GetAuthUser(r)

			// Check if user is authenticated? No => redirect to login
			if authUser == nil {
				helpers.ToFlashError(registry.GetCacheStore(), w, r, "Only authenticated users can access this page", loginURL, 15)
				return
			}

			// Check if user is active? No => redirect to website home
			if !authUser.IsActive() {
				helpers.ToFlashError(registry.GetCacheStore(), w, r, "User account not active", homeURL, 15)
				return
			}

			// Check if user has completed registration? No => redirect to profile to complete registration
			notOnProfilePage := strings.Trim(r.URL.Path, "/") != strings.Trim(links.USER_PROFILE, "/") &&
				strings.Trim(r.URL.Path, "/") != strings.Trim(links.AUTH_REGISTER, "/")

			if !authUser.IsRegistrationCompleted() && notOnProfilePage {
				if registrationEnabled {
					helpers.ToFlashInfo(registry.GetCacheStore(), w, r, "Please complete your registration to continue", registerURL, 15)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
