package middlewares

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
	"github.com/dracory/userstore"
)

// adminUserAdapter wraps userstore.UserInterface to satisfy rtrMiddleware.UserWithRole
type adminUserAdapter struct {
	user userstore.UserInterface
}

func (a *adminUserAdapter) IsActive() bool {
	if a.user == nil {
		return false
	}
	return a.user.IsActive()
}

func (a *adminUserAdapter) IsRegistrationCompleted() bool {
	if a.user == nil {
		return false
	}
	return a.user.IsRegistrationCompleted()
}

func (a *adminUserAdapter) HasRole(role string) bool {
	if a.user == nil {
		return false
	}
	switch role {
	case userstore.USER_ROLE_ADMINISTRATOR:
		return a.user.IsAdministrator()
	case userstore.USER_ROLE_SUPERUSER:
		return a.user.IsSuperuser()
	default:
		return false
	}
}

// NewAdminMiddleware checks if the user is an administrator or superuser
// before allowing access to the protected route.
//
// Business logic:
//  1. user must be authenticated
//  2. user must be active
//  3. user must be registered
//  4. user must be an admin or superuser
func NewAdminMiddleware(app app.AppInterface) rtr.MiddlewareInterface {
	return rtrMiddleware.UserMiddleware(rtrMiddleware.UserMiddlewareConfig{
		GetUser: func(r *http.Request) rtrMiddleware.UserMiddlewareUser {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return nil
			}
			return &adminUserAdapter{user: user}
		},
		RegistrationEnabled: true,
		OnNotAuthenticated: func(w http.ResponseWriter, r *http.Request) {
			returnURL := links.URL(r.URL.Path, map[string]string{})
			loginURL := links.Auth().Login(returnURL)
			helpers.ToFlashError(app.GetCacheStore(), w, r, "You must be logged in to access this page", loginURL, 15)
		},
		OnRegistrationIncomplete: func(w http.ResponseWriter, r *http.Request) {
			registerURL := links.Auth().Register(map[string]string{})
			helpers.ToFlashInfo(app.GetCacheStore(), w, r, "Please complete your registration to continue", registerURL, 15)
		},
		OnNotActive: func(w http.ResponseWriter, r *http.Request) {
			homeURL := links.Website().Home()
			helpers.ToFlash(app.GetCacheStore(), w, r, "error", "Your account is not active", homeURL, 15)
		},
		RequireRoles: []string{userstore.USER_ROLE_ADMINISTRATOR, userstore.USER_ROLE_SUPERUSER},
		OnNotAuthorized: func(w http.ResponseWriter, r *http.Request) {
			homeURL := links.Website().Home()
			helpers.ToFlash(app.GetCacheStore(), w, r, "error", "You must be an administrator to access this page", homeURL, 15)
		},
	})
}
