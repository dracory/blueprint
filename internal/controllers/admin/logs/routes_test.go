package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestLogsRoutesFunctionExists verifies Routes function is defined
func TestLogsRoutesFunctionExists(t *testing.T) {
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

// TestLogsRoutesNilRegistry verifies Routes handles nil registry
func TestLogsRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestLogsRoutesReturnsRoutes verifies Routes returns route slice
func TestLogsRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return logs route
	if len(routes) != 1 {
		t.Errorf("Expected 1 route (logs), got %d", len(routes))
	}
}
