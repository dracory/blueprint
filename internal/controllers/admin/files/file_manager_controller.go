package admin

import (
	"net/http"

	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"
	"project/pkg/fileadmin"
)

// FileManagerController wraps the pkg/fileadmin package
type FileManagerController struct {
	app app.AppInterface
}

// NewFileManagerController creates a new file manager controller
func NewFileManagerController(app app.AppInterface) *FileManagerController {
	return &FileManagerController{app: app}
}

// Handler processes file manager requests
func (c *FileManagerController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := fileadmin.New(fileadmin.AdminOptions{
		Registry:     c.app,
		AdminHomeURL: links.Admin().Home(),
		FileAdminURL: links.Admin().FileManager(),
		AuthUserID: func(r *http.Request) string {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
	})

	if err != nil {
		if logger := c.app.GetLogger(); logger != nil {
			logger.Error("At admin > FileManagerController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := c.app.GetLogger(); logger != nil {
				logger.Error("At admin > FileManagerController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	html := admin.Handle(w, r)

	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(html)); err != nil {
			if logger := c.app.GetLogger(); logger != nil {
				logger.Error("At FileManagerController > Handler", "write_error", err.Error())
			}
		}
	}
}
