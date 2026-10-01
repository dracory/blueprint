// Package forgot_password implements the "request a password reset link"
// page for the email/password login method.
//
// Flow:
//  1. GET links.AUTH_FORGOT_PASSWORD renders a self-contained page with an
//     email + first-name form (the first name acts as a lightweight
//     identity-verification step, like CourseThread).
//  2. POST links.AUTH_FORGOT_PASSWORD?action=forgot-password-ajax verifies
//     the details, generates a single-use reset token stored in the memory
//     cache, and enqueues EmailPasswordResetTask to deliver the link.
//
// All failure paths return an identical generic response (plus a
// constant-time delay) to prevent email enumeration via message or timing.
package forgot_password

import (
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"project/internal/app"
	"project/internal/controllers/auth/shared"
	"project/internal/ext"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"
	"project/internal/tasks/email_password_reset"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

const (
	tokenTTL            = 15 * time.Minute
	resetMaxSends       = 3
	sendsKeyPrefix      = "passwordreset_sent:"
	tokenCacheKeyPrefix = "passwordreset:token:"
	cacheKeyPrefix      = "passwordreset:"

	// genericSuccessMsg is returned for every non-system failure so the
	// endpoint cannot be used to enumerate registered emails.
	genericSuccessMsg = "If the account details match, a reset link has been sent."
)

// passwordResetCacheValue is the JSON payload stored in the memory cache.
// The token entry ("passwordreset:token:<token>") carries the email used at
// reset time; the nonce entry ("passwordreset:<nonce>") carries
// email+token for the email task to resolve.
type passwordResetCacheValue struct {
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

type forgotPasswordController struct {
	app app.AppInterface
	// counterMu serializes the read-modify-write on the send-throttle
	// counter in the memory cache (ttlcache is safe per-call, but
	// Get-then-Set is not atomic).
	counterMu sync.Mutex
}

func NewForgotPasswordController(application app.AppInterface) *forgotPasswordController {
	return &forgotPasswordController{app: application}
}

// PageHandler renders the forgot password page. Registered via rtr.GetHTML.
func (c *forgotPasswordController) PageHandler(w http.ResponseWriter, r *http.Request) string {
	return c.renderPage(r)
}

// AjaxHandler dispatches POST actions on the forgot password path.
// Registered via rtr.PostJSON.
func (c *forgotPasswordController) AjaxHandler(w http.ResponseWriter, r *http.Request) string {
	switch r.URL.Query().Get("action") {
	case "forgot-password-ajax":
		c.handleSubmit(w, r)
	default:
		c.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

func (c *forgotPasswordController) renderPage(r *http.Request) string {
	appName := "App"
	if c.app != nil && c.app.GetConfig() != nil {
		if c.app.GetConfig().GetAppName() != "" {
			appName = c.app.GetConfig().GetAppName()
		}
	}

	submitAjaxURL := links.AUTH_FORGOT_PASSWORD + "?action=forgot-password-ajax"

	// Replace placeholders in HTML with the actual values (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))
	htmlContent = strings.ReplaceAll(htmlContent, "{{ loginUrl }}", links.Auth().Login(""))

	script := `
	const SUBMIT_AJAX_URL = "` + submitAjaxURL + `";
	const LOGIN_URL = "` + links.Auth().Login("") + `";
	` + appJS

	layout := layouts.NewBlankLayout(c.app, r, baselayouts.Options{
		AppName: appName,
		Title:   "Reset Password",
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

// checkStores verifies the stores required by the forgot password flow are
// available.
func (c *forgotPasswordController) checkStores(w http.ResponseWriter) bool {
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
	if c.app.GetMemoryCache() == nil {
		c.sendErrorResponse(w, "Memory cache not initialized", http.StatusInternalServerError)
		return false
	}
	return true
}

// handleSubmit handles POST /auth/forgot-password?action=forgot-password-ajax
func (c *forgotPasswordController) handleSubmit(w http.ResponseWriter, r *http.Request) {
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
	firstName := strings.TrimSpace(r.FormValue("first_name"))

	if email == "" {
		c.sendErrorResponse(w, "Email is required", http.StatusBadRequest)
		return
	}

	if firstName == "" {
		c.sendErrorResponse(w, "First name is required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.sendErrorResponse(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Find user by email — any lookup failure returns the generic success
	// message so the endpoint cannot be used to enumerate emails.
	user, err := shared.UserFindByEmail(c.app, r.Context(), email)
	if err != nil {
		c.logError("Failed to find user", err)
		c.delayedSuccess(w)
		return
	}

	if user == nil {
		c.delayedSuccess(w)
		return
	}

	// Get the user's actual first name (handles vault decryption)
	_, userFirstName, _, _, _, err := ext.UserUntokenizeTransparently(r.Context(), c.app, user)
	if err != nil {
		c.logError("Failed to untokenize user", err)
		c.delayedSuccess(w)
		return
	}

	// Verify first name matches — a lightweight identity check so a reset
	// link is only sent to someone who knows more than just the email.
	if !strings.EqualFold(firstName, userFirstName) {
		c.delayedSuccess(w)
		return
	}

	// Per-email send throttle to prevent email-bombing. The quota is
	// reserved atomically and released if the send setup below fails.
	if !c.tryAcquireSend(email) {
		c.delayedSuccess(w)
		return
	}

	token, err := c.generateToken()
	if err != nil {
		c.releaseSend(email)
		c.logError("Failed to generate reset token", err)
		c.sendErrorResponse(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	nonce, err := c.generateNonce()
	if err != nil {
		c.releaseSend(email)
		c.logError("Failed to generate nonce", err)
		c.sendErrorResponse(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Note: the memory cache is process-local — an app restart between
	// enqueue and task execution drops the pending reset email, and issued
	// links expire before their 15-minute TTL. Accepted trade-off matching
	// the magic-link login.
	//
	// Two cache entries (JSON-encoded):
	// - "passwordreset:<nonce>" = {email, token} — read by the email task so
	//   the token is never persisted in the task queue.
	// - "passwordreset:token:<token>" = {email} — the reset path; deleted on
	//   use (single-use).
	nonceValue, _ := json.Marshal(passwordResetCacheValue{Email: email, Token: token})
	tokenValue, _ := json.Marshal(passwordResetCacheValue{Email: email})
	cache := c.app.GetMemoryCache()
	cache.Set(cacheKeyPrefix+nonce, string(nonceValue), tokenTTL)
	cache.Set(tokenCacheKeyPrefix+token, string(tokenValue), tokenTTL)

	resetTaskHandler := email_password_reset.NewEmailPasswordResetTask(c.app)
	resetTask, ok := resetTaskHandler.(*email_password_reset.EmailPasswordResetTask)
	if !ok {
		c.releaseSend(email)
		c.logError("Failed to cast email password reset task handler", nil)
		c.sendErrorResponse(w, "Failed to send reset link", http.StatusInternalServerError)
		return
	}

	if _, err = resetTask.Enqueue(nonce); err != nil {
		c.releaseSend(email)
		c.logError("Failed to enqueue password reset email task", err)
		c.sendErrorResponse(w, "Failed to send reset link", http.StatusInternalServerError)
		return
	}

	c.logInfo("Password reset email task enqueued for " + c.sendsKey(email))

	c.delayedSuccess(w)
}

// delayedSuccess sleeps briefly (constant-time response) and returns the
// generic success message. Used by every non-system failure path so the
// endpoint cannot be used for email enumeration.
func (c *forgotPasswordController) delayedSuccess(w http.ResponseWriter) {
	time.Sleep(300 * time.Millisecond)
	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "success",
		"message": genericSuccessMsg,
	}, http.StatusOK)
}

// tryAcquireSend atomically checks the per-email throttle and, when under
// the limit, reserves one send slot. Holding counterMu across the
// check-and-increment closes the check-then-act race between concurrent
// requests. Callers that subsequently fail must call releaseSend.
func (c *forgotPasswordController) tryAcquireSend(email string) bool {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()
	count := c.sendCountLocked(email)
	if count >= resetMaxSends {
		return false
	}
	c.app.GetMemoryCache().Set(c.sendsKey(email), fmt.Sprintf("%d", count+1), tokenTTL)
	return true
}

// releaseSend decrements the per-email send counter, undoing a reservation
// made by tryAcquireSend when the send setup fails.
func (c *forgotPasswordController) releaseSend(email string) {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()
	count := c.sendCountLocked(email)
	if count <= 0 {
		return
	}
	c.app.GetMemoryCache().Set(c.sendsKey(email), fmt.Sprintf("%d", count-1), tokenTTL)
}

func (c *forgotPasswordController) sendCountLocked(email string) int {
	count := 0
	if item := c.app.GetMemoryCache().Get(c.sendsKey(email)); item != nil {
		if s, ok := item.Value().(string); ok {
			fmt.Sscanf(s, "%d", &count)
		}
	}
	return count
}

func (c *forgotPasswordController) sendsKey(email string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(email)))
	return sendsKeyPrefix + hex.EncodeToString(sum[:])
}

// generateToken generates a cryptographically random 256-bit hex token.
func (c *forgotPasswordController) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateNonce generates a cryptographically random 128-bit hex nonce.
func (c *forgotPasswordController) generateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// sendJSONResponse sends a JSON response
func (c *forgotPasswordController) sendJSONResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (c *forgotPasswordController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "error",
		"message": message,
	}, statusCode)
}

// logError logs an error
func (c *forgotPasswordController) logError(message string, err error) {
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

// logInfo logs an info message
func (c *forgotPasswordController) logInfo(message string) {
	if logger := c.app.GetLogger(); logger != nil {
		logger.Info(message)
	}
}
