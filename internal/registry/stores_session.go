package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/sessionstore"
)

func sessionStoreInitialize(registry RegistryInterface) error {
	if !registry.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if store, err := newSessionStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetSessionStore(store)
	}

	return nil
}

func sessionStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if registry.GetSessionStore() == nil {
		return errors.New("session store is not initialized")
	}

	if err := registry.GetSessionStore().AutoMigrate(context.Background()); err != nil {
		return err
	}

	return nil
}

func newSessionStore(db *sql.DB) (sessionstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:               db,
		SessionTableName: "snv_sessions_session",
		TimeoutSeconds:   7200,
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("sessionstore.NewStore returned a nil store")
	}

	return st, nil
}
