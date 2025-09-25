package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestMetaStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetMetaStore() == nil {
		t.Error("expected meta store to be initialized")
	}
}

func TestMetaStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetMetaStore() != nil {
		t.Error("expected meta store to be nil when not used")
	}
}
