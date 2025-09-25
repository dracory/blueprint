package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestStatsStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetStatsStore() == nil {
		t.Error("expected stats store to be initialized")
	}
}

func TestStatsStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetStatsStore() != nil {
		t.Error("expected stats store to be nil when not used")
	}
}
