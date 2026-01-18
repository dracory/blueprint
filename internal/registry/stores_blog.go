package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/blogstore"
)

func blogStoreInitialize(registry RegistryInterface) error {
	if !registry.GetConfig().GetBlogStoreUsed() {
		return nil
	}

	if store, err := newBlogStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetBlogStore(store)
	}

	return nil
}

func blogStoreMigrate(registry RegistryInterface) error {
	if !registry.GetConfig().GetBlogStoreUsed() {
		return nil
	}

	if registry.GetBlogStore() == nil {
		return errors.New("blog store is not initialized")
	}

	if err := registry.GetBlogStore().AutoMigrate(); err != nil {
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
