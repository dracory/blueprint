package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestCustomStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetCustomStore() == nil {
		t.Error("expected custom store to be initialized")
	}
}

func TestCustomStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetCustomStore() != nil {
		t.Error("expected custom store to be nil when not used")
	}
}
