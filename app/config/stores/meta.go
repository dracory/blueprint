package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/metastore"
)

// MetaStoreInitialize initializes the meta store
func MetaStoreInitialize(db *sql.DB) (metastore.StoreInterface, error) {
	metaStoreInstance, err := metastore.NewStore(metastore.NewStoreOptions{
		DB:            db,
		MetaTableName: "snv_metas_meta",
	})

	if err != nil {
		return nil, errors.Join(errors.New("metastore.NewStore"), err)
	}

	if metaStoreInstance == nil {
		return nil, errors.New("MetaStore is nil")
	}

	return metaStoreInstance, nil
}

// MetaStoreAutoMigrate runs migrations for the meta store
func MetaStoreAutoMigrate(ctx context.Context, store metastore.StoreInterface) error {
	if store == nil {
		return errors.New("metastore.AutoMigrate: MetaStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("metastore.AutoMigrate"), err)
	}

	return nil
}
