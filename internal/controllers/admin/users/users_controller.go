package users

import (
	"net/http"

	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/tasks/constants"

	useradmin "github.com/dracory/useradmin"
)

// usersAdminController wraps the external github.com/dracory/useradmin
// package for integration with the blueprint admin interface.
type usersAdminController struct {
	app app.AppInterface
}

// NewUsersAdminController creates a new users admin controller
func NewUsersAdminController(app app.AppInterface) *usersAdminController {
	return &usersAdminController{app: app}
}

// Handler processes users admin requests
func (controller *usersAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetAuthUser(r)
	if user == nil {
		http.Redirect(w, r, links.Admin().Home(), http.StatusSeeOther)
		return
	}

	// Secure cookie is only false in development (HTTP). In production
	// (HTTPS) the impersonation cookie must be marked Secure.
	secure := true
	if controller.app.GetConfig() != nil && controller.app.GetConfig().IsEnvDevelopment() {
		secure = false
	}

	admin, err := useradmin.New(useradmin.AdminOptions{
		UserStore:       controller.app.GetUserStore(),
		GeoResolver:     adapters.NewGeoResolver(controller.app.GetGeoStore()),
		Logger:          controller.app.GetLogger(),
		SessionResolver: adapters.NewSessionResolver(controller.app.GetSessionStore()),
		OnUserSearch:    adapters.NewOnUserSearchFunc(controller.app),
		OnUserUpdate:    adapters.NewOnUserUpdateFunc(controller.app, constants.BlindIndexRebuildTaskAlias),
		VaultTokenizer:  adapters.NewVaultTokenizerAdapter(controller.app),
		FuncLayout:      adapters.NewLayoutFunc(controller.app),
		FlashRedirect:   adapters.NewFlashRedirectFunc(controller.app),
		SecureCookie:    secure,
		AdminHomeURL:    links.Admin().Home(),
		UserAdminURL:    links.Admin().Users(),
		UserHomeURL:     links.User().Home(),
		AuthUserID: func(r *http.Request) string {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
		AuthUser: helpers.GetAuthUser,
	})

	if err != nil {
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At admin > usersAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.app.GetLogger(); logger != nil {
				logger.Error("At admin > usersAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
