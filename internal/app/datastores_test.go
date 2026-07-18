package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestAuditStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetAuditStore() == nil {
		t.Error("expected audit store to be initialized")
	}
}

func TestAuditStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetAuditStore() != nil {
		t.Error("expected audit store to be nil when not used")
	}
}

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

func TestCacheStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetCacheStore() == nil {
		t.Error("expected cache store to be initialized")
	}
}

func TestCacheStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetCacheStore() != nil {
		t.Error("expected cache store to be nil when not used")
	}
}

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

func TestCustomStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetCustomStore() == nil {
		t.Error("expected custom store to be initialized")
	}
}

func TestCustomStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetCustomStore() != nil {
		t.Error("expected custom store to be nil when not used")
	}
}

func TestEntityStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetEntityStore() == nil {
		t.Error("expected entity store to be initialized")
	}
}

func TestEntityStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetEntityStore() != nil {
		t.Error("expected entity store to be nil when not used")
	}
}

func TestFeedStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetFeedStore() == nil {
		t.Error("expected feed store to be initialized")
	}
}

func TestFeedStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetFeedStore() != nil {
		t.Error("expected feed store to be nil when not used")
	}
}

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetLogStore() != nil {
		t.Error("expected log store to be nil when not used")
	}
}

func TestMetaStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetMetaStore() == nil {
		t.Error("expected meta store to be initialized")
	}
}

func TestMetaStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetMetaStore() != nil {
		t.Error("expected meta store to be nil when not used")
	}
}

func TestSessionStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetSessionStore() == nil {
		t.Error("expected session store to be initialized")
	}
}

func TestSessionStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetSessionStore() != nil {
		t.Error("expected session store to be nil when not used")
	}
}

func TestSettingStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetSettingStore() == nil {
		t.Error("expected setting store to be initialized")
	}
}

func TestSettingStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetSettingStore() != nil {
		t.Error("expected setting store to be nil when not used")
	}
}

func TestShopStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetShopStore() == nil {
		t.Error("expected shop store to be initialized")
	}
}

func TestShopStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetShopStore() != nil {
		t.Error("expected shop store to be nil when not used")
	}
}

func TestStatsStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetStatsStore() == nil {
		t.Error("expected stats store to be initialized")
	}
}

func TestStatsStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetStatsStore() != nil {
		t.Error("expected stats store to be nil when not used")
	}
}

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

func TestUserStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetUserStore() == nil {
		t.Error("expected user store to be initialized")
	}
}

func TestUserStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetUserStore() != nil {
		t.Error("expected user store to be nil when not used")
	}
}

func TestVaultStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.GetVaultStore() == nil {
		t.Fatal("expected vault store to be initialized, got nil")
	}
}

func TestVaultStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.GetVaultStore() != nil {
		t.Error("expected vault store to be nil when not used")
	}
}
