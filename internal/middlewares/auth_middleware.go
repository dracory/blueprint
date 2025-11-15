package middlewares

import (
	"context"
	"log"
	"net/http"
	"project/internal/config"
	"project/internal/types"
	"time"

	"github.com/dracory/auth"
	"github.com/dracory/rtr"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	"github.com/jellydator/ttlcache/v3"
)

const (
	sessionCacheTTL    = 5 * time.Minute
	userCacheTTL       = 5 * time.Minute
	sessionCachePrefix = "auth:session:"
	userCachePrefix    = "auth:user:"
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
			if _, err := w.Write([]byte("session store not enabled")); err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
			}
			return
		}

		if app.GetSessionStore() == nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("session store not initialized")); err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
			}
			return
		}

		if !app.GetConfig().GetUserStoreUsed() {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("user store not enabled")); err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
			}
			return
		}

		if app.GetUserStore() == nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("user store not initialized")); err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
			}
			return
		}

		sessionKey := authHandlerSessionKey(r)

		if sessionKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		sessionStore := app.GetSessionStore()
		memoryCache := app.GetMemoryCache()

		session := cacheGetSession(memoryCache, sessionKey)

		if session == nil {
			var err error
			session, err = sessionStore.SessionFindByKey(r.Context(), sessionKey)

			if err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
				next.ServeHTTP(w, r)
				return
			}

			if session == nil {
				next.ServeHTTP(w, r)
				return
			}

			cacheSetSession(memoryCache, sessionKey, session)
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

		user := cacheGetUser(memoryCache, userID)

		if user == nil {
			fetchedUser, err := app.GetUserStore().UserFindByID(r.Context(), userID)

			if err != nil {
				app.GetLogger().Error("auth_middleware", "error", err.Error())
				next.ServeHTTP(w, r)
				return
			}

			if fetchedUser == nil {
				next.ServeHTTP(w, r)
				return
			}

			cacheSetUser(memoryCache, userID, fetchedUser)
			user = fetchedUser
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

func cacheGetSession(cache *ttlcache.Cache[string, any], sessionKey string) sessionstore.SessionInterface {
	if cache == nil {
		return nil
	}

	item := cache.Get(sessionCachePrefix + sessionKey)

	if item == nil {
		return nil
	}

	session, ok := item.Value().(sessionstore.SessionInterface)

	if !ok || session == nil {
		return nil
	}

	if session.IsExpired() {
		return nil
	}

	return session
}

func cacheSetSession(cache *ttlcache.Cache[string, any], sessionKey string, session sessionstore.SessionInterface) {
	if cache == nil || session == nil {
		return
	}

	cache.Set(sessionCachePrefix+sessionKey, session, sessionCacheTTL)
}

func cacheGetUser(cache *ttlcache.Cache[string, any], userID string) userstore.UserInterface {
	if cache == nil {
		return nil
	}

	item := cache.Get(userCachePrefix + userID)

	if item == nil {
		return nil
	}

	user, ok := item.Value().(userstore.UserInterface)

	if !ok || user == nil {
		return nil
	}

	return user
}

func cacheSetUser(cache *ttlcache.Cache[string, any], userID string, user userstore.UserInterface) {
	if cache == nil || user == nil {
		return
	}

	cache.Set(userCachePrefix+userID, user, userCacheTTL)
}
