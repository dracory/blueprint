package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/customstore"
)

// customStoreInitialize initializes the custom store if enabled in the configuration.
func customStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetCustomStoreUsed() {
		return nil
	}

	if store, err := newCustomStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetCustomStore(store)
	}

	return nil
}

func customStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetCustomStoreUsed() {
		return nil
	}

	customStore := registry.GetCustomStore()
	if customStore == nil {
		return errors.New("custom store is not initialized")
	}

	err := customStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newCustomStore constructs the Custom store without running migrations
func newCustomStore(db *sql.DB) (customstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := customstore.NewStore(customstore.NewStoreOptions{
		DB:        db,
		TableName: "snv_custom_record",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("customstore.NewStore returned a nil store")
	}

	return st, nil
}
