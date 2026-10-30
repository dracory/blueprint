package middlewares

import (
	"context"
	"net/http"

	"project/internal/app"
	"project/internal/authadapters"
	"project/internal/config"

	"github.com/dracory/auth"
	authmiddlewares "github.com/dracory/auth/middlewares"
	"github.com/dracory/auth/types"
	"github.com/dracory/rtr"
)

// RememberMeMiddleware restores a logged-in state for requests that carry a
// valid remember cookie but no valid auth session. The exchange logic lives
// in the dracory/auth package — this wrapper wires the app's stores and
// config into it. It runs before AuthMiddleware (which has no fallback
// hook): on exchange it rewrites the request cookie header so
// AuthMiddleware picks up the minted session on the same request.
//
// Remember sessions are normal sessionstore rows marked with
// session_value = "remember" (see authadapters.IssueRememberSession). Two
// consequences:
//   - The middleware only accepts marked sessions, so a regular auth
//     session key planted in the remember cookie is ignored.
//   - A remember key planted in the AUTH cookie is stripped and expired —
//     AuthMiddleware accepts any valid session key, so an unmarked check
//     would let a copied remember token act as a 30-day auth credential
//     that never rotates.
//
// On every successful restore the used remember session is shortened to a
// grace window (not deleted) and a fresh one is issued — the grace window
// makes concurrent-tab restores safe.
func RememberMeMiddleware(a app.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Remember Me Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			if a.GetConfig() == nil || a.IsDisabledSessionStore() {
				return next
			}
			sessions := &authadapters.SessionStoreAdapter{Store: a.GetSessionStore()}
			cfg := a.GetConfig()
			secure := !(cfg.IsEnvDevelopment() || cfg.IsEnvLocal())
			return authmiddlewares.RememberMeMiddleware(authmiddlewares.RememberMeConfig{
				Enabled:              cfg.GetRememberMeEnabled(),
				SessionStore:         sessions,
				AuthCookieName:       auth.CookieName,
				RememberCookieName:   config.COOKIE_NAME_REMEMBER_TOKEN,
				IssueRememberSession: authadapters.IssueRememberSession(a, sessions),
				IssueAuthSession:     authadapters.IssueAuthSession(a, sessions),
				UserIsActive:         userIsActive(a),
				RememberMeDays:       cfg.GetRememberMeDays(),
				SecureCookies:        secure,
				AuthCookieSet: func(w http.ResponseWriter, r *http.Request, token string) {
					opts := []types.CookieOption{}
					if !secure {
						opts = append(opts, types.WithSecure(false))
					}
					auth.AuthCookieSet(w, r, token, opts...)
				},
				Logger: a.GetLogger(),
			}, next)
		})
}

// userIsActive returns the revocation check for the middleware, or nil
// (allow all) when the user store is disabled — without it a disabled
// store would fail every lookup and wipe remember sessions.
func userIsActive(a app.AppInterface) func(ctx context.Context, userID string) bool {
	if a.IsDisabledUserStore() {
		return nil
	}
	return func(ctx context.Context, userID string) bool {
		return authadapters.UserIsActiveHook(a)(ctx, userID, types.UserAuthOptions{})
	}
}
