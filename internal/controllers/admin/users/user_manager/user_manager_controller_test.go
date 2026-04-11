package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewUserManagerController verifies controller can be created
func TestNewUserManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewUserManagerController(app)

	if controller == nil {
		t.Error("NewUserManagerController() returned nil")
	}
}

// TestUserManagerControllerRegistry verifies controller has registry
func TestUserManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewUserManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestUserManagerControllerHandlerExists verifies Handler method exists
func TestUserManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewUserManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

func TestUserManagerController_NilRegistry(t *testing.T) {
	t.Parallel()
	controller := NewUserManagerController(nil)
	if controller == nil {
		t.Error("NewUserManagerController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}

func TestUserManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()
	app1 := testutils.Setup()
	app2 := testutils.Setup()

	controller1 := NewUserManagerController(app1)
	controller2 := NewUserManagerController(app2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != app1 {
		t.Error("Controller1 should have app1")
	}

	if controller2.registry != app2 {
		t.Error("Controller2 should have app2")
	}
}

func TestUserManagerController_Handler_Actions(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	controller := NewUserManagerController(app)

	// Verify handler exists - methods cannot be nil in Go
	// This test ensures the Handler method is callable
	_ = controller.Handler
}
