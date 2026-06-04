package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/metastore"
)

// metaStoreInitialize initializes the meta store if enabled in the configuration.
func metaStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetMetaStoreUsed() {
		return nil
	}

	if store, err := newMetaStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetMetaStore(store)
	}

	return nil
}

// newMetaStore constructs the Meta store without running migrations
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
