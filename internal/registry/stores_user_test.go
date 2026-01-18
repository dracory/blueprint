package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestUserStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetUserStore() == nil {
		t.Error("expected user store to be initialized")
	}
}

func TestUserStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetUserStore() != nil {
		t.Error("expected user store to be nil when not used")
	}
}
