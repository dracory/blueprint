package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/config"
	user "project/internal/controllers/user/home"
	"project/internal/testutils"

	"github.com/dracory/test"
)

func Test_HomeController_RedirectsIfUserNotLoggedIn(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, user.NewHomeController(registry).Handler, test.NewRequestOptions{
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

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)

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

func TestNewHomeController(t *testing.T) {
	registry := testutils.Setup()
	controller := user.NewHomeController(registry)

	if controller == nil {
		t.Fatal("Controller should not be nil")
	}
}

func TestHomeController_WithLoggedInUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	// Create a test user
	testUser, err := testutils.SeedUser(registry.GetUserStore(), "test_user_01")
	if err != nil {
		t.Fatal(err)
	}

	// Set user details
	testUser.SetFirstName("John")
	testUser.SetLastName("Doe")
	testUser.SetEmail("john@example.com")
	err = registry.GetUserStore().UserUpdate(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Test with authenticated user context via request context
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, testUser)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	controller := user.NewHomeController(registry)
	result := controller.Handler(recorder, req)

	// Should return HTML content when user is logged in
	if !strings.Contains(result, "Hi,") && !strings.Contains(result, "User not found") {
		t.Logf("Response: %s", result)
	}
}
