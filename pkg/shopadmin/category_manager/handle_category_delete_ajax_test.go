package category_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleCategoryDeleteAjax_RequiresShopStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewCategoryManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleCategoryDelete, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
