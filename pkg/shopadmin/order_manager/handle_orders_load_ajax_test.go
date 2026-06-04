package order_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleOrdersLoadAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewOrderManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleOrdersLoadAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleOrdersLoadAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewOrderManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleOrdersLoadAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleOrdersLoadAjax_LoadsOrders(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewOrderManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleOrdersLoadAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
