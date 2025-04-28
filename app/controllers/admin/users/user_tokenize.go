package admin

import (
	"context"
	"errors"
	"log/slog"
	"project/config"

	"github.com/gouniverse/userstore"
)

func userTokenize(user userstore.UserInterface, firstName string, lastName string, email string) (err error) {
	if config.VaultStore == nil {
		return errors.New("vault store is nil")
	}

	if user == nil {
		return errors.New("user is nil")
	}

	firstNameToken := user.FirstName()
	lastNameToken := user.LastName()
	emailToken := user.Email()

	ctx := context.Background()

	err = config.VaultStore.TokenUpdate(ctx, firstNameToken, firstName, config.VaultKey)

	if err != nil {
		config.Logger.Error("Error updating first name", slog.String("error", err.Error()))
		return err
	}

	err = config.VaultStore.TokenUpdate(ctx, lastNameToken, lastName, config.VaultKey)

	if err != nil {
		config.Logger.Error("Error updating last name", slog.String("error", err.Error()))
		return err
	}

	err = config.VaultStore.TokenUpdate(ctx, emailToken, email, config.VaultKey)

	if err != nil {
		config.Logger.Error("Error updating email", slog.String("error", err.Error()))
		return err
	}

	return nil
}
