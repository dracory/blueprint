package product_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHandleProductDeleteSelectedAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleProductDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleProductDeleteSelectedAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewProductManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleProductDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleProductDeleteSelectedAjax_DeletesSelectedProducts(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleProductDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}
