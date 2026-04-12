package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestBlogRoutesFunctionExists verifies Routes function is defined
func TestBlogRoutesFunctionExists(t *testing.T) {
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

// TestBlogRoutesNilRegistry verifies Routes handles nil registry
func TestBlogRoutesNilRegistry(t *testing.T) {
	routes, err := Routes(nil)

	if err == nil {
		t.Error("Routes(nil) should return error")
	}

	if routes != nil {
		t.Error("Routes(nil) should return nil routes")
	}
}

// TestBlogRoutesReturnsRoutes verifies Routes returns route slice
func TestBlogRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes() returned error: %v", err)
	}

	// Should return at least 1 route (blog)
	if len(routes) < 1 {
		t.Errorf("Expected at least 1 route (blog), got %d", len(routes))
	}
}
