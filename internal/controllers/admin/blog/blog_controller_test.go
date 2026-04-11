package admin

import (
	"testing"

	"project/internal/testutils"
)

func TestBlogRoutes(t *testing.T) {
	registry := testutils.Setup()

	routes, err := Routes(registry)
	if routes == nil {
		t.Error("Routes() should return non-nil routes")
	}
	if err != nil {
		t.Logf("Routes() error: %v", err)
	}
}

func TestBlogRoutesWithNilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil registry
		}
	}()

	routes, err := Routes(nil)
	if routes != nil || err == nil {
		t.Error("Routes(nil) should handle nil registry")
	}
}

func TestBlogRoutesMultipleCalls(t *testing.T) {
	registry := testutils.Setup()

	routes1, err1 := Routes(registry)
	routes2, err2 := Routes(registry)

	if routes1 == nil || routes2 == nil {
		t.Error("Routes() should return non-nil routes on multiple calls")
	}
	if err1 != nil || err2 != nil {
		t.Logf("Routes() errors: %v, %v", err1, err2)
	}
}

func TestBlogRoutesWithDifferentRegistries(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	routes1, _ := Routes(registry1)
	routes2, _ := Routes(registry2)

	if routes1 == nil || routes2 == nil {
		t.Error("Routes() should return non-nil routes for different registries")
	}
}

func TestBlogRoutesIntegration(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
		testutils.WithBlogStore(true),
	)

	routes, err := Routes(registry)
	if routes == nil {
		t.Fatal("Routes() should return non-nil routes")
	}
	if err != nil {
		t.Logf("Routes() error: %v", err)
	}
}

func TestBlogRoutesWithCacheStore(t *testing.T) {
	registry := testutils.Setup()

	routes, _ := Routes(registry)
	if routes == nil {
		t.Fatal("Routes() should return non-nil routes")
	}

	if registry.GetCacheStore() == nil {
		t.Error("Registry should have cache store")
	}
}

func TestBlogRoutesWithLogger(t *testing.T) {
	registry := testutils.Setup()

	routes, _ := Routes(registry)
	if routes == nil {
		t.Fatal("Routes() should return non-nil routes")
	}

	if registry.GetLogger() == nil {
		t.Error("Registry should have logger")
	}
}
