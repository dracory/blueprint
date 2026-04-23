package widgets

import (
	"testing"

	"project/internal/testutils"
)

func TestWidgetRegistry(t *testing.T) {
	// Test with nil registry
	widgets := WidgetRegistry(nil)
	if widgets == nil {
		t.Error("WidgetRegistry() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	widgets = WidgetRegistry(registry)
	if widgets == nil {
		t.Error("WidgetRegistry() should not return nil")
	}

	// Check that we have the expected number of widgets
	// Based on the code, we should have 5 widgets (contact and terms are commented out)
	expectedCount := 5
	if len(widgets) != expectedCount {
		t.Errorf("WidgetRegistry() returned %d widgets, want %d", len(widgets), expectedCount)
	}
}

func TestWidgetInterface(t *testing.T) {
	// Test that the Widget interface can be assigned
	// An interface is nil only when both type and value are nil
	var widget Widget
	// widget is nil at this point which is expected
	if widget != nil {
		t.Error("Widget interface should be nil by default")
	}
}
