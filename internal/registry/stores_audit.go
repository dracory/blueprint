package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/auditstore"
)

// auditStoreInitialize initializes the audit store if enabled in the configuration.
func auditStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetAuditStoreUsed() {
		return nil
	}

	if store, err := newAuditStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetAuditStore(store)
	}

	return nil
}


func auditStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetAuditStoreUsed() {
		return nil
	}

	auditStore := registry.GetAuditStore()
	if auditStore == nil {
		return errors.New("audit store is not initialized")
	}

	err := auditStore.AutoMigrate()
	if err != nil {
		return errors.New("error migrating audit store: " + err.Error())
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
