package tagscomponent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxSaveTags_Success tests successful tags saving
func TestHandleAjaxSaveTags_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	body := `{"tags": [{"id": "1", "tag": "featured"}, {"id": "2", "tag": "new"}]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-tags&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveTags(registry, req, product.GetID())

	if !strings.Contains(response, `"tags"`) {
		t.Error("expected response to contain tags field")
	}
	if !strings.Contains(response, `"featured"`) {
		t.Error("expected response to contain featured tag")
	}

	// Verify tags were actually saved
	updatedProduct, err := registry.GetShopStore().ProductFindByID(context.Background(), product.GetID())
	if err != nil {
		t.Fatalf("failed to find updated product: %v", err)
	}

	metas, err := updatedProduct.GetMetas()
	if err != nil {
		t.Fatalf("failed to get metas: %v", err)
	}

	if metas["tags"] != "featured,new" {
		t.Errorf("expected metas['tags'] = 'featured,new', got '%s'", metas["tags"])
	}
}

// TestHandleAjaxSaveTags_InvalidBody tests error with invalid JSON body
func TestHandleAjaxSaveTags_InvalidBody(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-tags&product_id=test123",
		strings.NewReader(body))

	response := HandleAjaxSaveTags(registry, req, "test123")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Invalid request body`) {
		t.Error("expected response to contain 'Invalid request body' message")
	}
}

// TestHandleAjaxSaveTags_ProductNotFound tests error when product not found
func TestHandleAjaxSaveTags_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `{"tags": [{"id": "1", "tag": "test"}]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-tags&product_id=nonexistent",
		strings.NewReader(body))

	response := HandleAjaxSaveTags(registry, req, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Product not found`) {
		t.Error("expected response to contain 'Product not found' message")
	}
}
