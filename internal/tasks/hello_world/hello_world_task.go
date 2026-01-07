package hello_world

import (
	"context"
	"errors"
	"project/internal/registry"

	"github.com/dracory/taskstore"
)

func NewHelloWorldTask(registry registry.RegistryInterface) *helloWorldTask {
	return &helloWorldTask{
		registry: registry,
	}
}

type helloWorldTask struct {
	taskstore.TaskHandlerBase
	registry registry.RegistryInterface
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
	if handler.registry == nil {
		return nil, errors.New("app is nil")
	}
	if handler.registry.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}
	return handler.registry.GetTaskStore().TaskDefinitionEnqueueByAlias(
		context.Background(),
		taskstore.DefaultQueueName,
		handler.Alias(),
		map[string]any{},
	)
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
