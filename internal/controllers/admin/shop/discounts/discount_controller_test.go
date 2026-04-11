package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestNewDiscountController verifies controller can be created
func TestNewDiscountController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	if controller == nil {
		t.Error("NewDiscountController() returned nil")
	}
}

// TestDiscountControllerRegistry verifies controller has registry
func TestDiscountControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestDiscountControllerAnyIndexExists verifies AnyIndex method exists
func TestDiscountControllerAnyIndexExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// Verify AnyIndex method exists (should compile without error)
	_ = controller.AnyIndex
}
