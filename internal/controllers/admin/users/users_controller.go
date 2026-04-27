package users

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/useradmin"
)

// usersAdminController wraps the pkg/useradmin package for integration
type usersAdminController struct {
	registry registry.RegistryInterface
}

// NewUsersAdminController creates a new users admin controller
func NewUsersAdminController(registry registry.RegistryInterface) *usersAdminController {
	return &usersAdminController{registry: registry}
}

// Handler processes users admin requests
func (controller *usersAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetAuthUser(r)
	if user == nil {
		http.Redirect(w, r, links.Admin().Home(), http.StatusSeeOther)
		return
	}

	admin, err := useradmin.New(useradmin.AdminOptions{
		Registry:     controller.registry,
		AdminHomeURL: links.Admin().Home(),
		UserAdminURL: links.Admin().Users(),
		AuthUserID: func(r *http.Request) string {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
	})

	if err != nil {
		if logger := controller.registry.GetLogger(); logger != nil {
			logger.Error("At admin > usersAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.registry.GetLogger(); logger != nil {
				logger.Error("At admin > usersAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	html := admin.Handle(w, r)

	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(html)); err != nil {
			if logger := controller.registry.GetLogger(); logger != nil {
				logger.Error("At usersAdminController > Handler", "write_error", err.Error())
			}
		}
	}
}
