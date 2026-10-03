package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"project/internal/app"
	"project/internal/config"
	"project/internal/controllers/auth/shared"

	"github.com/dracory/auth"
	"github.com/dracory/rtr"
	"github.com/dracory/sessionstore"
	"github.com/dromara/carbon/v2"
)

// rememberRotationGrace is how long a used remember session stays valid
// after rotation. Browser restores fire many concurrent requests carrying
// the same remember cookie; hard-deleting on first use would make the
// siblings fail (and expiring their cookie could clobber the rotated one).
// The short window absorbs concurrent replays; it is never extended, so a
// used token dies at most rememberRotationGrace after its first use.
const rememberRotationGrace = 60 * time.Second

// RememberMeMiddleware restores a logged-in state for requests that carry a
// valid remember cookie but no valid auth session. It runs before
// AuthMiddleware (which has no fallback hook): when it exchanges a remember
// session for a fresh auth session it rewrites the request cookie header so
// AuthMiddleware picks up the new session on the same request.
//
// Remember sessions are normal sessionstore rows marked with
// session_value = "remember" (see shared.IssueRememberSession). Two
// consequences:
//   - The middleware only accepts marked sessions, so a regular auth
//     session key planted in the remember cookie is ignored.
//   - A remember key planted in the AUTH cookie is stripped and expired —
//     AuthMiddleware accepts any valid session key, so an unmarked check
//     would let a copied remember token act as a 30-day auth credential
//     that never rotates.
//
// On every successful restore the used remember session is shortened to
// rememberRotationGrace (not deleted) and a fresh one is issued — the grace
// window makes concurrent-tab restores safe.
func RememberMeMiddleware(a app.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Remember Me Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return rememberMeHandler(a, next)
		})
}

func rememberMeHandler(a app.AppInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.IsDisabledSessionStore() {
			next.ServeHTTP(w, r)
			return
		}

		// The misuse check runs even when remember-me is disabled: remember
		// sessions issued before the feature was switched off stay in the
		// store and must still not act as auth credentials.
		switch checkAuthSession(a, r) {
		case authSessionValid:
			next.ServeHTTP(w, r)
			return
		case authSessionIsRemember:
			// A remember key used as an auth credential: strip it from the
			// request and expire the cookie, then let the normal remember
			// flow below do a proper restore + rotation.
			removeRequestCookie(r, auth.CookieName)
			expireCookie(w, r, auth.CookieName, a)
		}

		if a.GetConfig() == nil || !a.GetConfig().GetRememberMeEnabled() || a.IsDisabledUserStore() {
			next.ServeHTTP(w, r)
			return
		}

		rememberKey := rememberCookieValue(r)
		if rememberKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		rememberSession, err := a.GetSessionStore().SessionFindByKey(r.Context(), rememberKey)
		if err != nil {
			logRememberMeError(a, "session lookup error", err)
			next.ServeHTTP(w, r)
			return
		}

		// Unknown, expired, or unmarked (not a real remember token) — clear
		// the stale cookie and continue anonymously.
		if rememberSession == nil || rememberSession.IsExpired() ||
			rememberSession.GetValue() != shared.RememberSessionValue {
			shared.RemoveRememberCookie(a, w, r)
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.GetUserStore().UserFindByID(r.Context(), rememberSession.GetUserID())
		if err != nil {
			logRememberMeError(a, "user lookup error", err)
			next.ServeHTTP(w, r)
			return
		}

		if user == nil || !user.IsActive() {
			revokeRememberSession(a, w, r, rememberSession)
			next.ServeHTTP(w, r)
			return
		}

		session, err := shared.CreateSession(a, w, r, user)
		if err != nil {
			logRememberMeError(a, "failed to create session", err)
			next.ServeHTTP(w, r)
			return
		}

		// Rewrite the request cookie header so AuthMiddleware — which reads
		// the cookie — finds the new session on this same request. AddCookie
		// would leave the stale client cookie first in line (r.Cookie
		// returns the first match), so the header is rebuilt instead.
		setRequestCookie(r, auth.CookieName, session.GetKey())

		// Rotate: shorten the used remember session to the grace window and
		// issue a fresh one. A replay outside the window fails lookup. Only
		// ever shorten — re-extending on each replay would let a token
		// replayed within every window live forever.
		graceEnd := carbon.Now(carbon.UTC).AddSeconds(int(rememberRotationGrace / time.Second))
		if carbon.Parse(rememberSession.GetExpiresAt(), carbon.UTC).Gt(graceEnd) {
			rememberSession.SetExpiresAt(graceEnd.ToDateTimeString(carbon.UTC))
			if err := a.GetSessionStore().SessionUpdate(r.Context(), rememberSession); err != nil {
				logRememberMeError(a, "failed to rotate remember session", err)
			}
		}
		if err := shared.IssueRememberSession(a, w, r, user); err != nil {
			logRememberMeError(a, "failed to issue remember session", err)
		}

		next.ServeHTTP(w, r)
	})
}

