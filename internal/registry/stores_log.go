package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/logstore"
)

func logStoreInitialize(registry RegistryInterface) error {
	if !registry.GetConfig().GetLogStoreUsed() {
		return nil
	}

	if store, err := newLogStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetLogStore(store)
	}

	return nil
}

func logStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetLogStoreUsed() {
		return nil
	}

	if registry.GetLogStore() == nil {
		return errors.New("log store is not initialized")
	}

	if err := registry.GetLogStore().AutoMigrate(); err != nil {
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
