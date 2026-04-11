package admin

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewHomeController(t *testing.T) {
	registry := testutils.Setup()
	controller := NewHomeController(registry)

	if controller == nil {
		t.Error("NewHomeController() should not return nil")
	}
}

func TestHomeController_Handler(t *testing.T) {
	registry := testutils.Setup()
	controller := NewHomeController(registry)

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
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}

func TestNewOrderManagerController(t *testing.T) {
	controller := NewOrderManagerController()

	if controller == nil {
		t.Error("NewOrderManagerController() should not return nil")
	}
}

func TestOrderManagerController_Handler(t *testing.T) {
	controller := NewOrderManagerController()

	req := httptest.NewRequest("GET", "/admin/shop/orders", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result != "Order Manager" {
		t.Errorf("Handler() = %q, want %q", result, "Order Manager")
	}
}
