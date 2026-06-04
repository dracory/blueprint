package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestSettingStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetSettingStore() == nil {
		t.Error("expected setting store to be initialized")
	}
}

func TestSettingStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetSettingStore() != nil {
		t.Error("expected setting store to be nil when not used")
	}
}
