package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestSubscriptionStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSubscriptionStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetSubscriptionStore() == nil {
		t.Error("expected subscription store to be initialized")
	}
}

func TestSubscriptionStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSubscriptionStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	if registry.GetSubscriptionStore() != nil {
		t.Error("expected subscription store to be nil when not used")
	}
}
