package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileCloneAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		cloneFile    string
		currentDir   string
		newFile      string
		wantContains string
	}{
		{
			name:         "missing clone_file parameter",
			cloneFile:    "",
			currentDir:   "/uploads",
			wantContains: "clone_file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("clone_file", tt.cloneFile)
			form.Add("current_dir", tt.currentDir)
			form.Add("new_file", tt.newFile)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.fileCloneAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileCloneAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
