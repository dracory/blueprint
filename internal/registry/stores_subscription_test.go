package app_test

import (
    "project/internal/testutils"
    "testing"
)

func TestSubscriptionStoreInitialize_Success(t *testing.T) {
    cfg := testutils.DefaultConf()
    cfg.SetSubscriptionStoreUsed(true)
    app := testutils.Setup(testutils.WithCfg(cfg))

    if app.GetSubscriptionStore() == nil {
        t.Error("expected subscription store to be initialized")
    }
}

func TestSubscriptionStoreInitialize_NotUsed(t *testing.T) {
    cfg := testutils.DefaultConf()
    cfg.SetSubscriptionStoreUsed(false)
    app := testutils.Setup(testutils.WithCfg(cfg))

    if app.GetSubscriptionStore() != nil {
        t.Error("expected subscription store to be nil when not used")
    }
}
