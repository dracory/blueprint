package admin

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestNewMediaManagerController verifies controller can be created
func TestNewMediaManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	if controller == nil {
		t.Error("NewMediaManagerController() returned nil")
	}
}

// TestMediaManagerControllerRegistry verifies controller has registry
func TestMediaManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestMediaManagerControllerAnyIndexExists verifies AnyIndex method exists
func TestMediaManagerControllerAnyIndexExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	// Verify AnyIndex method exists (should compile without error)
	_ = controller.AnyIndex
}

// TestMediaManagerControllerHumanFilesize verifies HumanFilesize helper
func TestMediaManagerControllerHumanFilesize(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	tests := []struct {
		name string
		size int64
		want string
	}{
		{"bytes", 500, "500 B"},
		{"kilobytes", 1024, "1 KB"},
		{"megabytes", 1048576, "1 MB"},
		{"gigabytes", 1073741824, "1 GB"},
		{"zero", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := controller.HumanFilesize(tt.size)
			if result != tt.want {
				t.Errorf("HumanFilesize(%d) = %s, want %s", tt.size, result, tt.want)
			}
		})
	}
}

// TestMediaManagerControllerAnyIndexDefaultAction verifies default action
func TestMediaManagerControllerAnyIndexDefaultAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	result := controller.AnyIndex(w, r)

	// Should return JSON response
	if result == "" {
		t.Error("AnyIndex() should return response for default action")
	}
}

// TestMediaManagerControllerAnyIndexFileRenameAction verifies file_rename action
func TestMediaManagerControllerAnyIndexFileRenameAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("action", "file_rename")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.AnyIndex(w, r)

	if result == "" {
		t.Error("AnyIndex() should return response for file_rename action")
	}
}

// TestMediaManagerControllerAnyIndexFileDeleteAction verifies file_delete action
func TestMediaManagerControllerAnyIndexFileDeleteAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("action", "file_delete")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.AnyIndex(w, r)

	if result == "" {
		t.Error("AnyIndex() should return response for file_delete action")
	}
}

// TestMediaManagerControllerAnyIndexDirectoryCreateAction verifies directory_create action
func TestMediaManagerControllerAnyIndexDirectoryCreateAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("action", "directory_create")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.AnyIndex(w, r)

	if result == "" {
		t.Error("AnyIndex() should return response for directory_create action")
	}
}

// TestMediaManagerControllerAnyIndexDirectoryDeleteAction verifies directory_delete action
func TestMediaManagerControllerAnyIndexDirectoryDeleteAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("action", "directory_delete")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.AnyIndex(w, r)

	if result == "" {
		t.Error("AnyIndex() should return response for directory_delete action")
	}
}

// TestMediaManagerControllerAnyIndexFileUploadAction verifies file_upload action
func TestMediaManagerControllerAnyIndexFileUploadAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("action", "file_upload")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.AnyIndex(w, r)

	if result == "" {
		t.Error("AnyIndex() should return response for file_upload action")
	}
}

// TestMediaManagerControllerFileRenameAjaxValidation verifies validation errors
func TestMediaManagerControllerFileRenameAjaxValidation(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	tests := []struct {
		name         string
		paramName    string
		paramValue   string
		wantContains string
	}{
		{"missing rename_file", "rename_file", "", "rename_file is required"},
		{"missing new_file", "new_file", "", "new_file is required"},
		{"missing current_dir", "current_dir", "", "current_dir is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", nil)
			if tt.paramValue != "" {
				values := url.Values{}
				values.Add(tt.paramName, tt.paramValue)
				r = httptest.NewRequest("POST", "/?"+values.Encode(), nil)
			}

			result := controller.fileRenameAjax(r)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileRenameAjax() should contain %s, got %s", tt.wantContains, result)
			}
		})
	}
}

// TestMediaManagerControllerFileDeleteAjaxValidation verifies validation errors
func TestMediaManagerControllerFileDeleteAjaxValidation(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	tests := []struct {
		name         string
		paramName    string
		paramValue   string
		wantContains string
	}{
		{"missing delete_file", "delete_file", "", "delete_file is required"},
		{"missing current_dir", "current_dir", "", "current_dir is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", nil)
			if tt.paramValue != "" {
				values := url.Values{}
				values.Add(tt.paramName, tt.paramValue)
				r = httptest.NewRequest("POST", "/?"+values.Encode(), nil)
			}

			result := controller.fileDeleteAjax(r)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileDeleteAjax() should contain %s, got %s", tt.wantContains, result)
			}
		})
	}
}

