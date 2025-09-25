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
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

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
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

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
		`id="FormProfile"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="business_name"`,
		`name="phone"`,
		`name="country"`,
		`name="timezone"`,
	}
	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal("Response MUST contain", expected, " but was ", responseHTML)
		}
	}
}

func TestProfileController_RequiresFirstName(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":     {user.Email()},
			"last_name": {"LastName"},
			"country":   {"Country"},
			"timezone":  {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "First name is required field") {
		t.Fatal("Response MUST contain validation error, but was ", responseHTML)
	}
}

func TestProfileController_RequiresLastName(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
			"first_name": {"FirstName"},
			"country":    {"Country"},
			"timezone":   {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "Last name is required field") {
		t.Fatal("Response MUST contain validation error, but was ", responseHTML)
	}
}

func TestProfileController_RequiresEmail(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"first_name": {"FirstName"},
			"last_name":  {"LastName"},
			"country":    {"Country"},
			"timezone":   {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "Email is required field") {
		t.Fatal("Response MUST contain validation error, but was ", responseHTML)
	}
}

func TestProfileController_RequiresCountry(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {"user@example.com"},
			"first_name": {"FirstName"},
			"last_name":  {"LastName"},
			"timezone":   {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "Country is required field") {
		t.Fatal("Response MUST contain validation error, but was ", responseHTML)
	}
}

func TestProfileController_RequiresTimezone(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {"user@example.com"},
			"first_name": {"FirstName"},
			"last_name":  {"LastName"},
			"country":    {"Country"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "Timezone is required field") {
		t.Fatal("Response MUST contain validation error, but was ", responseHTML)
	}
}

func TestProfileController_Success_NoVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":         {"user@example.com"},
			"first_name":    {"FirstName"},
			"last_name":     {"LastName"},
			"business_name": {"Biz"},
			"phone":         {"123"},
			"country":       {"Country"},
			"timezone":      {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}

	expecteds := []string{
		`id="FormProfile"`,
		"Profile updated successfully",
		`<script>window.location.href = '/flash?message_id=`,
	}
	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal("Response MUST contain", expected, " but was ", responseHTML)
		}
	}
}

func TestProfileController_Success_WithVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewProfileController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":         {"user@example.com"},
			"first_name":    {"FirstName"},
			"last_name":     {"LastName"},
			"business_name": {"Biz"},
			"phone":         {"123"},
			"country":       {"Country"},
			"timezone":      {"Timezone"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("Response MUST not be nil")
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Response MUST be ", http.StatusOK, " but was: ", response.StatusCode)
	}

	expecteds := []string{
		`id="FormProfile"`,
		"Profile updated successfully",
		`<script>window.location.href = '/flash?message_id=`,
	}
	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal("Response MUST contain", expected, " but was ", responseHTML)
		}
	}
}
