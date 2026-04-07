package admin

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/base/cfmt"
	"github.com/dracory/req"
)

// directoryDeleteAjax handles directory deletion requests
func (c *FileManagerController) directoryDeleteAjax(r *http.Request) string {
	selectedDirName := req.GetStringTrimmed(r, "delete_dir")

	if selectedDirName == "" {
		return api.Error("delete_dir is required").ToString()
	}

	currentDir := req.GetStringTrimmed(r, "current_dir")

	if currentDir == "." || currentDir == ".." {
		return api.Error("current_dir is required").ToString()
	}

	if currentDir == "/" {
		currentDir = "" // eliminate double slashes
	}

	dirPath := currentDir + "/" + selectedDirName
	cfmt.Infoln("Deleting directory:", dirPath)
	dirPath = strings.ReplaceAll(dirPath, "//", "/") // remove double slashes
	dirPath = strings.TrimRight(dirPath, "/")        // remove trailing slashes

	if dirPath == "" || dirPath == "/" {
		return api.Error("root directory can not be deleted").ToString()
	}

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	errDeleted := c.storage.DeleteDirectory(dirPath)

	if errDeleted == nil {
		return api.Success("directory deleted successfully").ToString()
	}

	return api.Error(errDeleted.Error()).ToString()
}
