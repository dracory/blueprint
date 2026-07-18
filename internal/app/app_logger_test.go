package app_test

import (
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"project/internal/app"
	"project/internal/config"
)

func TestNew_SetsDefaultLogger(t *testing.T) {
	cfg := config.New()
	cfg.SetAppEnv("testing")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(fmt.Sprintf("file:mp_test_%d?mode=memory&cache=shared", time.Now().UnixNano()))

	app, err := app.New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}

	if app.GetLogger() == nil {
		t.Fatalf("expected app logger to be non-nil right after app.New; got nil")
	}
}
