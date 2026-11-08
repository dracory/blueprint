package theme

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/layouts"
)

func TestThemeControllerSetsCookieAndRedirects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/theme?theme=dark&redirect=/blog", nil)
	rec := httptest.NewRecorder()

	ThemeController(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if loc := rec.Header().Get("Location"); loc != "/blog" {
		t.Errorf("Location = %q, want %q", loc, "/blog")
	}

	cookies := rec.Result().Cookies()
	var themeCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == layouts.ThemeCookieKey {
			themeCookie = c
		}
	}
	if themeCookie == nil {
		t.Fatal("expected theme cookie to be set")
	}
	if themeCookie.Value != "dark" {
		t.Errorf("cookie value = %q, want %q", themeCookie.Value, "dark")
	}
}

func TestThemeControllerBlocksOpenRedirect(t *testing.T) {
	for _, redirect := range []string{"//evil.com", "https://evil.com", "evil.com"} {
		req := httptest.NewRequest(http.MethodGet, "/theme?theme=dark&redirect="+redirect, nil)
		rec := httptest.NewRecorder()

		ThemeController(rec, req)

		if loc := rec.Header().Get("Location"); loc != "/" {
			t.Errorf("redirect %q: Location = %q, want %q", redirect, loc, "/")
		}
	}
}

func TestThemeControllerNoThemeNoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/theme?redirect=/about", nil)
	rec := httptest.NewRecorder()

	ThemeController(rec, req)

	if loc := rec.Header().Get("Location"); loc != "/about" {
		t.Errorf("Location = %q, want %q", loc, "/about")
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == layouts.ThemeCookieKey {
			t.Error("theme cookie should not be set without theme param")
		}
	}
}

func TestThemeControllerDefaultRedirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/theme?theme=light", nil)
	rec := httptest.NewRecorder()

	ThemeController(rec, req)

	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("Location = %q, want %q", loc, "/")
	}
}
