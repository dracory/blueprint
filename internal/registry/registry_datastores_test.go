package registry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	app "project/internal/registry"
	"project/internal/types"

	"github.com/stretchr/testify/require"
)

// newTestApp creates a new Application with a unique in-memory SQLite DSN via cfg
func newTestApp(t *testing.T) types.RegistryInterface {
	t.Helper()
	cfg := &types.Config{}
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
	require.NoError(t, err)
	return a
}

func TestAppNew_InitializesStoresAndCreatesTables(t *testing.T) {
	a := newTestApp(t)
	ctx := context.Background()

	// Verify all stores are wired (non-nil)
	require.NotNil(t, a.GetAuditStore())
	require.NotNil(t, a.GetBlogStore())
	require.NotNil(t, a.GetChatStore())
	require.NotNil(t, a.GetCacheStore())
	require.NotNil(t, a.GetCmsStore())
	require.NotNil(t, a.GetCustomStore())
	require.NotNil(t, a.GetEntityStore())
	require.NotNil(t, a.GetFeedStore())
	require.NotNil(t, a.GetGeoStore())
	// MetaStore getter isn't exposed on AppInterface; table check below covers it
	require.NotNil(t, a.GetSessionStore())
	require.NotNil(t, a.GetShopStore())
	require.NotNil(t, a.GetSqlFileStorage())
	require.NotNil(t, a.GetStatsStore())
	require.NotNil(t, a.GetTaskStore())
	require.NotNil(t, a.GetUserStore())
	require.NotNil(t, a.GetVaultStore())
	require.NotNil(t, a.GetSubscriptionStore())
	require.NotNil(t, a.GetBlindIndexStoreEmail())
	require.NotNil(t, a.GetBlindIndexStoreFirstName())
	require.NotNil(t, a.GetBlindIndexStoreLastName())

	// Verify some key tables exist
	mustHaveTables := []string{
		"snv_chat_chats",
		"snv_chat_messages",
		"snv_users_user",
		"snv_sessions_session",
		"snv_caches_cache",
		"snv_blogs_post",
		"snv_files_file",
		"snv_logs_log",
		"snv_metas_meta",
		"snv_stats_visitor",
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
	require.NotNil(t, db)
	for _, tbl := range mustHaveTables {
		t.Run("has_"+tbl, func(t *testing.T) {
			var name string
			err := db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", tbl).Scan(&name)
			require.NoError(t, err, "expected table %s to exist", tbl)
			require.Equal(t, tbl, name)
		})
	}
}

func TestAppNew_IsIdempotent(t *testing.T) {
	a := newTestApp(t)

	// Second call should also succeed
	_, err := app.New(a.GetConfig())
	require.NoError(t, err)
}
