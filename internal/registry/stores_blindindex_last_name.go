package registry

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/blindindexstore"
)

func blindIndexLastNameStoreInitialize(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexLastNameStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetBlindIndexStoreLastName(store)
	}

	return nil
}

func blindIndexLastNameStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if app.GetBlindIndexStoreLastName() == nil {
		return errors.New("blind index last name store is not initialized")
	}

	if err := app.GetBlindIndexStoreLastName().AutoMigrate(); err != nil {
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
