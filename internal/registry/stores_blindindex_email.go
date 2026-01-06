package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/blindindexstore"
)

func blindIndexEmailStoreInitialize(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexEmailStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetBlindIndexStoreEmail(store)
	}

	return nil
}

func blindIndexEmailStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if app.GetBlindIndexStoreEmail() == nil {
		return errors.New("blind index email store is not initialized")
	}

	if err := app.GetBlindIndexStoreEmail().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newBlindIndexEmailStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	st, err := blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_email",
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
