package user_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHandleUserLoadAjax_RequiresPOST(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
}

func TestHandleUserLoadAjax_RequiresUserStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
}

func TestHandleUserLoadAjax_LoadsUsers(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUsersFetchAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
}
