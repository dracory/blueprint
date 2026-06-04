package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/logstore"
)

// logStoreInitialize initializes the log store if enabled in the configuration.
func logStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetLogStoreUsed() {
		return nil
	}

	if store, err := newLogStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetLogStore(store)
	}

	return nil
}

// newLogStore constructs the Log store without running migrations
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
