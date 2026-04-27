package admin

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewFileManagerController(t *testing.T) {
	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	if controller == nil {
		t.Error("NewFileManagerController() should not return nil")
	}
}

func TestFileManagerController_Handler(t *testing.T) {
	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	req := httptest.NewRequest("GET", "/admin/files", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)

	if w.Code == 0 {
		t.Error("Handler() should write a response")
	}
}

func TestFileManagerController_RegistryField(t *testing.T) {
	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
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

	if controller1.registry != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.registry != registry2 {
		t.Error("Controller2 should have registry2")
	}
}

func TestFileManagerController_HandlerMultipleCalls(t *testing.T) {
	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/admin/files", nil)
		w := httptest.NewRecorder()

		controller.Handler(w, req)
	}
}

func TestFileManagerController_HandlerWithDifferentMethods(t *testing.T) {
	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/admin/files", nil)
			w := httptest.NewRecorder()

			controller.Handler(w, req)
		})
	}
}
