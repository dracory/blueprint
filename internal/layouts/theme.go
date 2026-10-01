package layouts

import (
	"net/http"
	"regexp"
)

// ThemeCookieKey is the name of the cookie that stores the selected
// Bootlegance theme name.
const ThemeCookieKey = "theme"

// ThemeDefault is the Bootlegance theme used when no cookie is set.
const ThemeDefault = "vanguard"

// bootleganceVersion pins the Bootlegance release served via jsDelivr.
const bootleganceVersion = "0.5.0"

// BootleganceThemes lists the themes available in the pinned Bootlegance
// release. Each ships a theme.css under themes/{name}/.
var BootleganceThemes = []string{
	"brutalski",
	"carbon",
	"civic",
	"seneca",
	"vanguard",
	"whitehall",
	"yaru",
}

// themeNameSanitizer allows only lowercase letters in theme names, matching
// the Bootlegance theme directory naming.
var themeNameSanitizer = regexp.MustCompile(`[^a-z]`)

// ThemeName resolves the active Bootlegance theme for the request. It reads
// the theme cookie, sanitizes it, and falls back to ThemeDefault when empty.
func ThemeName(r *http.Request) string {
	if r == nil {
		return ThemeDefault
	}
	cookie, err := r.Cookie(ThemeCookieKey)
	if err != nil || cookie == nil {
		return ThemeDefault
	}
	name := themeNameSanitizer.ReplaceAllString(cookie.Value, "")
	if name == "" {
		return ThemeDefault
	}
	return name
}

// ThemeStyleURL returns the CDN URL of the theme's stylesheet.
func ThemeStyleURL(theme string) string {
	return "https://cdn.jsdelivr.net/gh/lesichkovm/bootlegance@v" +
		bootleganceVersion + "/themes/" + theme + "/theme.css"
}