// TestMediaManagerControllerDirectoryCreateAjaxValidation verifies validation errors
func TestMediaManagerControllerDirectoryCreateAjaxValidation(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	tests := []struct {
		name         string
		paramName    string
		paramValue   string
		wantContains string
	}{
		{"missing create_dir", "create_dir", "", "create_dir is required"},
		{"missing current_dir", "current_dir", "", "current_dir is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", nil)
			if tt.paramValue != "" {
				values := url.Values{}
				values.Add(tt.paramName, tt.paramValue)
				r = httptest.NewRequest("POST", "/?"+values.Encode(), nil)
			}

			result := controller.directoryCreateAjax(r)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("directoryCreateAjax() should contain %s, got %s", tt.wantContains, result)
			}
		})
	}
}

// TestMediaManagerControllerDirectoryDeleteAjaxValidation verifies validation errors
func TestMediaManagerControllerDirectoryDeleteAjaxValidation(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	tests := []struct {
		name         string
		paramName    string
		paramValue   string
		wantContains string
	}{
		{"missing delete_dir", "delete_dir", "", "delete_dir is required"},
		{"missing current_dir", "current_dir", "", "current_dir is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", nil)
			if tt.paramValue != "" {
				values := url.Values{}
				values.Add(tt.paramName, tt.paramValue)
				r = httptest.NewRequest("POST", "/?"+values.Encode(), nil)
			}

			result := controller.directoryDeleteAjax(r)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("directoryDeleteAjax() should contain %s, got %s", tt.wantContains, result)
			}
		})
	}
}

// TestMediaManagerControllerModalFileUpload verifies modal HTML generation
func TestMediaManagerControllerModalFileUpload(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalFileUpload("/test/dir")

	if html == "" {
		t.Error("modalFileUpload() should return HTML")
	}

	if !strings.Contains(html, "ModalUploadFile") {
		t.Error("modalFileUpload() should contain modal ID")
	}
}

// TestMediaManagerControllerModalDirectoryCreate verifies modal HTML generation
func TestMediaManagerControllerModalDirectoryCreate(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalDirectoryCreate("/test/dir")

	if html == "" {
		t.Error("modalDirectoryCreate() should return HTML")
	}

	if !strings.Contains(html, "ModalDirectoryCreate") {
		t.Error("modalDirectoryCreate() should contain modal ID")
	}
}

// TestMediaManagerControllerModalDirectoryDelete verifies modal HTML generation
func TestMediaManagerControllerModalDirectoryDelete(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalDirectoryDelete("/test/dir")

	if html == "" {
		t.Error("modalDirectoryDelete() should return HTML")
	}

	if !strings.Contains(html, "ModalDirectoryDelete") {
		t.Error("modalDirectoryDelete() should contain modal ID")
	}
}

// TestMediaManagerControllerModalFileDelete verifies modal HTML generation
func TestMediaManagerControllerModalFileDelete(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalFileDelete("/test/dir")

	if html == "" {
		t.Error("modalFileDelete() should return HTML")
	}

	if !strings.Contains(html, "ModalFileDelete") {
		t.Error("modalFileDelete() should contain modal ID")
	}
}

// TestMediaManagerControllerModalFileRename verifies modal HTML generation
func TestMediaManagerControllerModalFileRename(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalFileRename("/test/dir")

	if html == "" {
		t.Error("modalFileRename() should return HTML")
	}

	if !strings.Contains(html, "ModalFileRename") {
		t.Error("modalFileRename() should contain modal ID")
	}
}

// TestMediaManagerControllerModalFileView verifies modal HTML generation
func TestMediaManagerControllerModalFileView(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	html := controller.modalFileView()

	if html == "" {
		t.Error("modalFileView() should return HTML")
	}

	if !strings.Contains(html, "ModalFileView") {
		t.Error("modalFileView() should contain modal ID")
	}
}
