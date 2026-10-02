// Package password_reset implements the "choose a new password" page for
// the email/password login method.
//
// Flow:
//  1. GET links.AUTH_PASSWORD_RESET?token=... validates the single-use
//     reset token stored by the forgot-password controller in the memory
//     cache and renders the new-password form.
//  2. POST links.AUTH_PASSWORD_RESET?action=password-reset-ajax validates
//     the token again, applies the shared password policy rule, updates the
//     user's password hash, and consumes the token.
package password_reset

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"html"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/app"
	"project/internal/controllers/auth/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

const tokenCacheKeyPrefix = "passwordreset:token:"

// passwordResetCacheValue mirrors the JSON payload stored by the
// forgot-password controller under "passwordreset:token:<token>".
type passwordResetCacheValue struct {
	Email string `json:"email"`
}

type passwordResetController struct {
	app app.AppInterface
}

func NewPasswordResetController(application app.AppInterface) *passwordResetController {
	return &passwordResetController{app: application}
}

// PageHandler renders the password reset page. Registered via rtr.GetHTML.
func (c *passwordResetController) PageHandler(w http.ResponseWriter, r *http.Request) string {
	homeURL := links.Website().Home()

	if c.app.IsDisabledUserStore() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `user store is required`, homeURL, 5)
	}

	if c.app.GetConfig().GetUserStoreVaultEnabled() && c.app.IsDisabledVaultStore() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `vault store is required`, homeURL, 5)
	}

	if c.app.GetConfig().GetUserStoreVaultEnabled() && c.app.IsDisabledBlindIndexStoreEmail() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `blind index store is required`, homeURL, 5)
	}

	if c.app.GetMemoryCache() == nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `memory cache is required`, homeURL, 5)
	}

	token := strings.TrimSpace(req.GetStringTrimmed(r, "token"))
	if token == "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `Invalid or missing reset token`, homeURL, 5)
	}

	// Verify the token exists before rendering the form — the submit
	// handler validates it again (it may expire in between).
	if _, err := c.tokenEmail(r.Context(), token); err != nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `Invalid or expired reset link`, homeURL, 5)
	}

	return c.renderPage(r, token)
}

