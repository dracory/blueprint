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

// Test GetAuthSession with nil request
func TestGetAuthSession_NilRequest(t *testing.T) {
	result := GetAuthSession(nil)
	if result != nil {
		t.Errorf("GetAuthSession(nil) = %v, want nil", result)
	}
}

// Test GetAuthSession without context value
func TestGetAuthSession_NoContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := GetAuthSession(req)
	if result != nil {
		t.Errorf("GetAuthSession(req) without context = %v, want nil", result)
	}
}

// TestGetAuthSession_WithValidSession tests retrieving session from context
func TestGetAuthSession_WithValidSession(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("test@example.com").
		SetFirstName("Test").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a request with user agent and IP
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"

	// Create a session
	session := sessionstore.NewSession().
		SetUserID(user.GetID()).
		SetUserAgent(req.UserAgent()).
		SetIPAddress("127.0.0.1").
		SetExpiresAt("2099-12-31 23:59:59")

	err = registry.GetSessionStore().SessionCreate(req.Context(), session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Add session to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedSessionContextKey{}, session)
	req = req.WithContext(ctx)

	// Test GetAuthSession
	result := GetAuthSession(req)
	if result == nil {
		t.Fatal("GetAuthSession returned nil for valid session in context")
	}

	if result.GetUserID() != user.GetID() {
		t.Errorf("GetAuthSession user ID = %v, want %v", result.GetUserID(), user.GetID())
	}
}
