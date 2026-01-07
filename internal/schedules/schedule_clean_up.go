package schedules

import (
	"project/internal/registry"
	"project/internal/tasks/clean_up"

	"github.com/dracory/base/cfmt"
)

// scheduleCleanUpTask schedules the clean up task
func scheduleCleanUpTask(registry registry.RegistryInterface) {
	if registry == nil {
		cfmt.Errorln("CleanUp scheduling skipped; registry is nil")
		return
	}

	if registry.GetTaskStore() == nil {
		cfmt.Warningln("CleanUp scheduling skipped; task store not configured.")
		return
	}

	task := clean_up.NewCleanUpTask(registry)

	go func() {
		if handled := task.Handle(); !handled {
			cfmt.Warningln("CleanUp task handler reported failure")
		}
	}()
}
