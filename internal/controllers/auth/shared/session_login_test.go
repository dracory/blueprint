package shared

import (
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/auth"
	"github.com/dracory/userstore"
)

// TestCreateSession_LocalEnv_SetsInsecureCookie guards against a regression
// where APP_ENV=local emitted a Secure auth cookie, which browsers silently
// drop over plain HTTP, leaving every request unauthenticated.
func TestCreateSession_LocalEnv_SetsInsecureCookie(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAppEnv("local")

	application := testutils.Setup(
		testutils.WithCfg(cfg),
		testutils.WithSessionStore(true),
	)

	user := userstore.NewUser().SetID("test-user-id")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://example.com/", nil)

	session, err := CreateSession(application, recorder, request, user)

	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}

	cookie := findCookie(recorder, auth.CookieName)

	if cookie == nil {
		t.Fatal("expected auth cookie to be set")
	}
	if cookie.Secure {
		t.Error("expected auth cookie Secure=false under APP_ENV=local")
	}
}

func TestCreateSession_DevelopmentEnv_SetsInsecureCookie(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAppEnv("development")

	application := testutils.Setup(
		testutils.WithCfg(cfg),
		testutils.WithSessionStore(true),
	)

	user := userstore.NewUser().SetID("test-user-id")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://example.com/", nil)

	if _, err := CreateSession(application, recorder, request, user); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	cookie := findCookie(recorder, auth.CookieName)

	if cookie == nil {
		t.Fatal("expected auth cookie to be set")
	}
	if cookie.Secure {
		t.Error("expected auth cookie Secure=false under APP_ENV=development")
	}
}

func TestCreateSession_ProductionEnv_SetsSecureCookie(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAppEnv("production")

	application := testutils.Setup(
		testutils.WithCfg(cfg),
		testutils.WithSessionStore(true),
	)

	user := userstore.NewUser().SetID("test-user-id")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://example.com/", nil)

	if _, err := CreateSession(application, recorder, request, user); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	cookie := findCookie(recorder, auth.CookieName)

	if cookie == nil {
		t.Fatal("expected auth cookie to be set")
	}
	if !cookie.Secure {
		t.Error("expected auth cookie Secure=true under APP_ENV=production")
	}
}

func findCookie(recorder *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
