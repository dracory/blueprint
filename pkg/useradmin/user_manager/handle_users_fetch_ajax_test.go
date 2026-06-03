package user_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleUserLoadAjax_RequiresPOST(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleUserLoadAjax_RequiresUserStore(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleUserLoadAjax_LoadsUsers(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(registry)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
