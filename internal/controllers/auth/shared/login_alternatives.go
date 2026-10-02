package shared

import (
	"html"
	"strings"

	"project/internal/config"
	"project/internal/links"
)

// loginMethodLabels maps each login method to the label shown on the
// login page's "use a different method" buttons.
var loginMethodLabels = map[string]string{
	config.LOGIN_METHOD_MAGICLINK:  "Use magic link instead",
	config.LOGIN_METHOD_OTP:        "Use a one-time code instead",
	config.LOGIN_METHOD_PASSWORD:   "Use password instead",
	config.LOGIN_METHOD_AUTHKNIGHT: "Sign in with AuthKnight",
}

// loginMethodPath returns the login page path for a method — AUTH_LOGIN
// for the primary method, the dedicated AUTH_LOGIN_* path otherwise.
func loginMethodPath(method string, methods []string) string {
	if method == methods[0] {
		return links.AUTH_LOGIN
	}
	switch method {
	case config.LOGIN_METHOD_MAGICLINK:
		return links.AUTH_LOGIN_MAGICLINK
	case config.LOGIN_METHOD_OTP:
		return links.AUTH_LOGIN_OTP
	case config.LOGIN_METHOD_PASSWORD:
		return links.AUTH_LOGIN_PASSWORD
	case config.LOGIN_METHOD_AUTHKNIGHT:
		return links.AUTH_LOGIN_AUTHKNIGHT
	}
	return ""
}

// LoginAlternatives renders the "OR — use another method" section as
// trusted HTML to be injected into a login template's
// {{ alternatives }} placeholder. It returns an empty string when fewer
// than two methods are enabled.
//
// Output contains only fixed method labels and constant paths — no
// user-controlled input — so callers may inject it unescaped.
func LoginAlternatives(current string, methods []string) string {
	if len(methods) < 2 {
		return ""
	}

	var buttons strings.Builder
	for _, m := range methods {
		if m == current {
			continue
		}
		label, ok := loginMethodLabels[m]
		if !ok {
			continue
		}
		path := loginMethodPath(m, methods)
		if path == "" {
			continue
		}
		buttons.WriteString(`<a href="` + html.EscapeString(path) + `" class="btn btn-outline-secondary w-100 py-2 fw-bold mb-2">` +
			html.EscapeString(label) + `</a>`)
	}

	if buttons.Len() == 0 {
		return ""
	}

	return `<div class="mt-4">` +
		`<div class="d-flex align-items-center mb-3">` +
		`<hr class="flex-grow-1 my-0">` +
		`<span class="px-3 text-muted small fw-bold">OR</span>` +
		`<hr class="flex-grow-1 my-0">` +
		`</div>` +
		buttons.String() +
		`</div>`
}
