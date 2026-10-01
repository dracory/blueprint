// Package theme switches the active Bootlegance theme. It reads the
// requested theme name from the query string, stores it in a cookie, and
// redirects back to the referring page.
package theme

import (
	"net/http"
	"strings"
	"time"

	"project/internal/layouts"

	"github.com/dracory/req"
)

// ThemeController handles GET /theme?theme={name}&redirect={path}.
func ThemeController(w http.ResponseWriter, r *http.Request) {
	themeName := req.GetStringTrimmed(r, "theme")
	redirect := req.GetStringTrimmedOr(r, "redirect", "/")

	// Only allow local redirects (absolute paths) to avoid open redirects.
	if !strings.HasPrefix(redirect, "/") || strings.HasPrefix(redirect, "//") {
		redirect = "/"
	}

	if themeName != "" {
		cookie := http.Cookie{
			Name:    layouts.ThemeCookieKey,
			Value:   themeName,
			Path:    "/",
			Secure:  r.TLS != nil,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
