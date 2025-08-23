package auth

import (
	"net/http"
	"net/url"
	"project/internal/config"
	"project/internal/links"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/test"
	"github.com/gouniverse/auth"
)

func TestRegisterController_RequiresAuthenticatedUser_WithoutVault(t *testing.T) {
	application := testutils.Setup(testutils.WithVault(false))

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(application).Handler, test.NewRequestOptions{
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

	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), response)

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

func TestRegisterController_RequiresAuthenticatedUser_WithVault(t *testing.T) {
	application := testutils.Setup(testutils.WithVault(true))

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(application).Handler, test.NewRequestOptions{
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

	flashMessage, err := testutils.FlashMessageFindFromResponse(application.GetCacheStore(), response)

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
	application := testutils.Setup(testutils.WithVault(false))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(false))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(false))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(false))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(false))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
	application := testutils.Setup(testutils.WithVault(false))
	application.GetConfig().SetVaultStoreUsed(false)
	application.SetVaultStore(nil)

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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

func TestRegisterController_Success_WithVaultStore(t *testing.T) {
	application := testutils.Setup(testutils.WithVault(true))

	user, err := testutils.SeedUser(application.GetUserStore(), testutils.USER_01)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user should not be nil")
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodPost, NewRegisterController(application).Handler, test.NewRequestOptions{
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
