package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func TestEmailAllowlistMiddleware_UnauthenticatedRedirectsToLogin(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	body, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach next handler")
	}, test.NewRequestOptions{})

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response must not be nil")
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	msg, err := testutils.FlashMessageFindFromBody(registry.GetCacheStore(), body)
	if err != nil {
		t.Fatal(err)
	}

	if msg == nil {
		t.Fatal("flash message must not be nil")
	}

	if msg.Type != helpers.FLASH_ERROR {
		t.Fatalf("expected flash type %s, got %s", helpers.FLASH_ERROR, msg.Type)
	}

	if msg.Message != "Only authenticated users can access this page" {
		t.Fatalf("expected message %q, got %q", "Only authenticated users can access this page", msg.Message)
	}

	if msg.Url != links.AUTH_LOGIN {
		t.Fatalf("expected redirect %q, got %q", links.AUTH_LOGIN, msg.Url)
	}
}

func TestEmailAllowlistMiddleware_BlockedEmail(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetEmailsAllowedAccess([]string{"allowed@example.com"})
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, session, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		"blocked-user",
		httptest.NewRequest("GET", "/", nil),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetEmail("blocked@example.com")
	if err := registry.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	body, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach next handler")
	}, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	msg, err := testutils.FlashMessageFindFromBody(registry.GetCacheStore(), body)
	if err != nil {
		t.Fatal(err)
	}

	if msg == nil {
		t.Fatal("flash message must not be nil")
	}

	if msg.Message != "Access restricted to authorized emails only" {
		t.Fatalf("unexpected message: %s", msg.Message)
	}

	if msg.Url != links.Website().Home() {
		t.Fatalf("expected home redirect, got %s", msg.Url)
	}
}

func TestEmailAllowlistMiddleware_AllowedEmail(t *testing.T) {
	allowedEmail := "allowed@example.com"

	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetEmailsAllowedAccess([]string{allowedEmail})
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, session, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		test.USER_01,
		httptest.NewRequest("GET", "/", nil),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetEmail(allowedEmail)
	if err := registry.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	called := false
	body, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatalf("expected next handler to be called, got body: %s", body)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestEmailAllowlistMiddleware_AllEmailsAllowedWhenMapEmpty(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetEmailsAllowedAccess([]string{})
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, session, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		"any-user",
		httptest.NewRequest("GET", "/", nil),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetEmail("random@example.com")
	if err := registry.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	called := false

	_, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected OK, got %d", response.StatusCode)
	}

	if !called {
		t.Fatal("next handler should be called")
	}
}

func TestEmailAllowlistMiddleware_WithVaultCaching(t *testing.T) {
	allowedEmail := "vault@example.com"

	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetUserStoreVaultEnabled(true)
	cfg.SetVaultStoreUsed(true)
	cfg.SetVaultStoreKey("test-vault-key")
	cfg.SetEmailsAllowedAccess([]string{allowedEmail})
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, session, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		"vault-user",
		httptest.NewRequest("GET", "/", nil),
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	// Tokenize the email so the vault is used
	email := allowedEmail
	_, _, emailToken, _, _, err := ext.UserTokenize(
		context.Background(),
		registry.GetVaultStore(),
		registry.GetConfig().GetVaultStoreKey(),
		user,
		"Test", "User", email, "", "",
	)
	if err != nil {
		t.Fatalf("Failed to tokenize email: %v", err)
	}

	user.SetEmail(emailToken)
	if err := registry.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	// First call should hit the vault and cache
	user.SetStatus(userstore.USER_STATUS_ACTIVE)

	// Verify tokenization succeeded
	if !strings.HasPrefix(user.Email(), "tk_") {
		t.Fatalf("Expected tokenized email, got: %s", user.Email())
	}

	// First middleware call
	called1 := false
	body1, response1, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		called1 = true
		w.WriteHeader(http.StatusOK)
	}, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called1 {
		t.Fatalf("expected next handler to be called on first call, got body: %s", body1)
	}

	if response1.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d on first call, got %d", http.StatusOK, response1.StatusCode)
	}

	// Second call should hit the cache
	called2 := false
	body2, response2, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(registry).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
		called2 = true
		w.WriteHeader(http.StatusOK)
	}, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called2 {
		t.Fatalf("expected next handler to be called on second call, got body: %s", body2)
	}

	if response2.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d on second call, got %d", http.StatusOK, response2.StatusCode)
	}
}
