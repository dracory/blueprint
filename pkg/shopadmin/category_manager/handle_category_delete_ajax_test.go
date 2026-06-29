package category_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHandleCategoryDeleteAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewCategoryManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleCategoryDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleCategoryDeleteAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewCategoryManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleCategoryDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleCategoryDeleteAjax_DeletesCategory(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewCategoryManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleCategoryDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}
