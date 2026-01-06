package schedules

import (
	"project/internal/types"

	taskStats "project/internal/tasks/stats"

	"github.com/dracory/base/cfmt"
)

// scheduleStatsVisitorEnhanceTask schedules the stats visitor enhance task
func scheduleStatsVisitorEnhanceTask(app types.RegistryInterface) {
	if app == nil {
		cfmt.Errorln("StatsVisitorEnhance scheduling skipped; app is nil")
		return
	}

	if app.GetTaskStore() == nil {
		cfmt.Warningln("StatsVisitorEnhance scheduling skipped; task store not configured.")
		return
	}

	_, err := taskStats.NewStatsVisitorEnhanceTask(app).Enqueue()
	if err != nil {
		cfmt.Errorln(err.Error())
	}
}
