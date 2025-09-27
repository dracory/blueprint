package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/statsstore"
)

// statsStoreInitialize initializes the stats store if enabled in the configuration.
func statsStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetStatsStoreUsed() {
		return nil
	}

	if store, err := newStatsStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetStatsStore(store)
	}

	return nil
}

func statsStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetStatsStoreUsed() {
		return nil
	}

	if app.GetStatsStore() == nil {
		return errors.New("stats store is not initialized")
	}

	if err := app.GetStatsStore().AutoMigrate(); err != nil {
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
