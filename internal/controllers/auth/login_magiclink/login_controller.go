// Package login_magiclink implements the in-house email magic-link login.
//
// Flow:
//  1. GET links.AUTH_LOGIN renders a self-contained page with an email form.
//  2. POST links.AUTH_LOGIN?action=magiclink-send generates a single-use,
//     IP-bound token, stores it in the memory cache, and enqueues
//     EmailMagicLinkTask to deliver the link by email.
//  3. GET links.AUTH_AUTH?token=... verifies the token (IP match +
//     single-use), then hands off to shared.SessionLogin.
package login_magiclink

import (
	"context"
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
	"net/url"
	"strings"
	"sync"
	"time"

	"project/internal/app"
	"project/internal/controllers/auth/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"
	"project/internal/tasks/email_magic_link"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

const (
	tokenTTL            = 15 * time.Minute
	magicLinkMaxSends   = 3
	sendsKeyPrefix      = "magiclink_sent:"
	tokenCacheKeyPrefix = "magiclink:token:"
	cacheKeyPrefix      = "magiclink:"
)

type loginController struct {
	app app.AppInterface
	// counterMu serializes the read-modify-write on the send-throttle
	// counter in the memory cache (ttlcache is safe per-call, but
	// Get-then-Set is not atomic).
	counterMu sync.Mutex
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
	case "magiclink-send-ajax":
		c.handleSend(w, r)
	default:
		c.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

// Handler verifies a magic link token (GET links.AUTH_AUTH?token=...).
func (c *loginController) Handler(w http.ResponseWriter, r *http.Request) string {
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

	if c.app.IsDisabledSessionStore() {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, `session store is required`, homeURL, 5)
	}

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Invalid or missing login token", homeURL, 5)
	}

	// Look up "email|ip|return" stored under the token key.
	cache := c.app.GetMemoryCache()
	tokenKey := tokenCacheKeyPrefix + token

	item := cache.Get(tokenKey)
	if item == nil {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Invalid or expired login link", homeURL, 5)
	}

	stored, ok := item.Value().(string)
	if !ok {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Invalid or expired login link", homeURL, 5)
	}

	parts := strings.SplitN(stored, "|", 3)
	if len(parts) < 2 {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Invalid or expired login link", homeURL, 5)
	}

	email, boundIP, returnURL := parts[0], parts[1], ""
	if len(parts) == 3 {
		returnURL = parts[2]
	}

	// Verify IP binding BEFORE consuming the token — a link scanner or
	// prefetcher from a different IP must not burn the link for the real
	// user.
	if boundIP != "" && req.GetIP(r) != boundIP {
		c.logWarn("Magic link IP mismatch", map[string]string{
			"bound_ip":   boundIP,
			"request_ip": req.GetIP(r),
		})
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "This login link was issued from a different IP address and cannot be used here.", homeURL, 5)
	}

	// IP matches — consume the token (single-use).
	cache.Delete(tokenKey)

	redirectURL, _, errorMessage := shared.SessionLogin(c.app, w, r, email, "")
	if errorMessage != "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, errorMessage, homeURL, 5)
	}

	// Optional return URL: only relative paths are accepted (open-redirect
	// protection). The value stored at send time takes precedence; the
	// query parameter is honoured as a fallback.
	if qReturn := strings.TrimSpace(r.URL.Query().Get("return")); returnURL == "" && qReturn != "" {
		returnURL = qReturn
	}
	if returnURL != "" && strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//") {
		redirectURL = returnURL
	}

	return helpers.ToFlashSuccess(c.app.GetCacheStore(), w, r, "Login was successful", redirectURL, 5)
}

