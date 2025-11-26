package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/test"
)

func TestLoginControllerHandler_UserStoreNotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(false)
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	// Simulate user store not used/unavailable
	application.SetUserStore(nil)
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewLoginController(application).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'error', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "user store is required" {
		t.Fatal(`Response MUST contain 'user store is required', but got: `, flashMessage.Message)
	}
}

func TestLoginControllerHandler_VaultStoreNotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	// Ensure user store exists but vault store is required and missing
	if application.GetUserStore() == nil {
		t.Fatal("user store should be initialized in test setup")
	}
	application.GetConfig().SetVaultStoreUsed(true)
	application.SetVaultStore(nil)
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewLoginController(application).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'error', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "vault store is required" {
		t.Fatal(`Response MUST contain 'vault store is required', but got: `, flashMessage.Message)
	}
}

func TestLoginControllerHandler_ValidRedirectWithBackURL(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	backUrl := "http://localhost/home"
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{
		GetValues: url.Values{
			"back_url": {backUrl},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewLoginController(application).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}
}

func TestLoginControllerHandler_InvalidBackURL(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	invalidBackUrl := "http://evil.com"
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{
		GetValues: url.Values{
			"back_url": {invalidBackUrl},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewLoginController(application).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}
}
