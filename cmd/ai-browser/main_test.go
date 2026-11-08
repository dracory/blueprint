package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"project/internal/testutils"

	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

func TestSetHardcodedEnv(t *testing.T) {
	setHardcodedEnv()

	expected := map[string]string{
		"APP_NAME":     "Blueprint",
		"APP_URL":      "http://127.0.0.1:34756",
		"APP_HOST":     "127.0.0.1",
		"APP_PORT":     "34756",
		"APP_ENV":      "development",
		"APP_DEBUG":    "true",
		"DB_DRIVER":    "sqlite",
		"DB_DATABASE":  "tmp/ai-browser.db",
		"DB_SSL_MODE":  "disable",
		"MAIL_DRIVER":  "",
		"ENVENC_USED":  "no",
		"MCP_API_KEY":  "",
		"STRIPE_KEY_PRIVATE": "",
	}

	for key, want := range expected {
		if got := os.Getenv(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}

	// All external AI APIs must be disabled in the sandbox.
	for _, key := range []string{
		"ANTHROPIC_API_USED",
		"GEMINI_API_USED",
		"OPENAI_API_USED",
		"OPENROUTER_API_USED",
		"VERTEX_AI_API_USED",
	} {
		if got := os.Getenv(key); got != "no" {
			t.Errorf("%s = %q, want %q", key, got, "no")
		}
	}
}

func TestIsCliMode(t *testing.T) {
	original := os.Args
	defer func() { os.Args = original }()

	os.Args = []string{"ai-browser"}
	if isCliMode() {
		t.Error("isCliMode() = true, want false with no extra args")
	}

	os.Args = []string{"ai-browser", "-admin"}
	if !isCliMode() {
		t.Error("isCliMode() = false, want true with extra args")
	}
}

func TestBackgroundGroupRunsAndStops(t *testing.T) {
	group := newBackgroundGroup(context.Background())

	started := make(chan struct{})
	group.Go(func(ctx context.Context) {
		close(started)
		<-ctx.Done()
	})

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine did not start")
	}

	group.stop()

	select {
	case <-group.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("Done() channel was not closed after stop()")
	}
}

func TestBackgroundGroupNilParent(t *testing.T) {
	group := newBackgroundGroup(nil)
	if group == nil {
		t.Fatal("newBackgroundGroup(nil) returned nil")
	}
	group.stop()
}

func TestBackgroundGroupStopIsIdempotent(t *testing.T) {
	group := newBackgroundGroup(context.Background())
	group.stop()
	group.stop() // must not panic or deadlock
}

func TestStartBackgroundProcessesNilApp(t *testing.T) {
	group := newBackgroundGroup(context.Background())
	defer group.stop()

	err := startBackgroundProcesses(context.Background(), group, nil)
	if err == nil || !strings.Contains(err.Error(), "nil app") {
		t.Fatalf("expected nil app error, got %v", err)
	}
}

func TestStartBackgroundProcesses(t *testing.T) {
	app := testutils.Setup(
		testutils.WithSessionStore(true),
		testutils.WithTaskStore(true),
		testutils.WithCacheStore(true),
	)
	defer app.Close()

	group := newBackgroundGroup(context.Background())
	defer group.stop()

	if err := startBackgroundProcesses(context.Background(), group, app); err != nil {
		t.Fatalf("startBackgroundProcesses failed: %v", err)
	}
}

func TestSeedUserAndSessionCreatesUser(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
		testutils.WithSubscriptionStore(true),
	)
	defer app.Close()

	email := "test-seed@example.com"
	if err := seedUserAndSession(app, email, true); err != nil {
		t.Fatalf("seedUserAndSession failed: %v", err)
	}

	ctx := context.Background()

	user, err := app.GetUserStore().UserFindByEmail(ctx, email)
	if err != nil {
		t.Fatalf("UserFindByEmail failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected seeded user to exist")
	}
	if user.GetID() != aiBrowserUserID {
		t.Errorf("user ID = %q, want %q", user.GetID(), aiBrowserUserID)
	}

	roles, err := app.GetUserStore().UserRoles(ctx, user.GetID())
	if err != nil {
		t.Fatalf("UserRoles failed: %v", err)
	}
	foundAdmin := false
	for _, r := range roles {
		if r.GetHandle() == userstore.USER_ROLE_ADMINISTRATOR {
			foundAdmin = true
		}
	}
	if !foundAdmin {
		t.Error("expected administrator role to be assigned")
	}

	sessions, err := app.GetSessionStore().SessionList(ctx,
		sessionstore.NewSessionQuery().SetUserID(user.GetID()))
	if err != nil {
		t.Fatalf("SessionList failed: %v", err)
	}
	if len(sessions) == 0 {
		t.Fatal("expected at least one session for seeded user")
	}
	if sessions[0].IsExpired() {
		t.Error("expected seeded session to be unexpired")
	}
}

func TestSeedUserAndSessionIsIdempotent(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
	)
	defer app.Close()

	email := "test-idempotent@example.com"

	if err := seedUserAndSession(app, email, false); err != nil {
		t.Fatalf("first seedUserAndSession failed: %v", err)
	}
	if err := seedUserAndSession(app, email, false); err != nil {
		t.Fatalf("second seedUserAndSession failed: %v", err)
	}

	ctx := context.Background()
	user, err := app.GetUserStore().UserFindByEmail(ctx, email)
	if err != nil || user == nil {
		t.Fatalf("expected existing user to be found, err=%v", err)
	}

	sessions, err := app.GetSessionStore().SessionList(ctx,
		sessionstore.NewSessionQuery().SetUserID(user.GetID()))
	if err != nil {
		t.Fatalf("SessionList failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session (reused), got %d", len(sessions))
	}
}

func TestSeedUserAndSessionFailsWithoutStores(t *testing.T) {
	app := testutils.Setup() // all stores disabled
	defer app.Close()

	err := seedUserAndSession(app, "x@example.com", false)
	if err == nil {
		t.Fatal("expected error when user store is disabled")
	}
}
