package admin

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// fileRenameAjax handles file rename requests
func (c *FileManagerController) fileRenameAjax(r *http.Request) string {
	currentFileName := req.GetStringTrimmed(r, "rename_file")
	if currentFileName == "" {
		return api.Error("rename_file is required").ToString()
	}

	newFileName := req.GetStringTrimmed(r, "new_file")

	if newFileName == "" {
		return api.Error("new_file is required").ToString()
	}
	currentDir := req.GetStringTrimmed(r, "current_dir")

	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	if currentDir == "/" {
		currentDir = "" // eliminate double slashes
	}

	oldFilePath := currentDir + "/" + currentFileName
	newFilePath := currentDir + "/" + newFileName

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	err := c.storage.Move(oldFilePath, newFilePath)

	if err == nil {
		return api.Success("file renamed successfully").ToString()
	}

	return api.Error(err.Error()).ToString()
}
