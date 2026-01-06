package app

import (
	"database/sql"
	"errors"
	"project/internal/types"

	"github.com/dracory/logstore"
)

func logStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetLogStoreUsed() {
		return nil
	}

	if store, err := newLogStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetLogStore(store)
	}

	return nil
}

func logStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetLogStoreUsed() {
		return nil
	}

	if app.GetLogStore() == nil {
		return errors.New("log store is not initialized")
	}

	if err := app.GetLogStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

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
