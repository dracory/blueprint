package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestEntityStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetEntityStore() == nil {
		t.Error("expected entity store to be initialized")
	}
}

func TestEntityStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetEntityStore() != nil {
		t.Error("expected entity store to be nil when not used")
	}
}
