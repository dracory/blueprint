package user_update

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
	"github.com/stretchr/testify/assert"
)

func TestUserUpdateController_RequiresUserID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, response.StatusCode, "Should redirect with error")

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)
	assert.NoError(t, err, "Should find flash message")
	assert.Equal(t, "User ID is required", flash.Message, "Should show correct error message")
}

func TestUserUpdateController_InvalidUserID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {"invalid_id"},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, response.StatusCode, "Should redirect with error")

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)
	assert.NoError(t, err, "Should find flash message")
	assert.Equal(t, "User not found", flash.Message, "Should show correct error message")
}

func TestUserUpdateController_ShowsPage(t *testing.T) {
	registry, user := setupControllerAppAndUser(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {user.GetID()},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "Edit User", "Should show page heading")
	assert.Contains(t, responseHTML, "User:", "Should show user label")
}

func setupControllerAppAndUser(t *testing.T) (registry.RegistryInterface, userstore.UserInterface) {
	t.Helper()

	reg := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user := userstore.NewUser()
	user.SetFirstName("Test")
	user.SetLastName("User")
	user.SetEmail("test@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)

	if err := reg.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return reg, user
}

func TestActionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"actionUserFetch", actionUserFetch, "user-fetch-ajax"},
		{"actionGetTimezones", actionGetTimezones, "get-timezones-ajax"},
		{"actionUserUpdate", actionUserUpdate, "user-update-ajax"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestActionConstantsNotEmpty(t *testing.T) {
	constants := []string{
		actionUserFetch,
		actionGetTimezones,
		actionUserUpdate,
	}

	for _, c := range constants {
		if c == "" {
			t.Errorf("action constant should not be empty")
		}
		if !strings.HasSuffix(c, "-ajax") {
			t.Errorf("action constant %s should end with -ajax", c)
		}
	}
}
