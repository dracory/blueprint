package schedules

import (
	"project/internal/registry"
	"project/internal/tasks/blind_index_rebuild"

	"github.com/dracory/base/cfmt"
)

// scheduleBlindIndexRebuildTask schedules the blind index rebuild task
func scheduleBlindIndexRebuildTask(registry registry.RegistryInterface) {
	if registry == nil {
		cfmt.Errorln("BlindIndexRebuild scheduling skipped; registry is nil")
		return
	}

	_, err := blind_index_rebuild.NewBlindIndexRebuildTask(registry).
		Enqueue(blind_index_rebuild.BlindIndexAll)

	if err != nil {
		cfmt.Errorln(err.Error())
	}
}
