package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileManagerController(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		action       string
		wantContains string
	}{
		{
			name:         "load files action",
			action:       "load-files",
			wantContains: "",
		},
		{
			name:         "file clone action",
			action:       "file_clone",
			wantContains: "clone_file is required",
		},
		{
			name:         "file rename action",
			action:       "file_rename",
			wantContains: "rename_file is required",
		},
		{
			name:         "file delete action",
			action:       "file_delete",
			wantContains: "delete_file is required",
		},
		{
			name:         "directory create action",
			action:       "directory_create",
			wantContains: "create_dir is required",
		},
		{
			name:         "directory delete action",
			action:       "directory_delete",
			wantContains: "delete_dir is required",
		},
		{
			name:         "bulk move action",
			action:       "bulk_move",
			wantContains: "No items selected",
		},
		{
			name:         "bulk delete action",
			action:       "bulk_delete",
			wantContains: "No items selected",
		},
		{
			name:         "get move destinations action",
			action:       "get_move_destinations",
			wantContains: "No items selected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("action", tt.action)

			req, err := http.NewRequest("POST", "/file-manager", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.PostForm = form

			w := &testResponseWriter{}
			result := controller.Handler(w, req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("Handler() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
	}
}

type testResponseWriter struct {
	header http.Header
	body   []byte
}

func (w *testResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
}

func (w *testResponseWriter) String() string {
	return string(w.body)
}
