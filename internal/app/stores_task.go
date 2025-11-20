package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/taskstore"
)

// taskStoreInitialize initializes the task store if enabled in the configuration.
func taskStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetTaskStoreUsed() {
		return nil
	}

	if store, err := newTaskStore(app.GetDB()); err != nil {
		return err
	} else {
		app.SetTaskStore(store)
	}

	return nil
}

func taskStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetTaskStoreUsed() {
		return nil
	}

	if app.GetTaskStore() == nil {
		return errors.New("task store is not initialized")
	}

	if err := app.GetTaskStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

// newTaskStore constructs the Task store without running migrations
func newTaskStore(db *sql.DB) (taskstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := taskstore.NewStore(taskstore.NewStoreOptions{
		DB:                      db,
		TaskDefinitionTableName: "snv_tasks_task_definition",
		TaskQueueTableName:      "snv_tasks_task_queue",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("taskstore.NewStore returned a nil store")
	}
	return st, nil
}
