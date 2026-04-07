package register

import (
	"net/http"
	"net/url"
	"project/internal/config"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/auth"
	"github.com/dracory/test"
)

func TestRegisterController_RequiresAuthenticatedUser_WithoutVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context:   map[any]any{},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatal(`Response MUST be `, http.StatusSeeOther, ` but was: `, response.StatusCode)
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

	expected := "You must be logged in to access this page"
	if flashMessage.Message != expected {
		t.Fatal(`Response MUST contain '`, expected, `', but got: `, flashMessage.Message)
	}
}

func TestRegisterController_DisabledReturnsFlash(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)
	registry.GetConfig().SetRegistrationEnabled(false)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context:   map[any]any{},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("Response MUST be %d but was: %d", http.StatusSeeOther, response.StatusCode)
	}

	if !strings.Contains(responseHTML, `/flash?message_id=`) {
		t.Fatalf("Response MUST contain flash redirect, got: %s", responseHTML)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != helpers.FLASH_ERROR {
		t.Fatalf("Expected flash type %s but got %s", helpers.FLASH_ERROR, flashMessage.Type)
	}

	if flashMessage.Message != "Registrations are currently disabled" {
		t.Fatalf("Expected message 'Registrations are currently disabled' but got: %s", flashMessage.Message)
	}

	if flashMessage.Url != links.Website().Home() {
		t.Fatalf("Expected redirect URL %s but got %s", links.Website().Home(), flashMessage.Url)
	}
}

func TestRegisterController_RequiresAuthenticatedUser_WithVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
		Context:   map[any]any{},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusSeeOther {
		t.Fatal(`Response MUST be `, http.StatusSeeOther, ` but was: `, response.StatusCode)
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

	expected := "You must be logged in to access this page"
	if flashMessage.Message != expected {
		t.Fatal(`Response MUST contain '`, expected, `', but got: `, flashMessage.Message)
	}
}

func TestRegisterController_ShowsRegisterForm_WithoutVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(registry).Handler, test.NewRequestOptions{
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
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_ShowsRegisterForm_WithVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
		testutils.WithVaultStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(registry).Handler, test.NewRequestOptions{
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
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresFirstName_WithoutVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email": {user.Email()},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`First name is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresFirstName_WithVault(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
		testutils.WithVaultStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email": {user.Email()},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`First name is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresLastName_WithoutVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
			"first_name": {"FirstName"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Last name is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresLastName_WithVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
			"first_name": {"FirstName"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Last name is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresCountry_WithoutVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
			"first_name": {"FirstName"},
			"last_name":  {"LastName"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Country is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresCountry_WithVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
			"first_name": {"FirstName"},
			"last_name":  {"LastName"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Country is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresTimezone_WithoutVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
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
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Timezone is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_RequiresTimezone_WithVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
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
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Timezone is required field`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_Success_WithoutVault(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
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
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal("Response MUST NOT be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Your registration completed successfully. You can now continue browsing the website.`,
		`<script>window.location.href = '` + links.User().Home() + `'</script>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_SelectTimezoneByCountry_WithValidCountry(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	// Test the timezone selection AJAX endpoint
	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action":  {"on-country-selected-timezone-options"},
			"country": {"US"}, // United States
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	// Should return a timezone select element with actual timezone options
	expecteds := []string{
		`id="SelectTimezones"`,
		`name="timezone"`,
		`<select`,
		`</select>`,
		// Test for specific US timezone options
		`America/New_York`,
		`America/Los_Angeles`,
		`America/Chicago`,
		`America/Denver`,
		`Pacific/Honolulu`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}

	// Should NOT contain timezones from other countries
	notExpected := []string{
		`Europe/London`,    // UK timezone
		`Asia/Tokyo`,       // Japan timezone
		`Australia/Sydney`, // Australia timezone
		`Africa/Cairo`,     // Egypt timezone
		`America/Toronto`,  // Canada timezone (different country)
	}

	for _, unexpected := range notExpected {
		if strings.Contains(responseHTML, unexpected) {
			t.Fatal(`Response MUST NOT contain`, unexpected, ` but was found in: `, responseHTML)
		}
	}
}

func TestRegisterController_SelectTimezoneByCountry_WithEmptyCountry(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	// Test with empty country - should return all timezones
	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action":  {"on-country-selected-timezone-options"},
			"country": {""},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	// Should still return a timezone select element
	expecteds := []string{
		`id="SelectTimezones"`,
		`name="timezone"`,
		`<select`,
		`</select>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}

func TestRegisterController_SelectTimezoneByCountry_WithoutGeoStore(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
		// Note: No GeoStore - should return error
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	// Test without GeoStore configured
	_, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action":  {"on-country-selected-timezone-options"},
			"country": {"US"},
		},
		Context: map[any]any{
			auth.AuthenticatedUserID{}:           user.ID(),
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	// Should redirect to error page
	if response.StatusCode != http.StatusSeeOther {
		t.Fatal(`Response MUST be `, http.StatusSeeOther, ` but was: `, response.StatusCode)
	}

	location := response.Header.Get("Location")
	if !strings.Contains(location, `/flash?message_id=`) {
		t.Fatalf("Response Location MUST contain flash redirect, got: %s", location)
	}
}

func TestRegisterController_SelectTimezoneByCountry_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	// Test without authentication
	_, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action":  {"on-country-selected-timezone-options"},
			"country": {"US"},
		},
		Context: map[any]any{}, // No authenticated user
	})

	if err != nil {
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal(`Response MUST not be nil`)
	}

	// Should redirect to login
	if response.StatusCode != http.StatusSeeOther {
		t.Fatal(`Response MUST be `, http.StatusSeeOther, ` but was: `, response.StatusCode)
	}

	location := response.Header.Get("Location")
	if !strings.Contains(location, `/flash?message_id=`) {
		t.Fatalf("Response Location MUST contain flash redirect, got: %s", location)
	}
}

func TestRegisterController_Success_WithVaultStore(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"email":      {user.Email()},
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
		t.Fatal("Response MUST NOT trigger error, but was:", err)
	}

	if response == nil {
		t.Fatal("Response MUST NOT be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`id="FormRegister"`,
		`name="email"`,
		`name="first_name"`,
		`name="last_name"`,
		`name="country"`,
		`name="timezone"`,
		`Your registration completed successfully. You can now continue browsing the website.`,
		`<script>window.location.href = '` + links.User().Home() + `'</script>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(responseHTML, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, responseHTML)
		}
	}
}
