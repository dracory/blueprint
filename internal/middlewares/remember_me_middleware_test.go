package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/app"
	"project/internal/config"
	"project/internal/controllers/auth/shared"
	"project/internal/testutils"
	"testing"
	"time"

	"github.com/dracory/auth"
	"github.com/dracory/sessionstore"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// setupRememberMeApp returns an app with user+session stores and remember-me
// configured.
func setupRememberMeApp(t *testing.T, enabled bool) app.AppInterface {
	t.Helper()
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetRememberMeEnabled(enabled)
	cfg.SetRememberMeDays(30)
	return testutils.Setup(testutils.WithCfg(cfg))
}

// seedRememberSession seeds a session marked as a remember session
// (session_value = "remember") — required for the middleware to accept it.
func seedRememberSession(t *testing.T, application app.AppInterface, user userstore.UserInterface, expiresSeconds int) sessionstore.SessionInterface {
	t.Helper()
	seedReq := httptest.NewRequest("GET", "/", nil)
	session, err := testutils.SeedSession(application.GetSessionStore(), seedReq, user, expiresSeconds)
	if err != nil {
		t.Fatal(err)
	}
	session.SetValue(shared.RememberSessionValue)
	if err := application.GetSessionStore().SessionUpdate(context.Background(), session); err != nil {
		t.Fatal(err)
	}
	return session
}

func rememberMeRequest(rememberKey string) *http.Request {
	req := httptest.NewRequest("GET", "/", nil)
	if rememberKey != "" {
		req.AddCookie(&http.Cookie{Name: config.COOKIE_NAME_REMEMBER_TOKEN, Value: rememberKey})
	}
	return req
}

func runMiddleware(application app.AppInterface, req *http.Request, rec *httptest.ResponseRecorder, next http.HandlerFunc) {
	RememberMeMiddleware(application).GetHandler()(next).ServeHTTP(rec, req)
}

func TestRememberMeMiddleware_Disabled(t *testing.T) {
	app := setupRememberMeApp(t, false)

	req := rememberMeRequest("some_key")
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("no cookies should be set when remember-me is disabled")
	}
}

func TestRememberMeMiddleware_NoCookie(t *testing.T) {
	app := setupRememberMeApp(t, true)

	req := rememberMeRequest("")
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("auth cookie should not be injected without a remember cookie")
		}
		w.WriteHeader(http.StatusOK)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestRememberMeMiddleware_UnknownToken(t *testing.T) {
	app := setupRememberMeApp(t, true)

	req := rememberMeRequest("nonexistent_session_key")
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("auth cookie should not be injected for an unknown remember token")
		}
		w.WriteHeader(http.StatusOK)
	})

	// Stale remember cookie should be cleared
	cleared := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == config.COOKIE_NAME_REMEMBER_TOKEN && c.MaxAge < 0 {
			cleared = true
		}
	}
	if !cleared {
		t.Fatal("unknown remember token should expire the remember cookie")
	}
}

