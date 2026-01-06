package schedules

import (
	"context"
	"time"

	"project/internal/types"

	"github.com/dracory/base/cfmt"
	"github.com/go-co-op/gocron"
)

func newScheduler(app types.RegistryInterface) *gocron.Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule Building the Stats Every 2 Minutes
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		scheduleStatsVisitorEnhanceTask(app)
	}); err != nil {
		cfmt.Errorln("Error scheduling stats visitor enhance task:", err.Error())
	}
	
	// Blind index populate every hour
	if _, err := scheduler.Every(1).Hour().Do(func() {
		scheduleBlindIndexRebuildTask(app)
	}); err != nil {
		cfmt.Errorln("Error scheduling blind index rebuild task:", err.Error())
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
func StartAsync(ctx context.Context, app types.RegistryInterface) {
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
