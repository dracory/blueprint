package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/taskstore"
)

// TaskStoreInitialize initializes the task store
func TaskStoreInitialize(db *sql.DB) (*taskstore.Store, error) {
	taskStoreInstance, err := taskstore.NewStore(taskstore.NewStoreOptions{
		DB:             db,
		TaskTableName:  "snv_tasks_task",
		QueueTableName: "snv_tasks_queue",
	})

	if err != nil {
		return nil, errors.Join(errors.New("taskstore.NewStore"), err)
	}

	if taskStoreInstance == nil {
		return nil, errors.New("TaskStore is nil")
	}

	return taskStoreInstance, nil
}

// TaskStoreAutoMigrate runs migrations for the task store
func TaskStoreAutoMigrate(ctx context.Context, store *taskstore.Store) error {
	if store == nil {
		return errors.New("taskstore.AutoMigrate: TaskStore is nil")
	}

	err := store.AutoMigrate()

	if err != nil {
		return errors.Join(errors.New("taskstore.AutoMigrate"), err)
	}

	return nil
}
