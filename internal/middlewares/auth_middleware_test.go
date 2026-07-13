package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/auth"
	"github.com/dracory/test"
)

func TestAuthHandler_NoSessionKey(t *testing.T) {
	cfg := testutils.DefaultConf()
	// cfg.SetCacheStoreUsed(true)
	// cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	// Create a request without a session cookie
	req := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()

	// Create the middleware handler
	handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(config.AuthenticatedUserContextKey{}) != nil {
			t.Fatal("User should not be set in context")
		}
		if r.Context().Value(config.AuthenticatedSessionContextKey{}) != nil {
			t.Fatal("Session should not be set in context")
		}
		responseRecorder.WriteHeader(http.StatusOK)
	}))

	// Execute the handler
	handler.ServeHTTP(responseRecorder, req)

	// Assert that the next handler was called
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestAuthHandler_SessionNotFoundError(t *testing.T) {
	cfg := testutils.DefaultConf()
	// cfg.SetCacheStoreUsed(true)
	// cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	// Create a request with a session cookie
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: "some_session_key"})

	// Create a response recorder
	responseRecorder := httptest.NewRecorder()

	// Create the middleware handler
	handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(config.AuthenticatedUserContextKey{}) != nil {
			t.Fatal("User should not be set in context")
		}
		if r.Context().Value(config.AuthenticatedSessionContextKey{}) != nil {
			t.Fatal("Session should not be set in context")
		}
		responseRecorder.WriteHeader(http.StatusOK)
	}))

	// Execute the handler
	handler.ServeHTTP(responseRecorder, req)

	// Assert that the next handler was called
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestAuthHandler_SessionExpired(t *testing.T) {
	cfg := testutils.DefaultConf()
	// cfg.SetCacheStoreUsed(true)
	// cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	userStore := app.GetUserStore()
	sessionStore := app.GetSessionStore()

	if userStore == nil {
		t.Fatal("userStore should not be nil")
	}

	if sessionStore == nil {
		t.Fatal("sessionStore should not be nil")
	}

	user, err := testutils.SeedUser(userStore, test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	session, err := testutils.SeedSession(sessionStore, httptest.NewRequest("GET", "/", nil), user, -100)

	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("session should not be nil")
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: session.GetKey()})

	responseRecorder := httptest.NewRecorder()

	handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(config.AuthenticatedUserContextKey{}) != nil {
			t.Fatal("User should not be set in context")
		}
		if r.Context().Value(config.AuthenticatedSessionContextKey{}) != nil {
			t.Fatal("Session should not be set in context")
		}
		responseRecorder.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestAuthHandler_UserNotFound(t *testing.T) {
	cfg := testutils.DefaultConf()
	// cfg.SetCacheStoreUsed(true)
	// cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	userStore := app.GetUserStore()
	sessionStore := app.GetSessionStore()

	if userStore == nil {
		t.Fatal("userStore should not be nil")
	}

	if sessionStore == nil {
		t.Fatal("sessionStore should not be nil")
	}

	user, err := testutils.SeedUser(userStore, test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	session, err := testutils.SeedSession(sessionStore, httptest.NewRequest("GET", "/", nil), user, 1)

	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("session should not be nil")
	}

	err = userStore.UserDelete(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: session.GetKey()})

	responseRecorder := httptest.NewRecorder()

	handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(config.AuthenticatedUserContextKey{}) != nil {
			t.Fatal("User should not be set in context")
		}
		if r.Context().Value(config.AuthenticatedSessionContextKey{}) != nil {
			t.Fatal("Session should not be set in context")
		}

		responseRecorder.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestAuthHandler_SessionSuccess(t *testing.T) {
	cfg := testutils.DefaultConf()
	// cfg.SetCacheStoreUsed(true)
	// cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	userStore := app.GetUserStore()
	sessionStore := app.GetSessionStore()

	if userStore == nil {
		t.Fatal("userStore should not be nil")
	}

	if sessionStore == nil {
		t.Fatal("sessionStore should not be nil")
	}

	user, err := testutils.SeedUser(userStore, test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	session, err := testutils.SeedSession(sessionStore, httptest.NewRequest("GET", "/", nil), user, 1)

	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("session should not be nil")
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: session.GetKey()})

	responseRecorder := httptest.NewRecorder()

	handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(config.AuthenticatedUserContextKey{}) == nil {
			t.Fatal("User should be set in context")
		}
		if r.Context().Value(config.AuthenticatedSessionContextKey{}) == nil {
			t.Fatal("Session should be set in context")
		}
		responseRecorder.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestAuthHandler_SessionStoreNotEnabled(t *testing.T) {
	// When session store is disabled, it won't be initialized, so the middleware
	// will return a config error about SessionStore being required.
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(false)

	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := AuthMiddleware(app).GetHandler()(next)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "auth middleware: SessionStore is required") {
		t.Errorf("Expected body to contain 'auth middleware: SessionStore is required', got '%s'", rr.Body.String())
	}
}

func TestAuthHandler_SessionStoreEnabledButNotInitialized(t *testing.T) {
	// Configure app with session store enabled, then nil the store
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))
	// Simulate uninitialized session store
	app.SetSessionStore(nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := AuthMiddleware(app).GetHandler()(next)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "auth middleware: SessionStore is required") {
		t.Errorf("Expected body to contain 'auth middleware: SessionStore is required', got '%s'", rr.Body.String())
	}
}

func TestAuthHandler_UserStoreUsed_ReturnsUserStoreNotEnabledError(t *testing.T) {
	// When user store is disabled, it won't be initialized, so the middleware
	// will return a config error about UserStore being required.
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(false)
	cfg.SetSessionStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := AuthMiddleware(app).GetHandler()(next)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "auth middleware: UserStore is required") {
		t.Errorf("Expected body to contain 'auth middleware: UserStore is required', got '%s'", rr.Body.String())
	}
}

func TestAuthHandler_UserStoreNotInitialized(t *testing.T) {
	// Set session store used and initialized, but user store nil,
	// which should trigger the "UserStore is required" config error.
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))

	// Ensure user store is nil
	app.SetUserStore(nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := AuthMiddleware(app).GetHandler()(next)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "auth middleware: UserStore is required") {
		t.Errorf("Expected body to contain 'auth middleware: UserStore is required', got '%s'", rr.Body.String())
	}
}
