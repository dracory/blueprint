package file_manager

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// fileCloneAjax handles file clone (duplicate) requests
func (c *FileManagerController) fileCloneAjax(r *http.Request) string {
	selectedFileName := req.GetStringTrimmed(r, "clone_file")
	if selectedFileName == "" {
		return api.Error("clone_file is required").ToString()
	}

	currentDir := req.GetStringTrimmed(r, "current_dir")
	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	if currentDir == "/" {
		currentDir = ""
	}

	filePath := currentDir + "/" + selectedFileName

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	// Use provided new filename or auto-generate one
	newFileName := req.GetStringTrimmed(r, "new_file")
	if newFileName == "" {
		ext := filepath.Ext(selectedFileName)
		base := strings.TrimSuffix(selectedFileName, ext)
		newFileName = base + "_copy" + ext
	}
	newFilePath := currentDir + "/" + newFileName

	// Check if target already exists and append number if needed
	exists, _ := c.storage.Exists(newFilePath)
	counter := 2
	ext := filepath.Ext(selectedFileName)
	base := strings.TrimSuffix(newFileName, ext)
	for exists {
		newFileName = fmt.Sprintf("%s_copy_%d%s", base, counter, ext)
		newFilePath = currentDir + "/" + newFileName
		exists, _ = c.storage.Exists(newFilePath)
		counter++
	}

	errClone := c.storage.Copy(filePath, newFilePath)

	if errClone == nil {
		return api.SuccessWithData("file cloned successfully", map[string]any{
			"new_file_name": newFileName,
		}).ToString()
	}

	return api.Error(errClone.Error()).ToString()
}
