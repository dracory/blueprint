package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestGeoStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetGeoStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetGeoStore() == nil {
		t.Error("expected geo store to be initialized")
	}
}

func TestGeoStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetGeoStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetGeoStore() != nil {
		t.Error("expected geo store to be nil when not used")
	}
}
