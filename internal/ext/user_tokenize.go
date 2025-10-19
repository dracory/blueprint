package ext

import (
	"context"
	"errors"

	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
)

// UserTokenize updates or creates tokens for a user's details
//
// Business logic:
// - If the existingToken is empty, create a new token
// - If the existingToken is not empty, update the existing token
//
// Parameters:
// - ctx: The context
// - vaultStore: The vault store
// - vaultKey: The vault key
// - user: The user
// - firstNameValue: The first name value
// - lastNameValue: The last name value
// - emailValue: The email value
// - phoneValue: The phone value
// - businessNameValue: The business name value
//
// Returns:
// - firstNameTokenUpserted: The first name token upserted
// - lastNameTokenUpserted: The last name token upserted
// - emailTokenUpserted: The email token upserted
// - phoneTokenUpserted: The phone token upserted
// - businessNameTokenUpserted: The business name token upserted
// - err: The error
func UserTokenize(
	ctx context.Context,
	vaultStore vaultstore.StoreInterface,
	vaultKey string,
	user userstore.UserInterface,
	firstNameValue string,
	lastNameValue string,
	emailValue string,
	phoneValue string,
	businessNameValue string,
) (
	firstNameTokenUpserted string,
	lastNameTokenUpserted string,
	emailTokenUpserted string,
	phoneTokenUpserted string,
	businessNameTokenUpserted string,
	err error,
) {
	if vaultStore == nil {
		return "", "", "", "", "", errors.New("vault store is nil")
	}

	if user == nil {
		return "", "", "", "", "", errors.New("user is nil")
	}

	firstNameToken := user.FirstName()
	lastNameToken := user.LastName()
	emailToken := user.Email()
	phoneToken := user.Phone()
	businessNameToken := user.BusinessName()

	if firstNameTokenUpserted, err = VaultTokenUpsert(ctx, vaultStore, vaultKey, firstNameToken, firstNameValue); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating first name"))
	}

	if lastNameTokenUpserted, err = VaultTokenUpsert(ctx, vaultStore, vaultKey, lastNameToken, lastNameValue); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating last name"))
	}

	if emailTokenUpserted, err = VaultTokenUpsert(ctx, vaultStore, vaultKey, emailToken, emailValue); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating email"))
	}

	if phoneTokenUpserted, err = VaultTokenUpsert(ctx, vaultStore, vaultKey, phoneToken, phoneValue); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating phone"))
	}

	if businessNameTokenUpserted, err = VaultTokenUpsert(ctx, vaultStore, vaultKey, businessNameToken, businessNameValue); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating business name"))
	}

	return firstNameTokenUpserted,
		lastNameTokenUpserted,
		emailTokenUpserted,
		phoneTokenUpserted,
		businessNameTokenUpserted,
		nil
}

// func UserTokenize1(app types.AppInterface, user userstore.UserInterface) error {
// 	ctx := context.Background()
// 	vaultStore := app.GetVaultStore()
// 	vaultKey := app.GetConfig().GetVaultStoreKey()

// 	emailToken, err := VaultTokenUpsert(ctx, vaultStore, vaultKey, user.Email(), "john@example.com")
// 	if err != nil {
// 		return err
// 	}

// 	firstNameToken, err := VaultTokenUpsert(ctx, vaultStore, vaultKey, user.FirstName(), "John")
// 	if err != nil {
// 		return err
// 	}

// 	lastNameToken, err := VaultTokenUpsert(ctx, vaultStore, vaultKey, user.LastName(), "Doe")
// 	if err != nil {
// 		return err
// 	}

// 	businessNameToken, err := VaultTokenUpsert(ctx, vaultStore, vaultKey, user.BusinessName(), "JD Consulting")
// 	if err != nil {
// 		return err
// 	}

// 	phoneToken, err := VaultTokenUpsert(ctx, vaultStore, vaultKey, user.Phone(), "+44111222333")
// 	if err != nil {
// 		return err
// 	}

// 	user.SetEmail(emailToken)
// 	user.SetFirstName(firstNameToken)
// 	user.SetLastName(lastNameToken)
// 	user.SetBusinessName(businessNameToken)
// 	user.SetPhone(phoneToken)
// 	user.SetCountry("GB")
// 	user.SetTimezone("Europe/London")

// 	if err := app.GetUserStore().UserUpdate(ctx, user); err != nil {
// 		return errors.New("UserUpdate returned error: " + err.Error())
// 	}

// 	return nil
// }
