package registry_test

import (
	"project/internal/testutils"
	"testing"
)

// Log Store Tests
func TestLogStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetLogStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetLogStore() == nil {
		t.Error("expected log store to be initialized")
	}
}

func TestLogStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetLogStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetLogStore() != nil {
		t.Error("expected log store to be nil when not used")
	}
}
