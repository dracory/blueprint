package metadatacomponent

import (
	"context"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxLoadMetadata_Success tests successful metadata loading
func TestHandleAjaxLoadMetadata_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
	product.SetMetas(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	response := HandleAjaxLoadMetadata(registry, product.GetID())

	if !strings.Contains(response, `"metadata"`) {
		t.Error("expected response to contain metadata field")
	}
	if !strings.Contains(response, `"key1"`) {
		t.Error("expected response to contain key1")
	}
	if !strings.Contains(response, `"value1"`) {
		t.Error("expected response to contain value1")
	}
}

// TestHandleAjaxLoadMetadata_ProductNotFound tests error when product not found
func TestHandleAjaxLoadMetadata_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	response := HandleAjaxLoadMetadata(registry, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Product not found`) {
		t.Error("expected response to contain 'Product not found' message")
	}
}

// TestHandleAjaxLoadMetadata_NoMetadata tests product with no metadata
func TestHandleAjaxLoadMetadata_NoMetadata(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Product Without Metadata")
	product.SetPrice("50.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	response := HandleAjaxLoadMetadata(registry, product.GetID())

	if !strings.Contains(response, `"metadata"`) {
		t.Error("expected response to contain metadata field")
	}
}
