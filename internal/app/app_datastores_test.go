package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"project/database/migrations"
	"project/internal/app"
	"project/internal/config"
)

// newTestApp creates a new Application with a unique in-memory SQLite DSN via cfg
func newTestApp(t *testing.T) app.AppInterface {
	t.Helper()
	cfg := config.New()
	cfg.SetAppEnv("testing")
	cfg.SetAppDebug(true)
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")
	cfg.SetDatabaseName(fmt.Sprintf("file:mp_test_%d?mode=memory&cache=shared", time.Now().UnixNano()))

	cfg.SetAuditStoreUsed(true)
	cfg.SetBlogStoreUsed(true)
	cfg.SetChatStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCustomStoreUsed(true)
	cfg.SetEntityStoreUsed(true)
	cfg.SetFeedStoreUsed(true)
	cfg.SetGeoStoreUsed(true)
	cfg.SetLogStoreUsed(true)
	cfg.SetMetaStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetSettingStoreUsed(true)
	cfg.SetShopStoreUsed(true)
	cfg.SetSqlFileStoreUsed(true)
	cfg.SetStatsStoreUsed(true)
	cfg.SetSubscriptionStoreUsed(true)
	cfg.SetTaskStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetVaultStoreUsed(true)

	a, err := app.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if err := migrations.MigrateAll(a); err != nil {
		t.Fatalf("expected no migration error, got: %v", err)
	}

	return a
}

func TestAppNew_InitializesStoresAndCreatesTables(t *testing.T) {
	a := newTestApp(t)
	ctx := context.Background()

	// Verify all stores are wired (non-nil)
	if a.IsDisabledAuditStore() {
		t.Fatal("AuditStore should not be nil")
	}
	if a.IsDisabledBlogStore() {
		t.Fatal("BlogStore should not be nil")
	}
	if a.IsDisabledChatStore() {
		t.Fatal("ChatStore should not be nil")
	}
	if a.IsDisabledCacheStore() {
		t.Fatal("CacheStore should not be nil")
	}
	if a.IsDisabledCmsStore() {
		t.Fatal("CmsStore should not be nil")
	}
	if a.IsDisabledCustomStore() {
		t.Fatal("CustomStore should not be nil")
	}
	if a.IsDisabledEntityStore() {
		t.Fatal("EntityStore should not be nil")
	}
	if a.IsDisabledFeedStore() {
		t.Fatal("FeedStore should not be nil")
	}
	if a.IsDisabledGeoStore() {
		t.Fatal("GeoStore should not be nil")
	}
	// MetaStore getter isn't exposed on AppInterface; table check below covers it
	if a.IsDisabledSessionStore() {
		t.Fatal("SessionStore should not be nil")
	}
	if a.IsDisabledShopStore() {
		t.Fatal("ShopStore should not be nil")
	}
	if a.IsDisabledSqlFileStorage() {
		t.Fatal("SqlFileStorage should not be nil")
	}
	if a.IsDisabledStatsStore() {
		t.Fatal("StatsStore should not be nil")
	}
	if a.IsDisabledTaskStore() {
		t.Fatal("TaskStore should not be nil")
	}
	if a.IsDisabledUserStore() {
		t.Fatal("UserStore should not be nil")
	}
	if a.IsDisabledVaultStore() {
		t.Fatal("VaultStore should not be nil")
	}
	if a.IsDisabledSubscriptionStore() {
		t.Fatal("SubscriptionStore should not be nil")
	}
	if a.IsDisabledBlindIndexStoreEmail() {
		t.Fatal("BlindIndexStoreEmail should not be nil")
	}
	if a.IsDisabledBlindIndexStoreFirstName() {
		t.Fatal("BlindIndexStoreFirstName should not be nil")
	}
	if a.IsDisabledBlindIndexStoreLastName() {
		t.Fatal("BlindIndexStoreLastName should not be nil")
	}

	// Verify some key tables exist
	mustHaveTables := []string{
		"snv_chat_chats",
		"snv_chat_messages",
		"snv_users_user",
		"snv_sessions_session",
		"snv_caches_cache",
		"snv_blogs_post",
		"snv_blogs_version",
		"snv_files_file",
		"snv_logs_log",
		"snv_metas_meta",
		"snv_stats_visitor",
		"snv_tasks_schedule",
		"snv_tasks_task_definition",
		"snv_tasks_task_queue",
		"snv_subscriptions_plan",
		"snv_subscriptions_subscription",
		"snv_vault_vault",
		"snv_bindx_email",
		"snv_bindx_first_name",
		"snv_bindx_last_name",
	}

	db := a.GetDatabase()
	if db == nil {
		t.Fatal("Database should not be nil")
	}
	for _, tbl := range mustHaveTables {
		t.Run("has_"+tbl, func(t *testing.T) {
			var name string
			err := db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", tbl).Scan(&name)
			if err != nil {
				t.Fatalf("expected table %s to exist, got error: %v", tbl, err)
			}
			if name != tbl {
				t.Fatalf("expected table name %s, got %s", tbl, name)
			}
		})
	}
}

func TestAppNew_IsIdempotent(t *testing.T) {
	a := newTestApp(t)

	// Second call should also succeed
	_, err := app.New(a.GetConfig())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
