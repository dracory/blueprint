package registry_test

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
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetTaskStore() != nil {
		t.Error("expected task store to be nil when not used")
	}
}
