package app

import (
	"database/sql"
	"errors"
	"project/internal/types"

	"github.com/dracory/filesystem"
)

// sqlFileStorageInitialize initializes the SQL file storage if enabled in the configuration.
func sqlFileStorageInitialize(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetSqlFileStoreUsed() {
		return nil
	}

	store, err := newSqlFileStorage(app.GetDB())
	if err != nil {
		return err
	}

	app.SetSqlFileStorage(store)
	return nil
}

func sqlFileStorageMigrate(app types.RegistryInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetSqlFileStoreUsed() {
		return nil
	}

	if app.GetSqlFileStorage() == nil {
		return errors.New("sql file storage is not initialized")
	}

	// if err := app.GetSqlFileStorage().AutoMigrate(); err != nil {
	// 	return err
	// }

	return nil
}

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
		return nil, errors.New("filesystem.NewStorage returned a nil storage")
	}

	return st, nil
}
