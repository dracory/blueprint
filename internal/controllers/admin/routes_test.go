package admin

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

	routes := Routes(app)

	if routes == nil {
		t.Error("Routes() returned nil")
	}
}

// TestRoutesReturnsRoutes verifies Routes returns route slice
func TestRoutesReturnsRoutes(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	// Should return at least home and homeCatchAll routes
	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes (home, homeCatchAll), got %d", len(routes))
	}
}

// TestRoutesWithRequiredStores verifies routes with stores enabled
func TestRoutesWithRequiredStores(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithTaskStore(true),
	)
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	if routes == nil {
		t.Error("Routes() returned nil")
	}

	// Should have more routes with stores enabled
	t.Logf("Total routes with stores: %d", len(routes))
}
