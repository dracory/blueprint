package registry

import (
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

// newSettingStore constructs a SettingStore bound to our DB.
// Migration is handled separately in database/migrations.
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
