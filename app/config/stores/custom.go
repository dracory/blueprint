package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/customstore"
)

// CustomStoreInitialize initializes the custom store
func CustomStoreInitialize(db *sql.DB) (customstore.StoreInterface, error) {
	customStoreInstance, err := customstore.NewStore(customstore.NewStoreOptions{
		DB:        db,
		TableName: "snv_custom_record",
	})

	if err != nil {
		return nil, errors.Join(errors.New("customstore.NewStore"), err)
	}

	if customStoreInstance == nil {
		return nil, errors.New("CustomStore is nil")
	}

	return customStoreInstance, nil
}

// CustomStoreAutoMigrate runs migrations for the custom store
func CustomStoreAutoMigrate(ctx context.Context, store customstore.StoreInterface) error {
	if store == nil {
		return errors.New("customstore.AutoMigrate: CustomStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("customstore.AutoMigrate"), err)
	}

	return nil
}
