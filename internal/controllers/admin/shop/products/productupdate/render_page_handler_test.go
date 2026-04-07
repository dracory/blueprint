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

// TestHandleRenderPage_DetailsView tests rendering the details view
func TestHandleRenderPage_DetailsView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product for Details")
	product.SetDescription("Test Description")
	product.SetPrice("99.99")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "details", product.GetID())

	if !strings.Contains(result, "Test Product for Details") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Details") {
		t.Error("expected result to contain details section title")
	}
	if !strings.Contains(result, "Save Details") {
		t.Error("expected result to contain save button")
	}
}

// TestHandleRenderPage_MediaView tests rendering the media view
func TestHandleRenderPage_MediaView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product for Media")
	product.SetPrice("100.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "media", product.GetID())

	if !strings.Contains(result, "Test Product for Media") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Media") {
		t.Error("expected result to contain media section title")
	}
	if !strings.Contains(result, "Save Media") {
		t.Error("expected result to contain save button")
	}
}

// TestHandleRenderPage_TagsView tests rendering the tags view
func TestHandleRenderPage_TagsView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product for Tags")
	product.SetPrice("50.00")
	product.SetQuantity("3")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "tags", product.GetID())

	if !strings.Contains(result, "Test Product for Tags") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Tags") {
		t.Error("expected result to contain tags section title")
	}
	if !strings.Contains(result, "Save Tags") {
		t.Error("expected result to contain save button")
	}
}

// TestHandleRenderPage_MetadataView tests rendering the metadata view
func TestHandleRenderPage_MetadataView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product for Metadata")
	product.SetPrice("75.00")
	product.SetQuantity("2")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "metadata", product.GetID())

	if !strings.Contains(result, "Test Product for Metadata") {
		t.Error("expected result to contain product title")
	}
	if !strings.Contains(result, "Product Metadata") {
		t.Error("expected result to contain metadata section title")
	}
	if !strings.Contains(result, "Save Metadata") {
		t.Error("expected result to contain save button")
	}
}

// TestHandleRenderPage_DefaultView tests that default view is details
func TestHandleRenderPage_DefaultView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product Default")
	product.SetPrice("25.00")
	product.SetQuantity("1")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "", product.GetID())

	if !strings.Contains(result, "Test Product Default") {
		t.Error("expected result to contain product title")
	}
	// When view is empty, it defaults to details in the switch statement
	// The details component should be rendered
	if !strings.Contains(result, "FormProductDetailsUpdate") {
		t.Error("expected default view to render details form")
	}
}

// TestHandleRenderPage_InvalidView tests that invalid view defaults to details
func TestHandleRenderPage_InvalidView(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product Invalid View")
	product.SetPrice("30.00")
	product.SetQuantity("1")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "invalid-view", product.GetID())

	if !strings.Contains(result, "Test Product Invalid View") {
		t.Error("expected result to contain product title")
	}
	// When view is invalid, it defaults to details in the switch statement
	// The details component should be rendered
	if !strings.Contains(result, "FormProductDetailsUpdate") {
		t.Error("expected invalid view to default to details form")
	}
}

// TestHandleRenderPage_ContainsBreadcrumbs tests that breadcrumbs are rendered
func TestHandleRenderPage_ContainsBreadcrumbs(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product Breadcrumbs")
	product.SetPrice("40.00")
	product.SetQuantity("1")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "details", product.GetID())

	if !strings.Contains(result, "Home") {
		t.Error("expected result to contain breadcrumb for Home")
	}
	if !strings.Contains(result, "Shop") {
		t.Error("expected result to contain breadcrumb for Shop")
	}
	if !strings.Contains(result, "Product Manager") {
		t.Error("expected result to contain breadcrumb for Product Manager")
	}
	if !strings.Contains(result, "Edit Product") {
		t.Error("expected result to contain breadcrumb for Edit Product")
	}
}

// TestHandleRenderPage_ContainsTabs tests that navigation tabs are rendered
func TestHandleRenderPage_ContainsTabs(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
		testutils.WithCmsStore(true, "test-template"),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product Tabs")
	product.SetPrice("45.00")
	product.SetQuantity("1")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	controller := NewProductUpdateController(registry)
	req := httptest.NewRequest(http.MethodGet, "/admin/shop/products/update", nil)

	result := controller.handleRenderPage(req, product, "details", product.GetID())

	if !strings.Contains(result, "Details") {
		t.Error("expected result to contain Details tab")
	}
	if !strings.Contains(result, "Media") {
		t.Error("expected result to contain Media tab")
	}
	if !strings.Contains(result, "Tags") {
		t.Error("expected result to contain Tags tab")
	}
	if !strings.Contains(result, "Metadata") {
		t.Error("expected result to contain Metadata tab")
	}
}
