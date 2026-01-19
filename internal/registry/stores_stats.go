package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/statsstore"
)

// statsStoreInitialize initializes the stats store if enabled in the configuration.
func statsStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetStatsStoreUsed() {
		return nil
	}

	if store, err := newStatsStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetStatsStore(store)
	}

	return nil
}

func statsStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetStatsStoreUsed() {
		return nil
	}

	statsStore := registry.GetStatsStore()
	if statsStore == nil {
		return errors.New("stats store is not initialized")
	}

	err := statsStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newStatsStore constructs the Stats store without running migrations
func newStatsStore(db *sql.DB) (statsstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := statsstore.NewStore(statsstore.NewStoreOptions{
		DB:               db,
		VisitorTableName: "snv_stats_visitor",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("statsstore.NewStore returned a nil store")
	}
	return st, nil
}
