package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/geostore"
)

// GeoStoreInitialize initializes the geo store
func GeoStoreInitialize(db *sql.DB) (*geostore.Store, error) {
	geoStoreInstance, err := geostore.NewStore(geostore.NewStoreOptions{
		DB:                db,
		CountryTableName:  "snv_geo_country",
		StateTableName:    "snv_geo_state",
		TimezoneTableName: "snv_geo_timezone",
	})

	if err != nil {
		return nil, errors.Join(errors.New("geostore.NewStore"), err)
	}

	if geoStoreInstance == nil {
		return nil, errors.New("GeoStore is nil")
	}

	return geoStoreInstance, nil
}

// GeoStoreAutoMigrate runs migrations for the geo store
func GeoStoreAutoMigrate(ctx context.Context, store *geostore.Store) error {
	if store == nil {
		return errors.New("geostore.AutoMigrate: GeoStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("geostore.AutoMigrate"), err)
	}

	return nil
}
