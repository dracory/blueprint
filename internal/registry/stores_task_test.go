package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestTaskStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetTaskStore() == nil {
		t.Error("expected task store to be initialized")
	}
}

func TestTaskStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetTaskStore() != nil {
		t.Error("expected task store to be nil when not used")
	}
}
