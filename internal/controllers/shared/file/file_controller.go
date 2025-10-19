package file

import (
	"net/http"
	"strings"

	"github.com/dracory/str"
	"github.com/gouniverse/filesystem"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONTROLLER ==============================================================

type fileController struct {
	storage filesystem.StorageInterface
}

// == CONSTRUCTOR =============================================================

func NewFileController(storage filesystem.StorageInterface) *fileController {
	return &fileController{storage: storage}
}

// == PUBLIC METHODS ==========================================================

func (c *fileController) Handler(w http.ResponseWriter, r *http.Request) string {
	if c.storage == nil {
		return "File storage not configured"
	}

	filePath := lo.IfF(strings.HasPrefix(r.URL.Path, "/files"), func() string { return str.RightFrom(r.URL.Path, "/files") }).
		ElseIfF(strings.HasPrefix(r.URL.Path, "/file"), func() string { return str.RightFrom(r.URL.Path, "/file") }).
		ElseIfF(strings.HasPrefix(r.URL.Path, "/media"), func() string { return str.RightFrom(r.URL.Path, "/media") }).
		Else(r.URL.Path)

	exists, err := c.storage.Exists(filePath)

	if err != nil {
		return err.Error()
	}

	if !exists {
		return "File not found"
	}

	content, err := c.storage.ReadFile(filePath)

	if err != nil {
		return err.Error()
	}

	extension := c.findExtension(filePath)
	mimeType := c.findMIMEType(extension)

	if extension == "" {
		return "File not found"
	}

	w.Header().Set("Content-Type", mimeType)

	if mimeType == "application/octet-stream" {
		w.Header().Set("Content-Disposition", "attachment; filename="+r.URL.Path)
		w.Header().Set("Content-Length", cast.ToString(len(content)))
	}

	if _, err := w.Write(content); err != nil {
		return "Failed to write file content: " + err.Error()
	}

	return ""
}

// findExtension finds the file extension from a path.
//
// Parameter(s):
//   - path string - the path
//
// Return type(s):
//   - string - the file extension
func (controller fileController) findExtension(path string) string {
	fileName := controller.findFileName(path)

	if fileName == "" {
		return ""
	}

	nameParts := strings.Split(fileName, ".")

	if len(nameParts) < 2 {
		return ""
	}

	return nameParts[1]
}

func (controller fileController) findMIMEType(extension string) string {
	switch extension {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "svg":
		return "image/svg+xml"
	case "ico":
		return "image/x-icon"
	case "pdf":
		return "application/pdf"
	case "zip":
		return "application/zip"
	case "mp3":
		return "audio/mpeg"
	case "webm":
		return "video/webm"
	default:
		return "application/octet-stream"
	}
}

// findFileName finds the file name from a path.
//
// Parameter(s):
//   - path string - the path
//
// Return type(s):
//   - string - the file name
func (controller fileController) findFileName(path string) string {
	uriParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(uriParts) < 1 {
		return ""
	}

	return uriParts[len(uriParts)-1]
}
