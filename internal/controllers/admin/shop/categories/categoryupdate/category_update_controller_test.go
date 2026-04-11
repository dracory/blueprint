package categoryupdate

import (
	"testing"

	"project/internal/testutils"
)

// TestNewCategoryUpdateController verifies controller can be created
func TestNewCategoryUpdateController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	if controller == nil {
		t.Error("NewCategoryUpdateController() returned nil")
	}
}

// TestCategoryUpdateControllerRegistry verifies controller has registry
func TestCategoryUpdateControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryUpdateControllerHandlerExists verifies Handler method exists
func TestCategoryUpdateControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
