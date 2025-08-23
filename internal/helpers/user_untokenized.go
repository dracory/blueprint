package helpers

import (
	"context"
	"errors"
	"log/slog"
	"project/internal/types"

	"github.com/gouniverse/userstore"
)

func UserUntokenized(ctx context.Context, app types.AppInterface, vaultKey string, authUser userstore.UserInterface) (firstName string, lastName string, email string, err error) {
	if app.GetVaultStore() == nil {
		return "", "", "", errors.New("vaultstore is nil")
	}

	firstNameToken := authUser.FirstName()
	lastNameToken := authUser.LastName()
	emailToken := authUser.Email()

	keyFirstName := "first_name"
	keyLastName := "last_name"
	keyEmail := "email"

	keyTokenMap := map[string]string{
		keyFirstName: firstNameToken,
		keyLastName:  lastNameToken,
		keyEmail:     emailToken,
	}

	untokenized, err := Untokenize(ctx, app.GetVaultStore(), vaultKey, keyTokenMap) // use Untokenize as more resource optimized

	if err != nil {
		app.GetLogger().Error("Error reading tokens", slog.String("error", err.Error()))
		return "", "", "", err
	}

	firstName = untokenized[keyFirstName]
	lastName = untokenized[keyLastName]
	email = untokenized[keyEmail]

	return firstName, lastName, email, nil
}
