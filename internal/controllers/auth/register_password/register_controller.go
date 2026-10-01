// Package register_password implements the public sign-up page used when
// the login method is email/password.
//
// Unlike the post-authentication register controller (used by the external
// and passwordless login methods to complete a profile after first login),
// this controller creates the account itself: email, first/last name and
// password are collected up front, the user is created and logged in via
// shared.SessionLoginUser.
//
// Flow:
//  1. GET links.AUTH_REGISTER renders a self-contained sign-up form.
//  2. POST links.AUTH_REGISTER?action=register-ajax validates the form,
//     creates the user with a hashed password, and starts a session.
package register_password

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
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
	"github.com/dracory/str"
	"github.com/dracory/userstore"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

// genericExistsMsg is returned when the email is already registered so the
// endpoint cannot be used to enumerate accounts.
const genericExistsMsg = "If this email is not registered, your account has been created."

type registerController struct {
	app app.AppInterface
}

func NewRegisterController(application app.AppInterface) *registerController {
	return &registerController{app: application}
}

// PageHandler renders the registration page. Registered via rtr.GetHTML.
func (c *registerController) PageHandler(w http.ResponseWriter, r *http.Request) string {
	return c.renderPage(r)
}

// AjaxHandler dispatches POST actions on the register path. Registered via
// rtr.PostJSON — it returns a JSON string which rtr writes with the
// application/json content type.
func (c *registerController) AjaxHandler(w http.ResponseWriter, r *http.Request) string {
	switch r.URL.Query().Get("action") {
	case "register-ajax":
		c.handleSubmit(w, r)
	default:
		c.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

func (c *registerController) renderPage(r *http.Request) string {
	appName := "App"
	if c.app != nil && c.app.GetConfig() != nil {
		if c.app.GetConfig().GetAppName() != "" {
			appName = c.app.GetConfig().GetAppName()
		}
	}

	submitAjaxURL := links.AUTH_REGISTER + "?action=register-ajax"

	// Optional post-registration redirect target — only relative paths are
	// honoured (open-redirect protection); percent-encoded before being
	// embedded into the JS string literal below.
	returnURL := r.URL.Query().Get("return")
	if returnURL == "" {
		returnURL = r.URL.Query().Get("back_url")
	}
	if returnURL != "" && strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//") {
		returnURL = url.QueryEscape(returnURL)
	} else {
		returnURL = ""
	}

	loginURL := links.Auth().Login("")
	if returnURL != "" {
		loginURL = links.Auth().Login(r.URL.Query().Get("return"))
		if rawBackURL := r.URL.Query().Get("back_url"); rawBackURL != "" {
			loginURL = links.Auth().Login(rawBackURL)
		}
	}

	// Replace placeholders in HTML with the actual values (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))
	htmlContent = strings.ReplaceAll(htmlContent, "{{ loginUrl }}", loginURL)

	script := `
	const SUBMIT_AJAX_URL = "` + submitAjaxURL + `";
	const RETURN_URL = "` + returnURL + `";
	` + appJS

	layout := layouts.NewBlankLayout(c.app, r, baselayouts.Options{
		AppName: appName,
		Title:   "Register",
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

// checkStores verifies the stores required by the registration flow are
// available.
func (c *registerController) checkStores(w http.ResponseWriter) bool {
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

// handleSubmit handles POST /auth/register?action=register-ajax
func (c *registerController) handleSubmit(w http.ResponseWriter, r *http.Request) {
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

	// Trim inputs and strip HTML tags to prevent stored XSS
	firstName := str.StripHTMLTags(strings.TrimSpace(r.FormValue("first_name")))
	lastName := str.StripHTMLTags(strings.TrimSpace(r.FormValue("last_name")))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	// Only relative paths are accepted for the post-login redirect
	// (open-redirect protection).
	returnURL := strings.TrimSpace(r.FormValue("return"))
	if returnURL != "" && (!strings.HasPrefix(returnURL, "/") || strings.HasPrefix(returnURL, "//")) {
		returnURL = ""
	}

	if firstName == "" || lastName == "" {
		c.sendErrorResponse(w, "First and last name are required", http.StatusBadRequest)
		return
	}

	if email == "" {
		c.sendErrorResponse(w, "Email is required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.sendErrorResponse(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	if rule := authrules.NewCanRegisterRule(c.app, email); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	if rule := authrules.NewPasswordPolicyRule(authrules.PasswordPolicyData{
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusBadRequest)
		return
	}

	// Check if user already exists — return a generic message to prevent
	// email enumeration
	existingUser, err := shared.UserFindByEmail(c.app, r.Context(), email)
	if err != nil {
		c.logError("Error checking existing user", err)
		c.sendErrorResponse(w, "Registration failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		c.sendJSONResponse(w, map[string]interface{}{
			"status":   "success",
			"message":  genericExistsMsg,
			"redirect": links.Auth().Login(returnURL),
		}, http.StatusOK)
		return
	}

	// Create the user account
	user, err := shared.UserCreate(c.app, r.Context(), email, userstore.USER_STATUS_ACTIVE)
	if err != nil {
		c.logError("Error creating user", err)
		c.sendErrorResponse(w, "Registration failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	// From here on, any failure must roll back the created user — otherwise
	// an orphan record with no password would permanently block both
	// registration (email exists) and password recovery (no first name).
	rollback := func() {
		if delErr := c.app.GetUserStore().UserDelete(r.Context(), user); delErr != nil {
			c.logError("Failed to roll back partially created user", delErr)
		}
	}

	// Store the profile names (vault-tokenizing when the vault is enabled)
	if err := c.applyNames(r.Context(), user, firstName, lastName); err != nil {
		c.logError("Error saving user names", err)
		rollback()
		c.sendErrorResponse(w, "Registration failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Hash and store the password
	if err := user.SetPasswordAndHash(password); err != nil {
		c.logError("Error hashing password", err)
		rollback()
		c.sendErrorResponse(w, "Error securing password. Please try again.", http.StatusInternalServerError)
		return
	}

	if err := c.app.GetUserStore().UserUpdate(r.Context(), user); err != nil {
		c.logError("Error saving user", err)
		rollback()
		c.sendErrorResponse(w, "Registration failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Log the new user in
	redirectURL, _, errorMessage := shared.SessionLoginUser(c.app, w, r, user, returnURL)
	if errorMessage != "" {
		c.sendJSONResponse(w, map[string]interface{}{
			"status":   "success",
			"message":  "Account created! Please log in.",
			"redirect": links.Auth().Login(returnURL),
		}, http.StatusOK)
		return
	}

	c.sendJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"message":  "Account created successfully!",
		"redirect": redirectURL,
	}, http.StatusOK)
}

// applyNames sets the first and last name on the user, vault-tokenizing
// the values when the vault store is enabled.
func (c *registerController) applyNames(ctx context.Context, user userstore.UserInterface, firstName string, lastName string) error {
	if !c.app.GetConfig().GetUserStoreVaultEnabled() {
		user.SetFirstName(firstName)
		user.SetLastName(lastName)
		return nil
	}

	if c.app.IsDisabledVaultStore() {
		return errors.New("vault store is not configured")
	}

	firstNameToken, err := c.app.GetVaultStore().TokenCreate(ctx, firstName, c.app.GetConfig().GetVaultStoreKey(), 20)
	if err != nil {
		return err
	}
	lastNameToken, err := c.app.GetVaultStore().TokenCreate(ctx, lastName, c.app.GetConfig().GetVaultStoreKey(), 20)
	if err != nil {
		return err
	}

	user.SetFirstName(firstNameToken)
	user.SetLastName(lastNameToken)
	return nil
}

// sendJSONResponse sends a JSON response
func (c *registerController) sendJSONResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (c *registerController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "error",
		"message": message,
	}, statusCode)
}

// logError logs an error
func (c *registerController) logError(message string, err error) {
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
