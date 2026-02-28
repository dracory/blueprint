package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestAuditStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(true) // Only enable audit store for this test
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetAuditStore() == nil {
		t.Error("expected audit store to be initialized")
	}
}

func TestAuditStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetAuditStore() != nil {
		t.Error("expected audit store to be nil when not used")
	}
}
