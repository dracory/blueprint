package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/userstore"
)

func newUserStore(db *sql.DB) (userstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:            db,
		UserTableName: "snv_users_user",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("userstore.NewStore returned a nil store")
	}

	return st, nil
}
