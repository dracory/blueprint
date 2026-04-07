package registry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"project/internal/config"
	"project/internal/registry"
)

// newTestApp creates a new Application with a unique in-memory SQLite DSN via cfg
func newTestApp(t *testing.T) registry.RegistryInterface {
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

	a, err := registry.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	return a
}

func TestAppNew_InitializesStoresAndCreatesTables(t *testing.T) {
	a := newTestApp(t)
	ctx := context.Background()

	// Verify all stores are wired (non-nil)
	if a.GetAuditStore() == nil {
		t.Fatal("AuditStore should not be nil")
	}
	if a.GetBlogStore() == nil {
		t.Fatal("BlogStore should not be nil")
	}
	if a.GetChatStore() == nil {
		t.Fatal("ChatStore should not be nil")
	}
	if a.GetCacheStore() == nil {
		t.Fatal("CacheStore should not be nil")
	}
	if a.GetCmsStore() == nil {
		t.Fatal("CmsStore should not be nil")
	}
	if a.GetCustomStore() == nil {
		t.Fatal("CustomStore should not be nil")
	}
	if a.GetEntityStore() == nil {
		t.Fatal("EntityStore should not be nil")
	}
	if a.GetFeedStore() == nil {
		t.Fatal("FeedStore should not be nil")
	}
	if a.GetGeoStore() == nil {
		t.Fatal("GeoStore should not be nil")
	}
	// MetaStore getter isn't exposed on AppInterface; table check below covers it
	if a.GetSessionStore() == nil {
		t.Fatal("SessionStore should not be nil")
	}
	if a.GetShopStore() == nil {
		t.Fatal("ShopStore should not be nil")
	}
	if a.GetSqlFileStorage() == nil {
		t.Fatal("SqlFileStorage should not be nil")
	}
	if a.GetStatsStore() == nil {
		t.Fatal("StatsStore should not be nil")
	}
	if a.GetTaskStore() == nil {
		t.Fatal("TaskStore should not be nil")
	}
	if a.GetUserStore() == nil {
		t.Fatal("UserStore should not be nil")
	}
	if a.GetVaultStore() == nil {
		t.Fatal("VaultStore should not be nil")
	}
	if a.GetSubscriptionStore() == nil {
		t.Fatal("SubscriptionStore should not be nil")
	}
	if a.GetBlindIndexStoreEmail() == nil {
		t.Fatal("BlindIndexStoreEmail should not be nil")
	}
	if a.GetBlindIndexStoreFirstName() == nil {
		t.Fatal("BlindIndexStoreFirstName should not be nil")
	}
	if a.GetBlindIndexStoreLastName() == nil {
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
	_, err := registry.New(a.GetConfig())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
