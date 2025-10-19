package ext

import (
	"context"

	"github.com/dracory/vaultstore"
)

// VaultTokenUpsert updates or creates a token for a given value
//
// Business logic:
// - If the existingToken is empty, create a new token
// - If the existingToken is not empty, update the existing token
//
// Parameters:
// - ctx: The context
// - vaultStore: The vault store
// - vaultKey: The vault key
// - existingToken: The existing token
// - value: The value to store
// - field: The field name
//
// Returns:
// - newToken: The new token if created, or the existing token if updated
// - error: An error if something went wrong
func VaultTokenUpsert(
	ctx context.Context,
	vaultStore vaultstore.StoreInterface,
	vaultKey string,
	existingToken string,
	value string,
) (
	newToken string,
	error error,
) {
	if existingToken == "" {
		token, err := vaultStore.TokenCreate(ctx, value, vaultKey, 20)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	if err := vaultStore.TokenUpdate(ctx, existingToken, value, vaultKey); err != nil {
		return "", err
	}

	return existingToken, nil
}
