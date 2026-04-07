package metadatacomponent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxSaveMetadata_Success tests successful metadata saving
func TestHandleAjaxSaveMetadata_Success(t *testing.T) {
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

	body := `{"metadata": [{"id": "1", "key": "testkey", "value": "testvalue"}]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-metadata&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveMetadata(registry, req, product.GetID())

	if !strings.Contains(response, `"metadata"`) {
		t.Error("expected response to contain metadata field")
	}
	if !strings.Contains(response, `"testkey"`) {
		t.Error("expected response to contain testkey")
	}
	if !strings.Contains(response, `"testvalue"`) {
		t.Error("expected response to contain testvalue")
	}

	// Verify metadata was actually saved
	updatedProduct, err := registry.GetShopStore().ProductFindByID(context.Background(), product.GetID())
	if err != nil {
		t.Fatalf("failed to find updated product: %v", err)
	}

	metas, err := updatedProduct.GetMetas()
	if err != nil {
		t.Fatalf("failed to get metas: %v", err)
	}

	if metas["testkey"] != "testvalue" {
		t.Errorf("expected metas['testkey'] = 'testvalue', got '%s'", metas["testkey"])
	}
}

// TestHandleAjaxSaveMetadata_InvalidBody tests error with invalid JSON body
func TestHandleAjaxSaveMetadata_InvalidBody(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-metadata&product_id=test123",
		strings.NewReader(body))

	response := HandleAjaxSaveMetadata(registry, req, "test123")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Invalid request body`) {
		t.Error("expected response to contain 'Invalid request body' message")
	}
}

// TestHandleAjaxSaveMetadata_ProductNotFound tests error when product not found
func TestHandleAjaxSaveMetadata_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `{"metadata": [{"id": "1", "key": "test", "value": "value"}]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-metadata&product_id=nonexistent",
		strings.NewReader(body))

	response := HandleAjaxSaveMetadata(registry, req, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Product not found`) {
		t.Error("expected response to contain 'Product not found' message")
	}
}
