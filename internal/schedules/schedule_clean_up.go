package schedules

import (
	"project/internal/tasks"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
)

// scheduleCleanUpTask schedules the clean up task
func scheduleCleanUpTask(app types.AppInterface) {
	if app == nil {
		cfmt.Errorln("CleanUp scheduling skipped; app is nil")
		return
	}

	if app.GetTaskStore() == nil {
		cfmt.Warningln("CleanUp scheduling skipped; task store not configured.")
		return
	}

	task := tasks.NewCleanUpTask(app)

	go func() {
		if handled := task.Handle(); !handled {
			cfmt.Warningln("CleanUp task handler reported failure")
		}
	}()
}
