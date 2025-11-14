package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"project/internal/app"
	"project/internal/cli"
	"project/internal/config"
	"project/internal/emails"
	"project/internal/middlewares"
	"project/internal/routes"
	"project/internal/schedules"
	"project/internal/tasks"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/base/cfmt"

	"github.com/dracory/websrv"
)

// main starts the application
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
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Initialize application (logger, caches, database)
	application, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	defer closeResourcesDB(application.GetDB()) // Defer Closing the database

	tasks.RegisterTasks(application) // Register the task handlers

	if isCliMode() {
		if len(os.Args) < 2 {
			return
		}
		if err := cli.ExecuteCliCommand(application, os.Args[1:]); err != nil {
			fmt.Printf("Failed to execute CLI command: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Start background processes with explicit dependencies
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	background := newBackgroundGroup(ctx)
	startBackgroundProcesses(ctx, background, application)

	// Start the web server
	server, err := websrv.Start(websrv.Options{
		Host:    application.GetConfig().GetAppHost(),
		Port:    application.GetConfig().GetAppPort(),
		URL:     application.GetConfig().GetAppUrl(),
		Handler: routes.Routes(application).ServeHTTP,
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

// isCliMode checks if the application is running in CLI mode.
//
// Parameters:
// - none
//
// Returns:
// - bool: true if the application is running in CLI mode, false otherwise.
func isCliMode() bool {
	return len(os.Args) > 1
}

// startBackgroundProcesses starts the background processes.
//
// Parameters:
// - ctx: the context
// - group: the background group
// - app: the application
//
// Returns:
// - error: the error if any
func startBackgroundProcesses(ctx context.Context, group *backgroundGroup, app types.AppInterface) error {
	if app == nil {
		return errors.New("startBackgroundProcesses called with nil app")
	}

	if app.GetDB() == nil {
		return errors.New("startBackgroundProcesses called with nil db")
	}

	if app.GetTaskStore() == nil {
		return errors.New("startBackgroundProcesses called with nil task store")
	}

	if app.GetCacheStore() == nil {
		return errors.New("startBackgroundProcesses called with nil cache store")
	}

	if app.GetSessionStore() == nil {
		return errors.New("startBackgroundProcesses called with nil session store")
	}

	if ts := app.GetTaskStore(); ts != nil {
		group.Go(func(ctx context.Context) {
			ts.QueueRunGoroutine(ctx, 10, 2)
		})
	}
	if cs := app.GetCacheStore(); cs != nil {
		group.Go(func(ctx context.Context) {
			if err := cs.ExpireCacheGoroutine(ctx); err != nil {
				slog.Error("Cache expiration goroutine failed", "error", err)
			}
		})
	}
	if ss := app.GetSessionStore(); ss != nil {
		group.Go(func(ctx context.Context) {
			if err := ss.SessionExpiryGoroutine(ctx); err != nil {
				slog.Error("Session expiry goroutine failed", "error", err)
			}
		})
	}

	group.Go(func(ctx context.Context) {
		schedules.StartAsync(ctx, app)
	})

	// Initialize email sender
	emails.InitEmailSender(app)
	middlewares.CmsAddMiddlewares(app) // Add CMS middlewares
	widgets.CmsAddShortcodes(app)      // Add CMS shortcodes

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
