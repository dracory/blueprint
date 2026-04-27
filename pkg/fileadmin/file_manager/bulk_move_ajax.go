package file_manager

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// bulkMoveAjax handles bulk file/folder move requests
func (c *FileManagerController) bulkMoveAjax(r *http.Request) string {
	if r == nil || r.URL == nil {
		return api.Error("invalid request").ToString()
	}

	currentDir := req.GetStringTrimmed(r, "current_dir")
	if currentDir == "" {
		return api.Error("current_dir is required").ToString()
	}

	destinationDir := req.GetStringTrimmed(r, "destination_dir")

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

		// Extract the filename/directory name from the path
		itemName := filepath.Base(item.Path)

		// Build the new path
		var newPath string
		if destinationDir == "" || destinationDir == "/" {
			newPath = "/" + itemName
		} else {
			newPath = strings.TrimRight(destinationDir, "/") + "/" + itemName
		}

		// Check if trying to move into itself (for directories)
		if item.Type == "directory" && strings.HasPrefix(destinationDir, item.Path) {
			errors = append(errors, "Cannot move directory into itself: "+itemName)
			continue
		}

		// Perform the move using storage.Move
		err := c.storage.Move(item.Path, newPath)
		if err != nil {
			log.Printf("Error moving %s to %s: %v", item.Path, newPath, err)
			errors = append(errors, "Failed to move "+itemName+": "+err.Error())
			continue
		}

		successCount++
	}

	// Return appropriate response
	if successCount == 0 {
		return api.Error("Failed to move items: " + strings.Join(errors, "; ")).ToString()
	}

	message := "Successfully moved " + string(rune('0'+successCount)) + " item(s)"
	if len(errors) > 0 {
		message += ". Some items failed: " + strings.Join(errors, "; ")
	}

	return api.Success(message).ToString()
}
