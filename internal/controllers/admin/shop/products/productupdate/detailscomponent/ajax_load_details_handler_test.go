package detailscomponent

import (
	"context"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxLoadDetails_Success tests successful details loading
func TestHandleAjaxLoadDetails_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
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

	response := HandleAjaxLoadDetails(registry, product.GetID())

	if !strings.Contains(response, `"details"`) {
		t.Error("expected response to contain details field")
	}
	if !strings.Contains(response, `"Test Product"`) {
		t.Error("expected response to contain product title")
	}
	if !strings.Contains(response, `"99.99"`) {
		t.Error("expected response to contain price")
	}
}

// TestHandleAjaxLoadDetails_ProductNotFound tests error when product not found
func TestHandleAjaxLoadDetails_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	response := HandleAjaxLoadDetails(registry, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
}
