package users

import (
	"testing"

	"project/internal/testutils"
)

// TestUsersRoutesFunctionExists verifies Routes function is defined
func TestUsersRoutesFunctionExists(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	if routes == nil {
		t.Error("Routes() returned nil routes")
	}
}

// TestUsersRoutesNilRegistry verifies Routes handles nil registry
func TestUsersRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestUsersRoutesReturnsRoutes verifies Routes returns route slice
func TestUsersRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return users route
	if len(routes) != 1 {
		t.Errorf("Expected 1 route (users), got %d", len(routes))
	}
}
