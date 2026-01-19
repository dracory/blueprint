package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/metastore"
)

// metaStoreInitialize initializes the meta store if enabled in the configuration.
func metaStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetMetaStoreUsed() {
		return nil
	}

	if store, err := newMetaStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetMetaStore(store)
	}

	return nil
}

func metaStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetMetaStoreUsed() {
		return nil
	}

	metaStore := registry.GetMetaStore()
	if metaStore == nil {
		return errors.New("meta store is not initialized")
	}

	err := metaStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newMetaStore constructs the Meta store without running migrations
func newMetaStore(db *sql.DB) (metastore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := metastore.NewStore(metastore.NewStoreOptions{
		DB:            db,
		MetaTableName: "snv_metas_meta",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("metastore.NewStore returned a nil store")
	}

	return st, nil
}
