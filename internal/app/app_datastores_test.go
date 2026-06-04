package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"project/database/migrations"
	"project/internal/config"
	"project/internal/app"
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
	db := a.GetDatabase()
	if db == nil {
		t.Fatal("Database should not be nil")
	}

	var name string
	err := db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_chat_chats").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_chat_chats to exist, got error: %v", err)
	}
	if name != "snv_chat_chats" {
		t.Fatalf("expected table name snv_chat_chats, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_chat_messages").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_chat_messages to exist, got error: %v", err)
	}
	if name != "snv_chat_messages" {
		t.Fatalf("expected table name snv_chat_messages, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_users_user").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_users_user to exist, got error: %v", err)
	}
	if name != "snv_users_user" {
		t.Fatalf("expected table name snv_users_user, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_sessions_session").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_sessions_session to exist, got error: %v", err)
	}
	if name != "snv_sessions_session" {
		t.Fatalf("expected table name snv_sessions_session, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_caches_cache").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_caches_cache to exist, got error: %v", err)
	}
	if name != "snv_caches_cache" {
		t.Fatalf("expected table name snv_caches_cache, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_blogs_post").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_blogs_post to exist, got error: %v", err)
	}
	if name != "snv_blogs_post" {
		t.Fatalf("expected table name snv_blogs_post, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_blogs_version").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_blogs_version to exist, got error: %v", err)
	}
	if name != "snv_blogs_version" {
		t.Fatalf("expected table name snv_blogs_version, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_files_file").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_files_file to exist, got error: %v", err)
	}
	if name != "snv_files_file" {
		t.Fatalf("expected table name snv_files_file, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_logs_log").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_logs_log to exist, got error: %v", err)
	}
	if name != "snv_logs_log" {
		t.Fatalf("expected table name snv_logs_log, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_metas_meta").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_metas_meta to exist, got error: %v", err)
	}
	if name != "snv_metas_meta" {
		t.Fatalf("expected table name snv_metas_meta, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_stats_visitor").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_stats_visitor to exist, got error: %v", err)
	}
	if name != "snv_stats_visitor" {
		t.Fatalf("expected table name snv_stats_visitor, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_tasks_schedule").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_tasks_schedule to exist, got error: %v", err)
	}
	if name != "snv_tasks_schedule" {
		t.Fatalf("expected table name snv_tasks_schedule, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_tasks_task_definition").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_tasks_task_definition to exist, got error: %v", err)
	}
	if name != "snv_tasks_task_definition" {
		t.Fatalf("expected table name snv_tasks_task_definition, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_tasks_task_queue").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_tasks_task_queue to exist, got error: %v", err)
	}
	if name != "snv_tasks_task_queue" {
		t.Fatalf("expected table name snv_tasks_task_queue, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_subscriptions_plan").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_subscriptions_plan to exist, got error: %v", err)
	}
	if name != "snv_subscriptions_plan" {
		t.Fatalf("expected table name snv_subscriptions_plan, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_subscriptions_subscription").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_subscriptions_subscription to exist, got error: %v", err)
	}
	if name != "snv_subscriptions_subscription" {
		t.Fatalf("expected table name snv_subscriptions_subscription, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_vault_vault").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_vault_vault to exist, got error: %v", err)
	}
	if name != "snv_vault_vault" {
		t.Fatalf("expected table name snv_vault_vault, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_bindx_email").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_bindx_email to exist, got error: %v", err)
	}
	if name != "snv_bindx_email" {
		t.Fatalf("expected table name snv_bindx_email, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_bindx_first_name").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_bindx_first_name to exist, got error: %v", err)
	}
	if name != "snv_bindx_first_name" {
		t.Fatalf("expected table name snv_bindx_first_name, got %s", name)
	}

	err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", "snv_bindx_last_name").Scan(&name)
	if err != nil {
		t.Fatalf("expected table snv_bindx_last_name to exist, got error: %v", err)
	}
	if name != "snv_bindx_last_name" {
		t.Fatalf("expected table name snv_bindx_last_name, got %s", name)
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
