package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/cachestore"
)

func cacheStoreInitialize(app types.AppInterface) error {
	if !app.GetConfig().GetCacheStoreUsed() {
		return nil
	}

	if store, err := newCacheStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetCacheStore(store)
	}

	return nil
}

func cacheStoreMgrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetCacheStoreUsed() {
		return nil
	}

	if app.GetCacheStore() == nil {
		return errors.New("cache store is not initialized")
	}

	if err := app.GetCacheStore().AutoMigrate(); err != nil {
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
