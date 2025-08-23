package user_test

import (
	"net/http"
	"net/url"
	user "project/internal/controllers/user/home"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_HomeController_RedirectsIfUserNotLoggedIn(t *testing.T) {
	app := testutils.Setup()

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, user.NewHomeController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context:   map[any]any{},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	code := response.StatusCode

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "User not found" {
		t.Fatal(`Response MUST contain 'User not found', but got: `, flashMessage.Message)
	}

	expecteds := []string{
		`<a href="/flash?message_id=`,
		`">See Other</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}
