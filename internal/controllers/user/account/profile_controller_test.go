package account

import (
	"net/http"
	"net/url"
	"project/internal/config"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/auth"
	"github.com/dracory/test"
)

func TestProfileController_RequiresAuthenticatedUser(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewProfileController(app).Handler, test.NewRequestOptions{
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
		t.Fatal("Response MUST be ", http.StatusSeeOther, " but was: ", response.StatusCode)
	}

	expecteds := []string{
		`<a href="/flash?message_id=`,
		`">See Other</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal("Response MUST contain", expected, " but was ", responseHTML)
		}
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatal(err)
	}
	if flashMessage == nil {
		t.Fatal("Response MUST contain 'flash message'")
	}
	if flashMessage.Type != "error" {
		t.Fatal("Response be of type 'error', but got: ", flashMessage.Type, flashMessage.Message)
	}
	if flashMessage.Message != "User not found" {
		t.Fatal("Response MUST contain 'User not found', but got: ", flashMessage.Message)
	}
}

func TestProfileController_ShowsProfileForm(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewProfileController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}

	expecteds := []string{
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="business_name"`,
		`name="phone"`,
		`name="country"`,
		`name="timezone"`,
		`data-flux-action="apply"`,
	}
	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal("Response MUST contain", expected, " but was ", responseHTML)
		}
	}
}
