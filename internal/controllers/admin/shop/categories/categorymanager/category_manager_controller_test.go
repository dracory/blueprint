package categorymanager

import (
	"testing"

	"project/internal/testutils"
)

// TestNewCategoryManagerController verifies controller can be created
func TestNewCategoryManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	if controller == nil {
		t.Error("NewCategoryManagerController() returned nil")
	}
}

// TestCategoryManagerControllerRegistry verifies controller has registry
func TestCategoryManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryManagerControllerHandlerExists verifies Handler method exists
func TestCategoryManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
