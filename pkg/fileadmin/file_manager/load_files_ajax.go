package file_manager

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// handleLoadFilesAjax returns directory contents as JSON for the Vue app
func (controller *FileManagerController) handleLoadFilesAjax(r *http.Request) string {
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
