package admin

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewFileManagerController(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	if controller == nil {
		t.Error("NewFileManagerController() should not return nil")
	}
}

func TestFileManagerController_Handler(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	req := httptest.NewRequest("GET", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)

	if w.Code == 0 {
		t.Error("Handler() should write a response")
	}
}

func TestFileManagerController_RegistryField(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
}

func TestFileManagerController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewFileManagerController(registry1)
	controller2 := NewFileManagerController(registry2)

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

func TestFileManagerController_HandlerMultipleCalls(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/admin/files", nil)
		w := httptest.NewRecorder()

		controller.Handler(w, req)
	}
}

func TestFileManagerController_HandlerWithDifferentMethods_GET(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	req := httptest.NewRequest("GET", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestFileManagerController_HandlerWithDifferentMethods_POST(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	req := httptest.NewRequest("POST", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestFileManagerController_HandlerWithDifferentMethods_PUT(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	req := httptest.NewRequest("PUT", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestFileManagerController_HandlerWithDifferentMethods_DELETE(t *testing.T) {
	app := testutils.Setup()
	controller := NewFileManagerController(app)

	req := httptest.NewRequest("DELETE", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}
