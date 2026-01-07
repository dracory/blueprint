package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestShopStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetShopStore() == nil {
		t.Error("expected shop store to be initialized")
	}
}

func TestShopStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetShopStore() != nil {
		t.Error("expected shop store to be nil when not used")
	}
}
