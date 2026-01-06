package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/userstore"
)

// userStoreInitialize initializes the user store if enabled in the configuration.
func userStoreInitialize(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetUserStoreUsed() {
		return nil
	}

	if store, err := newUserStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetUserStore(store)
	}

	return nil
}

func userStoreMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetUserStoreUsed() {
		return nil
	}

	if app.GetUserStore() == nil {
		return errors.New("user store is not initialized")
	}

	if err := app.GetUserStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newUserStore(db *sql.DB) (userstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:            db,
		UserTableName: "snv_users_user",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("userstore.NewStore returned a nil store")
	}

	return st, nil
}
