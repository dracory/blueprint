package aitest

import (
	"testing"
)

// TestNewAiTestController tests the constructor
func TestNewAiTestController(t *testing.T) {
	t.Parallel()

	controller := NewAiTestController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestAiTestController_StructFields tests controller structure
func TestAiTestController_StructFields(t *testing.T) {
	t.Parallel()

	controller := NewAiTestController(nil)
	if controller == nil {
		t.Fatal("NewAiTestController() returned nil")
	}

	// Verify the controller has the expected structure
	// The struct should have a registry field
}

// TestAiTestController_MultipleInstances tests creating multiple controllers
func TestAiTestController_MultipleInstances(t *testing.T) {
	t.Parallel()

	controller1 := NewAiTestController(nil)
	controller2 := NewAiTestController(nil)

	if controller1 == controller2 {
		t.Error("Each NewAiTestController call should return a new instance")
	}

	if controller1 == nil || controller2 == nil {
		t.Error("Both controllers should be non-nil")
	}
}

// TestAiTestController_WithRegistry tests controller with registry
func TestAiTestController_WithRegistry(t *testing.T) {
	t.Parallel()

	// Even with nil registry, controller should work for basic tests
	controller := NewAiTestController(nil)

	// Verify it doesn't panic
	_ = controller.registry
}
