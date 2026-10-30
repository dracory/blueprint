package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/links"
	"project/internal/testutils"
)

// containsIgnoreCase checks if substring is contained in s (case-insensitive)
func containsIgnoreCase(s, substring string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substring))
}

// routePaths flattens route paths for membership checks.
func routePaths(routes []interface {
	GetPath() string
}) map[string]bool {
	paths := map[string]bool{}
	for _, route := range routes {
		paths[route.GetPath()] = true
	}
	return paths
}

// TestRoutesConfiguration verifies the /auth surface: the package mux
// catch-all plus the app-owned register and logout routes.
func TestRoutesConfiguration(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes := Routes(app)
	if len(routes) < 4 {
		t.Fatalf("Expected at least 4 routes (register get+post, logout, package base+catch-all), got %d", len(routes))
	}

	var foundPackageBase, foundCatchAll, foundLogout, foundRegister bool
	for _, route := range routes {
		switch route.GetPath() {
		case "/auth":
			foundPackageBase = true
			if !containsIgnoreCase(route.GetName(), "package") {
				t.Errorf("Expected /auth route name to mention package, got '%s'", route.GetName())
			}
		case "/auth/*":
			foundCatchAll = true
			if len(route.GetBeforeMiddlewares()) == 0 {
				t.Error("Expected rate-limit middleware on /auth/* catch-all")
			}
		case links.AUTH_LOGOUT:
			foundLogout = true
			if !containsIgnoreCase(route.GetName(), "logout") {
				t.Errorf("Expected logout route name to contain 'logout', got '%s'", route.GetName())
			}
		case links.AUTH_REGISTER:
			if route.GetMethod() == "GET" {
				foundRegister = true
			}
		}
	}

	if !foundPackageBase {
		t.Error("Package base route /auth not found")
	}
	if !foundCatchAll {
		t.Error("Package catch-all route /auth/* not found")
	}
	if !foundLogout {
		t.Error("Logout route not found")
	}
	if !foundRegister {
		t.Error("Register route not found (registration is enabled in test config)")
	}
}

// TestRoutesNoLegacyPaths verifies the removed per-controller paths do
// not resurface: secondary methods now live under the package's
// /auth/login-<name> convention served by the mux, not as rtr routes.
func TestRoutesNoLegacyPaths(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()

	legacy := []string{
		"/auth/otp-login",
		"/auth/magiclink-login",
		"/auth/password-login",
		"/auth/authknight-login",
		"/auth/forgot-password",
	}

	for _, route := range Routes(app) {
		for _, old := range legacy {
			if route.GetPath() == old {
				t.Errorf("Legacy path %s should not be mounted", old)
			}
		}
	}
}

// TestPackageMuxServesLoginPage hits the package mux handler directly to
// confirm the composed mux serves the default OTP login page.
func TestPackageMuxServesLoginPage(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	var muxHandler func(http.ResponseWriter, *http.Request)
	for _, route := range packageAuthRoutes(app) {
		if route.GetPath() == "/auth/*" {
			muxHandler = route.GetHandler()
		}
	}
	if muxHandler == nil {
		t.Fatal("catch-all package route has no handler")
	}

	req := httptest.NewRequest("GET", "/auth/login", nil)
	rec := httptest.NewRecorder()
	muxHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 from /auth/login, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<html") {
		t.Error("Expected HTML login page from package mux")
	}
}

// TestRoutesUniqueRouteNames ensures no duplicate route names are
// registered across the auth route table.
func TestRoutesUniqueRouteNames(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()

	names := map[string]int{}
	for _, route := range Routes(app) {
		names[route.GetName()]++
	}
	for name, count := range names {
		if count > 1 {
			t.Errorf("Duplicate route name %q (%d registrations)", name, count)
		}
	}
}

// TestRegisterRoutesPresence verifies the register route mounts when
// registration is enabled (the test config default).
func TestRegisterRoutesPresence(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()

	var registerGet, registerPost bool
	for _, route := range Routes(app) {
		if route.GetPath() != links.AUTH_REGISTER {
			continue
		}
		switch route.GetMethod() {
		case "GET":
			registerGet = true
		case "POST":
			registerPost = true
		}
		if len(route.GetBeforeMiddlewares()) == 0 {
			t.Error("Expected rate-limit middleware on register route")
		}
	}
	if !registerGet {
		t.Error("Register GET route not found")
	}
	if !registerPost {
		t.Error("Register POST route not found")
	}
}

// TestLoginAlternatives verifies sibling cross-link generation for a
// multi-method config.
func TestLoginAlternatives(t *testing.T) {
	t.Parallel()

	methods := []string{config.LOGIN_METHOD_OTP, config.LOGIN_METHOD_MAGICLINK}
	alts := loginAlternatives(methods, 0)
	if len(alts) != 1 {
		t.Fatalf("Expected 1 alternative for primary method, got %d", len(alts))
	}
	if alts[0].URL != "/auth/login-magiclink" {
		t.Errorf("Expected /auth/login-magiclink, got %s", alts[0].URL)
	}
}
