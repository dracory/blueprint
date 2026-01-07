package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"project/internal/cli"
	"project/internal/config"
	"project/internal/emails"
	"project/internal/middlewares"
	"project/internal/registry"
	"project/internal/routes"
	"project/internal/schedules"
	"project/internal/tasks"
	"project/internal/widgets"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/taskstore"
	"github.com/dracory/websrv"
)

// main starts the registrylication
//
// Business Logic:
// 1. Initialize the environment
// 2. Defer Closing the database
// 3. Initialize the models
// 4. Register the task handlers
// 5. Executes the command if provided
// 6. Initialize the task queue
// 7. Initialize the scheduler
// 8. Starts the cache expiration goroutine
// 9. Starts the session expiration goroutine
// 10. Adds CMS shortcodes
// 11. Starts the web server
//
// Parameters:
// - none
//
// Returns:
// - none
func main() {
	// Set log flags to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Initialize registry (logger, caches, database)
	registry, err := registry.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize registry: %v\n", err)
		return
	}

	defer closeResourcesDB(registry.GetDatabase()) // Defer Closing the database

	tasks.RegisterTasks(registry) // Register the task handlers

	if isCliMode() {
		if len(os.Args) < 2 {
			return
		}
		if err := cli.ExecuteCliCommand(registry, os.Args[1:]); err != nil {
			fmt.Printf("Failed to execute CLI command: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Start background processes with explicit dependencies
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	background := newBackgroundGroup(ctx)
	if err := startBackgroundProcesses(ctx, background, registry); err != nil {
		cfmt.Errorln("Failed to start background processes:", err.Error())
		return
	}

	// Start the web server
	server, err := websrv.Start(websrv.Options{
		Host:    registry.GetConfig().GetAppHost(),
		Port:    registry.GetConfig().GetAppPort(),
		URL:     registry.GetConfig().GetAppUrl(),
		Handler: routes.Router(registry).ServeHTTP,
	})

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		background.stop()
		return
	}

	// Listen for OS signals to gracefully drain background work
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		fmt.Println("Shutdown signal received, draining background workers")
		cancel()
	case <-background.Done():
		cancel()
	}

	background.stop()
	if server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	}
}

// closeResources closes the database connection if it exists.
//
// Parameters:
// - db: the database handle
//
// Returns:
// - none
func closeResourcesDB(db *sql.DB) {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		cfmt.Errorf("Failed to close database connection: %v", err)
	}
}

// isCliMode checks if the registrylication is running in CLI mode.
//
// Parameters:
// - none
//
// Returns:
// - bool: true if the registrylication is running in CLI mode, false otherwise.
func isCliMode() bool {
	return len(os.Args) > 1
}

// startBackgroundProcesses starts the background processes.
//
// Parameters:
// - ctx: the context
// - group: the background group
// - registry: the registrylication
//
// Returns:
// - error: the error if any
func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, registry registry.RegistryInterface) error {
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

	return nil
}

type backgroundGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	doneCh chan struct{}
	once   sync.Once
}

func newBackgroundGroup(parent context.Context) *backgroundGroup {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &backgroundGroup{ctx: ctx, cancel: cancel, doneCh: make(chan struct{})}
}

func (g *backgroundGroup) Go(fn func(ctx context.Context)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fn(g.ctx)
	}()
}

func (g *backgroundGroup) stop() {
	g.once.Do(func() {
		g.cancel()
		g.wg.Wait()
		close(g.doneCh)
	})
}

func (g *backgroundGroup) Done() <-chan struct{} {
	return g.doneCh
}
