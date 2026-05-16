package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDirectoryDeleteAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		deleteDir    string
		currentDir   string
		wantContains string
	}{
		{
			name:         "missing delete_dir parameter",
			deleteDir:    "",
			currentDir:   "/uploads",
			wantContains: "delete_dir is required",
		},
		{
			name:         "invalid current_dir",
			deleteDir:    "test",
			currentDir:   ".",
			wantContains: "invalid directory name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("delete_dir", tt.deleteDir)
			form.Add("current_dir", tt.currentDir)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.directoryDeleteAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("directoryDeleteAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
