package stats

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewStatsController(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	if controller == nil {
		t.Error("NewStatsController() should not return nil")
	}
	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
	if controller.logger == nil {
		t.Error("Controller logger should not be nil")
	}
}

func TestStatsController_Handler_DefaultAction(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_Handler_WithDifferentMethods_GET(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_Handler_WithDifferentMethods_POST(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("POST", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_RegistryField(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
}

func TestStatsController_NilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil app
		}
	}()

	controller := NewStatsController(nil)
	if controller == nil {
		t.Error("NewStatsController(nil) should not return nil")
	}
}

func TestStatsController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewStatsController(registry1)
	controller2 := NewStatsController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.app != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.app != registry2 {
		t.Error("Controller2 should have registry2")
	}
}
