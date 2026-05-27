package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetMoveDestinationsAjax_MissingSelectedItemsParameter(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("current_dir", "/uploads")
	form.Add("selected_items", "")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.getMoveDestinationsAjax(req)

	if !strings.Contains(result, "No items selected") {
		t.Errorf("getMoveDestinationsAjax() result = %q, want to contain %q", result, "No items selected")
	}
}

func TestGetMoveDestinationsAjax_InvalidJSON(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("current_dir", "/uploads")
	form.Add("selected_items", "invalid json")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.getMoveDestinationsAjax(req)

	if !strings.Contains(result, "Invalid selected items data") {
		t.Errorf("getMoveDestinationsAjax() result = %q, want to contain %q", result, "Invalid selected items data")
	}
}
