package helpers

import (
	"context"
	"errors"
	"log/slog"
	"project/config"

	"github.com/gouniverse/userstore"
)

func UserUntokenized(ctx context.Context, authUser userstore.UserInterface) (firstName string, lastName string, email string, err error) {
	if config.VaultStore == nil {
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

	untokenized, err := Untokenize(ctx, keyTokenMap) // use Untokenize as more resource optimized

	if err != nil {
		config.Logger.Error("Error reading tokens", slog.String("error", err.Error()))
		return "", "", "", err
	}

	firstName = untokenized[keyFirstName]
	lastName = untokenized[keyLastName]
	email = untokenized[keyEmail]

	return firstName, lastName, email, nil
}
