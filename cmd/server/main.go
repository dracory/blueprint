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

	"project/internal/cli"
	// "project/internal/cmsblocks"
	"project/internal/config"
	"project/internal/registry"
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

	// Initialize registry (logger, caches, database)
	registry, err := registry.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize registry: %v\n", err)
		return
	}
	defer func() {
		if err := registry.Close(); err != nil {
			cfmt.Errorf("Failed to close registry: %v", err)
		}
	}()

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
