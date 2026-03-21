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

	if firstNameTokenUpserted, err = vaultStore.TokenUpsert(ctx, firstNameToken, firstNameValue, vaultKey); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating first name"))
	}

	if lastNameTokenUpserted, err = vaultStore.TokenUpsert(ctx, lastNameToken, lastNameValue, vaultKey); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating last name"))
	}

	if emailTokenUpserted, err = vaultStore.TokenUpsert(ctx, emailToken, emailValue, vaultKey); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating email"))
	}

	if phoneTokenUpserted, err = vaultStore.TokenUpsert(ctx, phoneToken, phoneValue, vaultKey); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating phone"))
	}

	if businessNameTokenUpserted, err = vaultStore.TokenUpsert(ctx, businessNameToken, businessNameValue, vaultKey); err != nil {
		return "", "", "", "", "", errors.Join(err, errors.New("error updating business name"))
	}

	return firstNameTokenUpserted,
		lastNameTokenUpserted,
		emailTokenUpserted,
		phoneTokenUpserted,
		businessNameTokenUpserted,
		nil
}
