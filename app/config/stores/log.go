package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/logstore"
)

// LogStoreInitialize initializes the log store
func LogStoreInitialize(db *sql.DB) (logstore.StoreInterface, error) {
	logStoreInstance, err := logstore.NewStore(logstore.NewStoreOptions{
		DB:           db,
		LogTableName: "snv_logs_log",
	})

	if err != nil {
		return nil, errors.Join(errors.New("logstore.NewStore"), err)
	}

	if logStoreInstance == nil {
		return nil, errors.New("LogStore is nil")
	}

	return logStoreInstance, nil
}

// LogStoreAutoMigrate runs migrations for the log store
func LogStoreAutoMigrate(ctx context.Context, store logstore.StoreInterface) error {
	if store == nil {
		return errors.New("logstore.AutoMigrate: LogStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("logstore.AutoMigrate"), err)
	}

	return nil
}
