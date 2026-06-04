package app

import (
	"database/sql"
	"errors"

	"github.com/dracory/auditstore"
)

// auditStoreInitialize initializes the audit store if enabled in the configuration.
func auditStoreInitialize(app AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

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

// newAuditStore constructs the Audit store without running migrations
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

	if store == nil {
		return nil, errors.New("auditstore.NewStore returned a nil store")
	}

	return store, nil
}
