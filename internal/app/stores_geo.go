package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/geostore"
)

func geoStoreInitialize(app types.AppInterface) error {
	if !app.GetConfig().GetGeoStoreUsed() {
		return nil
	}

	if store, err := newGeoStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetGeoStore(store)
	}

	return nil
}

func geoStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetGeoStoreUsed() {
		return nil
	}

	if app.GetGeoStore() == nil {
		return errors.New("geo store is not initialized")
	}

	if err := app.GetGeoStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

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
