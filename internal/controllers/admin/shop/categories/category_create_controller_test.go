package categories

import (
	"testing"

	"project/internal/testutils"
)

// TestNewCategoryCreateController verifies controller can be created
func TestNewCategoryCreateController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	if controller == nil {
		t.Error("NewCategoryCreateController() returned nil")
	}
}

// TestCategoryCreateControllerRegistry verifies controller has registry
func TestCategoryCreateControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryCreateControllerHandlerExists verifies Handler method exists
func TestCategoryCreateControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
