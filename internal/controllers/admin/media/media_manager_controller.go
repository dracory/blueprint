package admin

import (
	"net/http"
	"strings"

	basesession "github.com/dracory/base/session"
	"github.com/dracory/filesystem"

	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/links"

	fileadmin "github.com/dracory/fileadmin"
)

// MediaManagerController wraps the external github.com/dracory/fileadmin
// package, backed by the configured S3 media storage, for integration with
// the blueprint admin interface.
type MediaManagerController struct {
	app app.AppInterface
}

// NewMediaManagerController creates a new media manager controller
func NewMediaManagerController(app app.AppInterface) *MediaManagerController {
	return &MediaManagerController{app: app}
}

// Handler processes media manager requests
func (c *MediaManagerController) Handler(w http.ResponseWriter, r *http.Request) {
	cfg := c.app.GetConfig()

	// Derive root dir path from config (same logic as the file manager)
	rootDirPath := strings.TrimSpace(cfg.GetMediaRoot())
	rootDirPath = strings.Trim(rootDirPath, "/")
	rootDirPath = strings.Trim(rootDirPath, ".")
	rootDirPath = "/" + rootDirPath

	storage, err := filesystem.NewStorage(filesystem.Disk{
		DiskName:             "S3",
		Driver:               filesystem.DRIVER_S3,
		Url:                  cfg.GetMediaUrl(),
		Region:               cfg.GetMediaRegion(),
		Key:                  cfg.GetMediaKey(),
		Secret:               cfg.GetMediaSecret(),
		Bucket:               cfg.GetMediaBucket(),
		UsePathStyleEndpoint: true,
	})

	if err != nil {
		if logger := c.app.GetLogger(); logger != nil {
			logger.Error("At admin > MediaManagerController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := c.app.GetLogger(); logger != nil {
				logger.Error("At admin > MediaManagerController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin, err := fileadmin.New(fileadmin.AdminOptions{
		Storage:      storage,
		RootDirPath:  rootDirPath,
		FuncLayout:   adapters.NewLayoutFunc(c.app),
		AdminHomeURL: links.Admin().Home(),
		FileAdminURL: links.Admin().MediaManager(),
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
			logger.Error("At admin > MediaManagerController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := c.app.GetLogger(); logger != nil {
				logger.Error("At admin > MediaManagerController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
