package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	logadmin "project/pkg/logadmin"
)

// logsAdminController wraps the pkg/logadmin package for integration
type logsAdminController struct {
	app app.AppInterface
}

// NewLogsAdminController creates a new logs admin controller
func NewLogsAdminController(app app.AppInterface) *logsAdminController {
	return &logsAdminController{app: app}
}

// Handler processes logs admin requests
func (controller *logsAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := logadmin.New(logadmin.AdminOptions{
		Registry:     controller.app,
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
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At admin > logsAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.app.GetLogger(); logger != nil {
				logger.Error("At admin > logsAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
