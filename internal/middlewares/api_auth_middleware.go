package middlewares

import (
	"context"
	"net/http"

	"project/internal/config"
	"project/internal/types"

	"github.com/dracory/api"
	"github.com/dracory/rtr"
)

func NewAPIAuthMiddleware(app types.AppInterface) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("API Auth Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 1. Get token from Authorization header
				token := r.Header.Get("Authorization")
				if token == "" {
					w.Write([]byte(api.Error("Authorization token required").ToString()))
					return
				}

				// 2. Validate session token
				session, err := app.GetSessionStore().SessionFindByKey(context.Background(), token)
				if err != nil {
					w.Write([]byte(api.Error("Failed to validate token").ToString()))
					return
				}

				if session == nil || session.IsExpired() {
					w.Write([]byte(api.Error("Invalid or expired token").ToString()))
					return
				}

				// 3. Load user and add session + user to request context
				userID := session.GetUserID()
				if userID == "" {
					w.Write([]byte(api.Error("Session missing user").ToString()))
					return
				}

				user, err := app.GetUserStore().UserFindByID(r.Context(), userID)
				if err != nil {
					w.Write([]byte(api.Error("Failed to load user").ToString()))
					return
				}

				if user == nil {
					w.Write([]byte(api.Error("User not found").ToString()))
					return
				}

				ctx := context.WithValue(r.Context(), config.APIAuthenticatedSessionContextKey{}, session)
				ctx = context.WithValue(ctx, config.APIAuthenticatedUserContextKey{}, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
}
