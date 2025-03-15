package helpers

import (
	"context"
	"errors"
	"project/app/config"

	"github.com/samber/lo"
)

// Untokenize accepts a map of key token pairs and returns a map of key value pairs
//
// Example:
//
//	keyTokenMap := map[string]string{
//	  "key1": "token1",
//	  "key2": "token2",
//	}
//
//	untokenizedMap, err := Untokenize(keyTokenMap)
//	if err != nil {
//	  return
//	}
//
//	fmt.Println(untokenizedMap)
//	// map[key1:value1 key2:value2]
//
// Parameters:
// - keyTokenMap (map[string]string): A map of key token pairs
//
// Returns:
// - untokenizedMap (map[string]string): A map of key value pairs
// - err (error): An error if one occurred
func Untokenize(ctx context.Context, cfg config.Config, keyTokenMap map[string]string) (map[string]string, error) {
	if cfg.VaultStore == nil {
		return map[string]string{}, errors.New("vaultstore is nil")
	}

	tokens := lo.Values(keyTokenMap)
	values, err := cfg.VaultStore.TokensRead(ctx, tokens, cfg.VaultKey)

	if err != nil {
		cfg.Logger.Error("Error reading tokens", "error", err.Error())
		return map[string]string{}, err
	}

	untokenized := lo.MapValues(keyTokenMap, func(token string, key string) string {
		return values[token]
	})

	return untokenized, nil
}
