package main

import (
	"database/sql"
	"log/slog"
	"os"

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
		cfmt.Error("Failed to load config:", slog.Any("error", err))
		return
	}

	// Initialize application (logger, caches, database)
	application, err := app.New(cfg)
	if err != nil {
		cfmt.Error("Failed to initialize app:", slog.Any("error", err))
		return
	}

	defer closeResourcesDB(application.GetDB()) // Defer Closing the database

	tasks.RegisterTasks(application) // Register the task handlers

	if isCliMode() {
		if len(os.Args) < 2 {
			return
		}
		cli.ExecuteCliCommand(application, os.Args[1:]) // Execute the command
		return
	}

	// Start background processes with explicit dependencies
	startBackgroundProcesses(application)

	// Start the web server
	_, err = websrv.Start(websrv.Options{
		Host:    application.GetConfig().GetAppHost(),
		Port:    application.GetConfig().GetAppPort(),
		URL:     application.GetConfig().GetAppUrl(),
		Handler: routes.Routes(application).ServeHTTP,
	})

	if err != nil {
		cfmt.Errorf("Failed to start server: %v", err)
		return
	}
}

// closeResources closes the database connection if it exists.
//
// Parameters:
// - dbx: the database handle
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
// - db: the database handle
//
// Returns:
// - none
func startBackgroundProcesses(app types.AppInterface) {
	if app.GetDB() != nil {
		if ts := app.GetTaskStore(); ts != nil {
			go ts.QueueRunGoroutine(10, 2) // Initialize the task queue
		}
		if cs := app.GetCacheStore(); cs != nil {
			go cs.ExpireCacheGoroutine() // Initialize the cache expiration goroutine
		}
		if ss := app.GetSessionStore(); ss != nil {
			go ss.SessionExpiryGoroutine() // Initialize the session expiration goroutine
		}
	}

	schedules.StartAsync(app) // Initialize the scheduler

	// Initialize email sender
	emails.InitEmailSender()

	middlewares.CmsAddMiddlewares(app) // Add CMS middlewares
	widgets.CmsAddShortcodes(app)      // Add CMS shortcodes
}
