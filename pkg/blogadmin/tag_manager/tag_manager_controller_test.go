package tag_manager

import (
	"testing"
)

// TestNewTagManagerController tests the constructor
func TestNewTagManagerController(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestSlugify tests the slugify function
func TestSlugify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test Tag", "test-tag"},
		{"My_Tag-Name", "my-tag-name"},
		{"UPPERCASE", "uppercase"},
		{"123 Numbers", "123-numbers"},
		{"", ""},
		{"Already-Slugified", "already-slugified"},
		{"Multiple   Spaces", "multiple---spaces"},
	}

	for _, tt := range tests {
		result := slugify(tt.input)
		if result != tt.expected {
			t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
