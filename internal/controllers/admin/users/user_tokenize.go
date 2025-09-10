package admin

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dracory/userstore"
	"github.com/gouniverse/vaultstore"
)

func userTokenize(ctx context.Context, vaultStore vaultstore.StoreInterface, logger *slog.Logger, vaultKey string, user userstore.UserInterface, firstName string, lastName string, email string) (err error) {
	if vaultStore == nil {
		return errors.New("vault store is nil")
	}

	if user == nil {
		return errors.New("user is nil")
	}

	firstNameToken := user.FirstName()
	lastNameToken := user.LastName()
	emailToken := user.Email()

	err = vaultStore.TokenUpdate(ctx, firstNameToken, firstName, vaultKey)

	if err != nil {
		if logger != nil {
			logger.Error("Error updating first name", slog.String("error", err.Error()))
		}
		return err
	}

	err = vaultStore.TokenUpdate(ctx, lastNameToken, lastName, vaultKey)

	if err != nil {
		if logger != nil {
			logger.Error("Error updating last name", slog.String("error", err.Error()))
		}
		return err
	}

	err = vaultStore.TokenUpdate(ctx, emailToken, email, vaultKey)

	if err != nil {
		if logger != nil {
			logger.Error("Error updating email", slog.String("error", err.Error()))
		}
		return err
	}

	return nil
}
