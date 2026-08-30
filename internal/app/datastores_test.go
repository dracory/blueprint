package app_test

import (
	"project/internal/testutils"
	"testing"
)

func TestAuditStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledAuditStore() {
		t.Error("expected audit store to be initialized")
	}
}

func TestAuditStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAuditStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledAuditStore() {
		t.Error("expected audit store to be nil when not used")
	}
}

func TestBlogStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledBlogStore() {
		t.Error("expected blog store to be initialized")
	}
}

func TestBlogStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledBlogStore() {
		t.Error("expected blog store to be nil when not used")
	}
}

func TestCacheStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledCacheStore() {
		t.Error("expected cache store to be initialized")
	}
}

func TestCacheStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledCacheStore() {
		t.Error("expected cache store to be nil when not used")
	}
}

func TestChatStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetChatStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsDisabledChatStore() {
		t.Error("expected chat store to be initialized")
	}
}

func TestChatStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetChatStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledChatStore() {
		t.Error("expected chat store to be nil when not used")
	}
}

func TestCmsStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCmsStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledCmsStore() {
		t.Error("expected CMS store to be initialized")
	}
}

func TestCmsStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCmsStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledCmsStore() {
		t.Error("expected CMS store to be nil when not used")
	}
}

func TestCustomStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledCustomStore() {
		t.Error("expected custom store to be initialized")
	}
}

func TestCustomStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCustomStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledCustomStore() {
		t.Error("expected custom store to be nil when not used")
	}
}

func TestEntityStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledEntityStore() {
		t.Error("expected entity store to be initialized")
	}
}

func TestEntityStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEntityStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledEntityStore() {
		t.Error("expected entity store to be nil when not used")
	}
}

func TestFeedStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledFeedStore() {
		t.Error("expected feed store to be initialized")
	}
}

func TestFeedStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetFeedStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledFeedStore() {
		t.Error("expected feed store to be nil when not used")
	}
}

func TestGeoStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetGeoStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledGeoStore() {
		t.Error("expected geo store to be initialized")
	}
}

func TestGeoStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetGeoStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledGeoStore() {
		t.Error("expected geo store to be nil when not used")
	}
}

func TestLogStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetLogStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledLogStore() {
		t.Error("expected log store to be initialized")
	}
}

func TestLogStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetLogStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledLogStore() {
		t.Error("expected log store to be nil when not used")
	}
}

func TestMetaStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledMetaStore() {
		t.Error("expected meta store to be initialized")
	}
}

func TestMetaStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetMetaStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledMetaStore() {
		t.Error("expected meta store to be nil when not used")
	}
}

func TestSessionStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledSessionStore() {
		t.Error("expected session store to be initialized")
	}
}

func TestSessionStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledSessionStore() {
		t.Error("expected session store to be nil when not used")
	}
}

func TestSettingStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledSettingStore() {
		t.Error("expected setting store to be initialized")
	}
}

func TestSettingStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSettingStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledSettingStore() {
		t.Error("expected setting store to be nil when not used")
	}
}

func TestShopStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledShopStore() {
		t.Error("expected shop store to be initialized")
	}
}

func TestShopStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetShopStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledShopStore() {
		t.Error("expected shop store to be nil when not used")
	}
}

func TestStatsStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledStatsStore() {
		t.Error("expected stats store to be initialized")
	}
}

func TestStatsStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetStatsStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledStatsStore() {
		t.Error("expected stats store to be nil when not used")
	}
}

func TestSubscriptionStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSubscriptionStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsDisabledSubscriptionStore() {
		t.Error("expected subscription store to be initialized")
	}
}

func TestSubscriptionStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetSubscriptionStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledSubscriptionStore() {
		t.Error("expected subscription store to be nil when not used")
	}
}

func TestTaskStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledTaskStore() {
		t.Error("expected task store to be initialized")
	}
}

func TestTaskStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledTaskStore() {
		t.Error("expected task store to be nil when not used")
	}
}

func TestUserStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledUserStore() {
		t.Error("expected user store to be initialized")
	}
}

func TestUserStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledUserStore() {
		t.Error("expected user store to be nil when not used")
	}
}

func TestVaultStoreInitialize_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(true)
	a := testutils.Setup(testutils.WithCfg(cfg))

	if a.IsDisabledVaultStore() {
		t.Fatal("expected vault store to be initialized, got nil")
	}
}

func TestVaultStoreInitialize_NotUsed(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetVaultStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	if app.IsEnabledVaultStore() {
		t.Error("expected vault store to be nil when not used")
	}
}
