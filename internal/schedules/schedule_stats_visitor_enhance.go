package schedules

import (
	"project/internal/registry"
	taskStats "project/internal/tasks/stats"

	"github.com/dracory/base/cfmt"
)

// scheduleStatsVisitorEnhanceTask schedules the stats visitor enhance task
func scheduleStatsVisitorEnhanceTask(registry registry.RegistryInterface) {
	if registry == nil {
		cfmt.Errorln("StatsVisitorEnhance scheduling skipped; registry is nil")
		return
	}

	if registry.GetTaskStore() == nil {
		cfmt.Warningln("StatsVisitorEnhance scheduling skipped; task store not configured.")
		return
	}

	_, err := taskStats.NewStatsVisitorEnhanceTask(registry).Enqueue()
	if err != nil {
		cfmt.Errorln(err.Error())
	}
}
