package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/sessionstore"
)

// sessionStoreInitialize initializes the session store if enabled in the configuration.
func sessionStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if store, err := newSessionStore(app.GetDatabase(), app); err != nil {
		return err
	} else {
		store.EnableDebug(app.GetConfig().GetAppDebug())
		app.SetSessionStore(store)
	}

	return nil
}

// newSessionStore constructs the Session store without running migrations
func newSessionStore(db *sql.DB, app AppInterface) (sessionstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	timeoutSeconds := int64(7200) // 2 hours default
	if app.GetConfig() != nil && app.GetConfig().IsEnvDevelopment() {
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
