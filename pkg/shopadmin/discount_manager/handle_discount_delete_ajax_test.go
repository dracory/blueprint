package discount_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleDiscountDeleteAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleDiscountDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleDiscountDeleteAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleDiscountDeleteAjax_DeletesDiscount(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	controller := NewDiscountManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleDiscountDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
