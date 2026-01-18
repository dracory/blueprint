package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/blindindexstore"
)

func blindIndexLastNameStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !registry.GetConfig().GetUserStoreUsed() || !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexLastNameStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetBlindIndexStoreLastName(store)
	}

	return nil
}

func blindIndexLastNameStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !registry.GetConfig().GetUserStoreUsed() || !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if registry.GetBlindIndexStoreLastName() == nil {
		return errors.New("blind index last name store is not initialized")
	}

	if err := registry.GetBlindIndexStoreLastName().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newBlindIndexLastNameStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	st, err := blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_last_name",
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
