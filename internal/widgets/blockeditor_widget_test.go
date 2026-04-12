package widgets

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/registry"
)

// TestNewBlockeditorWidget tests the constructor
func TestNewBlockeditorWidget(t *testing.T) {
	t.Parallel()

	// Test with nil registry
	widget := NewBlockeditorWidget(nil)
	if widget == nil {
		t.Fatal("NewBlockeditorWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry interface
	var mockRegistry registry.RegistryInterface
	widget = NewBlockeditorWidget(mockRegistry)
	if widget == nil {
		t.Fatal("NewBlockeditorWidget(registry) should return non-nil")
	}
}

// TestBlockeditorWidget_Alias tests the Alias method
func TestBlockeditorWidget_Alias(t *testing.T) {
	t.Parallel()

	widget := NewBlockeditorWidget(nil)
	alias := widget.Alias()
	if alias != "x-blockeditor" {
		t.Errorf("Alias() = %q, want %q", alias, "x-blockeditor")
	}
}

// TestBlockeditorWidget_Description tests the Description method
func TestBlockeditorWidget_Description(t *testing.T) {
	t.Parallel()

	widget := NewBlockeditorWidget(nil)
	desc := widget.Description()
	// Description returns empty string
	if desc != "" {
		t.Errorf("Description() = %q, want empty string", desc)
	}
}

// TestBlockeditorWidget_Render tests the Render method
func TestBlockeditorWidget_Render(t *testing.T) {
	t.Parallel()

	widget := NewBlockeditorWidget(nil)
	testURL, _ := url.Parse("/test")
	req := &http.Request{
		URL: testURL,
	}

	// Test with no example param
	result := widget.Render(req, "content", map[string]string{})
	if result != "Example not found" {
		t.Errorf("Render() with no example = %q, want 'Example not found'", result)
	}

	// Test with unknown example
	result = widget.Render(req, "content", map[string]string{"example": "unknown"})
	if result != "Example not found" {
		t.Errorf("Render() with unknown example = %q, want 'Example not found'", result)
	}
}

// TestBlockeditorWidget_Struct tests the widget struct
func TestBlockeditorWidget_Struct(t *testing.T) {
	t.Parallel()

	widget := &blockeditorWidget{}

	// Test that registry field exists and can be set
	var reg registry.RegistryInterface
	widget.registry = reg

	if widget.registry != reg {
		t.Error("Should be able to set registry field")
	}
}

// TestBlockeditorWidget_MultipleInstances tests creating multiple instances
func TestBlockeditorWidget_MultipleInstances(t *testing.T) {
	t.Parallel()

	widget1 := NewBlockeditorWidget(nil)
	widget2 := NewBlockeditorWidget(nil)

	if widget1 == widget2 {
		t.Error("Multiple instances should be independent")
	}

	if widget1 == nil || widget2 == nil {
		t.Error("All widgets should be non-nil")
	}
}
