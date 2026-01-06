package app_test

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetUserStore() != nil {
		t.Error("expected user store to be nil when not used")
	}
}
