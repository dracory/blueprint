package admin

import (
	"net/http"
	"path/filepath"
	"project/internal/layouts"
	"project/internal/registry"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/filesystem"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

const (
	JSON_ACTION_FILE_RENAME      = "file_rename"
	JSON_ACTION_FILE_DELETE      = "file_delete"
	JSON_ACTION_FILE_UPLOAD      = "file_upload"
	JSON_ACTION_DIRECTORY_CREATE = "directory_create"
	JSON_ACTION_DIRECTORY_DELETE = "directory_delete"
	JSON_ACTION_BULK_MOVE        = "bulk_move"
	JSON_ACTION_BULK_DELETE      = "bulk_delete"
	MAX_UPLOAD_SIZE              = 50 * 1024 * 1024 // 50MB
)

// FileManagerController handles file management operations
type FileManagerController struct {
	registry    registry.RegistryInterface
	rootDirPath string
	funcLayout  func(content string) string
	storage     filesystem.StorageInterface
}

// NewFileManagerController creates a new file manager controller
func NewFileManagerController(registry registry.RegistryInterface) *FileManagerController {
	cfg := registry.GetConfig()
	rootDirPath := strings.TrimSpace(cfg.GetMediaRoot())
	rootDirPath = strings.Trim(rootDirPath, "/")
	rootDirPath = strings.Trim(rootDirPath, ".")
	rootDirPath = "/" + rootDirPath

	return &FileManagerController{
		registry:    registry,
		rootDirPath: rootDirPath,
		storage:     registry.GetSqlFileStorage(),
	}
}

// Handler handles all file manager requests
func (c *FileManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	c.init(r)

	action := strings.TrimSpace(req.GetStringTrimmed(r, "action"))

	if lo.Contains([]string{
		JSON_ACTION_DIRECTORY_CREATE,
		JSON_ACTION_DIRECTORY_DELETE,
		JSON_ACTION_FILE_RENAME,
		JSON_ACTION_FILE_DELETE,
		JSON_ACTION_FILE_UPLOAD,
		JSON_ACTION_BULK_DELETE,
	}, action) {
		w.Header().Set("Content-Type", "application/json")
		return c.anyIndex(w, r)
	}

	return c.anyIndex(w, r)
}

// anyIndex routes to the appropriate action handler
func (c *FileManagerController) anyIndex(_ http.ResponseWriter, r *http.Request) string {
	action := strings.TrimSpace(req.GetStringTrimmed(r, "action"))

	switch action {
	case JSON_ACTION_FILE_RENAME:
		return c.fileRenameAjax(r)
	case JSON_ACTION_FILE_DELETE:
		return c.fileDeleteAjax(r)
	case JSON_ACTION_DIRECTORY_CREATE:
		return c.directoryCreateAjax(r)
	case JSON_ACTION_DIRECTORY_DELETE:
		return c.directoryDeleteAjax(r)
	case JSON_ACTION_FILE_UPLOAD:
		return c.fileUploadAjax(r)
	case JSON_ACTION_BULK_MOVE:
		return c.bulkMoveAjax(r)
	case JSON_ACTION_BULK_DELETE:
		return c.bulkDeleteAjax(r)
	default:
		return c.getMediaManager(r)
	}
}

// init initializes the controller by setting the layout function
func (controller *FileManagerController) init(r *http.Request) string {
	controller.funcLayout = func(content string) string {
		return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
			Title:   "File Manager",
			Content: hb.Raw(content),
		}).ToHTML()
	}
	return ""
}

// getMediaManager returns the main file manager UI
func (controller *FileManagerController) getMediaManager(r *http.Request) string {
	if controller.registry == nil {
		return api.Error("app is required").ToString()
	}

	cfg := controller.registry.GetConfig()
	if cfg == nil {
		return api.Error("config is required").ToString()
	}

	if !cfg.GetSqlFileStoreUsed() {
		return api.Error("SQL file store is not enabled").ToString()
	}

	if controller.storage == nil {
		return api.Error("storage is required").ToString()
	}

	currentDirectory := req.GetStringTrimmed(r, "current_dir")
	currentDirectory = strings.Trim(currentDirectory, "/")
	currentDirectory = strings.Trim(currentDirectory, ".")

	parentDirectory := ""
	if currentDirectory != "" {
		parentDirectory = filepath.Dir(currentDirectory)
	}

	parentDirectory = strings.Trim(parentDirectory, "/")
	parentDirectory = strings.Trim(parentDirectory, ".")

	if currentDirectory == "" {
		currentDirectory = controller.rootDirPath
	}

	directories, err := controller.storage.Directories(currentDirectory)
	if err != nil {
		return api.Error(err.Error()).ToString()
	}

	files, err := controller.storage.Files(currentDirectory)
	if err != nil {
		return api.Error(err.Error()).ToString()
	}

	directoryList := []FileEntry{}
	for _, dir := range directories {
		size, _ := controller.storage.Size(dir)
		hSize := lo.If(size > 0, controller.HumanFilesize(size)).Else("-")
		modified, _ := controller.storage.LastModified(dir)
		hModified := lo.If(lo.IsEmpty(modified), "-").Else(carbon.CreateFromStdTime(modified).ToDateTimeString())
		directoryList = append(directoryList, FileEntry{
			Path:              dir,
			Name:              filepath.Base(dir),
			Size:              size,
			SizeHuman:         hSize,
			LastModified:      modified,
			LastModifiedHuman: hModified,
		})
	}

	fileList := []FileEntry{}
	for _, file := range files {
		size, _ := controller.storage.Size(file)
		hSize := controller.HumanFilesize(size)
		modified, _ := controller.storage.LastModified(file)
		hModified := carbon.CreateFromStdTime(modified).ToDateTimeString()
		url, _ := controller.storage.Url(file)

		fileList = append(fileList, FileEntry{
			Path:              file,
			URL:               url,
			Name:              filepath.Base(file),
			Size:              size,
			SizeHuman:         hSize,
			LastModified:      modified,
			LastModifiedHuman: hModified,
		})
	}

	page := controller.uiManager(currentDirectory, parentDirectory, directoryList, fileList)

	if controller.funcLayout != nil {
		style := hb.StyleURL(cdn.BootstrapIconsCss_1_10_2()).ToHTML()
		script := hb.ScriptURL(cdn.Jquery_3_6_4()).ToHTML()
		script += hb.ScriptURL(cdn.Notify_0_4_2()).ToHTML()
		page = style + script + page
		return controller.funcLayout(page)
	}

	layout := uiLayout("File Manager", page)
	return layout
}
