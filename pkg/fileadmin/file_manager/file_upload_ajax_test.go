package file_manager

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileUploadAjax_MissingUploadFile(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	req, err := http.NewRequest("POST", "/file-manager", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	form := url.Values{}
	form.Add("current_dir", "/uploads")
	req.PostForm = form

	result := controller.fileUploadAjax(req)

	if !strings.Contains(result, "multipart/form-data") {
		t.Errorf("fileUploadAjax() result = %q, want to contain %q", result, "multipart/form-data")
	}
}

func TestFileUploadAjax_InvalidMultipartBoundary(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	req, err := http.NewRequest("POST", "/file-manager", strings.NewReader("invalid body"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")

	form := url.Values{}
	form.Add("current_dir", "/uploads")
	req.PostForm = form

	result := controller.fileUploadAjax(req)

	if !strings.Contains(result, "boundary") {
		t.Errorf("fileUploadAjax() result = %q, want to contain %q", result, "boundary")
	}
}

func TestFileUploadAjaxWithFile(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	// Create a proper multipart request with a file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload_file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("test content"))
	writer.WriteField("current_dir", "/uploads")
	writer.Close()

	req, err := http.NewRequest("POST", "/file-manager", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	result := controller.fileUploadAjax(req)

	// The result should contain success or error, but not a boundary error
	if strings.Contains(result, "boundary") {
		t.Errorf("fileUploadAjax() should not have boundary error with proper multipart data, got: %q", result)
	}
}
