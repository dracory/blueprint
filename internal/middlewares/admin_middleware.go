package middlewares

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

// NewAdminMiddleware checks if the user is an administrator or superuser
// before allowing access to the protected route.
//
// Business logic:
//  1. user must be authenticated
//  2. user must be active
//  3. user must be registered
//  4. user must be an admin or superuser
func NewAdminMiddleware(registry registry.RegistryInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Admin Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Admin validation logic here. CHange with your own

				authUser := helpers.GetAuthUser(r)

				// Check if user is authenticated? No => redirect to login
				if authUser == nil {
					returnURL := links.URL(r.URL.Path, map[string]string{})
					loginURL := links.Auth().Login(returnURL)
					helpers.ToFlashError(registry.GetCacheStore(), w, r, "You must be logged in to access this page", loginURL, 15)
					return
				}

				if !authUser.IsRegistrationCompleted() {
					registerURL := links.Auth().Register(map[string]string{})
					helpers.ToFlashInfo(registry.GetCacheStore(), w, r, "Please complete your registration to continue", registerURL, 15)
					return
				}

				// Check if user is active? No => redirect to website home
				if !authUser.IsActive() {
					homeURL := links.Website().Home()
					helpers.ToFlash(registry.GetCacheStore(), w, r, "error", "Your account is not active", homeURL, 15)
					return
				}

				// Check if user is an admin? No => redirect to website home
				if !authUser.IsAdministrator() && !authUser.IsSuperuser() {
					homeURL := links.Website().Home()
					helpers.ToFlash(registry.GetCacheStore(), w, r, "error", "You must be an administrator to access this page", homeURL, 15)
					return
				}

				next.ServeHTTP(w, r)
			})
		})
}
