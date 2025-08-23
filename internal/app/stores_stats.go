package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/statsstore"
)

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
