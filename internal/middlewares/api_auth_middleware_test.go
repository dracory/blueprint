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
	"github.com/dracory/userstore"
	"github.com/stretchr/testify/assert"
)

func TestAPIAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	middleware := NewAPIAuthMiddleware(app).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	assert.False(t, nextCalled, "next handler should not be called when token is missing")
	assert.Equal(t, api.Error("Authorization token required").ToString(), res.Body.String())
}

func TestAPIAuthMiddleware_InvalidToken(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	middleware := NewAPIAuthMiddleware(app).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "invalid-token")
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	assert.False(t, nextCalled, "next handler should not be called for invalid token")
	assert.Equal(t, api.Error("Invalid or expired token").ToString(), res.Body.String())
}

func TestAPIAuthMiddleware_ExpiredSession(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(app.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, -60)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(app).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	assert.False(t, nextCalled, "next handler should not be called for expired session")
	assert.Equal(t, api.Error("Invalid or expired token").ToString(), res.Body.String())
}

func TestAPIAuthMiddleware_SessionMissingUser(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	baseReq := httptest.NewRequest(http.MethodGet, "/", nil)
	session := sessionstore.NewSession().
		SetUserID("").
		SetUserAgent(baseReq.UserAgent()).
		SetIPAddress("127.0.0.1").
		SetExpiresAt(time.Now().Add(time.Hour).UTC().Format("2006-01-02 15:04:05"))

	if err := app.GetSessionStore().SessionCreate(baseReq.Context(), session); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(app).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	assert.False(t, nextCalled, "next handler should not be called when session is missing user")
	assert.Equal(t, api.Error("Session missing user").ToString(), res.Body.String())
}

func TestAPIAuthMiddleware_UserNotFound(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	userStore := app.GetUserStore()
	user, err := testutils.SeedUser(userStore, testutils.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(app.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, 3600)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	if err := userStore.UserDelete(context.Background(), user); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	middleware := NewAPIAuthMiddleware(app).GetHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", session.GetKey())
	res := httptest.NewRecorder()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware(next).ServeHTTP(res, req)

	assert.False(t, nextCalled, "next handler should not be called when user is missing")
	assert.Equal(t, api.Error("User not found").ToString(), res.Body.String())
}

func TestAPIAuthMiddleware_Success(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	session, err := testutils.SeedSession(app.GetSessionStore(), httptest.NewRequest(http.MethodGet, "/", nil), user, 3600)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	middleware := NewAPIAuthMiddleware(app).GetHandler()
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
		assert.True(t, ok, "session should be in request context")
		ctxUser, ok = r.Context().Value(config.APIAuthenticatedUserContextKey{}).(userstore.UserInterface)
		assert.True(t, ok, "user should be in request context")
		w.WriteHeader(http.StatusOK)
	})

	middleware(next).ServeHTTP(res, req)

	assert.True(t, nextCalled, "next handler should be called when authentication succeeds")
	if ctxSession != nil {
		assert.Equal(t, session.GetKey(), ctxSession.GetKey(), "context session key should match")
	}
	if ctxUser != nil {
		assert.Equal(t, user.ID(), ctxUser.ID(), "context user id should match")
	}
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}
