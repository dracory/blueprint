package shared_test

import (
	shared "project/internal/controllers/shared"
	approutes "project/internal/routes"
	"testing"

	"project/internal/testutils"
)

// TestSharedRoutesCount verifies the number of shared routes registered.
func TestSharedRoutesCount(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))
	routes := shared.Routes(app)
	if len(routes) != 8 {
		t.Fatalf("expected 8 shared routes, got %d", len(routes))
	}
}

// TestSharedRoutesNotNil ensures no route entries are nil.
func TestSharedRoutesNotNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))
	routes := shared.Routes(app)
	for i, rt := range routes {
		if rt == nil {
			t.Fatalf("route at index %d is nil", i)
		}
	}
}

func TestSharedRoutesAreAdded(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))
	routes := shared.Routes(app)
	expectedPaths := []string{
		"/ads.txt",
		"/files/*",
		"/flash",
		"/media/*",
		"/resources",
		// "/th/{extension:[a-z]+}/{size:[0-9x]+}/{quality:[0-9]+}/*",
		"/th/:extension/:size/:quality/:path",
		"/th/:extension/:size/:quality/:path...",
	}
	for _, expectedPath := range expectedPaths {
		found := false
		for _, rt := range routes {
			t.Log(rt.GetPath(), expectedPath)
			if rt.GetPath() == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected route with path %s not found", expectedPath)
		}
	}
}

// TestSharedRoutesAreInGlobalRouter verifies that all shared routes are added to the application router.
func TestSharedRoutesAreInGlobalRouter(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	sharedRoutes := shared.Routes(app)
	_, allRoutes := approutes.RoutesList(app)

	for i, s := range sharedRoutes {
		found := false
		for _, r := range allRoutes {
			if r.GetPath() == s.GetPath() { // same underlying route instance
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("shared route at index %d is not present in the global router list", i)
		}
	}
}
