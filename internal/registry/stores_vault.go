package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/vaultstore"
)

// vaultStoreInitialize initializes the vault store if enabled in the configuration.
func vaultStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if store, err := newVaultStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetVaultStore(store)
	}

	return nil
}

func vaultStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetVaultStoreUsed() {
		return nil
	}

	if registry.GetVaultStore() == nil {
		return errors.New("vault store is not initialized")
	}

	if err := registry.GetVaultStore().AutoMigrate(); err != nil {
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
