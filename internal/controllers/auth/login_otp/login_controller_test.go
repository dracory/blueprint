package login_otp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/app"
	"project/internal/tasks/email_otp"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// setupOtpApp builds a test app with all stores the OTP flow requires and
// registers the EmailOTPTask handler so enqueues resolve by alias.
func setupOtpApp(t *testing.T) app.AppInterface {
	t.Helper()
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetTaskStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	err := application.GetTaskStore().TaskHandlerAdd(
		context.Background(), email_otp.NewEmailOTPTask(application), true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	return application
}

func serve(t *testing.T, application app.AppInterface, method, path string,
	query url.Values, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req, err := test.NewRequest(method, path, test.NewRequestOptions{
		QueryParams: query,
		FormValues:  form,
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	controller := NewLoginController(application)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string
		if r.Method == http.MethodPost {
			body = controller.AjaxHandler(w, r)
		} else {
			body = controller.PageHandler(w, r)
		}
		_, _ = w.Write([]byte(body))
	})
	handler.ServeHTTP(recorder, req)
	return recorder
}

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var data map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return data
}

// sendOtp performs a successful otp-send and returns the nonce plus the OTP
// value read directly from the memory cache.
func sendOtp(t *testing.T, application app.AppInterface, email string) (nonce, otp string) {
	t.Helper()
	form := url.Values{}
	form.Set("email", email)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"otp-send-ajax"}}, form)

	if recorder.Code != http.StatusOK {
		t.Fatalf("otp-send expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("otp-send expected success, got %v", data)
	}

	nonce, _ = data["nonce"].(string)
	if nonce == "" {
		t.Fatal("otp-send expected nonce in response")
	}

	item := application.GetMemoryCache().Get("otp:" + nonce)
	if item == nil {
		t.Fatal("expected OTP stored under nonce key in memory cache")
	}

	stored, _ := item.Value().(string)
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected stored OTP value %q", stored)
	}
	return nonce, parts[1]
}

func verifyOtp(t *testing.T, application app.AppInterface, email, otp, nonce string) *httptest.ResponseRecorder {
	t.Helper()
	form := url.Values{}
	form.Set("email", email)
	form.Set("otp", otp)
	form.Set("nonce", nonce)
	return serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"otp-verify-ajax"}}, form)
}

func TestLoginControllerHandler_GetRendersPage(t *testing.T) {
	application := setupOtpApp(t)
	recorder := serve(t, application, http.MethodGet, "/", nil, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "LOGIN_AJAX_URL") {
		t.Fatal("login page should contain LOGIN_AJAX_URL")
	}
	if !strings.Contains(body, "VERIFY_AJAX_URL") {
		t.Fatal("login page should contain VERIFY_AJAX_URL")
	}
	if !strings.Contains(body, "Test app") {
		t.Fatal("login page should contain the app name")
	}
}

func TestLoginControllerHandler_InvalidAction(t *testing.T) {
	application := setupOtpApp(t)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"bogus"}}, nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestOtpSend_MissingEmail(t *testing.T) {
	application := setupOtpApp(t)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"otp-send-ajax"}}, url.Values{})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestOtpSend_InvalidEmail(t *testing.T) {
	application := setupOtpApp(t)
	form := url.Values{}
	form.Set("email", "not-an-email")
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"otp-send-ajax"}}, form)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestOtpSend_Throttle(t *testing.T) {
	application := setupOtpApp(t)
	form := url.Values{}
	form.Set("email", "throttle@test.com")

	// otpMaxSends = 3: the 4th request must be rejected
	var last *httptest.ResponseRecorder
	for i := 0; i < 4; i++ {
		last = serve(t, application, http.MethodPost, "/",
			url.Values{"action": {"otp-send-ajax"}}, form)
	}

	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after %d sends, got %d", otpMaxSends+1, last.Code)
	}
}

func TestOtpVerify_MissingParams(t *testing.T) {
	application := setupOtpApp(t)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"otp-verify-ajax"}}, url.Values{})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestOtpVerify_InvalidOtp(t *testing.T) {
	application := setupOtpApp(t)
	nonce, _ := sendOtp(t, application, "user@test.com")

	recorder := verifyOtp(t, application, "user@test.com", "000000", nonce)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestOtpVerify_UnknownNonce(t *testing.T) {
	application := setupOtpApp(t)
	recorder := verifyOtp(t, application, "user@test.com", "123456", "deadbeef")

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestOtpVerify_TooManyAttempts(t *testing.T) {
	application := setupOtpApp(t)
	nonce, _ := sendOtp(t, application, "user@test.com")

	var last *httptest.ResponseRecorder
	for i := 0; i < otpMaxAttempts; i++ {
		last = verifyOtp(t, application, "user@test.com", "000000", nonce)
	}
	if last.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on attempt %d, got %d", otpMaxAttempts, last.Code)
	}

	last = verifyOtp(t, application, "user@test.com", "000000", nonce)
	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after max attempts, got %d", last.Code)
	}
}

