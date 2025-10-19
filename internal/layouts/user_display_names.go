package layouts

import (
	"errors"
	"net/http"
	"project/internal/ext"
	"project/internal/types"

	"github.com/dracory/userstore"
)

// userDisplayNames returns the user's display names
//
// Business logic:
// - if authUser is nil, it will return "n/a"
// - if names are empty, it will return the email
// - If userstore vault is enabled, it will use the vault to get the display names
// - If userstore vault is disabled, it will use the userstore to get the display names
//
// Parameters:
// - app: the app interface
// - r: the http request
// - authUser: the authenticated user
// - vaultKey: the vault key
//
// Returns:
// - firstName: the user's first name
// - lastName: the user's last name
// - err: the error
func userDisplayNames(
	app types.AppInterface,
	r *http.Request,
	authUser userstore.UserInterface,
	vaultKey string,
) (
	firstName string,
	lastName string,
	err error,
) {
	if authUser == nil {
		return "n/a", "", errors.New("user is nil")
	}

	firstName = authUser.FirstName()
	lastName = authUser.LastName()
	email := authUser.Email()

	if app.GetConfig().GetUserStoreVaultEnabled() {
		firstName, lastName, _, _, _, err = ext.UserUntokenize(r.Context(), app, vaultKey, authUser)
		if err != nil {
			return "", "", err
		}
	}

	if firstName == "" && lastName == "" {
		return email, "", nil
	}

	return firstName, lastName, err
}
