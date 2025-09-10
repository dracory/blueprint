package helpers

import (
	"net/http"
	"project/internal/config"

	"github.com/dracory/userstore"
)

// GetAuthUser returns the authenticated user
func GetAuthUser(r *http.Request) userstore.UserInterface {
	if r == nil {
		return nil
	}

	value := r.Context().Value(config.AuthenticatedUserContextKey{})
	if value == nil {
		return nil
	}

	user := value.(userstore.UserInterface)
	return user
}
