package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestBlogStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetBlogStore() == nil {
		t.Error("expected blog store to be initialized")
	}
}

func TestBlogStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetBlogStore() != nil {
		t.Error("expected blog store to be nil when not used")
	}
}
