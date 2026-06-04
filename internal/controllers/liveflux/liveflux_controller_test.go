package liveflux

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewController(t *testing.T) {
	// Test with nil app
	controller := NewController(nil)
	if controller == nil {
		t.Error("NewController() should not return nil")
	}

	// Test with valid app
	app := testutils.Setup()
	controller = NewController(app)
	if controller == nil {
		t.Error("NewController() should not return nil")
	}
	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
	if controller.Engine == nil {
		t.Error("Controller Engine should not be nil")
	}
}

func TestAppContextKey(t *testing.T) {
	if AppContextKey != "app" {
		t.Errorf("AppContextKey = %q, want %q", AppContextKey, "app")
	}
}

func TestControllerHandler(t *testing.T) {
	app := testutils.Setup()
	controller := NewController(app)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// Handler should return a string (HTML response)
	if result == "" {
		t.Error("Handler() returned empty string")
	}
}
