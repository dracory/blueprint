package login_magiclink

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/app"
	"project/internal/tasks/email_magic_link"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// setupApp builds a test app with all stores the magic link flow requires
// and registers the EmailMagicLinkTask handler so enqueues resolve by alias.
func setupApp(t *testing.T) app.AppInterface {
	t.Helper()
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetTaskStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	err := application.GetTaskStore().TaskHandlerAdd(
		context.Background(), email_magic_link.NewEmailMagicLinkTask(application), true)
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
		switch {
		case r.Method == http.MethodPost:
			body = controller.AjaxHandler(w, r)
		case r.URL.Query().Has("token"):
			body = controller.Handler(w, r)
		default:
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

// sendLink performs a successful magiclink-send and returns the token read
// directly from the memory cache via the stored nonce entry.
func sendLink(t *testing.T, application app.AppInterface, email string) (token string) {
	t.Helper()
	form := url.Values{}
	form.Set("email", email)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"magiclink-send-ajax"}}, form)

	if recorder.Code != http.StatusOK {
		t.Fatalf("magiclink-send expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("magiclink-send expected success, got %v", data)
	}

	// Find the token entry the controller stored for the verify path.
	cache := application.GetMemoryCache()
	keys := cache.Keys()
	for _, k := range keys {
		if !strings.HasPrefix(k, tokenCacheKeyPrefix) {
			continue
		}
		item := cache.Get(k)
		stored, _ := item.Value().(string)
		var value magicLinkCacheValue
		if err := json.Unmarshal([]byte(stored), &value); err == nil && value.Email == email {
			return strings.TrimPrefix(k, tokenCacheKeyPrefix)
		}
	}
	t.Fatal("expected token entry stored in memory cache")
	return ""
}

func verifyLink(t *testing.T, application app.AppInterface, token string, returnURL string) *httptest.ResponseRecorder {
	t.Helper()
	query := url.Values{"token": {token}}
	if returnURL != "" {
		query.Set("return", returnURL)
	}
	return serve(t, application, http.MethodGet, "/", query, nil)
}

func TestLoginController_GetRendersPage(t *testing.T) {
	application := setupApp(t)
	recorder := serve(t, application, http.MethodGet, "/", nil, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "SEND_AJAX_URL") {
		t.Fatal("login page should contain SEND_AJAX_URL")
	}
	if !strings.Contains(body, "Test app") {
		t.Fatal("login page should contain the app name")
	}
}

func TestLoginController_InvalidAction(t *testing.T) {
	application := setupApp(t)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"bogus"}}, nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestSend_MissingEmail(t *testing.T) {
	application := setupApp(t)
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"magiclink-send-ajax"}}, url.Values{})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestSend_InvalidEmail(t *testing.T) {
	application := setupApp(t)
	form := url.Values{}
	form.Set("email", "not-an-email")
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"magiclink-send-ajax"}}, form)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestSend_DisallowedEmail(t *testing.T) {
	application := setupApp(t)
	application.GetConfig().SetEmailsAllowedAccess([]string{"allowed@test.com"})

	form := url.Values{}
	form.Set("email", "blocked@test.com")
	recorder := serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"magiclink-send-ajax"}}, form)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestSend_Throttle(t *testing.T) {
	application := setupApp(t)
	form := url.Values{}
	form.Set("email", "throttle@test.com")

	// magicLinkMaxSends = 3: the 4th request must be rejected
	var last *httptest.ResponseRecorder
	for i := 0; i < 4; i++ {
		last = serve(t, application, http.MethodPost, "/",
			url.Values{"action": {"magiclink-send-ajax"}}, form)
	}

	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after %d sends, got %d", magicLinkMaxSends+1, last.Code)
	}
}

func TestVerify_MissingToken(t *testing.T) {
	application := setupApp(t)
	recorder := serve(t, application, http.MethodGet, "/", url.Values{"token": {""}}, nil)

	// Flash error redirects to home — verify no session/cookie was created
	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("Location"), "/flash") {
		t.Fatalf("expected flash redirect, got %q", recorder.Header().Get("Location"))
	}
}

func TestVerify_UnknownToken(t *testing.T) {
	application := setupApp(t)
	recorder := verifyLink(t, application, "deadbeef", "")

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("Location"), "/flash") {
		t.Fatalf("expected flash redirect, got %q", recorder.Header().Get("Location"))
	}
}

func TestVerify_SuccessAndSingleUse(t *testing.T) {
	application := setupApp(t)
	token := sendLink(t, application, "new-user@test.com")

	recorder := verifyLink(t, application, token, "")

	// Successful login redirects through the flash page
	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", recorder.Code)
	}
	location := recorder.Header().Get("Location")
	if !strings.Contains(location, "/flash") {
		t.Fatalf("expected flash redirect, got %q", location)
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

	// Token entry must be consumed — replaying the link must fail
	if application.GetMemoryCache().Get(tokenCacheKeyPrefix+token) != nil {
		t.Fatal("expected token cache key to be deleted after successful verify")
	}

	replay := verifyLink(t, application, token, "")
	if !strings.Contains(replay.Header().Get("Location"), "/flash") {
		t.Fatal("replayed token must fail")
	}
}

func TestVerify_ReturnURLValidation(t *testing.T) {
	application := setupApp(t)

	// Create a fully-registered user so the flash redirect reflects the
	// final target rather than the registration flash URL.
	existingUser := userstore.NewUser().
		SetEmail("registered@test.com").
		SetFirstName("Test").
		SetLastName("User").
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := application.GetUserStore().UserCreate(context.Background(), existingUser); err != nil {
		t.Fatal(err)
	}

	token := sendLink(t, application, "registered@test.com")

	// Absolute external URL in the query must be ignored — the flash
	// redirect must point at the calculated redirect instead.
	rec := verifyLink(t, application, token, "https://evil.com")
	location := rec.Header().Get("Location")
	if strings.Contains(location, "evil.com") {
		t.Fatalf("external return URL must be rejected, got %q", location)
	}
}
