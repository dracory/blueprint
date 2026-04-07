package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// bulkDeleteAjax handles bulk file/folder delete requests
func (c *FileManagerController) bulkDeleteAjax(r *http.Request) string {
	currentDir := req.GetStringTrimmed(r, "current_dir")
	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	// Parse selected items JSON
	selectedItemsJSON := req.GetStringTrimmed(r, "selected_items")
	if selectedItemsJSON == "" {
		return api.Error("No items selected").ToString()
	}

	var selectedItems []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal([]byte(selectedItemsJSON), &selectedItems); err != nil {
		log.Printf("Error parsing selected items JSON: %v", err)
		return api.Error("Invalid selected items data").ToString()
	}

	if len(selectedItems) == 0 {
		return api.Error("No items selected").ToString()
	}

	if c.storage == nil {
		return api.Error("Storage not initialized").ToString()
	}

	// Track success and failures
	successCount := 0
	var errors []string

	for _, item := range selectedItems {
		if item.Path == "" {
			continue
		}

		itemName := filepath.Base(item.Path)

		var err error
		if item.Type == "directory" {
			err = c.storage.DeleteDirectory(item.Path)
		} else {
			err = c.storage.DeleteFile([]string{item.Path})
		}

		if err != nil {
			log.Printf("Error deleting %s: %v", item.Path, err)
			errors = append(errors, "Failed to delete "+itemName+": "+err.Error())
			continue
		}

		successCount++
	}

	// Return appropriate response
	if successCount == 0 {
		return api.Error("Failed to delete items: " + strings.Join(errors, "; ")).ToString()
	}

	message := "Successfully deleted " + strconv.Itoa(successCount) + " item(s)"
	if len(errors) > 0 {
		message += ". Some items failed: " + strings.Join(errors, "; ")
	}

	return api.Success(message).ToString()
}
