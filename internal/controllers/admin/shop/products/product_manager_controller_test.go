package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewProductManagerController verifies controller can be created
func TestNewProductManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	if controller == nil {
		t.Error("NewProductManagerController() returned nil")
	}
}

// TestProductManagerControllerRegistry verifies controller has registry
func TestProductManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestProductManagerControllerHandlerExists verifies Handler method exists
func TestProductManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
