package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/blogstore"
)

// BlogStoreInitialize initializes the blog store
func BlogStoreInitialize(db *sql.DB) (blogstore.StoreInterface, error) {
	blogStoreInstance, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:            db,
		PostTableName: "snv_blog_post",
	})

	if err != nil {
		return nil, errors.Join(errors.New("blogstore.NewStore"), err)
	}

	if blogStoreInstance == nil {
		return nil, errors.New("BlogStore is nil")
	}

	return blogStoreInstance, nil
}

// BlogStoreAutoMigrate runs migrations for the blog store
func BlogStoreAutoMigrate(ctx context.Context, store blogstore.StoreInterface) error {
	if store == nil {
		return errors.New("blogstore.AutoMigrate: BlogStore is nil")
	}

	// Use type assertion to access the AutoMigrate method
	if storeImpl, ok := store.(blogstore.StoreInterface); ok {
		err := storeImpl.AutoMigrate()
		if err != nil {
			return errors.Join(errors.New("blogstore.AutoMigrate"), err)
		}
		return nil
	}

	return errors.New("blogstore.AutoMigrate: Failed to cast store to *blogstore.Store")
}
