package contact

import (
	"testing"

	"project/internal/testutils"
)

func TestNewContactController(t *testing.T) {
	// Test with valid registry
	registry := testutils.Setup()
	controller := NewContactController(registry)
	if controller == nil {
		t.Error("NewContactController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestNewFormContact(t *testing.T) {
	// Test with nil registry
	component := NewFormContact(nil)
	if component == nil {
		t.Error("NewFormContact() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	component = NewFormContact(registry)
	if component == nil {
		t.Error("NewFormContact() should not return nil")
	}
}