func TestRememberMeMiddleware_UnmarkedSessionRejected(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	// A plain (unmarked) session key planted in the remember cookie must not
	// be honoured — otherwise any stolen auth session could mint persistent
	// sessions.
	seedReq := httptest.NewRequest("GET", "/", nil)
	plainSession, err := testutils.SeedSession(app.GetSessionStore(), seedReq, user, 30*24*3600)
	if err != nil {
		t.Fatal(err)
	}

	req := rememberMeRequest(plainSession.GetKey())
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("unmarked session in the remember cookie must not restore auth")
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestRememberMeMiddleware_RestoresAndRotates(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	rememberSession := seedRememberSession(t, app, user, 30*24*3600)

	req := rememberMeRequest(rememberSession.GetKey())
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		// The new auth session key must be visible to downstream middleware
		// on this same request.
		c, err := r.Cookie(auth.CookieName)
		if err != nil || c.Value == "" {
			t.Fatal("auth session key should be injected into the request")
		}
		w.WriteHeader(http.StatusOK)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Response must carry a fresh auth cookie and a rotated remember cookie.
	var authSet, rememberSet bool
	for _, c := range rec.Result().Cookies() {
		switch c.Name {
		case auth.CookieName:
			authSet = c.Value != ""
		case config.COOKIE_NAME_REMEMBER_TOKEN:
			rememberSet = c.Value != "" && c.Value != rememberSession.GetKey()
		}
	}
	if !authSet {
		t.Fatal("a new auth cookie should be set")
	}
	if !rememberSet {
		t.Fatal("the remember cookie should be rotated to a new session key")
	}

	// The used remember session is not deleted — it is shortened to the
	// grace window so concurrent tab restores still succeed.
	old, err := app.GetSessionStore().SessionFindByID(req.Context(), rememberSession.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if old == nil {
		t.Fatal("used remember session should remain during the grace window")
	}
	// It must be shortened to the grace window — far below the original
	// 30 days.
	expiresAt, err := time.Parse("2006-01-02 15:04:05", old.GetExpiresAt())
	if err != nil {
		t.Fatalf("cannot parse session expiry %q: %v", old.GetExpiresAt(), err)
	}
	if time.Until(expiresAt) > 2*time.Minute {
		t.Fatalf("used remember session should expire within the grace window, expires in %s", time.Until(expiresAt))
	}
}

func TestRememberMeMiddleware_RememberKeyInAuthCookieStripped(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	rememberSession := seedRememberSession(t, app, user, 30*24*3600)

	// Attacker/user plants the remember key in the AUTH cookie — it must not
	// act as a 30-day auth credential.
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: rememberSession.GetKey()})
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("a remember session key must be stripped from the auth cookie")
		}
		w.WriteHeader(http.StatusOK)
	})

	expired := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == auth.CookieName && c.MaxAge < 0 {
			expired = true
		}
	}
	if !expired {
		t.Fatal("the misused auth cookie should be expired in the response")
	}
}

func TestRememberMeMiddleware_DisabledStillStripsRememberKeyInAuthCookie(t *testing.T) {
	app := setupRememberMeApp(t, false)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	// A remember session issued before the feature was switched off must
	// still not act as an auth credential.
	rememberSession := seedRememberSession(t, app, user, 30*24*3600)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: rememberSession.GetKey()})
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("a remember session key must be stripped even when remember-me is disabled")
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestRememberMeMiddleware_ReplayDoesNotExtendGrace(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	// Already inside the grace window (used once, 10s left).
	rememberSession := seedRememberSession(t, app, user, 10)

	req := rememberMeRequest(rememberSession.GetKey())
	rec := httptest.NewRecorder()
	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	old, err := app.GetSessionStore().SessionFindByID(req.Context(), rememberSession.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if old == nil {
		t.Fatal("used remember session should still exist inside the grace window")
	}
	expiresAt, err := time.Parse("2006-01-02 15:04:05", old.GetExpiresAt())
	if err != nil {
		t.Fatalf("cannot parse session expiry %q: %v", old.GetExpiresAt(), err)
	}
	if time.Until(expiresAt) > 20*time.Second {
		t.Fatalf("replay must not extend the grace window, expires in %s", time.Until(expiresAt))
	}
}

func TestRememberMeMiddleware_ExistingValidSessionSkips(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}

	seedReq := httptest.NewRequest("GET", "/", nil)
	authSession, err := testutils.SeedSession(app.GetSessionStore(), seedReq, user, 3600)
	if err != nil {
		t.Fatal(err)
	}
	rememberSession := seedRememberSession(t, app, user, 30*24*3600)

	req := rememberMeRequest(rememberSession.GetKey())
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: authSession.GetKey()})
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Valid auth session → middleware is a no-op: no cookies re-set.
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("no cookies should be re-set when the auth session is already valid")
	}
}

func TestRememberMeMiddleware_InactiveUserRevokes(t *testing.T) {
	app := setupRememberMeApp(t, true)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatal(err)
	}
	user.SetStatus(userstore.USER_STATUS_INACTIVE)
	if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	rememberSession := seedRememberSession(t, app, user, 30*24*3600)

	req := rememberMeRequest(rememberSession.GetKey())
	rec := httptest.NewRecorder()

	runMiddleware(app, req, rec, func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(auth.CookieName); err == nil {
			t.Fatal("auth cookie should not be injected for an inactive user")
		}
		w.WriteHeader(http.StatusOK)
	})

	old, _ := app.GetSessionStore().SessionFindByID(req.Context(), rememberSession.GetID())
	if old != nil {
		t.Fatal("remember session of an inactive user should be revoked")
	}
}
