package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/geostore"
)

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
