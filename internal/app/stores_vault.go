package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/vaultstore"
)

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


