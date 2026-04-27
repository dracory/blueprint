package admin

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/fileadmin"
)

// FileManagerController wraps the pkg/fileadmin package
type FileManagerController struct {
	registry registry.RegistryInterface
}

// NewFileManagerController creates a new file manager controller
func NewFileManagerController(registry registry.RegistryInterface) *FileManagerController {
	return &FileManagerController{registry: registry}
}

// Handler processes file manager requests
func (c *FileManagerController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := fileadmin.New(fileadmin.AdminOptions{
		Registry:     c.registry,
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
		if logger := c.registry.GetLogger(); logger != nil {
			logger.Error("At admin > FileManagerController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := c.registry.GetLogger(); logger != nil {
				logger.Error("At admin > FileManagerController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	html := admin.Handle(w, r)

	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(html)); err != nil {
			if logger := c.registry.GetLogger(); logger != nil {
				logger.Error("At FileManagerController > Handler", "write_error", err.Error())
			}
		}
	}
}
