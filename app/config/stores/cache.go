package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/cachestore"
)

// CacheStoreInitialize initializes the cache store
func CacheStoreInitialize(db *sql.DB) (cachestore.StoreInterface, error) {
	cacheStoreInstance, err := cachestore.NewStore(cachestore.NewStoreOptions{
		DB:             db,
		CacheTableName: "snv_caches_cache",
	})

	if err != nil {
		return nil, errors.Join(errors.New("cachestore.NewStore"), err)
	}

	if cacheStoreInstance == nil {
		return nil, errors.New("CacheStore is nil")
	}

	return cacheStoreInstance, nil
}

// CacheStoreAutoMigrate runs migrations for the cache store
func CacheStoreAutoMigrate(ctx context.Context, store cachestore.StoreInterface) error {
	if store == nil {
		return errors.New("cachestore.AutoMigrate: CacheStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("cachestore.AutoMigrate"), err)
	}

	return nil
}
