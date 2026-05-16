package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileDeleteAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		deleteFile   string
		currentDir   string
		wantContains string
	}{
		{
			name:         "missing delete_file parameter",
			deleteFile:   "",
			currentDir:   "/uploads",
			wantContains: "delete_file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("delete_file", tt.deleteFile)
			form.Add("current_dir", tt.currentDir)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.fileDeleteAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileDeleteAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
