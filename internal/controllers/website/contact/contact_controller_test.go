package contact

import (
	"testing"

	"project/internal/testutils"
)

func TestNewContactController(t *testing.T) {
	// Test with valid app
	app := testutils.Setup()
	controller := NewContactController(app)
	if controller == nil {
		t.Error("NewContactController() should not return nil")
	}
	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
}

func TestNewFormContact(t *testing.T) {
	// Test with nil app
	component := NewFormContact(nil)
	if component == nil {
		t.Error("NewFormContact() should not return nil")
	}

	// Test with valid app
	app := testutils.Setup()
	component = NewFormContact(app)
	if component == nil {
		t.Error("NewFormContact() should not return nil")
	}
}
