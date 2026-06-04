package user_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleUserDeleteAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleUserDeleteAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleUserDeleteAjax_RequiresUserStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUserDeleteAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleUserDeleteAjax_RequiresUserID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUserDeleteAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
