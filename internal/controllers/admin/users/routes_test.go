package users

import (
	"testing"

	"project/internal/testutils"
)

func TestUserRoutes(t *testing.T) {
	// Test with nil registry
	routes := UserRoutes(nil)
	if routes == nil {
		t.Error("UserRoutes() should not return nil")
	}
	if len(routes) != 7 {
		t.Errorf("UserRoutes(nil) should return 7 routes, got %d", len(routes))
	}

	// Test with valid registry
	registry := testutils.Setup()
	routes = UserRoutes(registry)
	if routes == nil {
		t.Error("UserRoutes() should not return nil")
	}
	if len(routes) != 7 {
		t.Errorf("UserRoutes(registry) should return 7 routes, got %d", len(routes))
	}
}
