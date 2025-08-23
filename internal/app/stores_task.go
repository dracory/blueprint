package app

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/taskstore"
)

// newTaskStore constructs the Task store without running migrations
func newTaskStore(db *sql.DB) (taskstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := taskstore.NewStore(taskstore.NewStoreOptions{
		DB:             db,
		TaskTableName:  "snv_tasks_task",
		QueueTableName: "snv_tasks_queue",
	})
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errors.New("taskstore.NewStore returned a nil store")
	}
	return st, nil
}
