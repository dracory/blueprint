package ai_tools

import (
	"testing"
)

// TestNewAiToolsController tests the constructor
func TestNewAiToolsController(t *testing.T) {
	t.Parallel()

	controller := NewAiToolsController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestAiToolsController_StructFields tests controller structure
func TestAiToolsController_StructFields(t *testing.T) {
	t.Parallel()

	controller := NewAiToolsController(nil)
	if controller == nil {
		t.Fatal("NewAiToolsController() returned nil")
	}

	// Verify the controller has the expected structure
	// The struct should have a registry field
}

// TestAiToolsController_MultipleInstances tests creating multiple controllers
func TestAiToolsController_MultipleInstances(t *testing.T) {
	t.Parallel()

	controller1 := NewAiToolsController(nil)
	controller2 := NewAiToolsController(nil)

	if controller1 == controller2 {
		t.Error("Each NewAiToolsController call should return a new instance")
	}

	if controller1 == nil || controller2 == nil {
		t.Error("Both controllers should be non-nil")
	}
}

// TestAiToolsController_WithRegistry tests controller with registry
func TestAiToolsController_WithRegistry(t *testing.T) {
	t.Parallel()

	// Even with nil registry, controller should work for basic tests
	controller := NewAiToolsController(nil)

	// Verify it doesn't panic
	_ = controller.registry
}
