package mediacomponent

import (
	"context"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxLoadMedia_Success tests successful media loading
func TestHandleAjaxLoadMedia_Success(t *testing.T) {
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

	// Create some media
	media := shopstore.NewMedia()
	media.SetEntityID(product.GetID())
	media.SetURL("https://example.com/image1.jpg")
	media.SetTitle("Image 1")
	media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	media.SetSequence(0)
	media.SetType("image")

	if err := registry.GetShopStore().MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	response := HandleAjaxLoadMedia(registry, product.GetID())

	if !strings.Contains(response, `"media"`) {
		t.Error("expected response to contain media field")
	}
	if !strings.Contains(response, `"Image 1"`) {
		t.Error("expected response to contain Image 1")
	}
	if !strings.Contains(response, `"url"`) {
		t.Error("expected response to contain url field")
	}
}

// TestHandleAjaxLoadMedia_ProductNotFound tests when product has no media
func TestHandleAjaxLoadMedia_NoMedia(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Product Without Media")
	product.SetPrice("50.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	response := HandleAjaxLoadMedia(registry, product.GetID())

	if !strings.Contains(response, `"media"`) {
		t.Error("expected response to contain media field")
	}
}
