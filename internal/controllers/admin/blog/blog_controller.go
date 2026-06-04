package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	blogadmin "project/pkg/blogadmin"
)

// blogAdminController wraps the pkg/blogadmin package for integration
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
		Store:        controller.app.GetBlogStore(),
		AdminHomeURL: links.Admin().Home(),
		BlogAdminURL: links.Admin().Blog(),
		AuthUserID: func(r *http.Request) string {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
		LLMEngine: nil, // TODO: Set up LLM engine if available
		BlogTopic: "",  // TODO: Get from settings
		Registry:  controller.app,
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
