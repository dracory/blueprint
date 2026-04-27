package file_manager

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// directoryCreateAjax handles directory creation requests
func (c *FileManagerController) directoryCreateAjax(r *http.Request) string {
	newDirName := req.GetStringTrimmed(r, "create_dir")

	if newDirName == "" {
		return api.Error("create_dir is required").ToString()
	}

	currentDir := req.GetStringTrimmed(r, "current_dir")

	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	if currentDir == "/" {
		currentDir = "" // to prevent double slashes
	}

	dirPath := currentDir + "/" + newDirName
	dirPath = strings.ReplaceAll(dirPath, "//", "/") // remove double slashes
	dirPath = strings.TrimRight(dirPath, "/")        // remove trailing slashes

	if dirPath == "" || dirPath == "/" {
		return api.Error("root directory can not be created").ToString()
	}

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	err := c.storage.MakeDirectory(dirPath)

	if err == nil {
		return api.Success("directory created successfully").ToString()
	}

	return api.Error(err.Error()).ToString()
}
