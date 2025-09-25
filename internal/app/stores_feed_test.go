package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestFeedStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetFeedStore() == nil {
		t.Error("expected feed store to be initialized")
	}
}

func TestFeedStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetFeedStore() != nil {
		t.Error("expected feed store to be nil when not used")
	}
}
