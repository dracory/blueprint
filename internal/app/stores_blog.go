package app

import (
	"database/sql"
	"errors"
	"project/internal/types"

	"github.com/dracory/blogstore"
)

func blogStoreInitialize(app types.AppInterface) error {
	if !app.GetConfig().GetBlogStoreUsed() {
		return nil
	}

	if store, err := newBlogStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetBlogStore(store)
	}

	return nil
}

func blogStoreMigrate(app types.AppInterface) error {
	if !app.GetConfig().GetBlogStoreUsed() {
		return nil
	}

	if app.GetBlogStore() == nil {
		return errors.New("blog store is not initialized")
	}

	if err := app.GetBlogStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newBlogStore(db *sql.DB) (blogstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:            db,
		PostTableName: "snv_blogs_post",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("blogstore.NewStore returned a nil store")
	}

	return st, nil
}
