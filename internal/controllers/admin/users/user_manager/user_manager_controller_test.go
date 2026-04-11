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
