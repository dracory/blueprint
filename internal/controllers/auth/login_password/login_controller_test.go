package login_password

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/app"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// setupApp builds a test app with the stores the password login flow
// requires.
func setupApp(t *testing.T) app.AppInterface {
	t.Helper()
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	return testutils.Setup(testutils.WithCfg(cfg))
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

// createUserWithPassword persists an active user with a password and
// returns it.
func createUserWithPassword(t *testing.T, application app.AppInterface, email, password string) userstore.UserInterface {
	t.Helper()
	user := userstore.NewUser().
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetEmail(email).
		SetFirstName("Test").
		SetLastName("User")
	if err := user.SetPasswordAndHash(password); err != nil {
		t.Fatalf("SetPasswordAndHash() expected nil error, got %q", err)
	}
	if err := application.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("UserCreate() expected nil error, got %q", err)
	}
	return user
}

func loginRequest(t *testing.T, application app.AppInterface, email, password string) *httptest.ResponseRecorder {
	t.Helper()
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	return serve(t, application, http.MethodPost, "/",
		url.Values{"action": {"password-login-ajax"}}, form)
}

func TestLoginController_GetRendersPage(t *testing.T) {
	application := setupApp(t)
	recorder := serve(t, application, http.MethodGet, "/", nil, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "LOGIN_AJAX_URL") {
		t.Fatal("login page should contain LOGIN_AJAX_URL")
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

func TestLoginController_MissingCredentials(t *testing.T) {
	application := setupApp(t)
	recorder := loginRequest(t, application, "", "")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestLoginController_InvalidEmail(t *testing.T) {
	application := setupApp(t)
	recorder := loginRequest(t, application, "not-an-email", "password123")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestLoginController_UnknownUserAndWrongPassword_SameMessage(t *testing.T) {
	application := setupApp(t)
	createUserWithPassword(t, application, "user@example.com", "password123")

	unknown := loginRequest(t, application, "ghost@example.com", "password123")
	wrong := loginRequest(t, application, "user@example.com", "wrongpass123")

	if unknown.Code != http.StatusUnauthorized || wrong.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for both, got %d and %d", unknown.Code, wrong.Code)
	}

	// The messages must be identical to prevent email enumeration
	u := decodeJSON(t, unknown)
	w := decodeJSON(t, wrong)
	if u["message"] != w["message"] || u["message"] != "Invalid email or password" {
		t.Fatalf("expected identical enumeration-safe messages, got %v and %v", u["message"], w["message"])
	}
}

func TestLoginController_SuccessfulLogin(t *testing.T) {
	application := setupApp(t)
	createUserWithPassword(t, application, "user@example.com", "password123")

	recorder := loginRequest(t, application, "user@example.com", "password123")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	data := decodeJSON(t, recorder)
	if data["status"] != "success" {
		t.Fatalf("expected success, got %v", data)
	}
	if data["redirect"] == "" {
		t.Fatal("expected a redirect URL")
	}
}
