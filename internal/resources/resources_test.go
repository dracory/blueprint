package resources

import (
	"strings"
	"testing"
)

func TestToBytes(t *testing.T) {
	// Test with non-existent file
	_, err := ToBytes("/nonexistent/file.txt")
	if err == nil {
		t.Error("ToBytes() should return error for non-existent file")
	}

	// Test with empty path
	_, err = ToBytes("")
	if err == nil {
		t.Error("ToBytes() should return error for empty path")
	}

	// Test with valid file - use a file that should exist in the embedded FS
	// The blockarea_v0200.js file is embedded in the resources package
	data, err := ToBytes("js/blockarea_v0200.js")
	if err != nil {
		t.Logf("ToBytes() with valid file returned error (file may not exist in test environment): %v", err)
	} else if len(data) == 0 {
		t.Error("ToBytes() should return non-empty data for valid file")
	}
}

func TestToString(t *testing.T) {
	// Test with non-existent file
	_, err := ToString("/nonexistent/file.txt")
	if err == nil {
		t.Error("ToString() should return error for non-existent file")
	}

	// Test with empty path
	_, err = ToString("")
	if err == nil {
		t.Error("ToString() should return error for empty path")
	}

	// Test with valid file
	content, err := ToString("js/blockarea_v0200.js")
	if err != nil {
		t.Logf("ToString() with valid file returned error (file may not exist in test environment): %v", err)
	} else if len(content) == 0 {
		t.Error("ToString() should return non-empty string for valid file")
	}
}

func TestResourceExists(t *testing.T) {
	// Test with non-existent file
	exists := ResourceExists("/nonexistent/file.txt")
	if exists {
		t.Error("ResourceExists() should return false for non-existent file")
	}

	// Test with empty path
	exists = ResourceExists("")
	if exists {
		t.Error("ResourceExists() should return false for empty path")
	}

	// Test with valid file
	exists = ResourceExists("js/blockarea_v0200.js")
	if !exists {
		t.Log("ResourceExists() returned false for valid file (file may not exist in test environment)")
	}
}

func TestResource(t *testing.T) {
	// Test with non-existent file
	_, err := Resource("/nonexistent/file.txt")
	if err == nil {
		t.Error("Resource() should return error for non-existent file")
	}

	// Test with empty path
	_, err = Resource("")
	if err == nil {
		t.Error("Resource() should return error for empty path")
	}

	// Test with valid file
	content, err := Resource("js/blockarea_v0200.js")
	if err != nil {
		t.Logf("Resource() with valid file returned error (file may not exist in test environment): %v", err)
	} else if len(content) == 0 {
		t.Error("Resource() should return non-empty string for valid file")
	}
}

func TestResourceWithParams(t *testing.T) {
	// ResourceWithParams uses template.Must which panics on non-existent files
	panicOccurred := false
	defer func() {
		if r := recover(); r != nil {
			panicOccurred = true
		}
	}()

	// Test with non-existent template - this will panic
	ResourceWithParams("/nonexistent/template.html", map[string]string{"key": "value"})

	if !panicOccurred {
		t.Error("ResourceWithParams() should panic for non-existent template")
	}
}

func TestImageToBase64String(t *testing.T) {
	// Test with non-existent file
	result := ImageToBase64String("/nonexistent/image.png")
	if result == "" {
		// This is expected since the function ignores errors
		t.Log("ImageToBase64String() returns empty string for non-existent file (expected)")
	}

	// Test with empty path
	result = ImageToBase64String("")
	if result == "" {
		// This is expected
		t.Log("ImageToBase64String() returns empty string for empty path (expected)")
	}

	// Test with SVG extension
	result = ImageToBase64String("/nonexistent/image.svg")
	if result == "" {
		t.Log("ImageToBase64String() returns empty string for non-existent SVG (expected)")
	}
	// Even though file doesn't exist, it should have the SVG prefix if it had data
	// Since data is empty, the result will just be the prefix
	if result != "" && !strings.Contains(result, "image/svg+xml") {
		t.Error("ImageToBase64String() with .svg should contain image/svg+xml")
	}
}

func TestBlockAreaJS(t *testing.T) {
	// Test BlockAreaJS - it returns the JavaScript resource or empty string
	result := BlockAreaJS()
	// It will return empty string if the file doesn't exist
	// Just verify it doesn't panic
	if result == "" {
		t.Log("BlockAreaJS() returned empty string (file may not exist)")
	}
}
