package helpers

import (
	"net/http"
	"project/internal/config"

	"github.com/dracory/sessionstore"
)

// GetAuthSession returns the authenticated session
func GetAuthSession(r *http.Request) sessionstore.SessionInterface {
	if r == nil {
		return nil
	}

	value := r.Context().Value(config.AuthenticatedSessionContextKey{})

	if value == nil {
		return nil
	}

	session := value.(sessionstore.SessionInterface)
	return session
}
