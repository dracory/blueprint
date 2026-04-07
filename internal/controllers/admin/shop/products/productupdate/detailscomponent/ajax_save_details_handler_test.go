package detailscomponent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxSaveDetails_Success tests successful details saving
func TestHandleAjaxSaveDetails_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Old Title")
	product.SetDescription("Old Description")
	product.SetPrice("50.00")
	product.SetQuantity("5")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	body := `{"details": {"id": "` + product.GetID() + `", "title": "New Title", "description": "New Description", "price": "99.99", "quantity": "20", "status": "active"}}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-details&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveDetails(registry, req, product.GetID())

	if !strings.Contains(response, `"details"`) {
		t.Error("expected response to contain details field")
	}

	// Verify product was actually updated
	updatedProduct, err := registry.GetShopStore().ProductFindByID(context.Background(), product.GetID())
	if err != nil {
		t.Fatalf("failed to find updated product: %v", err)
	}

	if updatedProduct.GetTitle() != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", updatedProduct.GetTitle())
	}
	if updatedProduct.GetPrice() != "99.99" {
		t.Errorf("expected price '99.99', got '%s'", updatedProduct.GetPrice())
	}
}

// TestHandleAjaxSaveDetails_InvalidBody tests error with invalid JSON body
func TestHandleAjaxSaveDetails_InvalidBody(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-details&product_id=test123",
		strings.NewReader(body))

	response := HandleAjaxSaveDetails(registry, req, "test123")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
}

// TestHandleAjaxSaveDetails_ProductNotFound tests error when product not found
func TestHandleAjaxSaveDetails_ProductNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `{"details": {"title": "New Title", "price": "99.99"}}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-details&product_id=nonexistent",
		strings.NewReader(body))

	response := HandleAjaxSaveDetails(registry, req, "nonexistent")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
}
