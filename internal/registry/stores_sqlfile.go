package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/filesystem"
)

// sqlFileStorageInitialize initializes the SQL file storage if enabled in the configuration.
func sqlFileStorageInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSqlFileStoreUsed() {
		return nil
	}

	store, err := newSqlFileStorage(registry.GetDatabase())
	if err != nil {
		return err
	}

	registry.SetSqlFileStorage(store)
	return nil
}

func sqlFileStorageMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetSqlFileStoreUsed() {
		return nil
	}

	if registry.GetSqlFileStorage() == nil {
		return errors.New("sql file storage is not initialized")
	}

    // SQL file storage doesn't need migration
	// if err := registry.GetSqlFileStorage().AutoMigrate(); err != nil {
	// 	return err
	// }

	return nil
}

// newSqlFileStorage constructs the SQL file storage without running migrations
func newSqlFileStorage(db *sql.DB) (filesystem.StorageInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := filesystem.NewStorage(filesystem.Disk{
		DiskName:  filesystem.DRIVER_SQL,
		Driver:    filesystem.DRIVER_SQL,
		Url:       "/files",
		DB:        db,
		TableName: "snv_files_file",
	})

	if err != nil {
		return nil, err
	}

	if st == nil {
		return nil, errors.New("filesystem.NewSqlFileStorage returned a nil storage")
	}

	return st, nil
}
