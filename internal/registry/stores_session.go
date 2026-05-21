package registry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/sessionstore"
)

// sessionStoreInitialize initializes the session store if enabled in the configuration.
func sessionStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if store, err := newSessionStore(registry.GetDatabase(), registry); err != nil {
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

	sessionStore := registry.GetSessionStore()
	if sessionStore == nil {
		return errors.New("session store is not initialized")
	}

	if err := sessionStore.MigrateUp(context.Background()); err != nil {
		return err
	}

	return nil
}

// newSessionStore constructs the Session store without running migrations
func newSessionStore(db *sql.DB, registry RegistryInterface) (sessionstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	timeoutSeconds := int64(7200) // 2 hours default
	if registry.GetConfig() != nil && registry.GetConfig().IsEnvDevelopment() {
		timeoutSeconds = 14400 // 4 hours in development
	}

	st, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:               db,
		SessionTableName: "snv_sessions_session",
		TimeoutSeconds:   timeoutSeconds,
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("sessionstore.NewStore returned a nil store")
	}

	return st, nil
}
