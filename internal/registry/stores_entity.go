package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/entitystore"
)

func entityStoreInitialize(app RegistryInterface) error {
	if !app.GetConfig().GetEntityStoreUsed() {
		return nil
	}

	if store, err := newEntityStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetEntityStore(store)
	}

	return nil
}

func entityStoreMigrate(app RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetEntityStoreUsed() {
		return nil
	}

	if app.GetEntityStore() == nil {
		return errors.New("entity store is not initialized")
	}

	if err := app.GetEntityStore().AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

func newEntityStore(db *sql.DB) (entitystore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                      db,
		EntityTableName:         "snv_entities_entity",
		EntityTrashTableName:    "snv_entities_entity_trash",
		AttributeTableName:      "snv_entities_attribute",
		AttributeTrashTableName: "snv_entities_attribute_trash",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("entitystore.NewStore returned a nil store")
	}

	return st, nil
}
