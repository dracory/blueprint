package ext

import (
	"context"
	"errors"
	"log/slog"
	"project/internal/helpers"
	"project/internal/registry"

	"github.com/dracory/userstore"
)

func UserUntokenize(
	ctx context.Context,
	registry registry.RegistryInterface,
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
	if registry.GetVaultStore() == nil {
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

	// Build token map only with non-empty tokens to avoid "missing tokens" errors
	// when some fields are not tokenized (common in focused use-cases like
	// email-only vaulting).
	keyTokenMap := map[string]string{}
	if firstNameToken != "" {
		keyTokenMap[keyFirstName] = firstNameToken
	}
	if lastNameToken != "" {
		keyTokenMap[keyLastName] = lastNameToken
	}
	if emailToken != "" {
		keyTokenMap[keyEmail] = emailToken
	}
	if businessNameToken != "" {
		keyTokenMap[keyBusinessName] = businessNameToken
	}
	if phoneToken != "" {
		keyTokenMap[keyPhone] = phoneToken
	}

	untokenized, err := helpers.Untokenize(ctx, registry.GetVaultStore(), vaultKey, keyTokenMap) // use Untokenize as more resource optimized

	if err != nil {
		registry.GetLogger().Error("Error reading tokens", slog.String("error", err.Error()))
		return "", "", "", "", "", err
	}

	firstName = untokenized[keyFirstName]
	lastName = untokenized[keyLastName]
	email = untokenized[keyEmail]
	businessName = untokenized[keyBusinessName]
	phone = untokenized[keyPhone]

	return firstName, lastName, email, businessName, phone, nil
}

// UserUntokenizeTransparently returns user fields as plaintext regardless of
// whether the user store is configured to use the vault. When vault is
// disabled, it directly returns the values from the user object. When vault
// is enabled, it delegates to UserUntokenize to read and decrypt the tokens.
func UserUntokenizeTransparently(
	ctx context.Context,
	registry registry.RegistryInterface,
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
		return "", "", "", "", "", errors.New("user_untokenized_transparently: user is nil")
	}

	// Vault disabled: treat fields as plain text
	if !registry.GetConfig().GetUserStoreVaultEnabled() {
		return user.Email(), user.FirstName(), user.LastName(), user.BusinessName(), user.Phone(), nil
	}

	// Vault enabled: ensure vault store is available and untokenize
	firstName, lastName, email, businessName, phone, err = UserUntokenize(
		ctx,
		registry,
		registry.GetConfig().GetVaultStoreKey(),
		user,
	)
	if err != nil {
		return "", "", "", "", "", err
	}

	return email, firstName, lastName, businessName, phone, nil
}
