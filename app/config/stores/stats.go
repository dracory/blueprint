package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/statsstore"
)

// StatsStoreInitialize initializes the stats store
func StatsStoreInitialize(db *sql.DB) (*statsstore.Store, error) {
	statsStoreInstance, err := statsstore.NewStore(statsstore.NewStoreOptions{
		VisitorTableName: "snv_stats_visitor",
		DB:               db,
	})

	if err != nil {
		return nil, errors.Join(errors.New("statsstore.NewStore"), err)
	}

	if statsStoreInstance == nil {
		return nil, errors.New("StatsStore is nil")
	}

	return statsStoreInstance, nil
}

// StatsStoreAutoMigrate runs migrations for the stats store
func StatsStoreAutoMigrate(ctx context.Context, store *statsstore.Store) error {
	if store == nil {
		return errors.New("statsstore.AutoMigrate: StatsStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("statsstore.AutoMigrate"), err)
	}

	return nil
}
