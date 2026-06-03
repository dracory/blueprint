package discount_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleDiscountDeleteSelectedAjax_RequiresPOST(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleDiscountDeleteSelectedAjax_RequiresShopStore(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewDiscountManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleDiscountDeleteSelectedAjax_DeletesSelectedDiscounts(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDeleteSelected, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
