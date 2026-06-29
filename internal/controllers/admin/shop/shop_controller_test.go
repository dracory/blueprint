package admin

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewHomeController(t *testing.T) {
	app := testutils.Setup()
	controller := NewHomeController(app)

	if controller == nil {
		t.Error("NewHomeController() should not return nil")
	}
}

func TestHomeController_Handler(t *testing.T) {
	app := testutils.Setup()
	controller := NewHomeController(app)

	req := httptest.NewRequest("GET", "/admin/shop", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler() should return non-empty string")
	}
}

func TestHomeController_NilRegistry(t *testing.T) {
	controller := NewHomeController(nil)
	if controller == nil {
		t.Error("NewHomeController(nil) should not return nil")
	}
	if controller.app != nil {
		t.Error("Controller app should be nil when passed nil")
	}
}

