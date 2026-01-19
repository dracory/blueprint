package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/feedstore"
)

// feedStoreInitialize initializes the feed store if enabled in the configuration.
func feedStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetFeedStoreUsed() {
		return nil
	}

	if store, err := newFeedStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetFeedStore(store)
	}

	return nil
}

func feedStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetFeedStoreUsed() {
		return nil
	}

	feedStore := registry.GetFeedStore()
	if feedStore == nil {
		return errors.New("feed store is not initialized")
	}

	err := feedStore.AutoMigrate()
	if err != nil {
		return err
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
