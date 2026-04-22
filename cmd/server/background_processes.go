package main

import (
	"context"
	"errors"
	"log"
	"log/slog"

	// "project/internal/cmsblocks"
	"project/internal/emails"
	"project/internal/middlewares"
	"project/internal/registry"
	"project/internal/schedules"
	"project/internal/widgets"

	"github.com/dracory/taskstore"
)

func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, registry registry.RegistryInterface) error {
	_ = ctx // Suppress unused parameter warning, use it when needed for context cancellation or propagation

	if registry == nil {
		return errors.New("startBackgroundProcesses called with nil registry")
	}

	if registry.GetConfig() == nil {
		return errors.New("startBackgroundProcesses called with nil config")
	}

	if registry.GetDatabase() == nil {
		return errors.New("startBackgroundProcesses called with nil db")
	}

	if registry.GetConfig().GetTaskStoreUsed() && registry.GetTaskStore() == nil {
		return errors.New("startBackgroundProcesses task store is enabled but not initialized")
	}

	if registry.GetConfig().GetCacheStoreUsed() && registry.GetCacheStore() == nil {
		return errors.New("startBackgroundProcesses cache store is enabled but not initialized")
	}

	if registry.GetConfig().GetSessionStoreUsed() && registry.GetSessionStore() == nil {
		return errors.New("startBackgroundProcesses session store is enabled but not initialized")
	}

	if registry.GetConfig().GetTaskStoreUsed() {
		ts := registry.GetTaskStore()
		if ts != nil {
			group.Go(func(ctx context.Context) {
				// Run the default task queue worker loop using the updated TaskQueue API
				// 10 workers, 2-second polling interval
				runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
					IntervalSeconds: 2,
					UnstuckMinutes:  2,
					MaxConcurrency:  10,
					Logger:          log.Default(),
				})
				runner.Start(ctx)
			})
		}
	}
	if registry.GetConfig().GetCacheStoreUsed() {
		cs := registry.GetCacheStore()
		if cs != nil {
			group.Go(func(ctx context.Context) {
				if err := cs.ExpireCacheGoroutine(ctx); err != nil {
					slog.Error("Cache expiration goroutine failed", "error", err)
				}
			})
		}
	}
	if registry.GetConfig().GetSessionStoreUsed() {
		ss := registry.GetSessionStore()
		if ss != nil {
			group.Go(func(ctx context.Context) {
				if err := ss.SessionExpiryGoroutine(ctx); err != nil {
					slog.Error("Session expiry goroutine failed", "error", err)
				}
			})
		}
	}

	group.Go(func(ctx context.Context) {
		schedules.StartAsync(ctx, registry)
	})

	// Initialize email sender
	emails.InitEmailSender(registry)
	middlewares.CmsAddMiddlewares(registry) // Add CMS middlewares
	widgets.CmsAddShortcodes(registry)      // Add CMS shortcodes
	// cmsblocks.CmsAddBlockTypes(registry)    // Add CMS block types

	return nil
}
