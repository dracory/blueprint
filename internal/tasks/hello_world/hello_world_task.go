package hello_world

import (
	"errors"
	"project/internal/types"

	"github.com/dracory/taskstore"
)

func NewHelloWorldTask(app types.AppInterface) *helloWorldTask {
	return &helloWorldTask{
		app: app,
	}
}

type helloWorldTask struct {
	taskstore.TaskHandlerBase
	app types.AppInterface
}

var _ taskstore.TaskHandlerInterface = (*helloWorldTask)(nil) // verify it extends the task interface

func (handler *helloWorldTask) Alias() string {
	return "HelloWorldTask"
}

func (handler *helloWorldTask) Title() string {
	return "Hello World"
}

func (handler *helloWorldTask) Description() string {
	return "Say hello world"
}

func (handler *helloWorldTask) Enqueue() (task taskstore.TaskQueueInterface, err error) {
	if handler.app == nil {
		return nil, errors.New("app is nil")
	}
	if handler.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}
	return handler.app.GetTaskStore().TaskEnqueueByAlias(handler.Alias(), map[string]any{})
}

func (handler *helloWorldTask) Handle() bool {
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue()

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	handler.LogInfo("Hello World!")
	return true
}
