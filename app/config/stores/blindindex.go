package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/blindindexstore"
)

// BlindIndexStoreInitialize initializes a blind index store
func BlindIndexStoreInitialize(db *sql.DB, tableName string) (blindindexstore.StoreInterface, error) {
	blindIndexStoreInstance, err := blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   tableName,
		Transformer: &blindindexstore.Sha256Transformer{},
	})

	if err != nil {
		return nil, errors.Join(errors.New("blindindexstore.NewStore"), err)
	}

	if blindIndexStoreInstance == nil {
		return nil, errors.New("blindindexstore.NewStore: blindIndexStoreInstance is nil")
	}

	return blindIndexStoreInstance, nil
}

// BlindIndexStoreAutoMigrate runs migrations for a blind index store
func BlindIndexStoreAutoMigrate(ctx context.Context, store blindindexstore.StoreInterface) error {
	if store == nil {
		return errors.New("blindindexstore.AutoMigrate: BlindIndexStore is nil")
	}

	err := store.AutoMigrate()
	if err != nil {
		return errors.Join(errors.New("blindindexstore.AutoMigrate"), err)
	}

	return nil
}

// MigrateBlindIndexStore migrates a blind index store with type assertion
func MigrateBlindIndexStore(ctx context.Context, store interface{}) error {
	if store == nil {
		return errors.New("blindindexstore.MigrateBlindIndexStore: store is nil")
	}

	if blindIndexStore, ok := store.(blindindexstore.StoreInterface); ok {
		return BlindIndexStoreAutoMigrate(ctx, blindIndexStore)
	}

	return errors.New("blindindexstore.MigrateBlindIndexStore: invalid store type")
}
