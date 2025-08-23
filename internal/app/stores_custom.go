package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/customstore"
)

func newCustomStore(db *sql.DB) (customstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := customstore.NewStore(customstore.NewStoreOptions{
		DB:        db,
		TableName: "snv_custom_record",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("customstore.NewStore returned a nil store")
	}

	return st, nil
}
