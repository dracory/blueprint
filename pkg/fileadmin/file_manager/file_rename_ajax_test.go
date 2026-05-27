package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileRenameAjax_MissingRenameFileParameter(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("rename_file", "")
	form.Add("new_file", "new.txt")
	form.Add("current_dir", "/uploads")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.fileRenameAjax(req)

	if !strings.Contains(result, "rename_file is required") {
		t.Errorf("fileRenameAjax() result = %q, want to contain %q", result, "rename_file is required")
	}
}

func TestFileRenameAjax_MissingNewFileParameter(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("rename_file", "old.txt")
	form.Add("new_file", "")
	form.Add("current_dir", "/uploads")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	result := controller.fileRenameAjax(req)

	if !strings.Contains(result, "new_file is required") {
		t.Errorf("fileRenameAjax() result = %q, want to contain %q", result, "new_file is required")
	}
}
