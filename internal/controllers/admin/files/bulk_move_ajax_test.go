package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBulkMoveAjax_RequiresCurrentDir(t *testing.T) {
	controller := &FileManagerController{
		storage: &mockStorage{},
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"selected_items": {`[{"path":"/test.txt","type":"file"}]`},
		},
	})

	if response != "" && !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "current_dir is required") {
		t.Errorf("Should require current_dir, got: %s", response)
	}
}

func TestBulkMoveAjax_RequiresSelectedItems(t *testing.T) {
	controller := &FileManagerController{
		storage: &mockStorage{},
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir": {"/source"},
		},
	})

	if response != "" && !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "No items selected") {
		t.Errorf("Should require selected items, got: %s", response)
	}
}

func TestBulkMoveAjax_RequiresStorage(t *testing.T) {
	controller := &FileManagerController{
		storage: nil,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":    {"/source"},
			"selected_items": {`[{"path":"/test.txt","type":"file"}]`},
		},
	})

	if response != "" && !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "Storage not initialized") {
		t.Errorf("Should require storage, got: %s", response)
	}
}

func TestBulkMoveAjax_HandlesInvalidJSON(t *testing.T) {
	controller := &FileManagerController{
		storage: &mockStorage{},
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":    {"/source"},
			"selected_items": {`invalid json`},
		},
	})

	if response != "" && !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "Invalid selected items data") {
		t.Errorf("Should handle invalid JSON, got: %s", response)
	}
}

func TestBulkMoveAjax_MovesSingleFile(t *testing.T) {
	mockStore := &mockStorage{
		files: map[string][]byte{
			"/source/test.txt": []byte("test content"),
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/dest"},
			"selected_items":  {`[{"path":"/source/test.txt","type":"file"}]`},
		},
	})

	if !strings.Contains(response, `"status":"success"`) {
		t.Errorf("Should return success status, got: %s", response)
	}
	if !strings.Contains(response, "Successfully moved") {
		t.Errorf("Should confirm move, got: %s", response)
	}
	if _, exists := mockStore.files["/dest/test.txt"]; !exists {
		t.Error("File should be moved to destination")
	}
}

func TestBulkMoveAjax_MovesMultipleFiles(t *testing.T) {
	mockStore := &mockStorage{
		files: map[string][]byte{
			"/source/file1.txt": []byte("content1"),
			"/source/file2.txt": []byte("content2"),
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/dest"},
			"selected_items": {`[
				{"path":"/source/file1.txt","type":"file"},
				{"path":"/source/file2.txt","type":"file"}
			]`},
		},
	})

	if !strings.Contains(response, `"status":"success"`) {
		t.Errorf("Should return success status, got: %s", response)
	}
	if !strings.Contains(response, "Successfully moved 2 item(s)") {
		t.Errorf("Should confirm moving 2 items, got: %s", response)
	}
	if _, exists := mockStore.files["/dest/file1.txt"]; !exists {
		t.Error("File1 should be moved")
	}
	if _, exists := mockStore.files["/dest/file2.txt"]; !exists {
		t.Error("File2 should be moved")
	}
}

func TestBulkMoveAjax_MovesDirectory(t *testing.T) {
	mockStore := &mockStorage{
		directories: map[string][]string{
			"/source/mydir": {},
		},
		files: map[string][]byte{
			"/source/mydir/file.txt": []byte("content"),
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/dest"},
			"selected_items":  {`[{"path":"/source/mydir","type":"directory"}]`},
		},
	})

	if !strings.Contains(response, `"status":"success"`) {
		t.Errorf("Should return success status, got: %s", response)
	}
	if _, exists := mockStore.directories["/dest/mydir"]; !exists {
		t.Error("Directory should be moved")
	}
}

func TestBulkMoveAjax_PreventsMovingIntoItself(t *testing.T) {
	mockStore := &mockStorage{
		directories: map[string][]string{
			"/source/mydir": {},
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/source/mydir"},
			"selected_items":  {`[{"path":"/source/mydir","type":"directory"}]`},
		},
	})

	if !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "Cannot move directory into itself") {
		t.Errorf("Should prevent moving into itself, got: %s", response)
	}
}

func TestBulkMoveAjax_MovesToRoot(t *testing.T) {
	mockStore := &mockStorage{
		files: map[string][]byte{
			"/source/test.txt": []byte("test content"),
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {},
			"selected_items":  {`[{"path":"/source/test.txt","type":"file"}]`},
		},
	})

	if !strings.Contains(response, `"status":"success"`) {
		t.Errorf("Should return success status, got: %s", response)
	}
	if _, exists := mockStore.files["/test.txt"]; !exists {
		t.Error("File should be moved to root")
	}
}

