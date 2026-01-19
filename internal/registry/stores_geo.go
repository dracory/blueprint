package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/geostore"
)

// geoStoreInitialize initializes the geo store if enabled in the configuration.
func geoStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetGeoStoreUsed() {
		return nil
	}

	if store, err := newGeoStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetGeoStore(store)
	}

	return nil
}

func geoStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetGeoStoreUsed() {
		return nil
	}

	geoStore := registry.GetGeoStore()
	if geoStore == nil {
		return errors.New("geo store is not initialized")
	}

	err := geoStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newGeoStore constructs the Geo store without running migrations
func newGeoStore(db *sql.DB) (geostore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := geostore.NewStore(geostore.NewStoreOptions{
		DB:                db,
		CountryTableName:  "snv_geo_country",
		StateTableName:    "snv_geo_state",
		TimezoneTableName: "snv_geo_timezone",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("geostore.NewStore returned a nil store")
	}

	return st, nil
}
