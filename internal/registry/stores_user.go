package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/userstore"
)

// userStoreInitialize initializes the user store if enabled in the configuration.
func userStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetUserStoreUsed() {
		return nil
	}

	if store, err := newUserStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetUserStore(store)
	}

	return nil
}

func userStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetUserStoreUsed() {
		return nil
	}

	if registry.GetUserStore() == nil {
		return errors.New("user store is not initialized")
	}

	if err := registry.GetUserStore().AutoMigrate(); err != nil {
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
