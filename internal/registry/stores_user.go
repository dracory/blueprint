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

	userStore := registry.GetUserStore()
	if userStore == nil {
		return errors.New("user store is not initialized")
	}

	err := userStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newUserStore constructs the User store without running migrations
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
