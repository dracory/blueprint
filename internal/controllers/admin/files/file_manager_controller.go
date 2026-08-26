package admin

import (
	"net/http"
	"strings"

	basesession "github.com/dracory/base/session"

	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/links"

	fileadmin "github.com/dracory/fileadmin"
)

// FileManagerController wraps the external github.com/dracory/fileadmin package
// for integration with the blueprint admin interface.
type FileManagerController struct {
	app app.AppInterface
}

// NewFileManagerController creates a new file manager controller
func NewFileManagerController(app app.AppInterface) *FileManagerController {
	return &FileManagerController{app: app}
}

// Handler processes file manager requests
func (c *FileManagerController) Handler(w http.ResponseWriter, r *http.Request) {
	cfg := c.app.GetConfig()

	// Derive root dir path from config (same logic as the old pkg/fileadmin)
	rootDirPath := strings.TrimSpace(cfg.GetMediaRoot())
	rootDirPath = strings.Trim(rootDirPath, "/")
	rootDirPath = strings.Trim(rootDirPath, ".")
	rootDirPath = "/" + rootDirPath

	admin, err := fileadmin.New(fileadmin.AdminOptions{
		Storage:      c.app.GetSqlFileStorage(),
		RootDirPath:  rootDirPath,
		FuncLayout:   adapters.NewLayoutFunc(c.app),
		AdminHomeURL: links.Admin().Home(),
		FileAdminURL: links.Admin().FileManager(),
		AuthUserID: func(r *http.Request) string {
			user := basesession.GetAuthUser(r)
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

	admin.Handle(w, r)
}
