package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestHandleLoadFilesAjax_NilStorage(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)
	controller.storage = nil

	form := url.Values{}
	form.Add("current_dir", "/uploads")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.handleLoadFilesAjax(req)

	if !strings.Contains(result, "storage is required") {
		t.Errorf("handleLoadFilesAjax() result = %q, want to contain %q", result, "storage is required")
	}
}

func TestHandleLoadFilesAjax_LoadsFromRoot(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	// Create the root directory in the test storage so that listing works
	if err := controller.storage.MakeDirectory(controller.rootDirPath); err != nil {
		t.Fatalf("Failed to create root directory: %v", err)
	}

	form := url.Values{}
	form.Add("current_dir", "")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.handleLoadFilesAjax(req)

	if !strings.Contains(result, "Files loaded successfully") {
		t.Errorf("handleLoadFilesAjax() result = %q, want to contain %q", result, "Files loaded successfully")
	}
}
