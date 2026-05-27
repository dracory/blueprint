package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToSnakeCase_Simple(t *testing.T) {
	result := toSnakeCase("Simple.go")
	if result != "simple.go" {
		t.Errorf("toSnakeCase(\"Simple.go\") = %q; want \"simple.go\"", result)
	}
}

func TestToSnakeCase_CamelCase(t *testing.T) {
	result := toSnakeCase("camelCase.go")
	if result != "camel_case.go" {
		t.Errorf("toSnakeCase(\"camelCase.go\") = %q; want \"camel_case.go\"", result)
	}
}

func TestToSnakeCase_PascalCase(t *testing.T) {
	result := toSnakeCase("PascalCase.go")
	if result != "pascal_case.go" {
		t.Errorf("toSnakeCase(\"PascalCase.go\") = %q; want \"pascal_case.go\"", result)
	}
}

func TestToSnakeCase_HTTPClient(t *testing.T) {
	result := toSnakeCase("HTTPClient.go")
	if result != "http_client.go" {
		t.Errorf("toSnakeCase(\"HTTPClient.go\") = %q; want \"http_client.go\"", result)
	}
}

func TestToSnakeCase_XMLParser(t *testing.T) {
	result := toSnakeCase("XMLParser.go")
	if result != "xml_parser.go" {
		t.Errorf("toSnakeCase(\"XMLParser.go\") = %q; want \"xml_parser.go\"", result)
	}
}

func TestToSnakeCase_MyXMLParser(t *testing.T) {
	result := toSnakeCase("MyXMLParser.go")
	if result != "my_xml_parser.go" {
		t.Errorf("toSnakeCase(\"MyXMLParser.go\") = %q; want \"my_xml_parser.go\"", result)
	}
}

func TestToSnakeCase_File2Go(t *testing.T) {
	result := toSnakeCase("file2Go.go")
	if result != "file2_go.go" {
		t.Errorf("toSnakeCase(\"file2Go.go\") = %q; want \"file2_go.go\"", result)
	}
}

func TestToSnakeCase_AlreadySnake(t *testing.T) {
	result := toSnakeCase("already_snake.go")
	if result != "already_snake.go" {
		t.Errorf("toSnakeCase(\"already_snake.go\") = %q; want \"already_snake.go\"", result)
	}
}

func TestToSnakeCase_IDs(t *testing.T) {
	result := toSnakeCase("IDs.go")
	if result != "ids.go" {
		t.Errorf("toSnakeCase(\"IDs.go\") = %q; want \"ids.go\"", result)
	}
}

func TestToSnakeCase_UserIDs(t *testing.T) {
	result := toSnakeCase("UserIDs.go")
	if result != "user_ids.go" {
		t.Errorf("toSnakeCase(\"UserIDs.go\") = %q; want \"user_ids.go\"", result)
	}
}

func TestToSnakeCase_Server123(t *testing.T) {
	result := toSnakeCase("Server123.go")
	if result != "server123.go" {
		t.Errorf("toSnakeCase(\"Server123.go\") = %q; want \"server123.go\"", result)
	}
}

func TestToSnakeCase_JSONToXML(t *testing.T) {
	result := toSnakeCase("JSONToXML.go")
	if result != "json_to_xml.go" {
		t.Errorf("toSnakeCase(\"JSONToXML.go\") = %q; want \"json_to_xml.go\"", result)
	}
}

func TestToSnakeCase_WithDash(t *testing.T) {
	result := toSnakeCase("with-dash.go")
	if result != "with-dash.go" {
		t.Errorf("toSnakeCase(\"with-dash.go\") = %q; want \"with-dash.go\"", result)
	}
}

func TestShouldIgnore_Git(t *testing.T) {
	result := shouldIgnore(".git")
	if result != true {
		t.Errorf("shouldIgnore(\".git\") = %v; want true", result)
	}
}

func TestShouldIgnore_Hidden(t *testing.T) {
	result := shouldIgnore(".hidden")
	if result != true {
		t.Errorf("shouldIgnore(\".hidden\") = %v; want true", result)
	}
}

func TestShouldIgnore_Vendor(t *testing.T) {
	result := shouldIgnore("vendor")
	if result != true {
		t.Errorf("shouldIgnore(\"vendor\") = %v; want true", result)
	}
}

func TestShouldIgnore_NodeModules(t *testing.T) {
	result := shouldIgnore("node_modules")
	if result != true {
		t.Errorf("shouldIgnore(\"node_modules\") = %v; want true", result)
	}
}

func TestShouldIgnore_Src(t *testing.T) {
	result := shouldIgnore("src")
	if result != false {
		t.Errorf("shouldIgnore(\"src\") = %v; want false", result)
	}
}

func TestShouldIgnore_MainGo(t *testing.T) {
	result := shouldIgnore("main.go")
	if result != false {
		t.Errorf("shouldIgnore(\"main.go\") = %v; want false", result)
	}
}

func TestShouldIgnore_PathToGit(t *testing.T) {
	result := shouldIgnore("/path/to/.git")
	if result != true {
		t.Errorf("shouldIgnore(\"/path/to/.git\") = %v; want true", result)
	}
}

func TestShouldIgnore_PathToVendor(t *testing.T) {
	result := shouldIgnore("/path/to/vendor")
	if result != true {
		t.Errorf("shouldIgnore(\"/path/to/vendor\") = %v; want true", result)
	}
}

func TestShouldIgnore_PathToSrc(t *testing.T) {
	result := shouldIgnore("/path/to/src")
	if result != false {
		t.Errorf("shouldIgnore(\"/path/to/src\") = %v; want false", result)
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
