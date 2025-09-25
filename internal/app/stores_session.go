package app

import (
	"context"
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/sessionstore"
)

func sessionStoreInitialize(app types.AppInterface) error {
	if !app.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if store, err := newSessionStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetSessionStore(store)
	}

	return nil
}

func sessionStoreMgrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetSessionStoreUsed() {
		return nil
	}

	if app.GetSessionStore() == nil {
		return errors.New("session store is not initialized")
	}

	if err := app.GetSessionStore().AutoMigrate(context.Background()); err != nil {
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
