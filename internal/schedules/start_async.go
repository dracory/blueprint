package schedules

import (
	"context"
	"time"

	"project/internal/tasks"
	taskStats "project/internal/tasks/stats"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
	"github.com/go-co-op/gocron"
)

// scheduleStatsVisitorEnhanceTask schedules the stats visitor enhance task
func scheduleStatsVisitorEnhanceTask(app types.AppInterface) {
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

// scheduleCleanUpTask schedules the clean up task
func scheduleCleanUpTask(app types.AppInterface) {
	tasks.NewCleanUpTask(app).Handle()
}

func newScheduler(app types.AppInterface) *gocron.Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)

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

	return scheduler
}

// StartAsync starts the scheduler and stops it when the context is cancelled.
func StartAsync(ctx context.Context, app types.AppInterface) {
	if app == nil {
		cfmt.Errorln("Scheduler StartAsync called with nil app; skipping job registration")
		return
	}

	scheduler := newScheduler(app)
	scheduler.StartAsync()

	<-ctx.Done()
	scheduler.Stop()
	scheduler.Clear()
}
