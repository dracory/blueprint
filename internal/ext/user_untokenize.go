package ext

import (
	"context"
	"errors"
	"log/slog"
	"project/internal/registry"

	"github.com/dracory/userstore"
)

// UserUntokenizeFieldStatus tracks which fields were successfully untokenized
type UserUntokenizeFieldStatus struct {
	FirstName    bool
	LastName     bool
	Email        bool
	BusinessName bool
	Phone        bool
}

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

	firstNameToken := user.GetFirstName()
	lastNameToken := user.GetLastName()
	emailToken := user.GetEmail()
	businessNameToken := user.GetBusinessName()
	phoneToken := user.GetPhone()

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

	// If no tokens exist, return empty values without calling vault store
	if len(keyTokenMap) == 0 {
		return "", "", "", "", "", nil
	}

	untokenized, err := registry.GetVaultStore().TokensReadToResolvedMap(
		ctx,
		keyTokenMap,
		vaultKey,
	) // use TokensReadToResolvedMap as more resource optimized

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

// UserUntokenizeFieldByField untokenizes each field individually and returns
// the decrypted values along with a status map indicating which fields failed.
// This allows for granular handling of corrupted tokens.
func UserUntokenizeFieldByField(
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
	status UserUntokenizeFieldStatus,
) {
	if registry.GetVaultStore() == nil {
		return user.GetFirstName(), user.GetLastName(), user.GetEmail(), user.GetBusinessName(), user.GetPhone(), UserUntokenizeFieldStatus{}
	}

	if user == nil {
		return "", "", "", "", "", UserUntokenizeFieldStatus{}
	}

	status = UserUntokenizeFieldStatus{
		FirstName:    true,
		LastName:     true,
		Email:        true,
		BusinessName: true,
		Phone:        true,
	}

	// Helper function to untokenize a single field
	// Returns the decrypted value on success, or the original token on failure
	untokenizeField := func(token, key string) (string, bool) {
		if token == "" {
			return "", true // Empty token is not an error
		}
		result, err := registry.GetVaultStore().TokensReadToResolvedMap(
			ctx,
			map[string]string{key: token},
			vaultKey,
		)
		if err != nil {
			registry.GetLogger().Error("Error untokenizing field", slog.String("field", key), slog.String("error", err.Error()))
			return token, false // Return original token on error, mark as failed
		}
		return result[key], true
	}

	firstName, status.FirstName = untokenizeField(user.GetFirstName(), "first_name")
	lastName, status.LastName = untokenizeField(user.GetLastName(), "last_name")
	email, status.Email = untokenizeField(user.GetEmail(), "email")
	businessName, status.BusinessName = untokenizeField(user.GetBusinessName(), "business_name")
	phone, status.Phone = untokenizeField(user.GetPhone(), "phone")

	return firstName, lastName, email, businessName, phone, status
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
		return user.GetEmail(), user.GetFirstName(), user.GetLastName(), user.GetBusinessName(), user.GetPhone(), nil
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
