package productupdate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestNewProductUpdateController tests the controller constructor
func TestNewProductUpdateController(t *testing.T) {
	t.Run("creates controller with registry", func(t *testing.T) {
		registry := testutils.Setup(
			testutils.WithCacheStore(true),
			testutils.WithShopStore(true),
		)

		controller := NewProductUpdateController(registry)

		if controller == nil {
			t.Error("expected controller to be non-nil")
		}
		if controller.registry == nil {
			t.Error("expected registry to be set")
		}
	})
}

// TestHandler_MissingProductID tests error when product_id is missing
func TestHandler_MissingProductID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// ToFlashError returns a link, so we check if it's a valid link
	if !strings.Contains(result, "See Other") && !strings.Contains(result, "/flash") {
		t.Error("expected result to contain flash redirect link")
	}
}

// TestHandler_ProductNotFound tests error when product doesn't exist
func TestHandler_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id=nonexistent", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	// ToFlashError returns a link, so we check if it's a valid link
	if !strings.Contains(result, "See Other") && !strings.Contains(result, "/flash") {
		t.Error("expected result to contain flash redirect link")
	}
}

// TestHandler_RenderDetailsView tests rendering the details view
func TestHandler_RenderDetailsView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetDescription("Test Description")
	product.SetPrice("99.99")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id="+product.GetID()+"&view=details", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	if !strings.Contains(result, "Test Product") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Edit Product") {
		t.Error("expected result to contain page title")
	}
}

// TestHandler_RenderMediaView tests rendering the media view
func TestHandler_RenderMediaView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id="+product.GetID()+"&view=media", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	if !strings.Contains(result, "Test Product") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Media") {
		t.Error("expected result to contain media section title")
	}
}

// TestHandler_RenderTagsView tests rendering the tags view
func TestHandler_RenderTagsView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("50.00")
	product.SetQuantity("3")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id="+product.GetID()+"&view=tags", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	if !strings.Contains(result, "Test Product") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Tags") {
		t.Error("expected result to contain tags section title")
	}
}

// TestHandler_RenderMetadataView tests rendering the metadata view
func TestHandler_RenderMetadataView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("75.00")
	product.SetQuantity("2")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id="+product.GetID()+"&view=metadata", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	if !strings.Contains(result, "Test Product") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Metadata") {
		t.Error("expected result to contain metadata section title")
	}
}

// TestHandler_DefaultView tests that default view is details
func TestHandler_DefaultView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("25.00")
	product.SetQuantity("1")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update?product_id="+product.GetID(), nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)

	if !strings.Contains(result, "Product Details") {
		t.Error("expected default view to be details")
	}
}
