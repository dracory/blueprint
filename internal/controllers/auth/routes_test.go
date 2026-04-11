package auth

import (
	"strings"
	"testing"

	"project/internal/links"
	"project/internal/testutils"
)

// containsIgnoreCase checks if substring is contained in s (case-insensitive)
func containsIgnoreCase(s, substring string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substring))
}

// TestRoutesConfiguration verifies Routes function returns correct routes with expected configuration
func TestRoutesConfiguration(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	// Should return at least 3 routes (auth, login, logout)
	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes, got %d", len(routes))
	}

	// Verify auth route exists with correct path
	foundAuth := false
	foundLogin := false
	foundLogout := false

	for _, route := range routes {
		path := route.GetPath()
		switch path {
		case links.AUTH_AUTH:
			foundAuth = true
			if !containsIgnoreCase(route.GetName(), "auth") {
				t.Errorf("Expected auth route name to contain 'auth', got '%s'", route.GetName())
			}
		case links.AUTH_LOGIN:
			foundLogin = true
			if !containsIgnoreCase(route.GetName(), "login") {
				t.Errorf("Expected login route name to contain 'login', got '%s'", route.GetName())
			}
		case links.AUTH_LOGOUT:
			foundLogout = true
			if !containsIgnoreCase(route.GetName(), "logout") {
				t.Errorf("Expected logout route name to contain 'logout', got '%s'", route.GetName())
			}
		}
	}

	if !foundAuth {
		t.Error("Auth route not found")
	}
	if !foundLogin {
		t.Error("Login route not found")
	}
	if !foundLogout {
		t.Error("Logout route not found")
	}
}

// TestRoutesWithRegistrationEnabled verifies register route is included when enabled
func TestRoutesWithRegistrationEnabled(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	// Should have at least 4 routes with registration enabled (auth, login, logout, register)
	if len(routes) < 4 {
		t.Errorf("Expected at least 4 routes with registration enabled, got %d", len(routes))
	}

	// Verify register route exists
	foundRegister := false
	for _, route := range routes {
		if route.GetPath() == links.AUTH_REGISTER {
			foundRegister = true
			if !containsIgnoreCase(route.GetName(), "register") {
				t.Errorf("Expected register route name to contain 'register', got '%s'", route.GetName())
			}
			break
		}
	}

	if !foundRegister {
		t.Error("Register route not found when registration is enabled")
	}
}

// TestRoutesWithRegistrationDisabled verifies register route is excluded when disabled
func TestRoutesWithRegistrationDisabled(t *testing.T) {
	t.Parallel()
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(false)
	app := testutils.Setup(testutils.WithCfg(cfg))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)

	// Verify core routes are still present
	foundAuth := false
	foundLogin := false
	foundLogout := false

	// Verify register route is NOT included
	for _, route := range routes {
		switch route.GetPath() {
		case links.AUTH_AUTH:
			foundAuth = true
		case links.AUTH_LOGIN:
			foundLogin = true
		case links.AUTH_LOGOUT:
			foundLogout = true
		case links.AUTH_REGISTER:
			t.Error("Register route should not be present when registration is disabled")
		}
	}

	if !foundAuth {
		t.Error("Auth route not found when registration is disabled")
	}
	if !foundLogin {
		t.Error("Login route not found when registration is disabled")
	}
	if !foundLogout {
		t.Error("Logout route not found when registration is disabled")
	}
}
