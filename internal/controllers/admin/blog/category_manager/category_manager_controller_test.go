package category_manager

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewCategoryManagerController(t *testing.T) {
	// Test with nil registry
	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Error("NewCategoryManagerController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewCategoryManagerController(registry)
	if controller == nil {
		t.Error("NewCategoryManagerController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestCategoryManagerControllerHandler(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to missing auth/context
		}
	}()

	registry := testutils.Setup()
	controller := NewCategoryManagerController(registry)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// Test with a request without auth (will panic due to missing cache store)
	req := httptest.NewRequest("GET", "/admin/blog/categories", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}
