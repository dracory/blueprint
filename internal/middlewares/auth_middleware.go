package middlewares

import (
	"context"
	"project/internal/app"
	"project/internal/config"

	"github.com/dracory/auth"
	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

// authSessionStoreAdapter wraps sessionstore.StoreInterface to satisfy rtrMiddleware.AuthSessionStore
type authSessionStoreAdapter struct {
	store sessionstore.StoreInterface
}

func (a *authSessionStoreAdapter) SessionFindByKey(ctx context.Context, key string) (rtrMiddleware.AuthSession, error) {
	return a.store.SessionFindByKey(ctx, key)
}

// authUserStoreAdapter wraps userstore.StoreInterface to satisfy rtrMiddleware.AuthUserStore
type authUserStoreAdapter struct {
	store userstore.StoreInterface
}

func (a *authUserStoreAdapter) UserFindByID(ctx context.Context, id string) (rtrMiddleware.AuthUser, error) {
	return a.store.UserFindByID(ctx, id)
}

// authLoggerAdapter wraps slog.Logger to satisfy rtrMiddleware.AuthLogger
type authLoggerAdapter struct {
	logger interface {
		Error(msg string, args ...any)
	}
}

func (l *authLoggerAdapter) Error(msg string, args ...any) {
	if l.logger != nil {
		l.logger.Error(msg, args...)
	}
}

// AuthMiddleware creates the auth middleware using the shared rtr implementation
func AuthMiddleware(app app.AppInterface) rtr.MiddlewareInterface {
	var sessionStore rtrMiddleware.AuthSessionStore
	if s := app.GetSessionStore(); s != nil {
		sessionStore = &authSessionStoreAdapter{store: s}
	}

	var userStore rtrMiddleware.AuthUserStore
	if u := app.GetUserStore(); u != nil {
		userStore = &authUserStoreAdapter{store: u}
	}

	cfg := rtrMiddleware.AuthMiddlewareConfig{
		SessionStore:      sessionStore,
		UserStore:         userStore,
		Logger:            &authLoggerAdapter{logger: app.GetLogger()},
		MemoryCache:       app.GetMemoryCache(),
		ContextKeyUser:    config.AuthenticatedUserContextKey{},
		ContextKeySession: config.AuthenticatedSessionContextKey{},
		CookieName:        auth.CookieName,
	}

	return rtrMiddleware.AuthMiddleware(cfg)
}
