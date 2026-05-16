package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileRenameAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		renameFile   string
		newFile      string
		currentDir   string
		wantContains string
	}{
		{
			name:         "missing rename_file parameter",
			renameFile:   "",
			newFile:      "new.txt",
			currentDir:   "/uploads",
			wantContains: "rename_file is required",
		},
		{
			name:         "missing new_file parameter",
			renameFile:   "old.txt",
			newFile:      "",
			currentDir:   "/uploads",
			wantContains: "new_file is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("rename_file", tt.renameFile)
			form.Add("new_file", tt.newFile)
			form.Add("current_dir", tt.currentDir)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.fileRenameAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileRenameAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
