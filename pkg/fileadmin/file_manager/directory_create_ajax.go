package file_manager

import (
	"net/http"

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

	dirPath, err := verifyAndNormalizeDirPath(currentDir, newDirName)
	if err != nil {
		return api.Error("invalid directory name: " + err.Error()).ToString()
	}

	if dirPath == "" || dirPath == "/" {
		return api.Error("root directory can not be created").ToString()
	}

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	err = c.storage.MakeDirectory(dirPath)

	if err == nil {
		return api.Success("directory created successfully").ToString()
	}

	return api.Error(err.Error()).ToString()
}
