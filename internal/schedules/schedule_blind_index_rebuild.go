package schedules

import (
	"project/internal/tasks/blind_index_rebuild"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
)

// scheduleBlindIndexRebuildTask schedules the blind index rebuild task
func scheduleBlindIndexRebuildTask(app types.AppInterface) {
	_, err := blind_index_rebuild.NewBlindIndexRebuildTask(app).
		Enqueue(blind_index_rebuild.BlindIndexAll)

	if err != nil {
		cfmt.Errorln(err.Error())
	}
}
