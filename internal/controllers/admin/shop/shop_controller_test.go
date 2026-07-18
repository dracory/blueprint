package admin

import (
	"testing"

	"project/internal/testutils"
)

func TestNewShopAdminController(t *testing.T) {
	app := testutils.Setup()
	controller := NewShopAdminController(app)

	if controller == nil {
		t.Error("NewShopAdminController() should not return nil")
	}
}

func TestNewShopAdminController_NilApp(t *testing.T) {
	controller := NewShopAdminController(nil)
	if controller == nil {
		t.Error("NewShopAdminController(nil) should not return nil")
	}
	if controller.app != nil {
		t.Error("Controller app should be nil when passed nil")
	}
}
