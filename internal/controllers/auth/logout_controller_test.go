package auth

import (
	"net/http"
	"net/http/httptest"
	"project/internal/links"
	"project/internal/testutils"
	"testing"
	"time"

	"github.com/dracory/auth"
	"github.com/dracory/test"
	"github.com/gouniverse/responses"
)

func TestLogoutControllerHandler_SuccessfulLogout(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// Add an auth cookie to simulate a logged-in user
	cookie := &http.Cookie{
		Name:  auth.CookieName,
		Value: "test-token",
	}
	req.AddCookie(cookie)

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLogoutController(application).AnyIndex))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	// Verify auth cookie was removed
	cookies := recorder.Result().Cookies()
	hasActiveAuthCookie := false
	for _, cookie := range cookies {
		if cookie.Name == auth.CookieName {
			if time.Now().Before(cookie.Expires) {
				t.Fatal(`Auth cookie should not be present after logout`)
			}
		}
	}

	if hasActiveAuthCookie {
		t.Fatal(`Auth cookie should not be present after logout`)
	}

	// Verify flash message
	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	// Verify flash message type
	if flashMessage.Type != "success" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "You have been logged out successfully" {
		t.Fatal(`Response MUST contain 'You have been logged out successfully', but got: `, flashMessage.Message)
	}

	// Verify redirect URL
	if flashMessage.Url != links.Website().Home() {
		t.Fatal(`Flash message MUST contain redirect to home page, but got: `, flashMessage.Url)
	}
}

func TestLogoutControllerHandler_LogoutWithoutCookie(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLogoutController(application).AnyIndex))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	// Verify flash message still works even without existing cookie
	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	// Verify flash message type
	if flashMessage.Type != "success" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "You have been logged out successfully" {
		t.Fatal(`Response MUST contain 'You have been logged out successfully', but got: `, flashMessage.Message)
	}

	// Verify redirect URL
	if flashMessage.Url != links.Website().Home() {
		t.Fatal(`Flash message MUST contain redirect to home page, but got: `, flashMessage.Url)
	}
}
