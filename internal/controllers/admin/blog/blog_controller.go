package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/links"

	basesession "github.com/dracory/base/session"

	blogadmin "github.com/dracory/blogadmin"
)

// blogAdminController wraps the external github.com/dracory/blogadmin package
// for integration with the blueprint admin interface.
type blogAdminController struct {
	app app.AppInterface
}

// NewBlogAdminController creates a new blog admin controller
func NewBlogAdminController(app app.AppInterface) *blogAdminController {
	return &blogAdminController{app: app}
}

// Handler processes blog admin requests
func (controller *blogAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := blogadmin.New(blogadmin.AdminOptions{
		Store:          controller.app.GetBlogStore(),
		Logger:         controller.app.GetLogger(),
		CustomStore:    controller.app.GetCustomStore(),
		SettingStore:   controller.app.GetSettingStore(),
		LlmFactory:     adapters.NewLlmFactory(controller.app),
		FuncLayout:     adapters.NewLayoutFunc(controller.app),
		AdminHomeURL:   links.Admin().Home(),
		BlogAdminURL:   links.Admin().Blog(),
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
			logger.Error("At admin > blogAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.app.GetLogger(); logger != nil {
				logger.Error("At admin > blogAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
