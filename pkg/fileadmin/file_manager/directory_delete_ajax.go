package file_manager

import (
	"net/http"

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

	dirPath, err := verifyAndNormalizeDirPath(currentDir, selectedDirName)
	if err != nil {
		return api.Error("invalid directory name: " + err.Error()).ToString()
	}
	cfmt.Infoln("Deleting directory:", dirPath)

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
