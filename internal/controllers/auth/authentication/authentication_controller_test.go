package authentication

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), recorder.Result())

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
	app := testutils.Setup(testutils.WithCfg(cfg))

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
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()

	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), recorder.Result())

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{
			"once": {testutils.TestKey(app.GetConfig())},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()
	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), recorder.Result())

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetUserStore() == nil {
		t.Fatal("UserStore should not be nil")
	}

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	user.SetEmail("test@test.com")

	err = app.GetUserStore().UserUpdate(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		PostValues: url.Values{
			"once": {testutils.TestKey(app.GetConfig())},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)
	// response := recorder.Body.String()
	code := recorder.Code

	if code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), recorder.Result())

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

func TestAuthController_NilUserStore(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(false) // No session store
	cfg.SetUserStoreUsed(false)    // No user store
	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf(`Expected status 303, got %d`, recorder.Code)
	}
}

func TestAuthController_NilSessionStore(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(false) // No session store
	cfg.SetUserStoreUsed(true)     // Has user store
	app := testutils.Setup(testutils.WithCfg(cfg))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = NewAuthenticationController(app).Handler(w, r)
	})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf(`Expected status 303, got %d`, recorder.Code)
	}
}

func TestAuthController_NewController(t *testing.T) {
	app := testutils.Setup()
	controller := NewAuthenticationController(app)

	if controller == nil {
		t.Fatal("Controller should not be nil")
	}

	if controller.app == nil {
		t.Fatal("Controller app should not be nil")
	}
}
