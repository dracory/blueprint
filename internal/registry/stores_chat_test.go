package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestChatStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetChatStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetChatStore() == nil {
		t.Error("expected chat store to be initialized")
	}
}

func TestChatStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetChatStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetChatStore() != nil {
		t.Error("expected chat store to be nil when not used")
	}
}
