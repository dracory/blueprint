package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestMediaRoutesFunctionExists verifies Routes function is defined
func TestMediaRoutesFunctionExists(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	if routes == nil {
		t.Error("Routes() returned nil routes")
	}
}

// TestMediaRoutesNilRegistry verifies Routes handles nil registry
func TestMediaRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestMediaRoutesReturnsRoutes verifies Routes returns route slice
func TestMediaRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return media manager route
	if len(routes) < 1 {
		t.Errorf("Expected at least 1 route (media manager), got %d", len(routes))
	}
}
