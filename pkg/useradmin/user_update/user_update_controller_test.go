package user_update

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/app"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func TestUserUpdateController_RequiresUserID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("Handler should not return error: %v", err) }
	if http.StatusSeeOther != response.StatusCode { t.Errorf("Should redirect with error: expected %v, got %v", http.StatusSeeOther, response.StatusCode) }

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil { t.Errorf("Should find flash message: %v", err) }
	if "User ID is required" != flash.Message { t.Errorf("Should show correct error message: expected %v, got %v", "User ID is required", flash.Message) }
}

func TestUserUpdateController_InvalidUserID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {"invalid_id"},
		},
	})

	if err != nil { t.Errorf("Handler should not return error: %v", err) }
	if http.StatusSeeOther != response.StatusCode { t.Errorf("Should redirect with error: expected %v, got %v", http.StatusSeeOther, response.StatusCode) }

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil { t.Errorf("Should find flash message: %v", err) }
	if "User not found" != flash.Message { t.Errorf("Should show correct error message: expected %v, got %v", "User not found", flash.Message) }
}

func TestUserUpdateController_ShowsPage(t *testing.T) {
	app, user := setupControllerAppAndUser(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {user.GetID()},
		},
	})

	if err != nil { t.Errorf("Handler should not return error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("Should return 200 status: expected %v, got %v", http.StatusOK, response.StatusCode) }
	if !strings.Contains(responseHTML, "Edit User") { t.Errorf("Should show page heading: expected %q to contain %q", responseHTML, "Edit User") }
	if !strings.Contains(responseHTML, "User:") { t.Errorf("Should show user label: expected %q to contain %q", responseHTML, "User:") }
}

func setupControllerAppAndUser(t *testing.T) (app.AppInterface, userstore.UserInterface) {
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
