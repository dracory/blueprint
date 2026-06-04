package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/feedstore"
)

// feedStoreInitialize initializes the feed store if enabled in the configuration.
func feedStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetFeedStoreUsed() {
		return nil
	}

	if store, err := newFeedStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetFeedStore(store)
	}

	return nil
}

// newFeedStore constructs the Feed store without running migrations
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
