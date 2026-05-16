package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDirectoryCreateAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		createDir    string
		currentDir   string
		wantContains string
	}{
		{
			name:         "missing create_dir parameter",
			createDir:    "",
			currentDir:   "/uploads",
			wantContains: "create_dir is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("create_dir", tt.createDir)
			form.Add("current_dir", tt.currentDir)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.directoryCreateAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("directoryCreateAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
