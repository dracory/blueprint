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

func TestProductManagerController_NilRegistry(t *testing.T) {
	t.Parallel()
	controller := NewProductManagerController(nil)
	if controller == nil {
		t.Error("NewProductManagerController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}

func TestProductManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()
	app1 := testutils.Setup(testutils.WithShopStore(true))
	app2 := testutils.Setup(testutils.WithShopStore(true))
	if app1 == nil || app2 == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app1.GetDatabase().Close() })
	t.Cleanup(func() { _ = app2.GetDatabase().Close() })

	controller1 := NewProductManagerController(app1)
	controller2 := NewProductManagerController(app2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != app1 {
		t.Error("Controller1 should have app1")
	}

	if controller2.registry != app2 {
		t.Error("Controller2 should have app2")
	}
}

func TestProductManagerController_Handler_Actions(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	// Verify handler exists - methods cannot be nil in Go
	_ = controller.Handler
}
