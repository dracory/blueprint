package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestSessionStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetSessionStore() == nil {
		t.Error("expected session store to be initialized")
	}
}

func TestSessionStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetSessionStore() != nil {
		t.Error("expected session store to be nil when not used")
	}
}
