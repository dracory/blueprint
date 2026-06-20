package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/database/migrations"
	"project/internal/app"
	"project/internal/config"
	"project/internal/routes"
	"project/internal/tasks"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/websrv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var withUser string
	var isAdmin bool
	flag.StringVar(&withUser, "with-user", "", "Auto-seed a test user with the given email")
	flag.BoolVar(&isAdmin, "admin", false, "Assign administrator role to the auto-seeded user")
	flag.Parse()

	setHardcodedEnv()

	cfg, err := config.NewFromEnv()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	if cfg.IsEnvProduction() {
		fmt.Println("ERROR: AI browser must not run in production environment. Exiting.")
		return
	}

	appInstance, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	defer func() {
		if err := appInstance.Close(); err != nil {
			cfmt.Errorf("Failed to close app: %v", err)
		}
	}()

	if err := migrations.MigrateAll(appInstance); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	email := withUser
	if email == "" {
		email = aiBrowserDefaultEmail
	}

	if err := seedUserAndSession(appInstance, email, isAdmin || withUser == ""); err != nil {
		fmt.Printf("Failed to seed user and session: %v\n", err)
		return
	}

	tasks.RegisterTasks(appInstance)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	background := newBackgroundGroup(ctx)
	if err := startBackgroundProcesses(ctx, background, appInstance); err != nil {
		cfmt.Errorln("Background processes failed to start:", err.Error())
		cfmt.Infoln("Continuing without background processes...")
	}

	server, err := websrv.Start(websrv.Options{
		Host:    appInstance.GetConfig().GetAppHost(),
		Port:    appInstance.GetConfig().GetAppPort(),
		URL:     appInstance.GetConfig().GetAppUrl(),
		Handler: routes.AiBrowserRouter(appInstance).ServeHTTP,
	})

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		background.stop()
		return
	}

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

func setHardcodedEnv() {
	os.Setenv("APP_NAME", "Blueprint")
	os.Setenv("APP_URL", "http://127.0.0.1:34756")
	os.Setenv("APP_HOST", "127.0.0.1")
	os.Setenv("APP_PORT", "34756")
	os.Setenv("APP_ENV", "development")
	os.Setenv("APP_DEBUG", "true")

	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_DATABASE", "tmp/ai-browser.db")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_PORT", "")
	os.Setenv("DB_USERNAME", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("DB_CHARSET", "utf8mb4")
	os.Setenv("DB_TIMEZONE", "UTC")

	os.Setenv("AUTH_REGISTRATION_ENABLED", "yes")

	os.Setenv("MAIL_DRIVER", "")
	os.Setenv("MAIL_FROM_ADDRESS", "noreply@blueprint.local")
	os.Setenv("MAIL_FROM_NAME", "Blueprint")
	os.Setenv("MAIL_HOST", "")
	os.Setenv("MAIL_PORT", "587")
	os.Setenv("MAIL_USERNAME", "")
	os.Setenv("MAIL_PASSWORD", "")

	os.Setenv("TRANSLATION_LANGUAGE_DEFAULT", "en")

	os.Setenv("SESSION_SECRET", "ai-browser-session-secret-change-me")

	os.Setenv("ENVENC_USED", "no")
	os.Setenv("ENVENC_KEY_PRIVATE", "")

	os.Setenv("VAULT_STORE_KEY", "ai-browser-vault-key-32-chars-long!!")

	os.Setenv("CMS_STORE_TEMPLATE_ID", "")

	os.Setenv("ANTHROPIC_API_USED", "no")
	os.Setenv("ANTHROPIC_API_KEY", "")
	os.Setenv("ANTHROPIC_API_DEFAULT_MODEL", "")
	os.Setenv("GEMINI_API_USED", "no")
	os.Setenv("GEMINI_API_KEY", "")
	os.Setenv("GEMINI_API_DEFAULT_MODEL", "")
	os.Setenv("OPENAI_API_USED", "no")
	os.Setenv("OPENAI_API_KEY", "")
	os.Setenv("OPENAI_API_DEFAULT_MODEL", "")
	os.Setenv("OPENROUTER_API_USED", "no")
	os.Setenv("OPENROUTER_API_KEY", "")
	os.Setenv("OPENROUTER_API_DEFAULT_MODEL", "")
	os.Setenv("VERTEX_AI_API_USED", "no")
	os.Setenv("VERTEX_AI_API_PROJECT_ID", "")
	os.Setenv("VERTEX_AI_API_REGION_ID", "")
	os.Setenv("VERTEX_AI_API_MODEL_ID", "")
	os.Setenv("VERTEX_AI_API_DEFAULT_MODEL", "")

	os.Setenv("STRIPE_KEY_PRIVATE", "")
	os.Setenv("STRIPE_KEY_PUBLIC", "")

	os.Setenv("MCP_API_KEY", "")
}
