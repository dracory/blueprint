package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestTradingStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTradingStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetTradingStore() == nil {
		t.Error("expected trading store to be initialized")
	}
}

func TestTradingStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTradingStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetTradingStore() != nil {
		t.Error("expected trading store to be nil when not used")
	}
}
