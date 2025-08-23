package layouts

import (
	"errors"
	"net/http"
	"project/internal/helpers"
	"project/internal/types"

	"github.com/gouniverse/userstore"
)

func getUserData(app types.AppInterface, r *http.Request, authUser userstore.UserInterface, vaultKey string) (firstName string, lastName string, err error) {
	if authUser == nil {
		return "n/a", "", errors.New("user is nil")
	}

	if !app.GetConfig().GetVaultStoreUsed() {
		firstName = authUser.FirstName()
		lastName = authUser.LastName()

		if firstName == "" && lastName == "" {
			return authUser.Email(), "", nil
		}

		return firstName, lastName, nil
	}

	firstName, lastName, _, err = helpers.UserUntokenized(r.Context(), app, vaultKey, authUser)

	return firstName, lastName, err
}
