package testutils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/registry"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func LoginAs(registry registry.RegistryInterface, r *http.Request, user userstore.UserInterface) (*http.Request, error) {
	session, err := SeedSession(registry.GetSessionStore(), r, user, 10)

	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Context(), config.AuthenticatedSessionContextKey{}, session)
	ctx = context.WithValue(ctx, config.AuthenticatedUserContextKey{}, user)
	return r.WithContext(ctx), nil
}

func CallStringHandlerAsUser(registry registry.RegistryInterface, method string, handler func(http.ResponseWriter, *http.Request) string, options test.NewRequestOptions, userID string) (body string, response *http.Response, err error) {
	user, session, err := SeedUserAndSession(registry.GetUserStore(), registry.GetSessionStore(), userID, httptest.NewRequest("GET", "/", nil), 1)
	if err != nil {
		return "", nil, err
	}

	options.Context = map[any]any{
		config.AuthenticatedUserContextKey{}:    user,
		config.AuthenticatedSessionContextKey{}: session,
	}

	return test.CallStringEndpoint(http.MethodGet, handler, options)
}
