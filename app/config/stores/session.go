package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/sessionstore"
)

// SessionStoreInitialize initializes the session store
func SessionStoreInitialize(db *sql.DB) (sessionstore.StoreInterface, error) {
	sessionStoreInstance, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:               db,
		SessionTableName: "snv_sessions_session",
		TimeoutSeconds:   7200,
	})

	if err != nil {
		return nil, errors.Join(errors.New("sessionstore.NewStore"), err)
	}

	if sessionStoreInstance == nil {
		return nil, errors.New("SessionStore is nil")
	}

	return sessionStoreInstance, nil
}

// SessionStoreAutoMigrate runs migrations for the session store
func SessionStoreAutoMigrate(ctx context.Context, store sessionstore.StoreInterface) error {
	if store == nil {
		return errors.New("sessionstore.AutoMigrate: SessionStore is nil")
	}

	err := store.AutoMigrate(ctx)

	if err != nil {
		return errors.Join(errors.New("sessionstore.AutoMigrate"), err)
	}

	return nil
}
