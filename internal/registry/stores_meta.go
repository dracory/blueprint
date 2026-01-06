package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/metastore"
)

func metaStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetMetaStoreUsed() {
		return nil
	}

	if store, err := newMetaStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetMetaStore(store)
	}

	return nil
}

func metaStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetMetaStoreUsed() {
		return nil
	}

	if app.GetMetaStore() == nil {
		return errors.New("meta store is not initialized")
	}

	if err := app.GetMetaStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newMetaStore(db *sql.DB) (metastore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := metastore.NewStore(metastore.NewStoreOptions{
		DB:            db,
		MetaTableName: "snv_metas_meta",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("metastore.NewStore returned a nil store")
	}

	return st, nil
}
