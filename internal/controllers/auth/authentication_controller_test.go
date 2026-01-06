package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/test"
)

func TestAuthControllerOnceIsRequired(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(registry).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), recorder.Result())

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "Authentication Provider Error. Once is required field" {
		t.Fatal(`Response MUST contain 'Authentication Provider Error. Once is required field', but got: `, flashMessage.Message)
	}
}

func TestAuthControllerOnceMustBeValid(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{
			"once": {"test"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(registry).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()

	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), recorder.Result())

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "Authentication Provider Error. Invalid authentication response status" {
		t.Fatal(`Response MUST contain 'Authentication Provider Error. Invalid authentication response status', but got: `, flashMessage.Message, flashMessage.Message)
	}
}

func TestAuthControllerOnceSuccessWithNewUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{
			"once": {testutils.TestKey(registry.GetConfig())},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(registry).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()
	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), recorder.Result())

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "success" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	expected := "Login was successful"
	if flashMessage.Message != "Login was successful" {
		t.Fatal(`Response MUST contain '`+expected+`', but got: `, flashMessage.Message)
	}
}

func TestAuthControllerOnceSuccessWithExistingUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetUserStore() == nil {
		t.Fatal("UserStore should not be nil")
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	user.SetEmail("test@test.com")

	err = registry.GetUserStore().UserUpdate(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{
			"once": {testutils.TestKey(registry.GetConfig())},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(registry).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()
	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), recorder.Result())

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "success" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	expected := "Login was successful"
	if flashMessage.Message != "Login was successful" {
		t.Fatal(`Response MUST contain '`+expected+`', but got: `, flashMessage.Message)
	}
}
