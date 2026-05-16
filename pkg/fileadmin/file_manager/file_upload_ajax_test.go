package file_manager

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFileUploadAjax(t *testing.T) {
	reg, cleanup := setupTestRegistry()
	defer cleanup()

	controller := NewFileManagerController(reg)

	tests := []struct {
		name         string
		currentDir   string
		setupReq     func() (*http.Request, error)
		wantContains string
	}{
		{
			name:       "missing upload_file",
			currentDir: "/uploads",
			setupReq: func() (*http.Request, error) {
				return http.NewRequest("POST", "/file-manager", nil)
			},
			wantContains: "multipart/form-data",
		},
		{
			name:       "invalid multipart boundary",
			currentDir: "/uploads",
			setupReq: func() (*http.Request, error) {
				req, err := http.NewRequest("POST", "/file-manager", strings.NewReader("invalid body"))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "multipart/form-data")
				return req, err
			},
			wantContains: "boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.setupReq()
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			form := url.Values{}
			form.Add("current_dir", tt.currentDir)
			req.PostForm = form

			result := controller.fileUploadAjax(req)

			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("fileUploadAjax() result = %q, want to contain %q", result, tt.wantContains)
			}
		})
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
