package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/userstore"
)

// UserStoreInitialize initializes the user store
func UserStoreInitialize(db *sql.DB) (userstore.StoreInterface, error) {
	userStoreInstance, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:            db,
		UserTableName: "snv_users_user",
	})

	if err != nil {
		return nil, errors.Join(errors.New("userstore.NewStore"), err)
	}

	if userStoreInstance == nil {
		return nil, errors.New("UserStore is nil")
	}

	return userStoreInstance, nil
}

// UserStoreAutoMigrate runs migrations for the user store
func UserStoreAutoMigrate(ctx context.Context, store userstore.StoreInterface) error {
	if store == nil {
		return errors.New("userstore.AutoMigrate: UserStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("userstore.AutoMigrate"), err)
	}

	return nil
}
