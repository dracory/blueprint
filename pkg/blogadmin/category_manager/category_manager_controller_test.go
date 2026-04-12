package category_manager

import (
	"testing"
)

// TestNewCategoryManagerController tests the constructor
func TestNewCategoryManagerController(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
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
		{"Test Category", "test-category"},
		{"My_Category-Name", "my-category-name"},
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
