package app_test

import (
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	registrypkg "project/internal/registry"
	"project/internal/types"
)

func TestNew_SetsDefaultLogger(t *testing.T) {
	cfg := &types.Config{}
	cfg.SetAppEnv("testing")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(fmt.Sprintf("file:mp_test_%d?mode=memory&cache=shared", time.Now().UnixNano()))

	registry, err := registrypkg.New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}

	if registry.GetLogger() == nil {
		t.Fatalf("expected registry logger to be non-nil right after app.New; got nil")
	}
}
