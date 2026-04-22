package widgets

import (
	"testing"

	"project/internal/registry"
	"project/internal/testutils"
)

// TestNewBlogPostListWidget tests the constructor
func TestNewBlogPostListWidget(t *testing.T) {
	t.Parallel()

	// Test with nil registry
	widget := NewBlogPostListWidget(nil)
	if widget == nil {
		t.Fatal("NewBlogPostListWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewBlogPostListWidget(registry)
	if widget == nil {
		t.Fatal("NewBlogPostListWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

// TestBlogPostListWidget_Alias tests the Alias method
func TestBlogPostListWidget_Alias(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostListWidget(nil)
	alias := widget.Alias()
	if alias != "x-blog-post-list" {
		t.Errorf("Alias() = %q, want %q", alias, "x-blog-post-list")
	}
}

// TestBlogPostListWidget_Description tests the Description method
func TestBlogPostListWidget_Description(t *testing.T) {
	t.Parallel()

	widget := NewBlogPostListWidget(nil)
	desc := widget.Description()
	expected := "Renders a list of the blog posts"
	if desc != expected {
		t.Errorf("Description() = %q, want %q", desc, expected)
	}
}

// TestBlogPostListWidget_Struct tests the widget struct fields
func TestBlogPostListWidget_Struct(t *testing.T) {
	t.Parallel()

	widget := &blogPostListWidget{}

	// Test that registry field exists and can be set
	var reg registry.RegistryInterface
	widget.registry = reg

	// Test data struct
	data := blogPostListWidgetData{
		page:      1,
		perPage:   12,
		postCount: 0,
	}
	if data.page != 1 {
		t.Error("blogPostListWidgetData.page should be set correctly")
	}
	if data.perPage != 12 {
		t.Error("blogPostListWidgetData.perPage should be set correctly")
	}
}

// TestBlogPostListWidget_MultipleInstances tests creating multiple instances
func TestBlogPostListWidget_MultipleInstances(t *testing.T) {
	t.Parallel()

	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	widget1 := NewBlogPostListWidget(registry1)
	widget2 := NewBlogPostListWidget(registry2)

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

// TestBlogPostListWidget_Interface tests that widget implements Widget interface
func TestBlogPostListWidget_Interface(t *testing.T) {
	t.Parallel()

	var _ Widget = (*blogPostListWidget)(nil)

	widget := NewBlogPostListWidget(nil)
	var widgetInterface Widget = widget

	if widgetInterface.Alias() != "x-blog-post-list" {
		t.Error("Widget interface should work correctly")
	}
}