func TestBulkMoveAjax_HandlesMoveError(t *testing.T) {
	mockStore := &mockStorage{
		files:   map[string][]byte{},
		moveErr: fmt.Errorf("mock move error"),
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/dest"},
			"selected_items":  {`[{"path":"/source/test.txt","type":"file"}]`},
		},
	})

	if !strings.Contains(response, `"status":"error"`) {
		t.Errorf("Should return error status, got: %s", response)
	}
	if !strings.Contains(response, "Failed to move") {
		t.Errorf("Should report move failure, got: %s", response)
	}
}

func TestBulkMoveAjax_HandlesPartialSuccess(t *testing.T) {
	mockStore := &mockStorage{
		files: map[string][]byte{
			"/source/file1.txt": []byte("content1"),
		},
		moveShouldFail: map[string]bool{
			"/source/file2.txt": true,
		},
	}
	controller := &FileManagerController{
		storage: mockStore,
	}

	response := controller.bulkMoveAjax(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
		PostForm: url.Values{
			"current_dir":     {"/source"},
			"destination_dir": {"/dest"},
			"selected_items": {`[
				{"path":"/source/file1.txt","type":"file"},
				{"path":"/source/file2.txt","type":"file"}
			]`},
		},
	})

	if !strings.Contains(response, `"status":"success"`) {
		t.Errorf("Should return success status (partial), got: %s", response)
	}
	if !strings.Contains(response, "Successfully moved 1 item(s)") {
		t.Errorf("Should confirm partial success, got: %s", response)
	}
	if !strings.Contains(response, "Some items failed") {
		t.Errorf("Should report failures, got: %s", response)
	}
}

// Mock storage implementation for testing
type mockStorage struct {
	files          map[string][]byte
	directories    map[string][]string
	moveErr        error
	moveShouldFail map[string]bool
}

func (m *mockStorage) Directories(path string) ([]string, error) {
	return m.directories[path], nil
}

func (m *mockStorage) Files(path string) ([]string, error) {
	var files []string
	for filePath := range m.files {
		// Simple prefix matching for mock
		if len(filePath) > len(path) && filePath[:len(path)] == path {
			files = append(files, filePath)
		}
	}
	return files, nil
}

func (m *mockStorage) Size(path string) (int64, error) {
	if data, ok := m.files[path]; ok {
		return int64(len(data)), nil
	}
	return 0, nil
}

func (m *mockStorage) LastModified(path string) (time.Time, error) {
	return time.Time{}, nil
}

func (m *mockStorage) Url(path string) (string, error) {
	return "http://example.com" + path, nil
}

func (m *mockStorage) Exists(path string) (bool, error) {
	_, ok := m.files[path]
	return ok, nil
}

func (m *mockStorage) Get(path string) ([]byte, error) {
	return m.files[path], nil
}

func (m *mockStorage) Put(path string, data []byte) error {
	if m.files == nil {
		m.files = make(map[string][]byte)
	}
	m.files[path] = data
	return nil
}

func (m *mockStorage) Delete(path string) error {
	delete(m.files, path)
	return nil
}

func (m *mockStorage) Move(source, destination string) error {
	if m.moveErr != nil {
		return m.moveErr
	}
	if m.moveShouldFail != nil && m.moveShouldFail[source] {
		return fmt.Errorf("mock move failure for %s", source)
	}

	// Move file
	if data, ok := m.files[source]; ok {
		if m.files == nil {
			m.files = make(map[string][]byte)
		}
		m.files[destination] = data
		delete(m.files, source)
	}

	// Move directory entry if exists
	if dirs, ok := m.directories[source]; ok {
		if m.directories == nil {
			m.directories = make(map[string][]string)
		}
		m.directories[destination] = dirs
		delete(m.directories, source)
	}

	return nil
}

func (m *mockStorage) Copy(source, destination string) error {
	return nil
}

func (m *mockStorage) Mkdir(path string) error {
	if m.directories == nil {
		m.directories = make(map[string][]string)
	}
	m.directories[path] = []string{}
	return nil
}

func (m *mockStorage) Rmdir(path string) error {
	delete(m.directories, path)
	return nil
}

func (m *mockStorage) DeleteDirectory(path string) error {
	delete(m.directories, path)
	return nil
}

func (m *mockStorage) DeleteFile(paths []string) error {
	for _, path := range paths {
		delete(m.files, path)
	}
	return nil
}

func (m *mockStorage) MakeDirectory(path string) error {
	if m.directories == nil {
		m.directories = make(map[string][]string)
	}
	m.directories[path] = []string{}
	return nil
}

func (m *mockStorage) ReadFile(path string) ([]byte, error) {
	return m.files[path], nil
}
