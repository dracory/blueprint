package stats

import (
	"testing"

	"project/internal/testutils"
)

// TestRoutesFunctionExists verifies Routes function is defined
func TestRoutesFunctionExists(t *testing.T) {
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

// TestRoutesNilRegistry verifies Routes handles nil registry
func TestRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestRoutesReturnsRoutes verifies Routes returns route slice
func TestRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return statsHome and statsCatchAll routes
	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes (statsHome, statsCatchAll), got %d", len(routes))
	}
}
