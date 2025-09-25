package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/blindindexstore"
)

func blindIndexFirstNameStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexFirstNameStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetBlindIndexStoreFirstName(store)
	}

	return nil
}

func blindIndexFirstNameStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if app.GetBlindIndexStoreFirstName() == nil {
		return errors.New("blind index first name store is not initialized")
	}

	if err := app.GetBlindIndexStoreFirstName().AutoMigrate(); err != nil {
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
