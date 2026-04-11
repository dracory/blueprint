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

// Test TimezoneFromRequest with nil request
func TestTimezoneFromRequest_NilRequest(t *testing.T) {
	result := TimezoneFromRequest(nil)
	if result != "UTC" {
		t.Errorf("TimezoneFromRequest(nil) = %v, want UTC", result)
	}
}

// Test TimezoneFromRequest without user in context
func TestTimezoneFromRequest_NoUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := TimezoneFromRequest(req)
	if result != "UTC" {
		t.Errorf("TimezoneFromRequest(req) without user = %v, want UTC", result)
	}
}

// TestTimezoneFromRequest_WithUserTimezone tests retrieving timezone from user
func TestTimezoneFromRequest_WithUserTimezone(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user with timezone
	user := userstore.NewUser().
		SetEmail("timezone@example.com").
		SetFirstName("TimeZone").
		SetLastName("User").
		SetTimezone("America/New_York")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Add user to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	// Test TimezoneFromRequest
	result := TimezoneFromRequest(req)
	if result != "America/New_York" {
		t.Errorf("TimezoneFromRequest = %v, want America/New_York", result)
	}
}

// TestTimezoneFromRequest_UserWithEmptyTimezone tests that empty timezone returns UTC
func TestTimezoneFromRequest_UserWithEmptyTimezone(t *testing.T) {
	registry := testutils.Setup(testutils.WithUserStore(true))
	defer registry.GetDatabase().Close()

	// Create a test user with empty timezone
	user := userstore.NewUser().
		SetEmail("notimezone@example.com").
		SetFirstName("NoTimeZone").
		SetLastName("User")

	err := registry.GetUserStore().UserCreate(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Add user to context
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user)
	req = req.WithContext(ctx)

	// Test TimezoneFromRequest - should return UTC when user has no timezone
	result := TimezoneFromRequest(req)
	if result != "UTC" {
		t.Errorf("TimezoneFromRequest = %v, want UTC", result)
	}
}
