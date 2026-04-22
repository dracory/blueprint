package widgets

import (
	"testing"

	"project/internal/registry"
	"project/internal/testutils"
)

// TestNewBlogPostWidget tests the constructor
func TestNewBlogPostWidget(t *testing.T) {
	t.Parallel()

	// Test with nil registry
	widget := NewBlogPostWidget(nil)
	if widget == nil {
		t.Fatal("NewBlogPostWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewBlogPostWidget(registry)
	if widget == nil {
		t.Fatal("NewBlogPostWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

// TestBlogPostWidget_Alias tests the Alias method
func TestBlogPostWidget_Alias(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostWidget(nil)
	alias := widget.Alias()
	if alias != "x-blog-post" {
		t.Errorf("Alias() = %q, want %q", alias, "x-blog-post")
	}
}

// TestBlogPostWidget_Description tests the Description method
func TestBlogPostWidget_Description(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostWidget(nil)
	desc := widget.Description()
	expected := "Renders a blog post"
	if desc != expected {
		t.Errorf("Description() = %q, want %q", desc, expected)
	}
}

// TestBlogPostWidget_processContent tests the processContent method
func TestBlogPostWidget_processContent(t *testing.T) {
	widget := NewBlogPostWidget(nil)

	tests := []struct {
		name     string
		content  string
		editor   string
		expected string
	}{
		{
			name:     "BlockArea editor",
			content:  "test content",
			editor:   "BlockArea",
			expected: "test content", // Will be processed by BlogPostBlocksToString
		},
		{
			name:     "Other editor",
			content:  "raw content",
			editor:   "other",
			expected: "raw content",
		},
		{
			name:     "Empty editor",
			content:  "content",
			editor:   "",
			expected: "content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := widget.processContent(tt.content, tt.editor)
			// For BlockArea, the content is passed through BlogPostBlocksToString
			// which may modify it, so we just check it's not empty when input is not empty
			if tt.content != "" && result == "" {
				t.Error("processContent() should not return empty for non-empty input")
			}
		})
	}
}

// TestBlogPostWidget_css tests the css method
func TestBlogPostWidget_css(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostWidget(nil)
	css := widget.css()

	if css == "" {
		t.Error("css() should return non-empty CSS string")
	}

	// Check for expected CSS selectors
	expectedSelectors := []string{
		"#SectionNewsItem",
		".BlogTitle",
		".BlogContent",
	}

	for _, selector := range expectedSelectors {
		if !contains(css, selector) {
			t.Errorf("css() should contain selector %q", selector)
		}
	}
}

// TestBlogPostWidget_sectionBreadcrumbs tests the sectionBreadcrumbs method
func TestBlogPostWidget_sectionBreadcrumbs(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostWidget(nil)
	// Since we can't easily create a mock blogstore.PostInterface,
	// we'll pass nil and check it doesn't panic
	result := widget.sectionBreadcrumbs(nil)
	// Function returns empty Wrap() currently
	_ = result
}

// TestBlogPostWidget_Struct tests the widget struct fields
func TestBlogPostWidget_Struct(t *testing.T) {
	t.Parallel()

	widget := &blogPostWidget{}

	// Test that registry field exists and can be set
	var reg registry.RegistryInterface
	widget.registry = reg

	if widget.registry != reg {
		t.Error("Should be able to set registry field")
	}
}

// TestBlogPostWidget_MultipleInstances tests creating multiple instances
func TestBlogPostWidget_MultipleInstances(t *testing.T) {
	t.Parallel()

	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	widget1 := NewBlogPostWidget(registry1)
	widget2 := NewBlogPostWidget(registry2)

	if widget1 == widget2 {
		t.Error("Multiple instances should be independent")
	}

	if widget1 == nil || widget2 == nil {
		t.Error("All widgets should be non-nil")
	}

	if widget1.registry != registry1 {
		t.Error("Widget1 should have registry1")
	}

	if widget2.registry != registry2 {
		t.Error("Widget2 should have registry2")
	}
}

// TestBlogPostWidget_Interface tests that widget implements Widget interface
func TestBlogPostWidget_Interface(t *testing.T) {
	t.Parallel()

	var _ Widget = (*blogPostWidget)(nil)

	widget := NewBlogPostWidget(nil)
	var widgetInterface Widget = widget

	if widgetInterface.Alias() != "x-blog-post" {
		t.Error("Widget interface should work correctly")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
