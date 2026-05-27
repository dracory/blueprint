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

// TestMediaManagerControllerHumanFilesize_Bytes verifies HumanFilesize for bytes
func TestMediaManagerControllerHumanFilesize_Bytes(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	result := controller.HumanFilesize(500)
	if result != "500 B" {
		t.Errorf("HumanFilesize(500) = %s, want 500 B", result)
	}
}

// TestMediaManagerControllerHumanFilesize_Kilobytes verifies HumanFilesize for kilobytes
func TestMediaManagerControllerHumanFilesize_Kilobytes(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	result := controller.HumanFilesize(1024)
	if result != "1.0 kB" {
		t.Errorf("HumanFilesize(1024) = %s, want 1.0 kB", result)
	}
}

// TestMediaManagerControllerHumanFilesize_Megabytes verifies HumanFilesize for megabytes
func TestMediaManagerControllerHumanFilesize_Megabytes(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	result := controller.HumanFilesize(1048576)
	if result != "1.0 MB" {
		t.Errorf("HumanFilesize(1048576) = %s, want 1.0 MB", result)
	}
}

// TestMediaManagerControllerHumanFilesize_Gigabytes verifies HumanFilesize for gigabytes
func TestMediaManagerControllerHumanFilesize_Gigabytes(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	result := controller.HumanFilesize(1073741824)
	if result != "1.0 GB" {
		t.Errorf("HumanFilesize(1073741824) = %s, want 1.0 GB", result)
	}
}

// TestMediaManagerControllerHumanFilesize_Zero verifies HumanFilesize for zero
func TestMediaManagerControllerHumanFilesize_Zero(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	result := controller.HumanFilesize(0)
	if result != "0 B" {
		t.Errorf("HumanFilesize(0) = %s, want 0 B", result)
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

// TestMediaManagerControllerFileRenameAjaxValidation_MissingRenameFile verifies missing rename_file error
func TestMediaManagerControllerFileRenameAjaxValidation_MissingRenameFile(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.fileRenameAjax(r)
	if !strings.Contains(result, "rename_file is required") {
		t.Errorf("fileRenameAjax() should contain 'rename_file is required', got %s", result)
	}
}

// TestMediaManagerControllerFileRenameAjaxValidation_MissingNewFile verifies missing new_file error
func TestMediaManagerControllerFileRenameAjaxValidation_MissingNewFile(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	values.Add("rename_file", "old.txt")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.fileRenameAjax(r)
	if !strings.Contains(result, "new_file is required") {
		t.Errorf("fileRenameAjax() should contain 'new_file is required', got %s", result)
	}
}

// TestMediaManagerControllerFileRenameAjaxValidation_MissingCurrentDir verifies missing current_dir error
func TestMediaManagerControllerFileRenameAjaxValidation_MissingCurrentDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	values.Add("rename_file", "old.txt")
	values.Add("new_file", "new.txt")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.fileRenameAjax(r)
	if !strings.Contains(result, "current_dir is required") {
		t.Errorf("fileRenameAjax() should contain 'current_dir is required', got %s", result)
	}
}

// TestMediaManagerControllerFileDeleteAjaxValidation_MissingDeleteFile verifies missing delete_file error
func TestMediaManagerControllerFileDeleteAjaxValidation_MissingDeleteFile(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.fileDeleteAjax(r)
	if !strings.Contains(result, "delete_file is required") {
		t.Errorf("fileDeleteAjax() should contain 'delete_file is required', got %s", result)
	}
}

// TestMediaManagerControllerFileDeleteAjaxValidation_MissingCurrentDir verifies missing current_dir error
func TestMediaManagerControllerFileDeleteAjaxValidation_MissingCurrentDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	values.Add("delete_file", "test.txt")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.fileDeleteAjax(r)
	if !strings.Contains(result, "current_dir is required") {
		t.Errorf("fileDeleteAjax() should contain 'current_dir is required', got %s", result)
	}
}

// TestMediaManagerControllerDirectoryCreateAjaxValidation_MissingCreateDir verifies missing create_dir error
func TestMediaManagerControllerDirectoryCreateAjaxValidation_MissingCreateDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.directoryCreateAjax(r)
	if !strings.Contains(result, "create_dir is required") {
		t.Errorf("directoryCreateAjax() should contain 'create_dir is required', got %s", result)
	}
}

// TestMediaManagerControllerDirectoryCreateAjaxValidation_MissingCurrentDir verifies missing current_dir error
func TestMediaManagerControllerDirectoryCreateAjaxValidation_MissingCurrentDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	values.Add("create_dir", "newdir")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.directoryCreateAjax(r)
	if !strings.Contains(result, "current_dir is required") {
		t.Errorf("directoryCreateAjax() should contain 'current_dir is required', got %s", result)
	}
}

// TestMediaManagerControllerDirectoryDeleteAjaxValidation_MissingDeleteDir verifies missing delete_dir error
func TestMediaManagerControllerDirectoryDeleteAjaxValidation_MissingDeleteDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.directoryDeleteAjax(r)
	if !strings.Contains(result, "delete_dir is required") {
		t.Errorf("directoryDeleteAjax() should contain 'delete_dir is required', got %s", result)
	}
}

// TestMediaManagerControllerDirectoryDeleteAjaxValidation_InvalidCurrentDir verifies invalid current_dir error
func TestMediaManagerControllerDirectoryDeleteAjaxValidation_InvalidCurrentDir(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)
	controller.init(httptest.NewRequest("GET", "/", nil))

	values := url.Values{}
	values.Add("delete_dir", "testdir")
	values.Add("current_dir", ".")
	r := httptest.NewRequest("POST", "/?"+values.Encode(), nil)

	result := controller.directoryDeleteAjax(r)
	if !strings.Contains(result, "current_dir is required") {
		t.Errorf("directoryDeleteAjax() should contain 'current_dir is required', got %s", result)
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
