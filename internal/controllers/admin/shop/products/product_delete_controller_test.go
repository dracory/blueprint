package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewProductDeleteController verifies controller can be created
func TestNewProductDeleteController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductDeleteController(app)

	if controller == nil {
		t.Error("NewProductDeleteController() returned nil")
	}
}

// TestProductDeleteControllerRegistry verifies controller has registry
func TestProductDeleteControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductDeleteController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestProductDeleteControllerHandlerExists verifies Handler method exists
func TestProductDeleteControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductDeleteController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
