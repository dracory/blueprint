package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/feedstore"
)

func feedStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetFeedStoreUsed() {
		return nil
	}

	if store, err := newFeedStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetFeedStore(store)
	}

	return nil
}

func feedStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetFeedStoreUsed() {
		return nil
	}

	if app.GetFeedStore() == nil {
		return errors.New("feed store is not initialized")
	}

	if err := app.GetFeedStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

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
