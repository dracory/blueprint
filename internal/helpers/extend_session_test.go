package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

// Test ExtendSession with nil session store
func TestExtendSession_NilStore(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := ExtendSession(nil, req, 3600)
	if err == nil {
		t.Error("ExtendSession(nil store) expected error, got nil")
	}
	if err.Error() != "session store is nil" {
		t.Errorf("ExtendSession(nil store) error = %v, want 'session store is nil'", err.Error())
	}
}

// Test ExtendSession with session not found in context
func TestExtendSession_NoSessionInContext(t *testing.T) {
	registry := testutils.Setup(testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := ExtendSession(registry.GetSessionStore(), req, 3600)
	if err == nil {
		t.Error("ExtendSession(no session in context) expected error, got nil")
	}
	if err.Error() != "session not found" {
		t.Errorf("ExtendSession error = %v, want 'session not found'", err.Error())
	}
}

// Test ExtendSession with IP mismatch
func TestExtendSession_IPMismatch(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("extend@example.com").
		SetFirstName("Extend").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"

	// Create a session with different IP
	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent(req.UserAgent()).
		SetIPAddress("192.168.1.1"). // Different IP
		SetExpiresAt("2099-12-31 23:59:59")

	err = registry.GetSessionStore().SessionCreate(req.Context(), session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Add session to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedSessionContextKey{}, session)
	req = req.WithContext(ctx)

	// Test ExtendSession - should fail due to IP mismatch
	err = ExtendSession(registry.GetSessionStore(), req, 3600)
	if err == nil {
		t.Error("ExtendSession(IP mismatch) expected error, got nil")
	}
	if err.Error() != "session ip address does not match request ip address" {
		t.Errorf("ExtendSession error = %v, want IP mismatch error", err.Error())
	}
}

// Test ExtendSession with UserAgent mismatch
func TestExtendSession_UserAgentMismatch(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("extend2@example.com").
		SetFirstName("Extend2").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"

	// Create a session with different UserAgent
	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent("different-agent"). // Different UserAgent
		SetIPAddress("127.0.0.1").
		SetExpiresAt("2099-12-31 23:59:59")

	err = registry.GetSessionStore().SessionCreate(req.Context(), session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Add session to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedSessionContextKey{}, session)
	req = req.WithContext(ctx)

	// Test ExtendSession - should fail due to UserAgent mismatch
	err = ExtendSession(registry.GetSessionStore(), req, 3600)
	if err == nil {
		t.Error("ExtendSession(UserAgent mismatch) expected error, got nil")
	}
	if err.Error() != "session user agent does not match request user agent" {
		t.Errorf("ExtendSession error = %v, want UserAgent mismatch error", err.Error())
	}
}
