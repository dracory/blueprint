package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewMediaManagerController verifies controller can be created
func TestNewMediaManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	if controller == nil {
		t.Error("NewMediaManagerController() returned nil")
	}
}

// TestMediaManagerControllerRegistry verifies controller has registry
func TestMediaManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestMediaManagerControllerAnyIndexExists verifies AnyIndex method exists
func TestMediaManagerControllerAnyIndexExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewMediaManagerController(app)

	// Verify AnyIndex method exists (should compile without error)
	_ = controller.AnyIndex
}
