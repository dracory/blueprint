package register

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/internal/app"
	"project/internal/config"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/testutils"
	"strings"
	"testing"

	basetestutils "github.com/dracory/base/testutils"

	"github.com/dracory/auth"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// setupApp builds a test app with all stores the register flow requires.
// When vault is true the vault store is enabled as well.
func setupApp(t *testing.T, vault bool) app.AppInterface {
	t.Helper()
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	if vault {
		cfg.SetVaultStoreUsed(true)
		cfg.SetUserStoreVaultEnabled(true)
		// vaultstore requires keys of at least 16 characters
		cfg.SetVaultStoreKey("test-vault-key-0123456789")
	}
	return testutils.Setup(testutils.WithCfg(cfg))
}

// seedAuthUser creates a user and returns it together with the request
// context map that marks it as authenticated.
func seedAuthUser(t *testing.T, application app.AppInterface) (userstore.UserInterface, map[any]any) {
	t.Helper()
	user, err := testutils.SeedUser(application.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}
	// Users reaching the register page were created via email login — seed
	// an email to reflect that.
	if user.GetEmail() == "" {
		user.SetEmail("user01@test.com")
		if err := application.GetUserStore().UserUpdate(context.Background(), user); err != nil {
			t.Fatal(err)
		}
	}
	return user, map[any]any{
		auth.AuthenticatedUserID{}:           user.GetID(),
		config.AuthenticatedUserContextKey{}: user,
	}
}

// serve dispatches a request through the register controller the same way the
// router does: GET → PageHandler, POST → AjaxHandler.
func serve(t *testing.T, application app.AppInterface, method, path string,
	query url.Values, form url.Values, ctx map[any]any) *httptest.ResponseRecorder {
	t.Helper()
	req, err := test.NewRequest(method, path, test.NewRequestOptions{
		QueryParams: query,
		FormValues:  form,
		Context:     ctx,
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	controller := NewRegisterController(application)
	if method == http.MethodPost {
		controller.AjaxHandler(recorder, req)
	} else {
		_, _ = recorder.WriteString(controller.PageHandler(recorder, req))
	}
	return recorder
}

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var data map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON response %q: %v", recorder.Body.String(), err)
	}
	return data
}

func saveForm(overrides map[string]string) url.Values {
	form := url.Values{
		"email":      {"user01@test.com"},
		"first_name": {"FirstName"},
		"last_name":  {"LastName"},
		"country":    {"US"},
		"timezone":   {"America/New_York"},
	}
	for k, v := range overrides {
		form.Set(k, v)
	}
	return form
}

// == PAGE HANDLER ============================================================

func TestRegisterPage_UnauthenticatedRedirectsToFlash(t *testing.T) {
	application := setupApp(t, false)

	recorder := serve(t, application, http.MethodGet, "/", nil, nil, map[any]any{})

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected %d, got %d", http.StatusSeeOther, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `/flash?message_id=`) {
		t.Fatalf("expected flash redirect, got: %s", recorder.Body.String())
	}

	flashMessage, err := basetestutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}
	if flashMessage == nil {
		t.Fatal("expected flash message")
	}
	if flashMessage.Type != helpers.FLASH_ERROR {
		t.Fatalf("expected flash type %s, got %s", helpers.FLASH_ERROR, flashMessage.Type)
	}
	if flashMessage.Message != "You must be logged in to access this page" {
		t.Fatalf("unexpected flash message: %s", flashMessage.Message)
	}
}

func TestRegisterPage_RegistrationDisabledRedirectsToFlash(t *testing.T) {
	application := setupApp(t, false)
	application.GetConfig().SetRegistrationEnabled(false)

	_, ctx := seedAuthUser(t, application)
	recorder := serve(t, application, http.MethodGet, "/", nil, nil, ctx)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected %d, got %d", http.StatusSeeOther, recorder.Code)
	}

	flashMessage, err := basetestutils.FlashMessageFindFromResponse(application.GetCacheStore(), recorder.Result())
	if err != nil {
		t.Fatal(err)
	}
	if flashMessage == nil || flashMessage.Message != "Registrations are currently disabled" {
		t.Fatalf("unexpected flash message: %+v", flashMessage)
	}
	if flashMessage.Url != links.Website().Home() {
		t.Fatalf("expected redirect to %s, got %s", links.Website().Home(), flashMessage.Url)
	}
}

func TestRegisterPage_RendersVueApp(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodGet, "/", nil, nil, ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	expecteds := []string{
		`id="app"`,
		`REGISTER_AJAX_URL`,
		`TIMEZONES_AJAX_URL`,
		`INITIAL_DATA`,
		`COUNTRIES`,
		`first_name`,
		`last_name`,
		`business_name`,
		`timezone`,
	}
	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("response MUST contain %q", expected)
		}
	}
}

