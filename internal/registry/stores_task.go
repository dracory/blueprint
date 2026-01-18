package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/taskstore"
)

// taskStoreInitialize initializes the task store if enabled in the configuration.
func taskStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetTaskStoreUsed() {
		return nil
	}

	if store, err := newTaskStore(registry.GetDatabase()); err != nil {
		return err
	} else {
		registry.SetTaskStore(store)
	}

	return nil
}

func taskStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetTaskStoreUsed() {
		return nil
	}

	if registry.GetTaskStore() == nil {
		return errors.New("task store is not initialized")
	}

	if err := registry.GetTaskStore().AutoMigrate(); err != nil {
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
		ScheduleTableName:       "snv_tasks_schedule",
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
