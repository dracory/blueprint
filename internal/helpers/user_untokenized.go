package helpers

import (
	"context"
	"errors"
	"project/app/config"

	"github.com/gouniverse/userstore"
)

func UserUntokenized(ctx context.Context, cfg config.Config, authUser userstore.UserInterface) (firstName string, lastName string, email string, err error) {
	if cfg.VaultStore == nil {
		return "", "", "", errors.New("vaultstore is nil")
	}

	firstNameToken := authUser.FirstName()
	lastNameToken := authUser.LastName()
	emailToken := authUser.Email()

	firstName, err = cfg.VaultStore.TokenRead(ctx, firstNameToken, cfg.VaultKey)

	if err != nil {
		cfg.Logger.Error("Error reading first name", "error", err.Error())
		return "", "", "", err
	}

	lastName, err = cfg.VaultStore.TokenRead(ctx, lastNameToken, cfg.VaultKey)

	if err != nil {
		cfg.Logger.Error("Error reading last name", "error", err.Error())
		return "", "", "", err
	}

	email, err = cfg.VaultStore.TokenRead(ctx, emailToken, cfg.VaultKey)

	if err != nil {
		cfg.Logger.Error("Error reading email", "error", err.Error())
		return "", "", "", err
	}

	return firstName, lastName, email, nil
}
