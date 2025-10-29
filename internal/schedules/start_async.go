package schedules

import (
	taskStats "project/internal/tasks/stats"

	"project/internal/tasks"
	"project/internal/types"
	"time"

	"github.com/dracory/base/cfmt"
	"github.com/go-co-op/gocron"
)

// scheduleStatsVisitorEnhanceTask schedules the stats visitor enhance task
func scheduleStatsVisitorEnhanceTask(app types.AppInterface) {
	_, err := taskStats.NewStatsVisitorEnhanceTask(app).Enqueue()
	if err != nil {
		cfmt.Errorln(err.Error())
	}
}

// scheduleCleanUpTask schedules the clean up task
func scheduleCleanUpTask(app types.AppInterface) {
	tasks.NewCleanUpTask(app).Handle()
}

// StartAsync starts the scheduler in the background without blocking the main thread
func StartAsync(app types.AppInterface) {
	scheduler := gocron.NewScheduler(time.UTC)

	// Example of daily scheduled task
	// scheduler.Every(1).Day().At("01:00").Do(func() {
	// 	_, err := taskhandlers.NewHelloWorldTaskHandler().Enqueue()
	// 	if err != nil {
	// 		cfmt.Errorln(err.Error())
	// 	}
	// })

	// Schedule Building the Stats Every 2 Minutes
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		scheduleStatsVisitorEnhanceTask(app)
	}); err != nil {
		cfmt.Errorln("Error scheduling stats visitor enhance task:", err.Error())
	}

	// Clean up every 20 minutes
	if _, err := scheduler.Every(20).Minutes().Do(func() {
		scheduleCleanUpTask(app)
	}); err != nil {
		cfmt.Errorln("Error scheduling clean up task:", err.Error())
	}

	// Schedule queue clear job every 2 minutes
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		queueClearJob(app)
	}); err != nil {
		cfmt.Errorln("Error scheduling queue clear job:", err.Error())
	}

	scheduler.StartAsync()
}
