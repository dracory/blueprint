package liveflux

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewController(t *testing.T) {
	// Test with nil registry
	controller := NewController(nil)
	if controller == nil {
		t.Error("NewController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewController(registry)
	if controller == nil {
		t.Error("NewController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
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
	registry := testutils.Setup()
	controller := NewController(registry)
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

	// Verify context was set correctly
	if req.Context().Value(AppContextKey) == nil {
		t.Error("Handler() did not set app context")
	}
}
