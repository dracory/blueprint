package ext

import (
	"context"
	"errors"
	"log/slog"
	"project/internal/helpers"
	"project/internal/types"

	"github.com/dracory/userstore"
)

func UserUntokenize(
	ctx context.Context,
	app types.AppInterface,
	vaultKey string,
	user userstore.UserInterface,
) (
	firstName string,
	lastName string,
	email string,
	businessName string,
	phone string,
	err error,
) {
	if app.GetVaultStore() == nil {
		return "", "", "", "", "", errors.New("user_untokenized: vaultstore is nil")
	}

	if user == nil {
		return "", "", "", "", "", errors.New("user_untokenized: user is nil")
	}

	firstNameToken := user.FirstName()
	lastNameToken := user.LastName()
	emailToken := user.Email()
	businessNameToken := user.BusinessName()
	phoneToken := user.Phone()

	keyFirstName := "first_name"
	keyLastName := "last_name"
	keyEmail := "email"
	keyBusinessName := "business_name"
	keyPhone := "phone"

	keyTokenMap := map[string]string{
		keyFirstName:    firstNameToken,
		keyLastName:     lastNameToken,
		keyEmail:        emailToken,
		keyBusinessName: businessNameToken,
		keyPhone:        phoneToken,
	}

	untokenized, err := helpers.Untokenize(ctx, app.GetVaultStore(), vaultKey, keyTokenMap) // use Untokenize as more resource optimized

	if err != nil {
		app.GetLogger().Error("Error reading tokens", slog.String("error", err.Error()))
		return "", "", "", "", "", err
	}

	firstName = untokenized[keyFirstName]
	lastName = untokenized[keyLastName]
	email = untokenized[keyEmail]
	businessName = untokenized[keyBusinessName]
	phone = untokenized[keyPhone]

	return firstName, lastName, email, businessName, phone, nil
}