func (c *loginController) renderLoginPage(r *http.Request) string {
	appName := "App"
	if c.app != nil && c.app.GetConfig() != nil {
		if c.app.GetConfig().GetAppName() != "" {
			appName = c.app.GetConfig().GetAppName()
		}
	}

	sendAjaxURL := links.AUTH_LOGIN + "?action=magiclink-send-ajax"

	// Optional post-login redirect target. Both `return` (OTP convention)
	// and `back_url` (links.Auth().Login convention) are accepted. Only
	// relative paths are honoured; the value is percent-encoded before
	// being embedded into the JS string literal below.
	returnURL := r.URL.Query().Get("return")
	if returnURL == "" {
		returnURL = r.URL.Query().Get("back_url")
	}
	if returnURL != "" && strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//") {
		returnURL = url.QueryEscape(returnURL)
	} else {
		returnURL = ""
	}

	// Replace placeholders in HTML with the actual app name (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))

	script := `
	const SEND_AJAX_URL = "` + sendAjaxURL + `";
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

// checkStores verifies the stores required by the magic link flow are
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
	if c.app.GetMemoryCache() == nil {
		c.sendErrorResponse(w, "Memory cache not initialized", http.StatusInternalServerError)
		return false
	}
	return true
}

// handleSend handles POST /auth/login?action=magiclink-send-ajax
func (c *loginController) handleSend(w http.ResponseWriter, r *http.Request) {
	if !c.checkStores(w) {
		return
	}

	if err := r.ParseForm(); err != nil {
		c.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		c.logError("Failed to parse form", err)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))

	if email == "" {
		c.sendErrorResponse(w, "Email is required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.sendErrorResponse(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Enforce the email allowlist before generating a token — a blocked
	// email should never receive a link or consume send quota.
	if rule := authrules.NewEmailAllowedRule(c.app, email); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	// Per-email send throttle to prevent email-bombing. Checked before the
	// fallible operations; recorded only after the email task is enqueued
	// so failed sends do not consume the quota.
	if !c.canSend(email) {
		c.sendErrorResponse(w, "Too many links requested, please try again later", http.StatusTooManyRequests)
		return
	}

	// Optional post-login redirect target — only relative paths accepted.
	returnURL := strings.TrimSpace(r.FormValue("return"))
	if returnURL != "" && (!strings.HasPrefix(returnURL, "/") || strings.HasPrefix(returnURL, "//")) {
		returnURL = ""
	}

	token, err := c.generateToken()
	if err != nil {
		c.sendErrorResponse(w, "Failed to generate login token", http.StatusInternalServerError)
		c.logError("Failed to generate token", err)
		return
	}

	nonce, err := c.generateNonce()
	if err != nil {
		c.sendErrorResponse(w, "Failed to generate nonce", http.StatusInternalServerError)
		c.logError("Failed to generate nonce", err)
		return
	}

	// Two cache entries:
	// - "magiclink:<nonce>" = "email|token|return" — read by the email task
	//   so the token is never persisted in the task queue.
	// - "magiclink:token:<token>" = "email|ip|return" — the verify path;
	//   deleted on use (single-use).
	cache := c.app.GetMemoryCache()
	cache.Set(cacheKeyPrefix+nonce, email+"|"+token+"|"+returnURL, tokenTTL)
	cache.Set(tokenCacheKeyPrefix+token, email+"|"+req.GetIP(r)+"|"+returnURL, tokenTTL)

	magicLinkTaskHandler := email_magic_link.NewEmailMagicLinkTask(c.app)
	magicLinkTask, ok := magicLinkTaskHandler.(*email_magic_link.EmailMagicLinkTask)
	if !ok {
		c.logError("Failed to cast email magic link task handler", nil)
		c.sendErrorResponse(w, "Failed to send login link", http.StatusInternalServerError)
		return
	}

	if _, err = magicLinkTask.Enqueue(nonce); err != nil {
		c.logError("Failed to enqueue magic link email task", err)
		c.sendErrorResponse(w, "Failed to send login link", http.StatusInternalServerError)
		return
	}

	c.recordSend(email)

	c.logInfo(fmt.Sprintf("Magic link email task enqueued for %s", email))

	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "success",
		"message": "Login link sent to email",
	}, http.StatusOK)
}

// canSend reports whether another link may be sent to the email within the
// current throttle window. The quota is consumed separately via recordSend
// after the email task has been enqueued successfully, so transient failures
// do not eat into the user's allowance.
func (c *loginController) canSend(email string) bool {
	return c.sendCount(email) < magicLinkMaxSends
}

// recordSend increments the per-email send counter.
func (c *loginController) recordSend(email string) {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()
	c.app.GetMemoryCache().Set(c.sendsKey(email), fmt.Sprintf("%d", c.sendCountLocked(email)+1), tokenTTL)
}

// sendCount returns the current per-email send count.
func (c *loginController) sendCount(email string) int {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()
	return c.sendCountLocked(email)
}

func (c *loginController) sendCountLocked(email string) int {
	count := 0
	if item := c.app.GetMemoryCache().Get(c.sendsKey(email)); item != nil {
		if s, ok := item.Value().(string); ok {
			fmt.Sscanf(s, "%d", &count)
		}
	}
	return count
}

func (c *loginController) sendsKey(email string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(email)))
	return sendsKeyPrefix + hex.EncodeToString(sum[:])
}

// generateToken generates a cryptographically random 256-bit hex token.
func (c *loginController) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateNonce generates a cryptographically random 128-bit hex nonce.
func (c *loginController) generateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

// logInfo logs an info message
func (c *loginController) logInfo(message string) {
	if logger := c.app.GetLogger(); logger != nil {
		logger.Info(message)
	}
}

// logWarn logs a warning message with optional structured attributes.
func (c *loginController) logWarn(message string, attrs map[string]string) {
	logger := c.app.GetLogger()
	if logger == nil {
		return
	}
	args := make([]slog.Attr, 0, len(attrs))
	for k, v := range attrs {
		args = append(args, slog.String(k, v))
	}
	logger.LogAttrs(context.Background(), slog.LevelWarn, message, args...)
}
