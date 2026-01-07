package registry_test

import (
	"project/internal/testutils"
	"testing"
)

func TestCmsStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCmsStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetCmsStore() == nil {
		t.Error("expected CMS store to be initialized")
	}
}

func TestCmsStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCmsStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetCmsStore() != nil {
		t.Error("expected CMS store to be nil when not used")
	}
}
