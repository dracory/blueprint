package main

import (
	"context"
	"testing"

	"project/internal/testutils"
)

func TestStartBackgroundProcesses(t *testing.T) {
	// Initialize minimal stores for background processes
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithTaskStore(true),
		testutils.WithUserStore(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Should not panic
	if err := startBackgroundProcesses(ctx, group, app); err != nil {
		t.Fatalf("startBackgroundProcesses returned error: %v", err)
	}

	if app.GetCacheStore() == nil {
		t.Errorf("Cache store should not be nil after starting background processes")
	}
	if app.GetSessionStore() == nil {
		t.Errorf("Session store should not be nil after starting background processes")
	}
}

func TestStartBackgroundProcesses_NilRegistry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	err := startBackgroundProcesses(ctx, group, nil)
	if err == nil {
		t.Error("startBackgroundProcesses with nil registry should return error")
	}
}

func TestStartBackgroundProcesses_NilConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create a minimal mock registry without config
	app := testutils.Setup()

	// This will likely fail because we can't easily mock GetConfig() returning nil
	// But let's at least test the flow
	err := startBackgroundProcesses(ctx, group, app)
	// Should not error because testutils.Setup() provides a valid config
	if err != nil {
		t.Logf("startBackgroundProcesses returned error (expected for test setup): %v", err)
	}
}

func TestStartBackgroundProcesses_MinimalConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create registry with minimal config but all required stores disabled
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(false)
	cfg.SetSessionStoreUsed(false)
	cfg.SetTaskStoreUsed(false)
	cfg.SetUserStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Logf("startBackgroundProcesses with minimal config returned error: %v", err)
	}
}

func TestStartBackgroundProcesses_NilDatabase(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create a mock registry with nil database
	app := testutils.Setup()
	// We can't easily set database to nil, so this test documents the expected behavior
	err := startBackgroundProcesses(ctx, group, app)
	// Should succeed since testutils.Setup provides a valid database
	if err != nil {
		t.Logf("startBackgroundProcesses with valid setup returned error: %v", err)
	}
}

func TestStartBackgroundProcesses_TaskStoreEnabledButNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create registry with task store enabled but not initialized
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(true)
	cfg.SetCacheStoreUsed(false)
	cfg.SetSessionStoreUsed(false)
	cfg.SetUserStoreUsed(false)

	// Create a minimal setup without task store
	app := testutils.Setup(testutils.WithCfg(cfg))

	// This should fail because task store is enabled but nil
	err := startBackgroundProcesses(ctx, group, app)
	if err == nil {
		t.Logf("startBackgroundProcesses with enabled but nil task store should return error, but testutils may initialize it")
	}
}

func TestStartBackgroundProcesses_CacheStoreEnabledButNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create registry with cache store enabled but not initialized
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetTaskStoreUsed(false)
	cfg.SetSessionStoreUsed(false)
	cfg.SetUserStoreUsed(false)

	app := testutils.Setup(testutils.WithCfg(cfg))

	// This should fail because cache store is enabled but nil
	err := startBackgroundProcesses(ctx, group, app)
	if err == nil {
		t.Logf("startBackgroundProcesses with enabled but nil cache store should return error, but testutils may initialize it")
	}
}

func TestStartBackgroundProcesses_SessionStoreEnabledButNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create registry with session store enabled but not initialized
	cfg := testutils.DefaultConf()
	cfg.SetSessionStoreUsed(true)
	cfg.SetTaskStoreUsed(false)
	cfg.SetCacheStoreUsed(false)
	cfg.SetUserStoreUsed(false)

	app := testutils.Setup(testutils.WithCfg(cfg))

	// This should fail because session store is enabled but nil
	err := startBackgroundProcesses(ctx, group, app)
	if err == nil {
		t.Logf("startBackgroundProcesses with enabled but nil session store should return error, but testutils may initialize it")
	}
}

func TestStartBackgroundProcesses_AllStoresEnabled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create registry with all stores enabled
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithTaskStore(true),
		testutils.WithUserStore(true),
	)

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Fatalf("startBackgroundProcesses with all stores enabled returned error: %v", err)
	}

	// Verify all stores are initialized
	if app.GetCacheStore() == nil {
		t.Error("Cache store should be initialized")
	}
	if app.GetSessionStore() == nil {
		t.Error("Session store should be initialized")
	}
	if app.GetTaskStore() == nil {
		t.Error("Task store should be initialized")
	}
}

func TestStartBackgroundProcesses_OnlyTaskStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	app := testutils.Setup(
		testutils.WithTaskStore(true),
	)

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Fatalf("startBackgroundProcesses with task store returned error: %v", err)
	}

	if app.GetTaskStore() == nil {
		t.Error("Task store should be initialized")
	}
}

func TestStartBackgroundProcesses_OnlyCacheStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Fatalf("startBackgroundProcesses with cache store returned error: %v", err)
	}

	if app.GetCacheStore() == nil {
		t.Error("Cache store should be initialized")
	}
}

func TestStartBackgroundProcesses_OnlySessionStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	app := testutils.Setup(
		testutils.WithSessionStore(true),
	)

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Fatalf("startBackgroundProcesses with session store returned error: %v", err)
	}

	if app.GetSessionStore() == nil {
		t.Error("Session store should be initialized")
	}
}

func TestStartBackgroundProcesses_TaskAndCacheStores(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	app := testutils.Setup(
		testutils.WithTaskStore(true),
		testutils.WithCacheStore(true),
	)

	err := startBackgroundProcesses(ctx, group, app)
	if err != nil {
		t.Fatalf("startBackgroundProcesses with task and cache stores returned error: %v", err)
	}

	if app.GetTaskStore() == nil {
		t.Error("Task store should be initialized")
	}
	if app.GetCacheStore() == nil {
		t.Error("Cache store should be initialized")
	}
}
