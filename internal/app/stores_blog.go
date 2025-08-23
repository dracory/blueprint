package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/blogstore"
)

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
