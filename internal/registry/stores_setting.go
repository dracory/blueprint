package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/settingstore"
)

// settingStoreInitialize initializes the setting store if enabled in the configuration.
func settingStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSettingStoreUsed() {
		return nil
	}

	if store, err := newSettingStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetSettingStore(store)
	}

	return nil
}

func settingStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSettingStoreUsed() {
		return nil
	}

	settingStore := registry.GetSettingStore()
	if settingStore == nil {
		return errors.New("setting store is not initialized")
	}

	if err := settingStore.AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

// newSettingStore constructs a SettingStore bound to our DB.
// Migration is handled separately in dataStoresMigrate.
func newSettingStore(db *sql.DB) (settingstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := settingstore.NewStore(settingstore.NewStoreOptions{
		DB:               db,
		SettingTableName: "snv_settings",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("settingstore.NewStore returned a nil store")
	}
	return st, nil
}
