package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/feedstore"
)

func newFeedStore(db *sql.DB) (feedstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := feedstore.NewStore(feedstore.NewStoreOptions{
		DB:            db,
		FeedTableName: "snv_feeds_feed",
		LinkTableName: "snv_feeds_link",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("feedstore.NewStore returned a nil store")
	}

	return st, nil
}
