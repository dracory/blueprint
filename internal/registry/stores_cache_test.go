package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestCacheStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetCacheStore() == nil {
		t.Error("expected cache store to be initialized")
	}
}

func TestCacheStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetCacheStore() != nil {
		t.Error("expected cache store to be nil when not used")
	}
}
