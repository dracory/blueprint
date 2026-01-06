package app

import (
	"database/sql"
	"errors"
	"project/internal/types"

	"github.com/dracory/auditstore"
)

func auditStoreInitialize(app types.RegistryInterface) error {
	if !app.GetConfig().GetAuditStoreUsed() {
		return nil
	}

	if store, err := newAuditStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetAuditStore(store)
	}

	return nil
}

func newAuditStore(db *sql.DB) (auditstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	store, err := auditstore.NewStore(auditstore.NewStoreOptions{
		DB:             db,
		AuditTableName: "snv_audit_record",
	})

	if err != nil {
		return nil, errors.New("error creating audit store: " + err.Error())
	}

	return store, nil
}

func auditStoreMigrate(app types.RegistryInterface) error {
	if !app.GetConfig().GetAuditStoreUsed() {
		return nil
	}

	if store := app.GetAuditStore(); store != nil {
		if migratable, ok := store.(interface{ Migrate() error }); ok {
			if err := migratable.Migrate(); err != nil {
				return errors.New("error migrating audit store: " + err.Error())
			}
		}
	}
	return nil
}
