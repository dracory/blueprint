package aitest

import (
	"testing"

	"project/internal/testutils"
)

func TestNewAiTestController(t *testing.T) {
	// Test with nil registry
	controller := NewAiTestController(nil)
	if controller == nil {
		t.Error("NewAiTestController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiTestController(registry)
	if controller == nil {
		t.Error("NewAiTestController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}
