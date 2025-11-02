package app_test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"

	apppkg "project/internal/app"
	"project/internal/types"
)

func TestNew_SetsDefaultLogger(t *testing.T) {
	cfg := &types.Config{}
	cfg.SetAppEnv("testing")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(fmt.Sprintf("file:mp_test_%d?mode=memory&cache=shared", time.Now().UnixNano()))

	application, err := apppkg.New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}

	if application.GetLogger() == nil {
		t.Fatalf("expected application logger to be non-nil right after app.New; got nil")
	}
}
