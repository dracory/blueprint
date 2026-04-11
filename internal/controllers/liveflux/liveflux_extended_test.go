package liveflux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/rtr"
)

func TestController_HandlerWithNilRegistry(t *testing.T) {
	// Test handler with nil registry
	controller := NewController(nil)
	if controller == nil {
		t.Fatal("NewController(nil) should not return nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// Handler should still return a result even with nil registry
	if result == "" {
		t.Error("Handler() with nil registry returned empty string")
	}
}

func TestController_HandlerWithHeaders(t *testing.T) {
	registry := testutils.Setup()
	controller := NewController(registry)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	// Test with GET and POST methods (the ones supported by Routes)
	methods := []string{"GET", "POST"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/liveflux", nil)
			w := httptest.NewRecorder()

			result := controller.Handler(w, req)

			// Handler should return a result for supported methods
			if result == "" {
				t.Errorf("Handler() returned empty string for %s method", method)
			}
		})
	}
}

func TestController_HandlerWithContext(t *testing.T) {
	registry := testutils.Setup()
	controller := NewController(registry)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	// The handler should add the registry to context
	result := controller.Handler(w, req)

	if result == "" {
		t.Error("Handler() returned empty string")
	}
}

func TestRoutes(t *testing.T) {
	// Test with nil registry
	routes := Routes(nil)
	if routes == nil {
		t.Error("Routes(nil) should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("Routes(nil) returned %d routes, want 2", len(routes))
	}

	// Test with valid registry
	registry := testutils.Setup()
	routes = Routes(registry)
	if routes == nil {
		t.Error("Routes(registry) should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("Routes(registry) returned %d routes, want 2", len(routes))
	}

	// Verify route properties
	for i, route := range routes {
		if route == nil {
			t.Errorf("Route %d is nil", i)
			continue
		}

		// Check that each route has the expected path
		if route.GetPath() != "/liveflux" {
			t.Errorf("Route %d path = %q, want %q", i, route.GetPath(), "/liveflux")
		}
	}
}

func TestRoutesMethods(t *testing.T) {
	registry := testutils.Setup()
	routes := Routes(registry)

	if len(routes) != 2 {
		t.Fatalf("Routes() returned %d routes, want 2", len(routes))
	}

	// First route should be POST
	if routes[0].GetMethod() != http.MethodPost {
		t.Errorf("Route 0 method = %q, want %q", routes[0].GetMethod(), http.MethodPost)
	}

	// Second route should be GET
	if routes[1].GetMethod() != http.MethodGet {
		t.Errorf("Route 1 method = %q, want %q", routes[1].GetMethod(), http.MethodGet)
	}
}

func TestRoutesHandlerAssignment(t *testing.T) {
	registry := testutils.Setup()
	routes := Routes(registry)

	if len(routes) != 2 {
		t.Fatalf("Routes() returned %d routes, want 2", len(routes))
	}

	// Verify both routes have handlers assigned
	for i, route := range routes {
		if route.GetHTMLHandler() == nil {
			t.Errorf("Route %d has no HTML handler assigned", i)
		}
	}
}

func TestContextKeyType(t *testing.T) {
	// Verify contextKey is a string type
	var key contextKey = "test"
	if string(key) != "test" {
		t.Error("contextKey should be convertible to string")
	}
}

func TestControllerEngine(t *testing.T) {
	registry := testutils.Setup()
	controller := NewController(registry)

	if controller.Engine == nil {
		t.Error("Controller Engine should not be nil")
	}

	// Test that the engine is an http.Handler
	var _ http.Handler = controller.Engine
}

func TestNewControllerWithNil(t *testing.T) {
	controller := NewController(nil)

	if controller == nil {
		t.Fatal("NewController(nil) should not return nil")
	}

	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}

	if controller.Engine == nil {
		t.Error("Controller Engine should not be nil even with nil registry")
	}
}

func TestRoutesReturnType(t *testing.T) {
	registry := testutils.Setup()
	routes := Routes(registry)

	// Verify return type is []rtr.RouteInterface
	var _ []rtr.RouteInterface = routes

	// Verify each route implements the interface
	for i, route := range routes {
		if route == nil {
			t.Errorf("Route %d is nil", i)
		}
	}
}
