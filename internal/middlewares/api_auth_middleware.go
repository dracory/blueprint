package middlewares

import (
	"context"
	"net/http"
	"time"

	"project/internal/config"
	"project/internal/types"

	"github.com/dracory/api"
	"github.com/dracory/rtr"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

type apiAuthCacheItem struct {
	session sessionstore.SessionInterface
	user    userstore.UserInterface
}

func NewAPIAuthMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
	memoryCache := app.GetMemoryCache()

	return rtr.NewMiddleware().
		SetName("API Auth Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 1. Get token from Authorization header
				token := r.Header.Get("Authorization")
				if token == "" {
					if _, err := w.Write([]byte(api.Error("Authorization token required").ToString())); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					return
				}

				cacheKey := "api-auth:" + token
				if item := memoryCache.Get(cacheKey); item != nil {
					if cacheData, ok := item.Value().(apiAuthCacheItem); ok {
						if !cacheData.session.IsExpired() {
							ctx := context.WithValue(r.Context(), config.APIAuthenticatedSessionContextKey{}, cacheData.session)
							ctx = context.WithValue(ctx, config.APIAuthenticatedUserContextKey{}, cacheData.user)
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}
				}

				// 2. Validate session token
				session, err := app.GetSessionStore().SessionFindByKey(context.Background(), token)
				if err != nil {
					if _, writeErr := w.Write([]byte(api.Error("Failed to validate token").ToString())); writeErr != nil {
						return
					}
					return
				}

				if session == nil || session.IsExpired() {
					if _, err := w.Write([]byte(api.Error("Invalid or expired token").ToString())); err != nil {
						return
					}
					return
				}

				// 3. Load user and add session + user to request context
				userID := session.GetUserID()
				if userID == "" {
					if _, err := w.Write([]byte(api.Error("Session missing user").ToString())); err != nil {
						return
					}
					return
				}

				user, err := app.GetUserStore().UserFindByID(r.Context(), userID)
				if err != nil {
					if _, writeErr := w.Write([]byte(api.Error("Failed to load user").ToString())); writeErr != nil {
						return
					}
					return
				}

				if user == nil {
					if _, err := w.Write([]byte(api.Error("User not found").ToString())); err != nil {
						return
					}
					return
				}

				// Cache the data for 1 minute
				cacheData := apiAuthCacheItem{session: session, user: user}
				memoryCache.Set(cacheKey, cacheData, 1*time.Minute)

				ctx := context.WithValue(r.Context(), config.APIAuthenticatedSessionContextKey{}, session)
				ctx = context.WithValue(ctx, config.APIAuthenticatedUserContextKey{}, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
}
