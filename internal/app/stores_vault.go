package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/vaultstore"
)

// vaultStoreInitialize initializes the vault store if enabled in the configuration.
func vaultStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newVaultStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetVaultStore(store)
	}

	return nil
}

func vaultStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if app.GetVaultStore() == nil {
		return errors.New("vault store is not initialized")
	}

	if err := app.GetVaultStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newVaultStore(db *sql.DB) (vaultstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	st, err := vaultstore.NewStore(vaultstore.NewStoreOptions{
		DB:             db,
		VaultTableName: "snv_vault_vault",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("vaultstore.NewStore returned a nil store")
	}
	return st, nil
}
