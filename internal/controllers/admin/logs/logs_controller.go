package admin

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"

	logadmin "project/pkg/logadmin"
)

// logsAdminController wraps the pkg/logadmin package for integration
type logsAdminController struct {
	registry registry.RegistryInterface
}

// NewLogsAdminController creates a new logs admin controller
func NewLogsAdminController(registry registry.RegistryInterface) *logsAdminController {
	return &logsAdminController{registry: registry}
}

// Handler processes logs admin requests
func (controller *logsAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := logadmin.New(logadmin.AdminOptions{
		Registry:     controller.registry,
		AdminHomeURL: links.Admin().Home(),
		LogAdminURL:  links.Admin().Logs(),
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
			logger.Error("At admin > logsAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.registry.GetLogger(); logger != nil {
				logger.Error("At admin > logsAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
