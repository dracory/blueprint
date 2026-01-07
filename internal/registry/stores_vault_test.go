package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestVaultStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetVaultStore() == nil {
		t.Fatal("expected vault store to be initialized, got nil")
	}
}

func TestVaultStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetVaultStore() != nil {
		t.Error("expected vault store to be nil when not used")
	}
}
