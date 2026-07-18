package user_manager

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

func TestUserManagerController_ShowsPage(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewUserManagerController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil { t.Errorf("Handler should not return error: %v", err) }
	if http.StatusOK != response.StatusCode { t.Errorf("Should return 200 status: expected %v, got %v", http.StatusOK, response.StatusCode) }
	if !strings.Contains(responseHTML, "Users") { t.Errorf("Should show page heading: expected %q to contain %q", responseHTML, "Users") }
}

func TestActionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"actionLoadUsers", actionLoadUsers, "load-users-ajax"},
		{"actionDeleteUser", actionDeleteUser, "delete-user-ajax"},
		{"actionCreateUser", actionCreateUser, "create-user-ajax"},
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
		actionLoadUsers,
		actionDeleteUser,
		actionCreateUser,
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

func setupControllerAppAndUser(t *testing.T) (app.AppInterface, userstore.UserInterface) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user := userstore.NewUser()
	user.SetFirstName("Test")
	user.SetLastName("User")
	user.SetEmail("test@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)

	if err := app.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return app, user
}
