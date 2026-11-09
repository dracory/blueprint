package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/cmd/server/background"
	"project/cmd/server/helpers"
	"project/database/migrations"
	"project/internal/cli"
	"project/internal/app"
	"project/internal/config"
	"project/internal/routes"
	"project/internal/tasks"

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
	// Set log flags to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.NewFromEnv()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Log version on startup
	fmt.Printf("Starting %s v%s\n", cfg.GetAppName(), config.GetVersion())

	// Initialize app (logger, caches, database)
	app, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}
	defer func() {
		if err := app.Close(); err != nil {
			cfmt.Errorf("Failed to close app: %v", err)
		}
	}()

	// Run all migrations (store-level + custom SQL)
	if err := migrations.MigrateAll(app); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		return
	}

	tasks.RegisterTasks(app) // Register the task handlers

	if helpers.IsCliMode() {
		if len(os.Args) < 2 {
			return
		}
		if err := cli.ExecuteCliCommand(app, os.Args[1:]); err != nil {
			fmt.Printf("Failed to execute CLI command: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Start background processes with explicit dependencies
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backgroundGroup := background.NewBackgroundGroup(ctx)
	if err := background.StartBackgroundProcesses(ctx, backgroundGroup, app); err != nil {
		cfmt.Errorln("Failed to start background processes:", err.Error())
		return
	}

	// Start the web server
	server, err := websrv.Start(websrv.Options{
		Host:    app.GetConfig().GetAppHost(),
		Port:    app.GetConfig().GetAppPort(),
		URL:     app.GetConfig().GetAppUrl(),
		Handler: routes.Router(app).ServeHTTP,
	})

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		backgroundGroup.Stop()
		return
	}

	// Listen for OS signals to gracefully drain background work
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		fmt.Println("Shutdown signal received, draining background workers")
		cancel()
	case <-backgroundGroup.Done():
		cancel()
	}

	backgroundGroup.Stop()
	if server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	}
}
