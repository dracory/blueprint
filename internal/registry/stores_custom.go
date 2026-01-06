package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/customstore"
)

func customStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetCustomStoreUsed() {
		return nil
	}

	if store, err := newCustomStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetCustomStore(store)
	}

	return nil
}

func customStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetCustomStoreUsed() {
		return nil
	}

	if app.GetCustomStore() == nil {
		return errors.New("custom store is not initialized")
	}

	if err := app.GetCustomStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

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
