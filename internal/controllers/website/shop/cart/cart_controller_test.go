package cart

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"project/internal/app"

	"github.com/dracory/shopstore"
	"github.com/dracory/userstore"
)

func createProduct(t *testing.T, app app.AppInterface) shopstore.ProductInterface {
	t.Helper()
	product := shopstore.NewProduct().
		SetTitle("Test Product").
		SetPrice("9.99").
		SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
	if err := app.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("ProductCreate failed: %v", err)
	}
	return product
}

func TestRoutes(t *testing.T) {
	if routes := Routes(nil); len(routes) != 0 {
		t.Errorf("Routes(nil) = %d routes, want 0", len(routes))
	}

	app := testutils.Setup()
	defer app.Close()

	if routes := Routes(app); len(routes) != 1 {
		t.Errorf("Routes(app) = %d routes, want 1", len(routes))
	}
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	app := testutils.Setup(testutils.WithCacheStore(true))
	defer app.Close()

	controller := NewCartController(app)
	req := httptest.NewRequest(http.MethodPatch, "/shop/cart/api", nil)
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandleGetCartGuest(t *testing.T) {
	app := testutils.Setup(testutils.WithCacheStore(true))
	defer app.Close()

	controller := NewCartController(app)
	req := httptest.NewRequest(http.MethodGet, "/shop/cart/api", nil)
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"status":"success"`) {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandleAddToCartGuest(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)
	defer app.Close()

	product := createProduct(t, app)

	controller := NewCartController(app)
	body := `{"product_id":"` + product.GetID() + `","quantity":2}`
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api", strings.NewReader(body))
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}

	// Cart should now be in cache — GET should return the item.
	req2 := httptest.NewRequest(http.MethodGet, "/shop/cart/api", nil)
	rec2 := httptest.NewRecorder()
	controller.Handler(rec2, req2)

	if !strings.Contains(rec2.Body.String(), product.GetID()) {
		t.Errorf("expected cart to contain product, body: %s", rec2.Body.String())
	}
	if !strings.Contains(rec2.Body.String(), "9.99") {
		t.Errorf("expected server-side price 9.99, body: %s", rec2.Body.String())
	}
}

func TestHandleAddToCartGuestInvalidBody(t *testing.T) {
	app := testutils.Setup(testutils.WithCacheStore(true))
	defer app.Close()

	controller := NewCartController(app)
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api", strings.NewReader("{invalid"))
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleAddToCartGuestMissingProduct(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)
	defer app.Close()

	controller := NewCartController(app)
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api",
		strings.NewReader(`{"quantity":1}`))
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleAddToCartGuestProductNotFound(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)
	defer app.Close()

	controller := NewCartController(app)
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api",
		strings.NewReader(`{"product_id":"nonexistent","quantity":1}`))
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleAddToCartGuestQuantityLimit(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)
	defer app.Close()

	product := createProduct(t, app)

	controller := NewCartController(app)
	body := `{"product_id":"` + product.GetID() + `","quantity":10000}`
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api", strings.NewReader(body))
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleCartAuthenticatedUser(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
	)
	defer app.Close()

	ctx := context.Background()
	user := userstore.NewUser().
		SetEmail("buyer@example.com").
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := app.GetUserStore().UserCreate(ctx, user); err != nil {
		t.Fatalf("UserCreate failed: %v", err)
	}

	product := createProduct(t, app)
	controller := NewCartController(app)

	// POST add to cart as authenticated user
	body := `{"product_id":"` + product.GetID() + `","quantity":3}`
	req := httptest.NewRequest(http.MethodPost, "/shop/cart/api", strings.NewReader(body))
	req, err := testutils.LoginAs(app, req, user)
	if err != nil {
		t.Fatalf("LoginAs failed: %v", err)
	}
	rec := httptest.NewRecorder()

	controller.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}

	// Verify cart persisted to user metadata
	fresh, err := app.GetUserStore().UserFindByID(ctx, user.GetID())
	if err != nil || fresh == nil {
		t.Fatalf("UserFindByID failed: %v", err)
	}
	cartJSON := fresh.GetMeta(config.USER_META_CART)
	if cartJSON == "" {
		t.Fatal("expected cart metadata on user")
	}
	var cart Cart
	if err := json.Unmarshal([]byte(cartJSON), &cart); err != nil {
		t.Fatalf("cart metadata is not valid JSON: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].ProductID != product.GetID() {
		t.Fatalf("unexpected cart items: %+v", cart.Items)
	}
	if cart.Items[0].Quantity != 3 {
		t.Errorf("quantity = %d, want 3", cart.Items[0].Quantity)
	}

	// GET as the same user returns the cart
	req2 := httptest.NewRequest(http.MethodGet, "/shop/cart/api", nil)
	req2, err = testutils.LoginAs(app, req2, fresh)
	if err != nil {
		t.Fatalf("LoginAs failed: %v", err)
	}
	rec2 := httptest.NewRecorder()
	controller.Handler(rec2, req2)

	if !strings.Contains(rec2.Body.String(), product.GetID()) {
		t.Errorf("expected cart to contain product, body: %s", rec2.Body.String())
	}
}
