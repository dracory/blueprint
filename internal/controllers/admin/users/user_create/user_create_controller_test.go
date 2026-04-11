package user_create

import (
	"testing"

	"project/internal/testutils"
)

func TestNewUserCreateController(t *testing.T) {
	// Test with nil registry
	controller := NewUserCreateController(nil)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewUserCreateController(registry)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}
}

func TestUserCreateController_RegistryField(t *testing.T) {
	registry := testutils.Setup()
	controller := NewUserCreateController(registry)

	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestUserCreateController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewUserCreateController(registry1)
	controller2 := NewUserCreateController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.registry != registry2 {
		t.Error("Controller2 should have registry2")
	}
}

func TestUserCreateController_Handler_Actions(t *testing.T) {
	registry := testutils.Setup()
	controller := NewUserCreateController(registry)

	// Verify handler exists - methods cannot be nil in Go
	_ = controller.Handler
}
