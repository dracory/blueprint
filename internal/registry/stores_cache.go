package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/cachestore"
)

func cacheStoreInitialize(registry RegistryInterface) error {
	if !registry.GetConfig().GetCacheStoreUsed() {
		return nil
	}

	if store, err := newCacheStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetCacheStore(store)
	}

	return nil
}

func cacheStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetCacheStoreUsed() {
		return nil
	}

	if registry.GetCacheStore() == nil {
		return errors.New("cache store is not initialized")
	}

	if err := registry.GetCacheStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newCacheStore(db *sql.DB) (cachestore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := cachestore.NewStore(cachestore.NewStoreOptions{
		DB:             db,
		CacheTableName: "snv_caches_cache",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("cachestore.NewStore returned a nil store")
	}

	return st, nil
}
