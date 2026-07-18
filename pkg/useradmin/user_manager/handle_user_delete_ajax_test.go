package user_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
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

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
}

func TestHandleUserDeleteAjax_RequiresUserStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserManagerController(app)
	_, response, err := test.CallStringEndpoint(http.MethodPost, controller.handleUserDeleteAjax, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
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

	if err != nil { t.Errorf("unexpected error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("expected %v, got %v", http.StatusOK, response.StatusCode) }
}