func TestRegisterPage_RendersVueApp_WithVault(t *testing.T) {
	application := setupApp(t, true)
	user, ctx := seedAuthUser(t, application)

	// With the vault enabled the stored email is a token — replace the
	// plaintext seed email with a real vault token.
	token, err := application.GetVaultStore().TokenCreate(context.Background(),
		user.GetEmail(), application.GetConfig().GetVaultStoreKey(), 20)
	if err != nil {
		t.Fatal(err)
	}
	user.SetEmail(token)
	if err := application.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	recorder := serve(t, application, http.MethodGet, "/", nil, nil, ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `id="app"`) {
		t.Fatal("expected Vue app container in response")
	}
}

// == AJAX HANDLER ============================================================

func TestRegisterAjax_Unauthenticated(t *testing.T) {
	application := setupApp(t, false)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(nil), map[any]any{})

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if data["status"] != "error" {
		t.Fatalf("expected error status, got %v", data)
	}
}

func TestRegisterAjax_InvalidAction(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"bogus"}}, nil, ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestRegisterSave_RequiresFirstName(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(map[string]string{"first_name": ""}), ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if !strings.Contains(data["message"].(string), "First name") {
		t.Fatalf("unexpected message: %v", data["message"])
	}
}

func TestRegisterSave_RequiresLastName(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(map[string]string{"last_name": ""}), ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if !strings.Contains(data["message"].(string), "Last name") {
		t.Fatalf("unexpected message: %v", data["message"])
	}
}

func TestRegisterSave_RequiresCountry(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(map[string]string{"country": ""}), ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if !strings.Contains(data["message"].(string), "Country") {
		t.Fatalf("unexpected message: %v", data["message"])
	}
}

func TestRegisterSave_RequiresTimezone(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(map[string]string{"timezone": ""}), ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if !strings.Contains(data["message"].(string), "Timezone") {
		t.Fatalf("unexpected message: %v", data["message"])
	}
}

func TestRegisterSave_Success(t *testing.T) {
	application := setupApp(t, false)
	user, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(map[string]string{
			"business_name": "Acme Ltd",
			"phone":         "+1234567890",
		}), ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}
	if data["redirect"] != links.User().Home() {
		t.Fatalf("expected redirect to %s, got %v", links.User().Home(), data["redirect"])
	}

	// Verify the profile was persisted
	updated, err := application.GetUserStore().UserFindByID(context.Background(), user.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if updated.GetFirstName() != "FirstName" || updated.GetLastName() != "LastName" {
		t.Fatalf("name not saved: %q %q", updated.GetFirstName(), updated.GetLastName())
	}
	if updated.GetBusinessName() != "Acme Ltd" {
		t.Fatalf("business name not saved: %q", updated.GetBusinessName())
	}
	if updated.GetCountry() != "US" || updated.GetTimezone() != "America/New_York" {
		t.Fatalf("country/timezone not saved: %q %q", updated.GetCountry(), updated.GetTimezone())
	}
}

func TestRegisterSave_Success_WithVault(t *testing.T) {
	application := setupApp(t, true)
	user, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"save"}}, saveForm(nil), ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}

	// Sensitive fields must be stored as vault tokens, not plaintext
	updated, err := application.GetUserStore().UserFindByID(context.Background(), user.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if updated.GetFirstName() == "FirstName" {
		t.Fatal("expected first name to be vault-tokenized")
	}
	if updated.GetCountry() != "US" || updated.GetTimezone() != "America/New_York" {
		t.Fatalf("country/timezone not saved: %q %q", updated.GetCountry(), updated.GetTimezone())
	}
}

func TestRegisterTimezones_FiltersByCountry(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"timezones"}}, url.Values{"country": {"US"}}, ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}

	raw, _ := json.Marshal(data["timezones"])
	list := string(raw)

	for _, expected := range []string{"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles"} {
		if !strings.Contains(list, expected) {
			t.Fatalf("expected timezone %q in list", expected)
		}
	}
	for _, unexpected := range []string{"Europe/London", "Asia/Tokyo", "Australia/Sydney"} {
		if strings.Contains(list, unexpected) {
			t.Fatalf("unexpected timezone %q in list", unexpected)
		}
	}
}

func TestRegisterTimezones_EmptyCountryReturnsAll(t *testing.T) {
	application := setupApp(t, false)
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"timezones"}}, url.Values{"country": {""}}, ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}
	raw, _ := json.Marshal(data["timezones"])
	if !strings.Contains(string(raw), "Europe/London") {
		t.Fatal("expected unfiltered timezone list")
	}
}

func TestRegisterTimezones_NoGeoStore(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))
	_, ctx := seedAuthUser(t, application)

	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"timezones"}}, url.Values{"country": {"US"}}, ctx)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}
