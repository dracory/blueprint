package auth

import (
	"testing"

	"project/internal/testutils"
)

// TestRoutesFunctionExists verifies Routes function is defined
func TestRoutesFunctionExists(t *testing.T) {
	// Create a test registry using testutils.Setup
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	// Routes should return a slice of routes
	routes := Routes(app)

	// Should return routes (at minimum auth, login, logout)
	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes (auth, login, logout), got %d", len(routes))
	}
}

// TestRoutesWithRegistrationEnabled verifies routes with registration enabled
func TestRoutesWithRegistrationEnabled(t *testing.T) {
	// Setup with user store enabled (registration needs user store)
	app := testutils.Setup(testutils.WithUserStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	// The Routes function checks GetRegistrationEnabled()
	routes := Routes(app)

	// Should have at least auth, login, logout (and possibly register if enabled)
	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes, got %d", len(routes))
	}
}

// TestRoutesReturnType verifies Routes returns slice
func TestRoutesReturnType(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	// Verify routes is not nil and is a slice
	if routes == nil {
		t.Error("Routes() returned nil")
	}
}
