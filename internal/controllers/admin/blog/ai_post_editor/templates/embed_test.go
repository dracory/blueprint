package templates

import (
	"testing"
)

// TestResourceExists verifies ResourceExists returns true for existing files
func TestResourceExists(t *testing.T) {
	t.Parallel()
	exists := ResourceExists("app.html")

	if !exists {
		t.Error("ResourceExists should return true for app.html")
	}
}

// TestResourceExistsNotFound verifies ResourceExists returns false for non-existent files
func TestResourceExistsNotFound(t *testing.T) {
	t.Parallel()
	exists := ResourceExists("nonexistent.html")

	if exists {
		t.Error("ResourceExists should return false for non-existent files")
	}
}

// TestToString verifies ToString returns content for existing files
func TestToString(t *testing.T) {
	t.Parallel()
	content, err := ToString("app.html")

	if err != nil {
		t.Errorf("ToString should not return error for app.html: %v", err)
	}
	if content == "" {
		t.Error("ToString should return non-empty content for app.html")
	}
}

// TestToStringNotFound verifies ToString returns error for non-existent files
func TestToStringNotFound(t *testing.T) {
	t.Parallel()
	_, err := ToString("nonexistent.html")

	if err == nil {
		t.Error("ToString should return error for non-existent files")
	}
}

// TestToBytes verifies ToBytes returns content for existing files
func TestToBytes(t *testing.T) {
	t.Parallel()
	content, err := ToBytes("app.html")

	if err != nil {
		t.Errorf("ToBytes should not return error for app.html: %v", err)
	}
	if len(content) == 0 {
		t.Error("ToBytes should return non-empty content for app.html")
	}
}

// TestTpl verifies Tpl returns non-empty string for valid template
func TestTpl(t *testing.T) {
	t.Parallel()
	result := Tpl("app.html", map[string]any{})

	if result == "" {
		t.Error("Tpl should return non-empty string for app.html")
	}
}

// TestTplInvalidFile verifies Tpl returns empty string for invalid file
func TestTplInvalidFile(t *testing.T) {
	t.Parallel()
	result := Tpl("nonexistent.html", map[string]any{})

	if result != "" {
		t.Error("Tpl should return empty string for non-existent files")
	}
}
