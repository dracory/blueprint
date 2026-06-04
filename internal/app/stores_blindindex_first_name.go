package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/blindindexstore"
)

// blindIndexFirstNameStoreInitialize initializes the blind index first name store if enabled in the configuration.
func blindIndexFirstNameStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	// Blind index stores: create and set only if user store is enabled and vault store is enabled
	if !app.GetConfig().GetUserStoreUsed() || !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newBlindIndexFirstNameStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetBlindIndexStoreFirstName(store)
	}

	return nil
}

// newBlindIndexFirstNameStore constructs the Blind Index First Name store without running migrations
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
