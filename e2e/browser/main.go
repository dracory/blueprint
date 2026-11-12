// e2e-browser runs a sandboxed instance of the app on a dedicated port with a
// disposable database, intended to be driven by the Playwright suite in /e2e.
//
// Unlike cmd/ai-browser it uses the regular router (no auto-login middleware),
// so tests exercise the real authentication flows. A seeded user is always
// created and its credentials printed on startup.
//
// Never use in production — it refuses to start when APP_ENV=production.
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

	"project/cmd/server/background"
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
	var reset bool
	flag.StringVar(&withUser, "with-user", "", "Auto-seed a test user with the given email")
	flag.BoolVar(&isAdmin, "admin", false, "Assign administrator role to the auto-seeded user")
	flag.BoolVar(&reset, "reset", false, "Delete the e2e database before starting (clean state)")
	flag.Parse()

	setHardcodedEnv()

	// The disposable DB lives under tmp/, which `task dev` cleanup deletes.
	if err := os.MkdirAll("tmp", 0o755); err != nil {
		fmt.Printf("Failed to create tmp dir: %v\n", err)
		return
	}

	if reset {
		resetDatabase()
	}

	email := withUser
	if email == "" {
		email = e2eDefaultEmail
	}

	// Restrict the allowlist to the seeded email so the sandbox only ever
	// authenticates its own fixture user.
	_ = os.Setenv("AUTH_EMAILS_ALLOWED_ACCESS", email)

	cfg, err := config.NewFromEnv()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	if cfg.IsEnvProduction() {
		fmt.Println("ERROR: e2e-browser must not run in production environment. Exiting.")
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

	if err := seedUserAndSession(appInstance, email, isAdmin || withUser == ""); err != nil {
		fmt.Printf("Failed to seed user and session: %v\n", err)
		return
	}

	tasks.RegisterTasks(appInstance)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backgroundGroup := background.NewBackgroundGroup(ctx)
	if err := background.StartBackgroundProcesses(ctx, backgroundGroup, appInstance); err != nil {
		cfmt.Errorln("Failed to start background processes:", err.Error())
		cfmt.Infoln("Continuing without background processes...")
	}

	server, err := websrv.Start(websrv.Options{
		Host:    appInstance.GetConfig().GetAppHost(),
		Port:    appInstance.GetConfig().GetAppPort(),
		URL:     appInstance.GetConfig().GetAppUrl(),
		Handler: routes.Router(appInstance).ServeHTTP,
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

// resetDatabase deletes the disposable e2e SQLite files so the next start
// runs migrations and seeding from scratch.
func resetDatabase() {
	for _, f := range []string{"tmp/e2e-browser.db", "tmp/e2e-browser.db-shm", "tmp/e2e-browser.db-wal"} {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to remove %s: %v\n", f, err)
			return
		}
	}
	cfmt.Successln("Reset e2e database")
}

func setHardcodedEnv() {
	_ = os.Setenv("APP_NAME", "Blueprint")
	_ = os.Setenv("APP_URL", "http://127.0.0.1:34757")
	_ = os.Setenv("APP_HOST", "127.0.0.1")
	_ = os.Setenv("APP_PORT", "34757")
	_ = os.Setenv("APP_ENV", "development")
	_ = os.Setenv("APP_DEBUG", "true")

	_ = os.Setenv("DB_DRIVER", "sqlite")
	_ = os.Setenv("DB_DATABASE", "tmp/e2e-browser.db")
	_ = os.Setenv("DB_HOST", "")
	_ = os.Setenv("DB_PORT", "")
	_ = os.Setenv("DB_USERNAME", "")
	_ = os.Setenv("DB_PASSWORD", "")
	_ = os.Setenv("DB_SSL_MODE", "disable")
	_ = os.Setenv("DB_CHARSET", "utf8mb4")
	_ = os.Setenv("DB_TIMEZONE", "UTC")

	_ = os.Setenv("AUTH_REGISTRATION_ENABLED", "yes")
	// Password login so tests can sign in deterministically (no mail/OTP).
	_ = os.Setenv("AUTH_LOGIN_METHODS", "password")

	_ = os.Setenv("MAIL_DRIVER", "")
	_ = os.Setenv("MAIL_FROM_ADDRESS", "noreply@blueprint.local")
	_ = os.Setenv("MAIL_FROM_NAME", "Blueprint")
	_ = os.Setenv("MAIL_HOST", "")
	_ = os.Setenv("MAIL_PORT", "587")
	_ = os.Setenv("MAIL_USERNAME", "")
	_ = os.Setenv("MAIL_PASSWORD", "")

	_ = os.Setenv("TRANSLATION_LANGUAGE_DEFAULT", "en")

	_ = os.Setenv("SESSION_SECRET", "e2e-browser-session-secret-change-me")

	_ = os.Setenv("ENVENC_USED", "no")
	_ = os.Setenv("ENVENC_KEY_PRIVATE", "")

	_ = os.Setenv("VAULT_STORE_KEY", "e2e-browser-vault-key-32-chars!!")

	_ = os.Setenv("CMS_STORE_TEMPLATE_ID", "")

	_ = os.Setenv("ANTHROPIC_API_USED", "no")
	_ = os.Setenv("ANTHROPIC_API_KEY", "")
	_ = os.Setenv("ANTHROPIC_API_DEFAULT_MODEL", "")
	_ = os.Setenv("GEMINI_API_USED", "no")
	_ = os.Setenv("GEMINI_API_KEY", "")
	_ = os.Setenv("GEMINI_API_DEFAULT_MODEL", "")
	_ = os.Setenv("OPENAI_API_USED", "no")
	_ = os.Setenv("OPENAI_API_KEY", "")
	_ = os.Setenv("OPENAI_API_DEFAULT_MODEL", "")
	_ = os.Setenv("OPENROUTER_API_USED", "no")
	_ = os.Setenv("OPENROUTER_API_KEY", "")
	_ = os.Setenv("OPENROUTER_API_DEFAULT_MODEL", "")
	_ = os.Setenv("VERTEX_AI_API_USED", "no")
	_ = os.Setenv("VERTEX_AI_API_PROJECT_ID", "")
	_ = os.Setenv("VERTEX_AI_API_REGION_ID", "")
	_ = os.Setenv("VERTEX_AI_API_MODEL_ID", "")
	_ = os.Setenv("VERTEX_AI_API_DEFAULT_MODEL", "")

	_ = os.Setenv("STRIPE_KEY_PRIVATE", "")
	_ = os.Setenv("STRIPE_KEY_PUBLIC", "")

	_ = os.Setenv("MCP_API_KEY", "")
}