// AjaxHandler dispatches POST actions on the password reset path.
// Registered via rtr.PostJSON.
func (c *passwordResetController) AjaxHandler(w http.ResponseWriter, r *http.Request) string {
	switch r.URL.Query().Get("action") {
	case "password-reset-ajax":
		c.handleSubmit(w, r)
	default:
		c.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

func (c *passwordResetController) renderPage(r *http.Request, token string) string {
	appName := "App"
	if c.app != nil && c.app.GetConfig() != nil {
		if c.app.GetConfig().GetAppName() != "" {
			appName = c.app.GetConfig().GetAppName()
		}
	}

	submitAjaxURL := links.AUTH_PASSWORD_RESET + "?action=password-reset-ajax"

	// Replace placeholders in HTML with the actual values (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))

	script := `
	const SUBMIT_AJAX_URL = "` + submitAjaxURL + `";
	const LOGIN_URL = "` + links.Auth().Login("") + `";
	const TOKEN = "` + token + `";
	` + appJS

	layout := layouts.NewBlankLayout(c.app, r, baselayouts.Options{
		AppName: appName,
		Title:   "Choose a New Password",
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

// handleSubmit handles POST /auth/password-reset?action=password-reset-ajax
func (c *passwordResetController) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if c.app.IsDisabledUserStore() || c.app.GetMemoryCache() == nil {
		c.sendErrorResponse(w, "Required stores not initialized", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		c.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		c.logError("Failed to parse form", err)
		return
	}

	token := strings.TrimSpace(r.FormValue("token"))
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	if token == "" {
		c.sendErrorResponse(w, "Reset token is required", http.StatusBadRequest)
		return
	}

	if rule := authrules.NewPasswordPolicyRule(authrules.PasswordPolicyData{
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusBadRequest)
		return
	}

	// Consume the token atomically (GetAndDelete) so two concurrent
	// requests cannot both pass the token check — the link is single-use.
	// If the password update below fails the user must request a fresh
	// reset link, which is the acceptable trade-off.
	email, err := c.consumeToken(token)
	if err != nil {
		c.sendErrorResponse(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := shared.UserFindByEmail(c.app, r.Context(), email)
	if err != nil {
		c.logError("Failed to find user", err)
		c.sendErrorResponse(w, "Failed to reset password. Please try again.", http.StatusInternalServerError)
		return
	}

	if user == nil {
		c.sendErrorResponse(w, "Account not found", http.StatusBadRequest)
		return
	}

	// Set the new password
	if err := user.SetPasswordAndHash(password); err != nil {
		c.logError("Failed to hash password", err)
		c.sendErrorResponse(w, "Failed to reset password. Please try again.", http.StatusInternalServerError)
		return
	}

	if err := c.app.GetUserStore().UserUpdate(r.Context(), user); err != nil {
		c.logError("Failed to save password", err)
		c.sendErrorResponse(w, "Failed to save password. Please try again.", http.StatusInternalServerError)
		return
	}

	// Invalidate all existing sessions — an attacker holding a session must
	// lose access once the legitimate user resets the password.
	c.invalidateUserSessions(r.Context(), user.GetID())

	c.sendJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"message":  "Password reset successfully! You can now log in.",
		"redirect": links.Auth().Login(""),
	}, http.StatusOK)
}

// tokenEmail returns the email a reset token was issued for, or an error
// when the token does not exist or has expired.
func (c *passwordResetController) tokenEmail(ctx context.Context, token string) (string, error) {
	cache := c.app.GetMemoryCache()
	if cache == nil {
		return "", errors.New("memory cache is required")
	}

	item := cache.Get(tokenCacheKeyPrefix + token)
	if item == nil {
		return "", errors.New("token not found or expired")
	}

	stored, ok := item.Value().(string)
	if !ok {
		return "", errors.New("token not found or expired")
	}

	var value passwordResetCacheValue
	if err := json.Unmarshal([]byte(stored), &value); err != nil || value.Email == "" {
		return "", errors.New("token not found or expired")
	}

	return value.Email, nil
}

// consumeToken atomically reads and deletes a reset token so it cannot be
// used twice.
func (c *passwordResetController) consumeToken(token string) (string, error) {
	cache := c.app.GetMemoryCache()
	if cache == nil {
		return "", errors.New("memory cache is required")
	}

	item, found := cache.GetAndDelete(tokenCacheKeyPrefix + token)
	if !found || item == nil {
		return "", errors.New("token not found or expired")
	}

	stored, ok := item.Value().(string)
	if !ok {
		return "", errors.New("token not found or expired")
	}

	var value passwordResetCacheValue
	if err := json.Unmarshal([]byte(stored), &value); err != nil || value.Email == "" {
		return "", errors.New("token not found or expired")
	}

	return value.Email, nil
}

// invalidateUserSessions deletes all sessions belonging to the user so that
// stolen sessions lose access once the password is reset.
func (c *passwordResetController) invalidateUserSessions(ctx context.Context, userID string) {
	if c.app.IsDisabledSessionStore() {
		return
	}

	sessions, err := c.app.GetSessionStore().SessionList(ctx, sessionstore.NewSessionQuery().
		SetUserID(userID))
	if err != nil {
		c.logError("Failed to list user sessions for invalidation", err)
		return
	}

	for _, session := range sessions {
		if err := c.app.GetSessionStore().SessionDelete(ctx, session); err != nil {
			c.logError("Failed to delete session during password reset", err)
		}
	}
}

// sendJSONResponse sends a JSON response
func (c *passwordResetController) sendJSONResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (c *passwordResetController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "error",
		"message": message,
	}, statusCode)
}

// logError logs an error
func (c *passwordResetController) logError(message string, err error) {
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
