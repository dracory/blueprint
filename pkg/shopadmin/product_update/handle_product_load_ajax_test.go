package product_update

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHandleProductLoadAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleLoadProduct, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleProductLoadAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewProductUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleLoadProduct, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleProductLoadAjax_LoadsProduct(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewProductUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleLoadProduct, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}
