package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/cachestore"
)

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
