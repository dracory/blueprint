package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetMoveDestinationsAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name          string
		currentDir    string
		selectedItems string
		wantContains  string
	}{
		{
			name:          "missing selected_items parameter",
			currentDir:    "/uploads",
			selectedItems: "",
			wantContains:  "No items selected",
		},
		{
			name:          "invalid JSON",
			currentDir:    "/uploads",
			selectedItems: "invalid json",
			wantContains:  "Invalid selected items data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("current_dir", tt.currentDir)
			form.Add("selected_items", tt.selectedItems)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			result := controller.getMoveDestinationsAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("getMoveDestinationsAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}
