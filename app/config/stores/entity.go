package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/entitystore"
)

// EntityStoreInitialize initializes the entity store
func EntityStoreInitialize(db *sql.DB) (entitystore.StoreInterface, error) {
	entityStoreInstance, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                      db,
		EntityTableName:         "snv_entity",
		EntityTrashTableName:    "snv_entity_trash",
		AttributeTableName:      "snv_entity_attribute",
		AttributeTrashTableName: "snv_entity_attribute_trash",
	})

	if err != nil {
		return nil, errors.Join(errors.New("entitystore.NewStore"), err)
	}

	if entityStoreInstance == nil {
		return nil, errors.New("EntityStore is nil")
	}

	return entityStoreInstance, nil
}

// EntityStoreAutoMigrate runs migrations for the entity store
func EntityStoreAutoMigrate(ctx context.Context, store entitystore.StoreInterface) error {
	if store == nil {
		return errors.New("entitystore.AutoMigrate: EntityStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("entitystore.AutoMigrate"), err)
	}

	return nil
}
