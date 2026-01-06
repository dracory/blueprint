package schedules

import (
	"project/internal/tasks/clean_up"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
)

// scheduleCleanUpTask schedules the clean up task
func scheduleCleanUpTask(app types.RegistryInterface) {
	if app == nil {
		cfmt.Errorln("CleanUp scheduling skipped; app is nil")
		return
	}

	if app.GetTaskStore() == nil {
		cfmt.Warningln("CleanUp scheduling skipped; task store not configured.")
		return
	}

	task := clean_up.NewCleanUpTask(app)

	go func() {
		if handled := task.Handle(); !handled {
			cfmt.Warningln("CleanUp task handler reported failure")
		}
	}()
}
