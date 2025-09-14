package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/entitystore"
)

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
