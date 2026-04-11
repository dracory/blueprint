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

// Test GetAuthUser with nil request
func TestGetAuthUser_NilRequest(t *testing.T) {
	result := GetAuthUser(nil)
	if result != nil {
		t.Errorf("GetAuthUser(nil) = %v, want nil", result)
	}
}

// Test GetAuthUser without context value
func TestGetAuthUser_NoContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := GetAuthUser(req)
	if result != nil {
		t.Errorf("GetAuthUser(req) without context = %v, want nil", result)
	}
}

// TestGetAuthUser_WithValidUser tests retrieving user from context
func TestGetAuthUser_WithValidUser(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("test2@example.com").
		SetFirstName("Test2").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Add user to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	// Test GetAuthUser
	result := GetAuthUser(req)
	if result == nil {
		t.Fatal("GetAuthUser returned nil for valid user in context")
	}

	if result.ID() != user.ID() {
		t.Errorf("GetAuthUser ID = %v, want %v", result.ID(), user.ID())
	}
}

// Test GetAPIAuthUser with nil request
func TestGetAPIAuthUser_NilRequest(t *testing.T) {
	result := GetAPIAuthUser(nil)
	if result != nil {
		t.Errorf("GetAPIAuthUser(nil) = %v, want nil", result)
	}
}

// Test GetAPIAuthUser without context value
func TestGetAPIAuthUser_NoContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := GetAPIAuthUser(req)
	if result != nil {
		t.Errorf("GetAPIAuthUser(req) without context = %v, want nil", result)
	}
}

// TestGetAPIAuthUser_WithValidUser tests retrieving API user from context
func TestGetAPIAuthUser_WithValidUser(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user
	user := userstore.NewUser().
		SetEmail("apiuser@example.com").
		SetFirstName("API").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	// Add API user to context
	ctx := context.WithValue(req.Context(), config.APIAuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	// Test GetAPIAuthUser
	result := GetAPIAuthUser(req)
	if result == nil {
		t.Fatal("GetAPIAuthUser returned nil for valid user in context")
	}

	if result.ID() != user.ID() {
		t.Errorf("GetAPIAuthUser ID = %v, want %v", result.ID(), user.ID())
	}
}
