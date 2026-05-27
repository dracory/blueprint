package file_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileManagerController_LoadFilesAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "load-files")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	_ = controller.Handler(w, req)
}

func TestFileManagerController_FileCloneAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "file_clone")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "clone_file is required") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "clone_file is required")
	}
}

func TestFileManagerController_FileRenameAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "file_rename")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "rename_file is required") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "rename_file is required")
	}
}

func TestFileManagerController_FileDeleteAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "file_delete")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "delete_file is required") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "delete_file is required")
	}
}

func TestFileManagerController_DirectoryCreateAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "directory_create")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "create_dir is required") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "create_dir is required")
	}
}

func TestFileManagerController_DirectoryDeleteAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "directory_delete")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "delete_dir is required") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "delete_dir is required")
	}
}

func TestFileManagerController_BulkMoveAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "bulk_move")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "No items selected") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "No items selected")
	}
}

func TestFileManagerController_BulkDeleteAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "bulk_delete")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "No items selected") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "No items selected")
	}
}

func TestFileManagerController_GetMoveDestinationsAction(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	form := url.Values{}
	form.Add("action", "get_move_destinations")

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.PostForm = form

	w := &testResponseWriter{}
	result := controller.Handler(w, req)

	if !strings.Contains(result, "No items selected") {
		t.Errorf("Handler() result = %q, want to contain %q", result, "No items selected")
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
