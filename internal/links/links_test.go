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

func TestQuery(t *testing.T) {
	tests := []struct {
		name      string
		queryData map[string]string
		contains  string
	}{
		{
			name:      "Empty query",
			queryData: map[string]string{},
			contains:  "",
		},
		{
			name:      "Single parameter",
			queryData: map[string]string{"key": "value"},
			contains:  "key=value",
		},
		{
			name:      "Multiple parameters",
			queryData: map[string]string{"key1": "value1", "key2": "value2"},
			contains:  "key", // Just check it contains something
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := query(tt.queryData)
			if tt.contains != "" {
				if result == "" {
					t.Error("query() should return non-empty string for non-empty input")
				}
				if !strings.Contains(result, tt.contains) {
					t.Errorf("query() = %q, should contain %q", result, tt.contains)
				}
			}
			if tt.contains == "" && result != "" {
				t.Errorf("query() = %q, want empty string", result)
			}
		})
	}
}

func TestHttpBuildQuery(t *testing.T) {
	tests := []struct {
		name      string
		queryData url.Values
		contains  string
	}{
		{
			name:      "Empty values",
			queryData: url.Values{},
			contains:  "",
		},
		{
			name:      "Single value",
			queryData: url.Values{"key": []string{"value"}},
			contains:  "key=value",
		},
		{
			name:      "Multiple values",
			queryData: url.Values{"key1": []string{"value1"}, "key2": []string{"value2"}},
			contains:  "key", // Order may vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := httpBuildQuery(tt.queryData)
			if tt.contains != "" {
				if result == "" {
					t.Error("httpBuildQuery() should return non-empty string for non-empty input")
				}
				if !strings.Contains(result, tt.contains) {
					t.Errorf("httpBuildQuery() = %q, should contain %q", result, tt.contains)
				}
			}
			if len(tt.queryData) == 0 && result != "" {
				t.Errorf("httpBuildQuery() = %q, want empty string for empty input", result)
			}
		})
	}
}

func TestURL(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	tests := []struct {
		name   string
		path   string
		params map[string]string
	}{
		{
			name:   "Simple path",
			path:   "/test",
			params: nil,
		},
		{
			name:   "Path with parameters",
			path:   "/test",
			params: map[string]string{"key": "value"},
		},
		{
			name:   "Empty path",
			path:   "",
			params: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := URL(tt.path, tt.params)
			if result == "" && tt.path != "" {
				t.Error("URL() should return non-empty string for non-empty path")
			}
		})
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
