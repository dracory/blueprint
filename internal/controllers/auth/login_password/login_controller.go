// Package login_password implements the in-house email/password login.
//
// Flow:
//  1. GET links.AUTH_LOGIN renders a self-contained page with an
//     email + password form.
//  2. POST links.AUTH_LOGIN?action=password-login-ajax verifies the
//     credentials against the user store (vault-aware lookup, prefers the
//     user record that has a password set) and hands off to
//     shared.SessionLoginUser to create the session.
//
// A valid email/password combination never creates a user — password login
// requires an account registered beforehand (e.g. via the register page).
package login_password

import (
	_ "embed"
	"encoding/json"
	"html"
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"project/internal/app"
	"project/internal/controllers/auth/shared"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

// dummyPasswordUser carries a valid bcrypt hash so a login attempt for an
// unknown email still pays the same PasswordCompare cost as a wrong
// password — preventing timing-based email enumeration.
var dummyPasswordUser = func() userstore.UserInterface {
	u := userstore.NewUser()
	_ = u.SetPasswordAndHash("dummy-password-for-timing-parity")
	return u
}()

type loginController struct {
	app app.AppInterface
}

func NewLoginController(application app.AppInterface) *loginController {
	return &loginController{app: application}
}

// PageHandler renders the login page. Registered via rtr.GetHTML.
func (c *loginController) PageHandler(w http.ResponseWriter, r *http.Request) string {
	return c.renderLoginPage(r)
}

// AjaxHandler dispatches POST actions on the login path. Registered via
// rtr.PostJSON — it returns a JSON string which rtr writes with the
// application/json content type.
func (c *loginController) AjaxHandler(w http.ResponseWriter, r *http.Request) string {
	switch r.URL.Query().Get("action") {
	case "password-login-ajax":
		c.handleLogin(w, r)
	default:
		c.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

func (c *loginController) renderLoginPage(r *http.Request) string {
	appName := "App"
	if c.app != nil && c.app.GetConfig() != nil {
		if c.app.GetConfig().GetAppName() != "" {
			appName = c.app.GetConfig().GetAppName()
		}
	}

	loginAjaxURL := links.AUTH_LOGIN + "?action=password-login-ajax"

	// Optional post-login redirect target. Both `return` (OTP convention)
	// and `back_url` (links.Auth().Login convention) are accepted. Only
	// relative paths are honoured; the value is percent-encoded before
	// being embedded into the JS string literal below.
	returnURL := r.URL.Query().Get("return")
	if returnURL == "" {
		returnURL = r.URL.Query().Get("back_url")
	}
	if returnURL != "" && (!strings.HasPrefix(returnURL, "/") || strings.HasPrefix(returnURL, "//")) {
		returnURL = ""
	}
	rawReturnURL := returnURL
	returnURL = url.QueryEscape(returnURL)

	registerURL := links.Auth().Register()
	if rawReturnURL != "" {
		// Pass the raw path — links.URL() URL-encodes the param itself;
		// the escaped copy is only for the JS string literal below.
		registerURL = links.Auth().Register(map[string]string{"back_url": rawReturnURL})
	}

	// Replace placeholders in HTML with the actual app name (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))
	htmlContent = strings.ReplaceAll(htmlContent, "{{ registerUrl }}", registerURL)
	htmlContent = strings.ReplaceAll(htmlContent, "{{ forgotPasswordUrl }}", links.Auth().ForgotPassword())

	script := `
	const LOGIN_AJAX_URL = "` + loginAjaxURL + `";
	const RETURN_URL = "` + returnURL + `";
	` + appJS

	layout := layouts.NewBlankLayout(c.app, r, baselayouts.Options{
		AppName: appName,
		Title:   "Login",
		Content: hb.Div().HTML(htmlContent),
		StyleURLs: []string{
			cdn.BootstrapIconsCss_1_11_3(),
			cdn.Notiflix_3_2_8_CSS(),
		},
		ScriptURLs: []string{
			cdn.VueJs_3_5_32(),
			cdn.Notiflix_3_2_8(),
		},
		Scripts: []string{
			script,
		},
		Styles: []string{
			appCSS,
		},
	})
	return layout.ToHTML()
}

// checkStores verifies the stores required by the password login flow are
// available.
func (c *loginController) checkStores(w http.ResponseWriter) bool {
	if c.app.IsDisabledUserStore() {
		c.sendErrorResponse(w, "User store not initialized", http.StatusInternalServerError)
		return false
	}
	if c.app.GetConfig().GetUserStoreVaultEnabled() && c.app.IsDisabledVaultStore() {
		c.sendErrorResponse(w, "Vault store not initialized", http.StatusInternalServerError)
		return false
	}
	if c.app.GetConfig().GetUserStoreVaultEnabled() && c.app.IsDisabledBlindIndexStoreEmail() {
		c.sendErrorResponse(w, "Blind index store not initialized", http.StatusInternalServerError)
		return false
	}
	if c.app.IsDisabledSessionStore() {
		c.sendErrorResponse(w, "Session store not initialized", http.StatusInternalServerError)
		return false
	}
	return true
}

// handleLogin handles POST /auth/login?action=password-login-ajax
func (c *loginController) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !c.checkStores(w) {
		return
	}

	if rule := authrules.NewCanUsePasswordAuthRule(c.app); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		c.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		c.logError("Failed to parse form", err)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	// Only relative paths are accepted for the post-login redirect
	// (open-redirect protection).
	returnURL := strings.TrimSpace(r.FormValue("return"))
	if returnURL != "" && (!strings.HasPrefix(returnURL, "/") || strings.HasPrefix(returnURL, "//")) {
		returnURL = ""
	}

	if email == "" || password == "" {
		c.sendErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.sendErrorResponse(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Enforce the email allowlist before the credential check — a blocked
	// email should never reach the user store.
	if rule := authrules.NewEmailAllowedRule(c.app, email); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	user, err := shared.UserFindByEmail(c.app, r.Context(), email)
	if err != nil {
		c.logError("Error finding user", err)
		c.sendErrorResponse(w, "Error finding user", http.StatusInternalServerError)
		return
	}

	// Unknown email: run a bcrypt compare against a dummy hash anyway so the
	// response time is indistinguishable from a wrong-password attempt.
	if user == nil {
		dummyPasswordUser.PasswordCompare(password)
		c.sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Use an identical message for both missing user and wrong password to
	// prevent email enumeration attacks.
	if !user.PasswordCompare(password) {
		c.sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	redirectURL, _, errorMessage := shared.SessionLoginUser(c.app, w, r, user, returnURL)
	if errorMessage != "" {
		status := http.StatusForbidden
		if errorMessage == shared.MsgSessionError {
			// Session-creation failures are internal errors, not auth denials
			status = http.StatusInternalServerError
		}
		c.sendErrorResponse(w, errorMessage, status)
		return
	}

	c.sendJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"message":  "Login successful",
		"redirect": redirectURL,
	}, http.StatusOK)
}

// sendJSONResponse sends a JSON response
func (c *loginController) sendJSONResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (c *loginController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "error",
		"message": message,
	}, statusCode)
}

// logError logs an error
func (c *loginController) logError(message string, err error) {
	logger := c.app.GetLogger()
	if logger == nil {
		return
	}
	if err != nil {
		logger.Error(message, slog.String("error", err.Error()))
	} else {
		logger.Error(message)
	}
}
