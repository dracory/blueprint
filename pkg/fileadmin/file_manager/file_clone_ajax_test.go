package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileCloneAjax_MissingCloneFileParameter(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("clone_file", "")
	form.Add("current_dir", "/uploads")
	form.Add("new_file", "")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.fileCloneAjax(req)

	if !strings.Contains(result, "clone_file is required") {
		t.Errorf("fileCloneAjax() result = %q, want to contain %q", result, "clone_file is required")
	}
}
