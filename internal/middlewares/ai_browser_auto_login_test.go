package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/auth"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
)

func TestAiBrowserAutoLoginHandler_NoCookie_CreatesSession(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	devUser, err := testutils.SeedUser(app.GetUserStore(), aiBrowserUserID)
	if err != nil {
		t.Fatal("failed to seed dev user:", err)
	}
	_ = devUser

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := aiBrowserAutoLoginHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var authCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == auth.CookieName {
			authCookie = c
			break
		}
	}

	if authCookie == nil {
		t.Fatal("expected auth cookie to be set, got none")
	}
	if authCookie.Value == "" {
		t.Error("expected auth cookie to have a non-empty value")
	}
}

func TestAiBrowserAutoLoginHandler_WithCookie_PassesThrough(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), "test-user-passthrough")
	if err != nil {
		t.Fatal("failed to seed user:", err)
	}

	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent("test").
		SetIPAddress("127.0.0.1").
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(24).ToDateTimeString(carbon.UTC))

	if err := app.GetSessionStore().SessionCreate(context.Background(), session); err != nil {
		t.Fatal("failed to create session:", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: session.GetKey()})
	rr := httptest.NewRecorder()

	handler := aiBrowserAutoLoginHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	sessionList, err := app.GetSessionStore().SessionList(context.Background(),
		sessionstore.NewSessionQuery().SetUserID(user.GetID()))
	if err != nil {
		t.Fatal(err)
	}

	if len(sessionList) != 1 {
		t.Errorf("expected 1 session (no duplicate created), got %d", len(sessionList))
	}
}

func TestAiBrowserAutoLoginHandler_DevUserNotFound_PassesThrough(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handlerCalled := false
	handler := aiBrowserAutoLoginHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("expected next handler to be called even when dev user is not found")
	}

	for _, c := range rr.Result().Cookies() {
		if c.Name == auth.CookieName {
			t.Error("expected no auth cookie to be set when dev user is not found")
		}
	}
}

func TestAiBrowserAutoLoginHandler_NilStores_PassesThrough(t *testing.T) {
	cfg := testutils.DefaultConf()
	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handlerCalled := false
	handler := aiBrowserAutoLoginHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("expected next handler to be called when stores are nil")
	}
}

func TestAiBrowserAutoLoginHandler_SetsContextUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	_, err := testutils.SeedUser(app.GetUserStore(), aiBrowserUserID)
	if err != nil {
		t.Fatal("failed to seed dev user:", err)
	}

	var contextUser userstore.UserInterface
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := aiBrowserAutoLoginHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextUser, _ = r.Context().Value(config.AuthenticatedUserContextKey{}).(userstore.UserInterface)
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if contextUser == nil {
		t.Error("expected authenticated user to be set in request context")
	}
}
