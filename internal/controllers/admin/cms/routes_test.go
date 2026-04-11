package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestCmsRoutesFunctionExists verifies Routes function is defined
func TestCmsRoutesFunctionExists(t *testing.T) {
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

// TestCmsRoutesNilRegistry verifies Routes handles nil registry
func TestCmsRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestCmsRoutesReturnsRoutes verifies Routes returns route slice
func TestCmsRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return cms manager route
	if len(routes) < 1 {
		t.Errorf("Expected at least 1 route (cms manager), got %d", len(routes))
	}
}
