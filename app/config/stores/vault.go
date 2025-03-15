package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/vaultstore"
)

// VaultStoreInitialize initializes the vault store
func VaultStoreInitialize(db *sql.DB) (*vaultstore.Store, error) {
	vaultStoreInstance, err := vaultstore.NewStore(vaultstore.NewStoreOptions{
		DB:             db,
		VaultTableName: "snv_vault_vault",
	})

	if err != nil {
		return nil, errors.Join(errors.New("vaultstore.NewStore"), err)
	}

	if vaultStoreInstance == nil {
		return nil, errors.New("VaultStore is nil")
	}

	return vaultStoreInstance, nil
}

// VaultStoreAutoMigrate runs migrations for the vault store
func VaultStoreAutoMigrate(ctx context.Context, store *vaultstore.Store) error {
	if store == nil {
		return errors.New("vaultstore.AutoMigrate: VaultStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("vaultstore.AutoMigrate"), err)
	}

	return nil
}
