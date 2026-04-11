package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple.go", "simple.go"},
		{"camelCase.go", "camel_case.go"},
		{"PascalCase.go", "pascal_case.go"},
		{"HTTPClient.go", "http_client.go"},
		{"XMLParser.go", "xml_parser.go"},
		{"MyXMLParser.go", "my_xml_parser.go"},
		{"file2Go.go", "file2_go.go"},
		{"already_snake.go", "already_snake.go"},
		{"IDs.go", "ids.go"},          // Updated expectation
		{"UserIDs.go", "user_ids.go"}, // Updated expectation
		{"Server123.go", "server123.go"},
		{"JSONToXML.go", "json_to_xml.go"},
		{"with-dash.go", "with-dash.go"},
	}

	for _, test := range tests {
		result := toSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("toSnakeCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{".git", true},
		{".hidden", true},
		{"vendor", true},
		{"node_modules", true},
		{"src", false},
		{"main.go", false},
		{"/path/to/.git", true},
		{"/path/to/vendor", true},
		{"/path/to/src", false},
	}

	for _, test := range tests {
		result := shouldIgnore(test.path)
		if result != test.expected {
			t.Errorf("shouldIgnore(%q) = %v; want %v", test.path, result, test.expected)
		}
	}
}

func TestRenameGoFiles(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "snakecase_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"SimpleFile.go",
		"AnotherFile.go",
		"already_snake.go",
	}

	for _, filename := range testFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte("package main"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Run renameGoFiles in dry-run mode first
	err = renameGoFiles(tempDir, true, false)
	if err != nil {
		t.Errorf("renameGoFiles(dry-run) failed: %v", err)
	}

	// Run renameGoFiles for real
	err = renameGoFiles(tempDir, false, false)
	if err != nil {
		t.Errorf("renameGoFiles failed: %v", err)
	}

	// Check that files were renamed
	expectedFiles := []string{
		"simple_file.go",
		"another_file.go",
		"already_snake.go",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(tempDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", filename)
		}
	}
}

func TestRenameGoFilesNonExistentDir(t *testing.T) {
	err := renameGoFiles("/nonexistent/path/12345", false, false)
	if err == nil {
		t.Error("renameGoFiles with non-existent dir should return error")
	}
}

func TestRenameGoFilesNotADirectory(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "snakecase_test_file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	err = renameGoFiles(tempFile.Name(), false, false)
	if err == nil {
		t.Error("renameGoFiles with file (not dir) should return error")
	}
}
