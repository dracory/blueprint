package main

import (
	"fmt"
	"log"

	"project/app/config"
	"project/app/routes"
	"project/internal/platform/database"

	"github.com/dracory/base/server"

	_ "modernc.org/sqlite"
)

func main() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to the database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	err = db.MigrateDatabase()
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// Setup routes
	r := routes.SetupRoutes(cfg, db)

	// Start the server
	fmt.Println("Starting Dracory web server...")
	fmt.Printf("Server running at http://%s:%s\n", cfg.ServerHost, cfg.ServerPort)

	// Use the server package from the base library
	serverOptions := server.Options{
		Host:    cfg.ServerHost,
		Port:    cfg.ServerPort,
		URL:     cfg.AppURL,
		Handler: r.ServeHTTP,
		Mode:    server.ProductionMode,
	}

	if cfg.IsDevelopment() {
		serverOptions.Mode = "development"
	}

	_, err = server.Start(serverOptions)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
