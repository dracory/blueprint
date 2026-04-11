package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/userstore"
)

// Test UserSettingGet with nil session store
func TestUserSettingGet_NilStore(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := UserSettingGet(nil, req, "key", "default")
	if result != "default" {
		t.Errorf("UserSettingGet(nil store) = %v, want 'default'", result)
	}
}

// Test UserSettingGet with nil request
func TestUserSettingGet_NilRequest(t *testing.T) {
	result := UserSettingGet(nil, nil, "key", "default")
	if result != "default" {
		t.Errorf("UserSettingGet(nil request) = %v, want 'default'", result)
	}
}

// Test UserSettingSet with nil session store
func TestUserSettingSet_NilStore(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := UserSettingSet(nil, req, "key", "value")
	if err == nil {
		t.Error("UserSettingSet(nil store) expected error, got nil")
	}
	if err.Error() != "session store is nil" {
		t.Errorf("UserSettingSet(nil store) error = %v, want 'session store is nil'", err.Error())
	}
}

// Test UserSettingSet with nil request (checks store first)
func TestUserSettingSet_NilRequest(t *testing.T) {
	err := UserSettingSet(nil, nil, "key", "value")
	if err == nil {
		t.Error("UserSettingSet(nil store and request) expected error, got nil")
	}
	// UserSettingSet checks for nil store first
	if err.Error() != "session store is nil" {
		t.Errorf("UserSettingSet(nil store) error = %v, want 'session store is nil'", err.Error())
	}
}

// Test UserSettingGet with no user in context
func TestUserSettingGet_NoUser(t *testing.T) {
	registry := testutils.Setup(testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := UserSettingGet(registry.GetSessionStore(), req, "key", "default")
	if result != "default" {
		t.Errorf("UserSettingGet(no user) = %v, want 'default'", result)
	}
}

// Test UserSettingSet with no user in context
func TestUserSettingSet_NoUser(t *testing.T) {
	registry := testutils.Setup(testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := UserSettingSet(registry.GetSessionStore(), req, "key", "value")
	if err == nil {
		t.Error("UserSettingSet(no user) expected error, got nil")
	}
	if err.Error() != "auth user is nil" {
		t.Errorf("UserSettingSet error = %v, want 'auth user is nil'", err.Error())
	}
}

// Test UserSettingGet with no session found
func TestUserSettingGet_NoSession(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("settings@example.com").
		SetFirstName("Settings").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Add user to context but no session
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	result := UserSettingGet(registry.GetSessionStore(), req, "nonexistent-key", "default")
	if result != "default" {
		t.Errorf("UserSettingGet(no session) = %v, want 'default'", result)
	}
}

// Test UserSettingSet with no session found
func TestUserSettingSet_NoSession(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("settings2@example.com").
		SetFirstName("Settings2").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Add user to context but no session
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	err = UserSettingSet(registry.GetSessionStore(), req, "nonexistent-key", "value")
	if err == nil {
		t.Error("UserSettingSet(no session) expected error, got nil")
	}
	if err.Error() != "session is nil" {
		t.Errorf("UserSettingSet error = %v, want 'session is nil'", err.Error())
	}
}
