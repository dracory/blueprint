package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/filesystem"
)

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
