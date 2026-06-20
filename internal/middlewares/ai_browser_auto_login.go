package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"project/internal/app"
	"project/internal/config"

	"github.com/dracory/auth"
	"github.com/dracory/rtr"
	"github.com/dracory/sessionstore"
	"github.com/dromara/carbon/v2"
)

const aiBrowserUserID = "ai-browser-user"

// AiBrowserAutoLoginMiddleware automatically authenticates every request as the
// ai-browser dev user. It checks for an existing auth cookie; if none is present
// it creates a new 24-hour session and sets the cookie so that AuthMiddleware
// (which runs after this middleware) can pick it up normally.
//
// IMPORTANT: This middleware must NEVER be registered in the main server's
// middleware chain. It is only for use with routes.AiBrowserRouter.
func AiBrowserAutoLoginMiddleware(a app.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("AI Browser Auto-Login Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return aiBrowserAutoLoginHandler(a, next)
		})
}

func aiBrowserAutoLoginHandler(a app.AppInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for an existing auth cookie directly — do NOT use helpers.GetAuthUser(r)
		// because GetAuthUser reads from the request context, which is populated by
		// AuthMiddleware. Since this middleware runs before AuthMiddleware, the context
		// value will always be nil here.
		cookie, err := r.Cookie(auth.CookieName)
		if err == nil && cookie != nil && cookie.Value != "" {
			// A cookie exists — let AuthMiddleware validate it normally.
			next.ServeHTTP(w, r)
			return
		}

		if a.GetUserStore() == nil || a.GetSessionStore() == nil {
			slog.Warn("AiBrowserAutoLoginMiddleware: user or session store not initialized, skipping auto-login")
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.GetUserStore().UserFindByID(r.Context(), aiBrowserUserID)
		if err != nil {
			slog.Error("AiBrowserAutoLoginMiddleware: failed to find dev user", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		if user == nil {
			slog.Warn("AiBrowserAutoLoginMiddleware: dev user not found, skipping auto-login", "userID", aiBrowserUserID)
			next.ServeHTTP(w, r)
			return
		}

		session := sessionstore.NewSession().
			SetUserID(user.GetID()).
			SetUserAgent(r.UserAgent()).
			SetIPAddress(r.RemoteAddr).
			SetExpiresAt(carbon.Now(carbon.UTC).AddHours(24).ToDateTimeString(carbon.UTC))

		if err := a.GetSessionStore().SessionCreate(r.Context(), session); err != nil {
			slog.Error("AiBrowserAutoLoginMiddleware: failed to create session", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     auth.CookieName,
			Value:    session.GetKey(),
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		// Inject the session key into the request cookie header so that
		// AuthMiddleware (which also reads the cookie) can find it on this
		// same request without requiring a redirect.
		r.AddCookie(&http.Cookie{
			Name:  auth.CookieName,
			Value: session.GetKey(),
		})

		ctx := context.WithValue(r.Context(), config.AuthenticatedUserContextKey{}, user)
		ctx = context.WithValue(ctx, config.AuthenticatedSessionContextKey{}, session.GetKey())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