type authSessionState int

const (
	authSessionNone authSessionState = iota
	authSessionValid
	authSessionIsRemember
)

// checkAuthSession classifies the request's auth cookie: valid normal
// session, a misused remember session, or absent/invalid.
func checkAuthSession(a app.AppInterface, r *http.Request) authSessionState {
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil || cookie.Value == "" {
		return authSessionNone
	}

	session, err := a.GetSessionStore().SessionFindByKey(r.Context(), cookie.Value)
	if err != nil || session == nil || session.IsExpired() {
		return authSessionNone
	}

	if session.GetValue() == shared.RememberSessionValue {
		return authSessionIsRemember
	}

	return authSessionValid
}

// rememberCookieValue extracts the remember session key from the request.
func rememberCookieValue(r *http.Request) string {
	cookie, err := r.Cookie(config.COOKIE_NAME_REMEMBER_TOKEN)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}

// revokeRememberSession deletes the remember session server-side and
// expires the cookie — used when the owning user is gone or deactivated.
func revokeRememberSession(a app.AppInterface, w http.ResponseWriter, r *http.Request, rememberSession sessionstore.SessionInterface) {
	if err := a.GetSessionStore().SessionDelete(r.Context(), rememberSession); err != nil {
		logRememberMeError(a, "failed to revoke remember session", err)
	}
	shared.RemoveRememberCookie(a, w, r)
}

// setRequestCookie rebuilds the request Cookie header with name=value,
// replacing any existing entry of the same name — r.Cookie resolves the
// first match, so appending would leave a stale duplicate in place.
func setRequestCookie(r *http.Request, name, value string) {
	removeRequestCookie(r, name)
	r.AddCookie(&http.Cookie{Name: name, Value: value}) // #nosec G124 -- in-request cookie injection, not sent to browser
}

// removeRequestCookie rebuilds the request Cookie header without the named
// cookie.
func removeRequestCookie(r *http.Request, name string) {
	var kept []string
	for _, c := range r.Cookies() {
		if c.Name != name {
			kept = append(kept, c.Name+"="+c.Value)
		}
	}
	r.Header.Set("Cookie", strings.Join(kept, "; "))
}

// expireCookie emits an expiring Set-Cookie for name, matching the Secure
// convention used for the auth cookie (plain HTTP in development).
func expireCookie(w http.ResponseWriter, r *http.Request, name string, a app.AppInterface) {
	secure := true
	if a.GetConfig() != nil && a.GetConfig().IsEnvDevelopment() {
		secure = false
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "none",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
}

// logRememberMeError logs an error with the middleware prefix when a logger
// is configured.
func logRememberMeError(a app.AppInterface, message string, err error) {
	if logger := a.GetLogger(); logger != nil {
		logger.Error("RememberMeMiddleware: "+message, slog.String("error", err.Error()))
	}
}
