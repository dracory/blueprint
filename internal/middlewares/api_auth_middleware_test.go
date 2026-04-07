package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/api"
	"github.com/dracory/sessionstore"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func TestAPIAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	if nextCalled {
		t.Error("next handler should not be called when token is missing")
	}
	expected := api.Error("Authorization token required").ToString()
	if res.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, res.Body.String())
	}
}

func TestAPIAuthMiddleware_InvalidToken(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "invalid-token")
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	if nextCalled {
		t.Error("next handler should not be called for invalid token")
	}
	expected := api.Error("Invalid or expired token").ToString()
	if res.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, res.Body.String())
	}
}

func TestAPIAuthMiddleware_ExpiredSession(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(registry.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, -60)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	if nextCalled {
		t.Error("next handler should not be called for expired session")
	}
	expected := api.Error("Invalid or expired token").ToString()
	if res.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, res.Body.String())
	}
}

func TestAPIAuthMiddleware_SessionMissingUser(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	baseReq := httptest.NewRequest(http.MethodGet, "/", nil)
	session := sessionstore.NewSession().
		SetUserID("").
		SetUserAgent(baseReq.UserAgent()).
		SetIPAddress("127.0.0.1").
		SetExpiresAt(time.Now().Add(time.Hour).UTC().Format("2006-01-02 15:04:05"))

	if err := registry.GetSessionStore().SessionCreate(baseReq.Context(), session); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	if nextCalled {
		t.Error("next handler should not be called when session is missing user")
	}
	expected := api.Error("Session missing user").ToString()
	if res.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, res.Body.String())
	}
}

func TestAPIAuthMiddleware_UserNotFound(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	userStore := registry.GetUserStore()
	user, err := testutils.SeedUser(userStore, test.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(registry.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, 3600)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	if err := userStore.UserDelete(context.Background(), user); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	if nextCalled {
		t.Error("next handler should not be called when user is missing")
	}
	expected := api.Error("User not found").ToString()
	if res.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, res.Body.String())
	}
}

func TestAPIAuthMiddleware_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(registry.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, 3600)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(registry).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false
	var ctxSession sessionstore.SessionInterface
	var ctxUser userstore.UserInterface

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		var ok bool
		ctxSession, ok = r.Context().Value(config.APIAuthenticatedSessionContextKey{}).(sessionstore.SessionInterface)
		if !ok {
			t.Error("session should be in request context")
		}
		ctxUser, ok = r.Context().Value(config.APIAuthenticatedUserContextKey{}).(userstore.UserInterface)
		if !ok {
			t.Error("user should be in request context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware(next).ServeHTTP(res, req)

	if !nextCalled {
		t.Error("next handler should be called when authentication succeeds")
	}
	if ctxSession != nil {
		if session.GetKey() != ctxSession.GetKey() {
			t.Errorf("context session key should match: expected %s, got %s", session.GetKey(), ctxSession.GetKey())
		}
	}
	if ctxUser != nil {
		if user.ID() != ctxUser.ID() {
			t.Errorf("context user id should match: expected %s, got %s", user.ID(), ctxUser.ID())
		}
	}
	if res.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, res.Result().StatusCode)
	}
}
