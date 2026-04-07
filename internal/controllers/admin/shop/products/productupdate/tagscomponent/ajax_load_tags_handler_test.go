package tagscomponent

import (
	"context"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxLoadTags_Success tests successful tags loading
func TestHandleAjaxLoadTags_Success(t *testing.T) {
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
		"tags": "featured, new, sale",
	})

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	response := HandleAjaxLoadTags(registry, product.GetID())

	if !strings.Contains(response, `"tags"`) {
		t.Error("expected response to contain tags field")
	}
	if !strings.Contains(response, `"featured"`) {
		t.Error("expected response to contain featured tag")
	}
	if !strings.Contains(response, `"new"`) {
		t.Error("expected response to contain new tag")
	}
	if !strings.Contains(response, `"sale"`) {
		t.Error("expected response to contain sale tag")
	}
}

// TestHandleAjaxLoadTags_ProductNotFound tests error when product not found
func TestHandleAjaxLoadTags_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	response := HandleAjaxLoadTags(registry, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Product not found`) {
		t.Error("expected response to contain 'Product not found' message")
	}
}

// TestHandleAjaxLoadTags_NoTags tests product with no tags
func TestHandleAjaxLoadTags_NoTags(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Product Without Tags")
	product.SetPrice("50.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	response := HandleAjaxLoadTags(registry, product.GetID())

	if !strings.Contains(response, `"tags"`) {
		t.Error("expected response to contain tags field")
	}
}
