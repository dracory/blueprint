package liveflux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/rtr"
)

func TestController_HandlerWithNilRegistry(t *testing.T) {
	// Test handler with nil app
	controller := NewController(nil)
	if controller == nil {
		t.Fatal("NewController(nil) should not return nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// Handler should still return a result even with nil app
	if result == "" {
		t.Error("Handler() with nil app returned empty string")
	}
}

func TestController_HandlerWithHeaders_GET(t *testing.T) {
	app := testutils.Setup()
	controller := NewController(app)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// Handler should return a result for supported methods
	if result == "" {
		t.Error("Handler() returned empty string for GET method")
	}
}

func TestController_HandlerWithHeaders_POST(t *testing.T) {
	app := testutils.Setup()
	controller := NewController(app)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	req := httptest.NewRequest("POST", "/liveflux", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// Handler should return a result for supported methods
	if result == "" {
		t.Error("Handler() returned empty string for POST method")
	}
}

func TestController_HandlerWithContext(t *testing.T) {
	app := testutils.Setup()
	controller := NewController(app)
	if controller == nil {
		t.Fatal("NewController() returned nil")
	}

	req := httptest.NewRequest("GET", "/liveflux", nil)
	w := httptest.NewRecorder()

	// The handler should add the app to context
	result := controller.Handler(w, req)

	if result == "" {
		t.Error("Handler() returned empty string")
	}
}

func TestRoutes(t *testing.T) {
	// Test with nil app
	routes := Routes(nil)
	if routes == nil {
		t.Error("Routes(nil) should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("Routes(nil) returned %d routes, want 2", len(routes))
	}

	// Test with valid app
	app := testutils.Setup()
	routes = Routes(app)
	if routes == nil {
		t.Error("Routes(app) should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("Routes(app) returned %d routes, want 2", len(routes))
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
	app := testutils.Setup()
	routes := Routes(app)

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
	app := testutils.Setup()
	routes := Routes(app)

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
	app := testutils.Setup()
	controller := NewController(app)

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

	if controller.app != nil {
		t.Error("Controller app should be nil when passed nil")
	}

	if controller.Engine == nil {
		t.Error("Controller Engine should not be nil even with nil app")
	}
}

func TestRoutesReturnType(t *testing.T) {
	app := testutils.Setup()
	routes := Routes(app)

	// Verify return type is []rtr.RouteInterface
	var _ []rtr.RouteInterface = routes

	// Verify each route implements the interface
	for i, route := range routes {
		if route == nil {
			t.Errorf("Route %d is nil", i)
		}
	}
}
