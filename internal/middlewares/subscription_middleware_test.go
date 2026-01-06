package middlewares

import (
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/test"
)

func TestSubscriptionOnlyMiddleware_AdminUserPassesThrough(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)

	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, session, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		testutils.ADMIN_01,
		httptest.NewRequest("GET", "/", nil),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	body, response, err := test.CallMiddleware(
		"GET",
		NewSubscriptionOnlyMiddleware(registry).GetHandler(),
		func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("ok")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		},
		test.NewRequestOptions{
			Context: map[any]any{
				config.AuthenticatedUserContextKey{}:    user,
				config.AuthenticatedSessionContextKey{}: session,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("response should not be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if body != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", body)
	}
}
