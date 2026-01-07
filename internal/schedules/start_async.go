package schedules

import (
	"context"
	"project/internal/registry"
	"time"

	"github.com/dracory/base/cfmt"
	"github.com/go-co-op/gocron"
)

func newScheduler(registry registry.RegistryInterface) *gocron.Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule Building the Stats Every 2 Minutes
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		scheduleStatsVisitorEnhanceTask(registry)
	}); err != nil {
		cfmt.Errorln("Error scheduling stats visitor enhance task:", err.Error())
	}

	// Blind index populate every hour
	if _, err := scheduler.Every(1).Hour().Do(func() {
		scheduleBlindIndexRebuildTask(registry)
	}); err != nil {
		cfmt.Errorln("Error scheduling blind index rebuild task:", err.Error())
	}

	// Clean up every 20 minutes
	if _, err := scheduler.Every(20).Minutes().Do(func() {
		scheduleCleanUpTask(registry)
	}); err != nil {
		cfmt.Errorln("Error scheduling clean up task:", err.Error())
	}

	// Schedule queue clear job every 2 minutes
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		queueClearJob(registry)
	}); err != nil {
		cfmt.Errorln("Error scheduling queue clear job:", err.Error())
	}

	return scheduler
}

// StartAsync starts the scheduler and stops it when the context is cancelled.
func StartAsync(ctx context.Context, registry registry.RegistryInterface) {
	if registry == nil {
		cfmt.Errorln("Scheduler StartAsync called with nil registry; skipping job registration")
		return
	}

	scheduler := newScheduler(registry)
	scheduler.StartAsync()

	<-ctx.Done()
	scheduler.Stop()
	scheduler.Clear()
}
