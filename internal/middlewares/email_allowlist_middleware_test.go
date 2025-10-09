package middlewares

import (
    "context"
    "net/http"
    "net/http/httptest"
    "project/internal/config"
    "project/internal/helpers"
    "project/internal/links"
    "project/internal/testutils"
    "testing"

    "github.com/dracory/test"
    "github.com/dracory/userstore"
)

func TestEmailAllowlistMiddleware_UnauthenticatedRedirectsToLogin(t *testing.T) {
    cfg := testutils.DefaultConf()
    cfg.SetCacheStoreUsed(true)
    cfg.SetSessionStoreUsed(true)
    cfg.SetUserStoreUsed(true)
    app := testutils.Setup(testutils.WithCfg(cfg))

    body, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(app).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
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

    msg, err := testutils.FlashMessageFindFromBody(app.GetCacheStore(), body)
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
    original := allowedEmails
    allowedEmails = map[string]struct{}{
        "allowed@example.com": {},
    }
    defer func() { allowedEmails = original }()

    cfg := testutils.DefaultConf()
    cfg.SetCacheStoreUsed(true)
    cfg.SetSessionStoreUsed(true)
    cfg.SetUserStoreUsed(true)
    app := testutils.Setup(testutils.WithCfg(cfg))

    user, session, err := testutils.SeedUserAndSession(
        app.GetUserStore(),
        app.GetSessionStore(),
        "blocked-user",
        httptest.NewRequest("GET", "/", nil),
        1,
    )

    if err != nil {
        t.Fatal(err)
    }

    user.SetStatus(userstore.USER_STATUS_ACTIVE)
    user.SetEmail("blocked@example.com")
    if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
        t.Fatal(err)
    }

    body, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(app).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
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

    msg, err := testutils.FlashMessageFindFromBody(app.GetCacheStore(), body)
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
    original := allowedEmails
    defer func() { allowedEmails = original }()

    cfg := testutils.DefaultConf()
    cfg.SetCacheStoreUsed(true)
    cfg.SetSessionStoreUsed(true)
    cfg.SetUserStoreUsed(true)
    app := testutils.Setup(testutils.WithCfg(cfg))

    user, session, err := testutils.SeedUserAndSession(
        app.GetUserStore(),
        app.GetSessionStore(),
        testutils.USER_01,
        httptest.NewRequest("GET", "/", nil),
        1,
    )

    if err != nil {
        t.Fatal(err)
    }

    allowedEmail := "allowed@example.com"
    allowedEmails = map[string]struct{}{allowedEmail: {}}
    user.SetStatus(userstore.USER_STATUS_ACTIVE)
    user.SetEmail(allowedEmail)
    if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
        t.Fatal(err)
    }

    called := false

    _, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(app).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
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

func TestEmailAllowlistMiddleware_AllEmailsAllowedWhenMapEmpty(t *testing.T) {
    original := allowedEmails
    allowedEmails = map[string]struct{}{}
    defer func() { allowedEmails = original }()

    cfg := testutils.DefaultConf()
    cfg.SetCacheStoreUsed(true)
    cfg.SetSessionStoreUsed(true)
    cfg.SetUserStoreUsed(true)
    app := testutils.Setup(testutils.WithCfg(cfg))

    user, session, err := testutils.SeedUserAndSession(
        app.GetUserStore(),
        app.GetSessionStore(),
        "any-user",
        httptest.NewRequest("GET", "/", nil),
        1,
    )

    if err != nil {
        t.Fatal(err)
    }

    user.SetStatus(userstore.USER_STATUS_ACTIVE)
    user.SetEmail("random@example.com")
    if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
        t.Fatal(err)
    }

    called := false

    _, response, err := test.CallMiddleware("GET", NewEmailAllowlistMiddleware(app).GetHandler(), func(w http.ResponseWriter, r *http.Request) {
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
