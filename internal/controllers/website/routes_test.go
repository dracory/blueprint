package website

import (
	"testing"

	"project/internal/testutils"
)

func TestRoutes(t *testing.T) {
	// Test with nil registry
	routes := Routes(nil)
	if routes == nil {
		t.Error("Routes() should not return nil")
	}
	if len(routes) != 0 {
		t.Errorf("Routes(nil) should return empty slice, got %d routes", len(routes))
	}

	// Test with registry with config
	registry := testutils.Setup()
	routes = Routes(registry)
	if routes == nil {
		t.Error("Routes() should not return nil")
	}
	// Should have favicon, blog, contact, seo, swagger routes (minimum 5)
	if len(routes) < 5 {
		t.Errorf("Routes(registry with config) should return at least 5 routes, got %d", len(routes))
	}
}
