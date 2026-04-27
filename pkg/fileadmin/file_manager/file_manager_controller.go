package file_manager

import (
	"embed"
	"net/http"
	"path/filepath"
	"project/internal/layouts"
	"project/internal/registry"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/filesystem"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"

	"project/pkg/fileadmin/shared"
)

//go:embed *.html
//go:embed *.js
var filesEmbed embed.FS

const (
	actionLoadFiles                   = "load-files"
	JSON_ACTION_FILE_CLONE            = "file_clone"
	JSON_ACTION_FILE_RENAME           = "file_rename"
	JSON_ACTION_FILE_DELETE           = "file_delete"
	JSON_ACTION_FILE_UPLOAD           = "file_upload"
	JSON_ACTION_DIRECTORY_CREATE      = "directory_create"
	JSON_ACTION_DIRECTORY_DELETE      = "directory_delete"
	JSON_ACTION_BULK_MOVE             = "bulk_move"
	JSON_ACTION_BULK_DELETE           = "bulk_delete"
	JSON_ACTION_GET_MOVE_DESTINATIONS = "get_move_destinations"
	MAX_UPLOAD_SIZE                   = 50 * 1024 * 1024 // 50MB
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

	if action == actionLoadFiles ||
		action == JSON_ACTION_DIRECTORY_CREATE ||
		action == JSON_ACTION_DIRECTORY_DELETE ||
		action == JSON_ACTION_FILE_CLONE ||
		action == JSON_ACTION_FILE_RENAME ||
		action == JSON_ACTION_FILE_DELETE ||
		action == JSON_ACTION_FILE_UPLOAD ||
		action == JSON_ACTION_BULK_MOVE ||
		action == JSON_ACTION_BULK_DELETE ||
		action == JSON_ACTION_GET_MOVE_DESTINATIONS {
		w.Header().Set("Content-Type", "application/json")
		return c.anyIndex(w, r)
	}

	return c.anyIndex(w, r)
}

// anyIndex routes to the appropriate action handler
func (c *FileManagerController) anyIndex(_ http.ResponseWriter, r *http.Request) string {
	action := strings.TrimSpace(req.GetStringTrimmed(r, "action"))

	switch action {
	case actionLoadFiles:
		return c.handleLoadFiles(r)
	case JSON_ACTION_FILE_CLONE:
		return c.fileCloneAjax(r)
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
	case JSON_ACTION_GET_MOVE_DESTINATIONS:
		return c.getMoveDestinationsAjax(r)
	default:
		return c.renderPage(r)
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

// renderPage renders the file manager Vue.js application
func (controller *FileManagerController) renderPage(r *http.Request) string {
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

	htmlContent, err := filesEmbed.ReadFile("files.html")
	if err != nil {
		return api.Error("Failed to read files HTML template: " + err.Error()).ToString()
	}

	jsContent, err := filesEmbed.ReadFile("files.js")
	if err != nil {
		return api.Error("Failed to read files JavaScript file: " + err.Error()).ToString()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlFileManager = '` + shared.NewLinks("/admin/file-manager").FileManager() + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(vueContainer)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:      "File Manager",
		Content:    content,
		ScriptURLs: []string{},
		Styles:     []string{},
	}).ToHTML()
}

// handleLoadFiles returns directory contents as JSON for the Vue app
func (controller *FileManagerController) handleLoadFiles(r *http.Request) string {
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

	return api.SuccessWithData("Files loaded successfully", map[string]any{
		"current_directory": currentDirectory,
		"parent_directory":  parentDirectory,
		"directories":       directoryList,
		"files":             fileList,
	}).ToString()
}
