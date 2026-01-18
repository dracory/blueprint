package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/blindindexstore"
)

func blindIndexFirstNameStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !registry.GetConfig().GetUserStoreUsed() || !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexFirstNameStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetBlindIndexStoreFirstName(store)
	}

	return nil
}

func blindIndexFirstNameStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !registry.GetConfig().GetUserStoreUsed() || !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if registry.GetBlindIndexStoreFirstName() == nil {
		return errors.New("blind index first name store is not initialized")
	}

	if err := registry.GetBlindIndexStoreFirstName().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newBlindIndexFirstNameStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	st, err := blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_first_name",
		Transformer: &blindindexstore.Sha256Transformer{},
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("blindindexstore.NewStore returned a nil store")
	}
	return st, nil
}
