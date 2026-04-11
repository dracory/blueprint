package admin

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestRoutesNilRegistry verifies Routes handles nil registry
func TestRoutesNilRegistry(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	routes, err := Routes(app)

	if err != nil {
		t.Errorf("Routes(app) should not return error, got: %v", err)
	}
	if routes == nil {
		t.Error("Routes(app) should return non-nil routes")
	}
	if len(routes) == 0 {
		t.Error("Routes(app) should return at least one route")
	}
}

// TestRoutesRouteName verifies route has correct name
func TestRoutesRouteName(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	routes, err := Routes(app)
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) == 0 {
		t.Fatal("Routes should return at least one route")
	}

	routeName := routes[0].GetName()
	if !strings.Contains(strings.ToLower(routeName), "file") {
		t.Errorf("Route name should contain 'file', got: %s", routeName)
	}
}
