package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestShopRoutesFunctionExists verifies ShopRoutes function is defined
func TestShopRoutesFunctionExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := ShopRoutes(app)

	if err != nil {
		t.Errorf("ShopRoutes() returned error: %v", err)
	}

	if routes == nil {
		t.Error("ShopRoutes() returned nil routes")
	}
}

// TestShopRoutesNilRegistry verifies ShopRoutes handles nil registry
func TestShopRoutesNilRegistry(t *testing.T) {
	t.Parallel()
	routes, err := ShopRoutes(nil)

	if err == nil {
		t.Error("ShopRoutes(nil) should return error")
	}

	if routes != nil {
		t.Error("ShopRoutes(nil) should return nil routes")
	}
}

// TestShopRoutesReturnsRoutes verifies ShopRoutes returns route slice
func TestShopRoutesReturnsRoutes(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	routes, err := ShopRoutes(app)

	if err != nil {
		t.Errorf("ShopRoutes() returned error: %v", err)
	}

	// Should return shopOrders and shopCatchAll routes
	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes (shopOrders, shopCatchAll), got %d", len(routes))
	}
}
