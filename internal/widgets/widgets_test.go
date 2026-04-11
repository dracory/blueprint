package widgets

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"
)

func TestNewAuthenticatedWidget(t *testing.T) {
	// Test with nil registry
	widget := NewAuthenticatedWidget(nil)
	if widget == nil {
		t.Fatal("NewAuthenticatedWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewAuthenticatedWidget(registry)
	if widget == nil {
		t.Fatal("NewAuthenticatedWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

func TestAuthenticatedWidget_Alias(t *testing.T) {
	registry := testutils.Setup()
	widget := NewAuthenticatedWidget(registry)

	alias := widget.Alias()
	if alias != "x-authenticated" {
		t.Errorf("Alias() = %q, want %q", alias, "x-authenticated")
	}
}

func TestAuthenticatedWidget_Description(t *testing.T) {
	registry := testutils.Setup()
	widget := NewAuthenticatedWidget(registry)

	desc := widget.Description()
	if desc == "" {
		t.Error("Description() should return non-empty string")
	}
	if desc != "Renders the content if the user is authenticated" {
		t.Errorf("Description() = %q, want %q", desc, "Renders the content if the user is authenticated")
	}
}

func TestAuthenticatedWidget_Render(t *testing.T) {
	registry := testutils.Setup()
	widget := NewAuthenticatedWidget(registry)

	// Test with nil request (no auth user)
	req := &http.Request{}
	result := widget.Render(req, "Test content", map[string]string{})
	if result != "" {
		t.Error("Render() with unauthenticated request should return empty string")
	}
}

func TestNewUnauthenticatedWidget(t *testing.T) {
	// Test with nil registry
	widget := NewUnauthenticatedWidget(nil)
	if widget == nil {
		t.Fatal("NewUnauthenticatedWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewUnauthenticatedWidget(registry)
	if widget == nil {
		t.Fatal("NewUnauthenticatedWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

func TestUnauthenticatedWidget_Alias(t *testing.T) {
	registry := testutils.Setup()
	widget := NewUnauthenticatedWidget(registry)

	alias := widget.Alias()
	if alias != "x-unauthenticated" {
		t.Errorf("Alias() = %q, want %q", alias, "x-unauthenticated")
	}
}

func TestUnauthenticatedWidget_Description(t *testing.T) {
	registry := testutils.Setup()
	widget := NewUnauthenticatedWidget(registry)

	desc := widget.Description()
	if desc == "" {
		t.Error("Description() should return non-empty string")
	}
	if desc != "Renders the content if the user is not authenticated" {
		t.Errorf("Description() = %q, want %q", desc, "Renders the content if the user is not authenticated")
	}
}

func TestUnauthenticatedWidget_Render(t *testing.T) {
	registry := testutils.Setup()
	widget := NewUnauthenticatedWidget(registry)

	// Test with nil request (no auth user)
	req := &http.Request{}
	result := widget.Render(req, "Test content", map[string]string{})
	if result != "Test content" {
		t.Errorf("Render() = %q, want %q", result, "Test content")
	}
}

func TestNewPrintWidget(t *testing.T) {
	// Test with nil registry
	widget := NewPrintWidget(nil)
	if widget == nil {
		t.Fatal("NewPrintWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewPrintWidget(registry)
	if widget == nil {
		t.Fatal("NewPrintWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

func TestPrintWidget_Alias(t *testing.T) {
	registry := testutils.Setup()
	widget := NewPrintWidget(registry)

	alias := widget.Alias()
	if alias != "x-print" {
		t.Errorf("Alias() = %q, want %q", alias, "x-print")
	}
}

func TestPrintWidget_Description(t *testing.T) {
	registry := testutils.Setup()
	widget := NewPrintWidget(registry)

	desc := widget.Description()
	if desc == "" {
		t.Error("Description() should return non-empty string")
	}
	if desc != "Renders the result of the provided content" {
		t.Errorf("Description() = %q, want %q", desc, "Renders the result of the provided content")
	}
}

func TestPrintWidget_Render(t *testing.T) {
	registry := testutils.Setup()
	widget := NewPrintWidget(registry)

	// Test with valid request
	testURL, err := url.Parse("/test")
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}
	req := &http.Request{
		URL: testURL,
	}
	result := widget.Render(req, "'hello'", map[string]string{})
	if result == "" {
		t.Error("Render() should return non-empty result")
	}

	// Test with path variable
	result = widget.Render(req, "path", map[string]string{})
	if result == "" {
		t.Error("Render() with path should return non-empty result")
	}
}

func TestNewVisibleWidget(t *testing.T) {
	// Test with nil registry
	widget := NewVisibleWidget(nil)
	if widget == nil {
		t.Fatal("NewVisibleWidget(nil) should return non-nil")
	}
	if widget.registry != nil {
		t.Error("widget.registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widget = NewVisibleWidget(registry)
	if widget == nil {
		t.Fatal("NewVisibleWidget(registry) should return non-nil")
	}
	if widget.registry == nil {
		t.Error("widget.registry should not be nil when passed valid registry")
	}
}

func TestVisibleWidget_Alias(t *testing.T) {
	registry := testutils.Setup()
	widget := NewVisibleWidget(registry)

	alias := widget.Alias()
	if alias != "x-visible" {
		t.Errorf("Alias() = %q, want %q", alias, "x-visible")
	}
}

func TestVisibleWidget_Description(t *testing.T) {
	registry := testutils.Setup()
	widget := NewVisibleWidget(registry)

	desc := widget.Description()
	if desc == "" {
		t.Error("Description() should return non-empty string")
	}
	if desc != "Renders the content if the condition is met" {
		t.Errorf("Description() = %q, want %q", desc, "Renders the content if the condition is met")
	}
}

func TestVisibleWidget_Render(t *testing.T) {
	registry := testutils.Setup()
	widget := NewVisibleWidget(registry)

	// Test with nil request
	req := &http.Request{}
	result := widget.Render(req, "Test content", map[string]string{})
	// Result depends on environment matching, but should not panic
	_ = result
}

func TestRoutes(t *testing.T) {
	// Test with nil registry
	routes := Routes(nil)
	if routes == nil {
		t.Error("Routes(nil) should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	routes = Routes(registry)
	if routes == nil {
		t.Error("Routes(registry) should not return nil")
	}
}
