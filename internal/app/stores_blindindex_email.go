package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/blindindexstore"
)

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
