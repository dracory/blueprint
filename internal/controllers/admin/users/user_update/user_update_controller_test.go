package admin

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/auth"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func TestUserUpdateController_RequiresUserID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context:   map[any]any{},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response MUST not be nil")
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("Response MUST be %d but was %d", http.StatusSeeOther, response.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal("Response MUST contain flash message")
	}

	if flashMessage.Message != "User ID is required" {
		t.Fatalf("Flash message MUST be 'User ID is required', got %q", flashMessage.Message)
	}

	if !strings.Contains(responseHTML, "See Other") {
		t.Fatalf("Response MUST contain redirect notice, got: %s", responseHTML)
	}
}

func TestUserUpdateController_ShowsForm(t *testing.T) {
	app, user := setupControllerAppAndUser(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewUserUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {user.ID()},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Handler MUST NOT return error, but got:", err)
	}

	if response == nil {
		t.Fatal("Response MUST not be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Response MUST be %d but was %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="user_status"`,
		`name="user_first_name"`,
		`name="user_last_name"`,
		`name="user_email"`,
		`name="user_business_name"`,
		`name="user_phone"`,
		`name="user_country"`,
		`name="user_timezone"`,
		`data-flux-action="apply"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatalf("Response MUST contain %q but was %s", expected, responseHTML)
		}
	}

	if !strings.Contains(responseHTML, "User: "+user.FirstName()+" "+user.LastName()) {
		t.Fatalf("Response MUST contain user title for %s", user.ID())
	}
}

func setupControllerAppAndUser(t *testing.T) (registry.RegistryInterface, userstore.UserInterface) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetFirstName("John")
	user.SetLastName("Doe")
	user.SetEmail("john@example.com")
	user.SetBusinessName("JD Consulting")
	user.SetPhone("+44111222333")
	user.SetMemo("Initial memo")
	user.SetCountry("GB")
	user.SetTimezone("Europe/London")

	if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatalf("UserUpdate returned error: %v", err)
	}

	return app, user
}
