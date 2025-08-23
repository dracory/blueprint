package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/settingstore"
)

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
