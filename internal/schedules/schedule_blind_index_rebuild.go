package schedules

import (
	"project/internal/app"
	"project/internal/tasks/blind_index_rebuild"

	"github.com/dracory/base/cfmt"
)

// scheduleBlindIndexRebuildTask schedules the blind index rebuild task
func scheduleBlindIndexRebuildTask(app app.AppInterface) {
	if app == nil {
		cfmt.Errorln("BlindIndexRebuild scheduling skipped; app is nil")
		return
	}

	_, err := blind_index_rebuild.NewBlindIndexRebuildTask(app).
		Enqueue(blind_index_rebuild.BlindIndexAll)

	if err != nil {
		cfmt.Errorln(err.Error())
	}
}
