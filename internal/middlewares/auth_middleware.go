package middlewares

import (
	"context"
	"log"
	"net/http"
	"project/internal/config"
	"project/internal/types"

	"github.com/dracory/auth"
	"github.com/dracory/rtr"
)

func AuthMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Auth Middleware").
		SetHandler(func(next http.Handler) http.Handler { return authHandler(app, next) })
}

// authHandler adds the user and session to the context.
//
//  1. Checks if the user session key exists in the incoming request.
//
//  2. Retrieves the session using the session key..
//
//  3. Checks the session is not expired.
//
//  4. Retrieves the user using the user ID from the session.
//
// Params:
//   - next http.Handler. The `next` handler is the next handler in the middleware chain.
//
// Returns
// - an http.Handler which represents the modified handler with the user.
func authHandler(app types.AppInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.GetConfig().GetSessionStoreUsed() {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("session store not enabled"))
			return
		}

		if app.GetSessionStore() == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("session store not initialized"))
			return
		}

		if !app.GetConfig().GetUserStoreUsed() {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("user store not enabled"))
			return
		}

		if app.GetUserStore() == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("user store not initialized"))
			return
		}

		sessionKey := authHandlerSessionKey(r)

		if sessionKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		session, err := app.GetSessionStore().SessionFindByKey(r.Context(), sessionKey)

		if err != nil {
			app.GetLogger().Error("auth_middleware", "error", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		if session == nil {
			next.ServeHTTP(w, r)
			return
		}

		if session.IsExpired() {
			next.ServeHTTP(w, r)
			return
		}

		userID := session.GetUserID()

		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.GetUserStore().UserFindByID(r.Context(), userID)

		if err != nil {
			app.GetLogger().Error("auth_middleware", "error", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		if user == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), config.AuthenticatedUserContextKey{}, user)
		ctx = context.WithValue(ctx, config.AuthenticatedSessionContextKey{}, sessionKey)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// authHandlerSessionKey returns the session key from the incoming request.
func authHandlerSessionKey(r *http.Request) string {
	authTokenFromCookie, err := r.Cookie(auth.CookieName)

	if err != nil {
		if err != http.ErrNoCookie {
			log.Println(err.Error())
		}
	}

	if authTokenFromCookie == nil {
		return ""
	}

	sessionKey := authTokenFromCookie.Value

	return sessionKey
}
