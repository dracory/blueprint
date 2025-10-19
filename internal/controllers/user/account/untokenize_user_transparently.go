package account

import (
	"context"
	"errors"
	"project/internal/ext"
	"project/internal/types"

	"github.com/dracory/userstore"
)

func UntokenizeUserTranparently(
	app types.AppInterface,
	ctx context.Context,
	user userstore.UserInterface,
) (
	email string,
	firstName string,
	lastName string,
	businessName string,
	phone string,
	err error,
) {
	if user == nil {
		return "", "", "", "", "", errors.New("user is nil")
	}

	// If vault is not used, treat fields as plaintext
	if app.GetConfig().GetUserStoreVaultEnabled() {
		if app.GetUserStore() == nil {
			return "", "", "", "", "", errors.New("UserStore is not initialized")
		}
		return ext.UserUntokenize(
			ctx,
			app,
			app.GetConfig().GetVaultStoreKey(),
			user)
	}

	email = user.Email()
	firstName = user.FirstName()
	lastName = user.LastName()
	businessName = user.BusinessName()
	phone = user.Phone()

	return email, firstName, lastName, businessName, phone, nil
}
