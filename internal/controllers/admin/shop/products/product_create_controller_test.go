package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewProductCreateController verifies controller can be created
func TestNewProductCreateController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductCreateController(app)

	if controller == nil {
		t.Error("NewProductCreateController() returned nil")
	}
}

// TestProductCreateControllerRegistry verifies controller has registry
func TestProductCreateControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductCreateController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestProductCreateControllerHandlerExists verifies Handler method exists
func TestProductCreateControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductCreateController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
