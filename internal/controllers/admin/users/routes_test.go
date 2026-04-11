package users

import (
	"testing"

	"project/internal/testutils"
)

func TestUserRoutes(t *testing.T) {
	// Test with nil registry
	routes, err := UserRoutes(nil)
	if err == nil {
		t.Error("UserRoutes(nil) should return an error")
	}
	if routes != nil {
		t.Error("UserRoutes(nil) should return nil routes")
	}

	// Test with valid registry
	registry := testutils.Setup()
	routes, err = UserRoutes(registry)
	if err != nil {
		t.Errorf("UserRoutes(registry) should not return an error, got: %v", err)
	}
	if routes == nil {
		t.Error("UserRoutes() should not return nil")
	}
	if len(routes) != 7 {
		t.Errorf("UserRoutes(registry) should return 7 routes, got %d", len(routes))
	}
}
