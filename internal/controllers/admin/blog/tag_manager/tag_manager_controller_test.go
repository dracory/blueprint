package tag_manager

import (
	"testing"

	"project/internal/testutils"
)

// TestNewTagManagerController verifies controller can be created
func TestNewTagManagerController(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewTagManagerController(app)

	if controller == nil {
		t.Error("NewTagManagerController() returned nil")
	}
}

// TestTagManagerControllerRegistry verifies controller has registry
func TestTagManagerControllerRegistry(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewTagManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestTagManagerControllerHandlerExists verifies Handler method exists
func TestTagManagerControllerHandlerExists(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewTagManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
