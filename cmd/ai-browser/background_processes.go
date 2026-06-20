package main

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"project/internal/app"
	"project/internal/cmsblocks"
	"project/internal/emails"
	"project/internal/middlewares"
	"project/internal/schedules"
	"project/internal/widgets"

	"github.com/dracory/taskstore"
)

func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, app app.AppInterface) error {
	_ = ctx

	if app == nil {
		return errors.New("startBackgroundProcesses called with nil app")
	}

	if app.GetConfig() == nil {
		return errors.New("startBackgroundProcesses called with nil config")
	}

	if app.GetDatabase() == nil {
		return errors.New("startBackgroundProcesses called with nil db")
	}

	if app.GetConfig().GetTaskStoreUsed() && app.GetTaskStore() == nil {
		return errors.New("startBackgroundProcesses task store is enabled but not initialized")
	}

	if app.GetConfig().GetCacheStoreUsed() && app.GetCacheStore() == nil {
		return errors.New("startBackgroundProcesses cache store is enabled but not initialized")
	}

	if app.GetConfig().GetSessionStoreUsed() && app.GetSessionStore() == nil {
		return errors.New("startBackgroundProcesses session store is enabled but not initialized")
	}

	if app.GetConfig().GetTaskStoreUsed() {
		ts := app.GetTaskStore()
		if ts != nil {
			group.Go(func(ctx context.Context) {
				runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
					IntervalSeconds: 2,
					UnstuckMinutes:  2,
					QueueName:       taskstore.DefaultQueueName,
					MaxConcurrency:  10,
					Logger:          log.Default(),
				})
				runner.Start(ctx)
			})
		}
	}

	if app.GetConfig().GetCacheStoreUsed() {
		cs := app.GetCacheStore()
		if cs != nil {
			group.Go(func(ctx context.Context) {
				if err := cs.ExpireCacheGoroutine(ctx); err != nil {
					slog.Error("Cache expiration goroutine failed", "error", err)
				}
			})
		}
	}

	if app.GetConfig().GetSessionStoreUsed() {
		ss := app.GetSessionStore()
		if ss != nil {
			group.Go(func(ctx context.Context) {
				if err := ss.SessionExpiryGoroutine(ctx); err != nil {
					slog.Error("Session expiry goroutine failed", "error", err)
				}
			})
		}
	}

	group.Go(func(ctx context.Context) {
		schedules.StartAsync(ctx, app)
	})

	emails.InitEmailSender(app)
	middlewares.CmsAddMiddlewares(app)
	widgets.CmsAddShortcodes(app)
	cmsblocks.CmsAddBlockTypes(app)

	return nil
}
