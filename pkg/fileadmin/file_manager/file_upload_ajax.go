package file_manager

import (
	"log"
	"net/http"
	"os"

	"github.com/dracory/api"
	"github.com/dracory/base/files"
	"github.com/dracory/req"
)

// fileUploadAjax handles file upload requests
func (c *FileManagerController) fileUploadAjax(r *http.Request) string {
	if r.ContentLength > MAX_UPLOAD_SIZE {
		return api.Error("The uploaded image is too big. Please use an file less than 50MB in size").ToString()
	}

	currentDir := req.GetStringTrimmed(r, "current_dir")
	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("upload_file")
	if err != nil {
		return api.Error(err.Error()).ToString()
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close uploaded file: %v", err)
		}
	}()

	filePath, errSave := files.SaveToTempDir(fileHeader.Filename, file)
	if errSave != nil {
		log.Println(errSave.Error())
		return api.Error(errSave.Error()).ToString()
	}
	defer func() {
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: failed to remove temp file %s: %v", filePath, err)
		}
	}()

	remoteFilePath := currentDir + "/" + fileHeader.Filename

	data, err := os.ReadFile(filePath)
	if err != nil {
		return api.Error(err.Error()).ToString()
	}

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	err = c.storage.Put(remoteFilePath, data)

	if err != nil {
		return api.Error(err.Error()).ToString()
	}

	return api.Success("File uploaded successfully").ToString()
}
