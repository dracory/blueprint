package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/links"

	basesession "github.com/dracory/base/session"

	logadmin "github.com/dracory/logadmin"
)

// logsAdminController wraps the external github.com/dracory/logadmin package
// for integration with the blueprint admin interface.
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
		Store:          controller.app.GetLogStore(),
		Logger:         controller.app.GetLogger(),
		FuncLayout:     adapters.NewLayoutFunc(controller.app),
		AdminHomeURL:   links.Admin().Home(),
		LogAdminURL:    links.Admin().Logs(),
		FileManagerURL: links.Admin().FileManager(),
		AuthUserID: func(r *http.Request) string {
			user := basesession.GetAuthUser(r)
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
