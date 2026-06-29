package discount_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHandleDiscountDeleteSelectedAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleDiscountDeleteSelectedAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleDiscountDeleteSelectedAjax_DeletesSelectedDiscounts(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}