// TestOtpSend_FailureDoesNotConsumeQuota verifies that when the email task
// cannot be enqueued the per-email send quota is not consumed.
func TestOtpSend_FailureDoesNotConsumeQuota(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetTaskStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	// Note: EmailOTPTask handler intentionally NOT registered — Enqueue fails
	controller := NewLoginController(application)

	form := url.Values{}
	form.Set("email", "quota@test.com")

	for i := 0; i < otpMaxSends+2; i++ {
		req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
			QueryParams: url.Values{"action": {"otp-send-ajax"}},
			FormValues:  form,
		})
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		controller.handleOtpSend(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 (enqueue failure), got %d", rec.Code)
		}
	}

	// Register the task so a send can succeed — quota must not be exhausted
	if err := application.GetTaskStore().TaskHandlerAdd(
		context.Background(), email_otp.NewEmailOTPTask(application), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	req, err := test.NewRequest(http.MethodPost, "/", test.NewRequestOptions{
		QueryParams: url.Values{"action": {"otp-send-ajax"}},
		FormValues:  form,
	})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	controller.handleOtpSend(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 after failures (quota must not be consumed), got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestOtpVerify_ReturnURLValidation verifies the open-redirect protection on
// the `return` parameter for a fully-registered user.
// TestLoginPage_ReturnURLCannotBreakJSString verifies that a `return`
// parameter containing a quote is percent-encoded before being embedded
// into the VERIFY_AJAX_URL JavaScript string literal (XSS regression).
func TestLoginPage_ReturnURLCannotBreakJSString(t *testing.T) {
	application := setupOtpApp(t)
	recorder := serve(t, application, http.MethodGet, "/",
		url.Values{"return": {`/";alert(1);//`}}, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if strings.Contains(body, `/";alert(1);//`) {
		t.Fatal("raw return URL leaked into page — JS string break-out possible")
	}
	if !strings.Contains(body, "return=%2F%22") {
		t.Fatal("expected percent-encoded return URL in page")
	}
}

func TestOtpVerify_ReturnURLValidation(t *testing.T) {
	application := setupOtpApp(t)

	// Create a fully-registered user so needs_registration is false
	existingUser := userstore.NewUser().
		SetEmail("registered@test.com").
		SetFirstName("Test").
		SetLastName("User").
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := application.GetUserStore().UserCreate(context.Background(), existingUser); err != nil {
		t.Fatal(err)
	}

	nonce, otp := sendOtp(t, application, "registered@test.com")

	verify := func(returnURL string) *httptest.ResponseRecorder {
		form := url.Values{}
		form.Set("email", "registered@test.com")
		form.Set("otp", otp)
		form.Set("nonce", nonce)
		query := url.Values{"action": {"otp-verify-ajax"}}
		if returnURL != "" {
			query.Set("return", returnURL)
		}
		return serve(t, application, http.MethodPost, "/", query, form)
	}

	// Absolute external URL must be ignored
	rec := verify("https://evil.com")
	data := decodeJSON(t, rec)
	if redirect, _ := data["redirect"].(string); redirect == "https://evil.com" {
		t.Fatal("absolute external return URL must be rejected")
	}
}

func TestOtpVerify_Success(t *testing.T) {
	application := setupOtpApp(t)
	nonce, otp := sendOtp(t, application, "new-user@test.com")

	recorder := verifyOtp(t, application, "new-user@test.com", otp, nonce)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}

	// New user has no names -> registration must be completed
	if data["needs_registration"] != true {
		t.Fatalf("expected needs_registration=true, got %v", data["needs_registration"])
	}

	// Redirect goes through the flash page which forwards to /auth/register
	redirect, _ := data["redirect"].(string)
	if redirect == "" {
		t.Fatal("expected redirect in response")
	}

	// Auth cookie must be set
	foundCookie := false
	for _, c := range recorder.Result().Cookies() {
		if c.Value != "" {
			foundCookie = true
		}
	}
	if !foundCookie {
		t.Fatal("expected auth cookie to be set")
	}

	// OTP keys must be consumed
	if application.GetMemoryCache().Get("otp:"+nonce) != nil {
		t.Fatal("expected OTP cache key to be deleted after successful verify")
	}
}
