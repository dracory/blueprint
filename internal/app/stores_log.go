package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/logstore"
)

func newLogStore(db *sql.DB) (logstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := logstore.NewStore(logstore.NewStoreOptions{
		DB:           db,
		LogTableName: "snv_logs_log",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("logstore.NewStore returned a nil store")
	}

	return st, nil
}
