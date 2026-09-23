package login_otp

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"math/big"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"sync"
	"time"

	"project/internal/app"
	"project/internal/controllers/auth/shared"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"
	"project/internal/tasks/email_otp"

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
	otpTTL            = 15 * time.Minute
	otpMaxAttempts    = 5
	otpMaxSends       = 3
	otpSendsKeyPrefix = "otp_sent:"
)

type loginController struct {
	app app.AppInterface
	// counterMu serializes the read-modify-write on the send-throttle and
	// attempt counters in the memory cache (ttlcache is safe per-call, but
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
	action := r.URL.Query().Get("action")
	switch action {
	case "otp-send-ajax":
		c.handleOtpSend(w, r)
	case "otp-verify-ajax":
		c.handleOtpVerify(w, r)
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

	loginAjaxURL := links.AUTH_LOGIN + "?action=otp-send-ajax"
	verifyAjaxURL := links.AUTH_LOGIN + "?action=otp-verify-ajax"
	returnURL := r.URL.Query().Get("return")
	if returnURL != "" && strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//") {
		// QueryEscape keeps the value intact for the verify handler while
		// encoding quotes so it cannot break out of the JS string literal
		// the URL is embedded into below.
		verifyAjaxURL += "&return=" + url.QueryEscape(returnURL)
	}

	// Replace placeholders in HTML with the actual app name (escaped —
	// appName is developer-controlled config, but defence in depth)
	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))

	script := `
	const LOGIN_AJAX_URL = "` + loginAjaxURL + `";
	const VERIFY_AJAX_URL = "` + verifyAjaxURL + `";
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

// checkStores verifies the stores required by the OTP flow are available.
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

// handleOtpSend handles POST /auth/login?action=otp-send-ajax
func (c *loginController) handleOtpSend(w http.ResponseWriter, r *http.Request) {
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

	// Validate email format
	if _, err := mail.ParseAddress(email); err != nil {
		c.sendErrorResponse(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Enforce the email allowlist before generating a code — a blocked email
	// should never receive an OTP, consume send quota, or reach verify.
	if rule := authrules.NewEmailAllowedRule(c.app, email); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	// Per-email send throttle to prevent email-bombing. Checked before the
	// fallible operations; recorded only after the email task is enqueued
	// so failed sends do not consume the quota.
	if !c.canSend(email) {
		c.sendErrorResponse(w, "Too many codes requested, please try again later", http.StatusTooManyRequests)
		return
	}

	// Generate 6-digit OTP
	otp, err := c.generateOTP()
	if err != nil {
		c.sendErrorResponse(w, "Failed to generate OTP", http.StatusInternalServerError)
		c.logError("Failed to generate OTP", err)
		return
	}

	// Generate a random nonce to use as the cache key. The nonce is returned to
	// the client and must be presented alongside the OTP on verification. This
	// prevents brute-forcing OTP values directly against the cache because an
	// attacker would also need to guess the 128-bit nonce.
	nonce, err := c.generateNonce()
	if err != nil {
		c.sendErrorResponse(w, "Failed to generate nonce", http.StatusInternalServerError)
		c.logError("Failed to generate nonce", err)
		return
	}

	// Store OTP under the nonce key. Value: "email:otp" — verified with
	// constant-time compare on retrieval, and read by EmailOTPTask so the
	// plaintext code is never persisted in the task queue.
	cache := c.app.GetMemoryCache()
	cache.Set("otp:"+nonce, email+":"+otp, otpTTL)
	cache.Set("otp_attempts:"+nonce, "0", otpTTL)

	// Enqueue OTP email task for background processing. Only the nonce is
	// passed — the task resolves email+OTP from the memory cache.
	emailOTPTaskHandler := email_otp.NewEmailOTPTask(c.app)
	emailOTPTask, ok := emailOTPTaskHandler.(*email_otp.EmailOTPTask)
	if !ok {
		c.logError("Failed to cast email OTP task handler", nil)
		c.sendErrorResponse(w, "Failed to send OTP email", http.StatusInternalServerError)
		return
	}

	_, err = emailOTPTask.Enqueue(nonce)
	if err != nil {
		c.logError("Failed to enqueue OTP email task", err)
		c.sendErrorResponse(w, "Failed to send OTP email", http.StatusInternalServerError)
		return
	}

	c.recordSend(email)

	c.logInfo(fmt.Sprintf("OTP email task enqueued for %s", email))

	c.sendJSONResponse(w, map[string]interface{}{
		"status":  "success",
		"message": "OTP sent to email",
		"nonce":   nonce,
	}, http.StatusOK)
}

// handleOtpVerify handles POST /auth/login?action=otp-verify-ajax
func (c *loginController) handleOtpVerify(w http.ResponseWriter, r *http.Request) {
	if !c.checkStores(w) {
		return
	}

	if err := r.ParseForm(); err != nil {
		c.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		c.logError("Failed to parse form", err)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")

	if email == "" || otp == "" || nonce == "" {
		c.sendErrorResponse(w, "email, otp, and nonce are required", http.StatusBadRequest)
		return
	}

	cache := c.app.GetMemoryCache()
	cacheKey := "otp:" + nonce
	attemptsKey := "otp_attempts:" + nonce

	// Enforce per-nonce attempt limit before any cache lookup to prevent
	// enumeration. The read-modify-write is serialized via counterMu.
	c.counterMu.Lock()
	attemptsItem := cache.Get(attemptsKey)
	if attemptsItem != nil {
		attempts, _ := attemptsItem.Value().(string)
		count := 0
		fmt.Sscanf(attempts, "%d", &count)
		if count >= otpMaxAttempts {
			c.counterMu.Unlock()
			cache.Delete(cacheKey)
			cache.Delete(attemptsKey)
			c.sendErrorResponse(w, "Too many attempts, please request a new OTP", http.StatusTooManyRequests)
			return
		}
		cache.Set(attemptsKey, fmt.Sprintf("%d", count+1), otpTTL)
	}
	c.counterMu.Unlock()

	// Look up stored "email:otp" value by nonce
	item := cache.Get(cacheKey)
	if item == nil {
		c.sendErrorResponse(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	stored, ok := item.Value().(string)
	if !ok {
		c.sendErrorResponse(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Constant-time compare to prevent timing attacks
	expected := email + ":" + otp
	if subtle.ConstantTimeCompare([]byte(stored), []byte(expected)) != 1 {
		c.sendErrorResponse(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// OTP is valid — remove both keys from cache
	cache.Delete(cacheKey)
	cache.Delete(attemptsKey)

	// The allowlist may have changed since the OTP was sent — re-check here
	// so a newly-blocked email cannot complete login. SessionLogin also
	// enforces this as defense-in-depth for every login method.
	if rule := authrules.NewEmailAllowedRule(c.app, email); rule.Fails() {
		c.sendErrorResponse(w, rule.FailMessageFirst(), http.StatusForbidden)
		return
	}

	// Shared post-auth pipeline: find-or-create user, session, cookie, redirect
	redirectURL, needsRegistration, errorMessage := shared.SessionLogin(c.app, w, r, email, "")
	if errorMessage != "" {
		c.sendErrorResponse(w, errorMessage, http.StatusUnauthorized)
		return
	}

	// Optional return URL: only relative paths are accepted (open-redirect
	// protection), and only when registration is already complete.
	returnURL := r.URL.Query().Get("return")
	if !needsRegistration && returnURL != "" &&
		strings.HasPrefix(returnURL, "/") && !strings.HasPrefix(returnURL, "//") {
		redirectURL = returnURL
	}

	c.sendJSONResponse(w, map[string]interface{}{
		"status":             "success",
		"message":            "Authentication successful",
		"needs_registration": needsRegistration,
		"redirect":           redirectURL,
	}, http.StatusOK)
}

// canSend reports whether another OTP may be sent to the email within the
// current throttle window. The quota is consumed separately via recordSend
// after the email task has been enqueued successfully, so transient failures
// do not eat into the user's allowance.
func (c *loginController) canSend(email string) bool {
	return c.sendCount(email) < otpMaxSends
}

// recordSend increments the per-email send counter.
func (c *loginController) recordSend(email string) {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()
	c.app.GetMemoryCache().Set(c.sendsKey(email), fmt.Sprintf("%d", c.sendCountLocked(email)+1), otpTTL)
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
	return otpSendsKeyPrefix + hex.EncodeToString(sum[:])
}

// generateOTP generates a 6-digit OTP
func (c *loginController) generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
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
