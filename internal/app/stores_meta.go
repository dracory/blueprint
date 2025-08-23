package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/metastore"
)

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
