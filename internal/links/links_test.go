package links

import (
	"net/url"
	"strings"
	"testing"
)

func TestInitializeURLBuilder(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	initializeURLBuilder()
}

func TestRootURL(t *testing.T) {
	// In testing mode, RootURL always returns empty string
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	result := RootURL()
	if result == "" {
		t.Error("RootURL should return a URL even with empty APP_URL in testing mode")
	}
}

func TestQuery_EmptyQuery(t *testing.T) {
	result := query(map[string]string{})
	if result != "" {
		t.Errorf("query() = %q, want empty string", result)
	}
}

func TestQuery_SingleParameter(t *testing.T) {
	result := query(map[string]string{"key": "value"})
	if result == "" {
		t.Error("query() should return non-empty string for non-empty input")
	}
	if !strings.Contains(result, "key=value") {
		t.Errorf("query() = %q, should contain key=value", result)
	}
}

func TestQuery_MultipleParameters(t *testing.T) {
	result := query(map[string]string{"key1": "value1", "key2": "value2"})
	if result == "" {
		t.Error("query() should return non-empty string for non-empty input")
	}
	if !strings.Contains(result, "key") {
		t.Errorf("query() = %q, should contain key", result)
	}
}

func TestHttpBuildQuery_EmptyValues(t *testing.T) {
	result := httpBuildQuery(url.Values{})
	if result != "" {
		t.Errorf("httpBuildQuery() = %q, want empty string for empty input", result)
	}
}

func TestHttpBuildQuery_SingleValue(t *testing.T) {
	result := httpBuildQuery(url.Values{"key": []string{"value"}})
	if result == "" {
		t.Error("httpBuildQuery() should return non-empty string for non-empty input")
	}
	if !strings.Contains(result, "key=value") {
		t.Errorf("httpBuildQuery() = %q, should contain key=value", result)
	}
}

func TestHttpBuildQuery_MultipleValues(t *testing.T) {
	result := httpBuildQuery(url.Values{"key1": []string{"value1"}, "key2": []string{"value2"}})
	if result == "" {
		t.Error("httpBuildQuery() should return non-empty string for non-empty input")
	}
	if !strings.Contains(result, "key") {
		t.Errorf("httpBuildQuery() = %q, should contain key", result)
	}
}

func TestURL_SimplePath(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	result := URL("/test", nil)
	if result == "" {
		t.Error("URL() should return non-empty string for non-empty path")
	}
}

func TestURL_PathWithParameters(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	result := URL("/test", map[string]string{"key": "value"})
	if result == "" {
		t.Error("URL() should return non-empty string for non-empty path")
	}
}

func TestURL_EmptyPath(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	result := URL("", nil)
	if result == "" {
		t.Error("URL() should return non-empty string for empty path")
	}
}

func TestAuthLinks(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	auth := Auth()

	// Test Auth()
	result := auth.Auth()
	if result == "" {
		t.Error("auth.Auth() should return non-empty string")
	}

	// Test Auth with params
	result = auth.Auth(map[string]string{"key": "value"})
	if result == "" {
		t.Error("auth.Auth(params) should return non-empty string")
	}

	// Test Login()
	result = auth.Login("/back")
	if result == "" {
		t.Error("auth.Login() should return non-empty string")
	}

	// Test Login with params
	result = auth.Login("/back", map[string]string{"extra": "param"})
	if result == "" {
		t.Error("auth.Login(params) should return non-empty string")
	}

	// Test Logout()
	result = auth.Logout()
	if result == "" {
		t.Error("auth.Logout() should return non-empty string")
	}

	// Test Register()
	result = auth.Register()
	if result == "" {
		t.Error("auth.Register() should return non-empty string")
	}

	// Test AuthKnightLogin()
	result = auth.AuthKnightLogin("/back-url")
	if result == "" {
		t.Error("auth.AuthKnightLogin() should return non-empty string")
	}
	if !strings.Contains(result, "authknight.com") {
		t.Error("auth.AuthKnightLogin() should contain authknight.com")
	}
	if !strings.Contains(result, "back_url") {
		t.Error("auth.AuthKnightLogin() should contain back_url parameter")
	}
	if !strings.Contains(result, "next_url") {
		t.Error("auth.AuthKnightLogin() should contain next_url parameter")
	}
}

func TestUserLinks(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	user := User()

	// Test Home()
	result := user.Home()
	if result == "" {
		t.Error("user.Home() should return non-empty string")
	}

	// Test Home with params
	result = user.Home(map[string]string{"key": "value"})
	if result == "" {
		t.Error("user.Home(params) should return non-empty string")
	}

	// Test Profile()
	result = user.Profile()
	if result == "" {
		t.Error("user.Profile() should return non-empty string")
	}

	// Test Profile with params
	result = user.Profile(map[string]string{"key": "value"})
	if result == "" {
		t.Error("user.Profile(params) should return non-empty string")
	}
}
