package authadapters

import (
	"errors"

	"github.com/dracory/cachestore"
)

// TempKeyAdapter exposes cachestore as the auth package's temporary-key
// store (login codes, attempt counters, delivery keys). cachestore.Get
// returns the stored value or the provided default; an empty default + ""
// result is treated as "missing" — matching the package's fail-open
// semantics for absent keys.
type TempKeyAdapter struct {
	Store cachestore.StoreInterface
}

// Get returns the value for key, or an error when the key is absent or the
// read fails.
func (a *TempKeyAdapter) Get(key string) (string, error) {
	if a.Store == nil {
		return "", errors.New("authadapters: cache store is nil")
	}
	v, err := a.Store.Get(key, "")
	if err != nil {
		return "", err
	}
	if v == "" {
		return "", errors.New("key not found")
	}
	return v, nil
}

// Set stores value under key with the given TTL in seconds.
func (a *TempKeyAdapter) Set(key string, value string, expiresSeconds int) error {
	if a.Store == nil {
		return errors.New("authadapters: cache store is nil")
	}
	return a.Store.Set(key, value, int64(expiresSeconds))
}

// Delete removes key — used by code invalidation and delivery consume.
func (a *TempKeyAdapter) Delete(key string) error {
	if a.Store == nil {
		return errors.New("authadapters: cache store is nil")
	}
	return a.Store.Remove(key)
}
