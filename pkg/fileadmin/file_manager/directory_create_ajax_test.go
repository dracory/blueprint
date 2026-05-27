package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDirectoryCreateAjax_MissingCreateDirParameter(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("create_dir", "")
	form.Add("current_dir", "/uploads")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.directoryCreateAjax(req)

	if !strings.Contains(result, "create_dir is required") {
		t.Errorf("directoryCreateAjax() result = %q, want to contain %q", result, "create_dir is required")
	}
}
