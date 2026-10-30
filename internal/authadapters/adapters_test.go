package authadapters

import (
	"context"
	"testing"
	"time"

	"project/internal/testutils"
)

// TestTempKeyAdapter exercises Get/Set/Delete against the test cachestore.
func TestTempKeyAdapter(t *testing.T) {
	application := testutils.Setup(testutils.WithCacheStore(true))
	adapter := &TempKeyAdapter{Store: application.GetCacheStore()}

	if err := adapter.Set("k1", "v1", 60); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, err := adapter.Get("k1")
	if err != nil || v != "v1" {
		t.Fatalf("Get = (%q, %v), want (v1, nil)", v, err)
	}
	if err := adapter.Delete("k1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := adapter.Get("k1"); err == nil {
		t.Fatal("Get after delete expected error (missing key), got nil")
	}
}

// TestTempKeyAdapter_NilStore verifies all methods error on nil store.
func TestTempKeyAdapter_NilStore(t *testing.T) {
	adapter := &TempKeyAdapter{}
	if _, err := adapter.Get("k"); err == nil {
		t.Error("Get expected error on nil store")
	}
	if err := adapter.Set("k", "v", 1); err == nil {
		t.Error("Set expected error on nil store")
	}
	if err := adapter.Delete("k"); err == nil {
		t.Error("Delete expected error on nil store")
	}
}

// TestSessionStoreAdapter exercises create/find/delete round-trips plus the
// value-name mapping (remember -> session_value, ip -> session_ip).
func TestSessionStoreAdapter(t *testing.T) {
	application := testutils.Setup(testutils.WithSessionStore(true))
	adapter := &SessionStoreAdapter{Store: application.GetSessionStore()}
	ctx := context.Background()

	s := adapter.NewSession()
	s.SetUserID("u1")
	s.SetKey("test-key")
	s.SetExpiresAt(time.Now().Add(time.Hour))
	s.SetValue(ValueRemember, "remember")
	s.SetValue(valueIP, "127.0.0.1")
	s.SetValue(valueUA, "test-agent")

	if err := adapter.SessionCreate(ctx, s); err != nil {
		t.Fatalf("SessionCreate: %v", err)
	}

	found, err := adapter.SessionFindByKey(ctx, "test-key")
	if err != nil || found == nil {
		t.Fatalf("SessionFindByKey = (%v, %v)", found, err)
	}
	if found.GetUserID() != "u1" {
		t.Errorf("GetUserID = %q, want u1", found.GetUserID())
	}
	if got := found.GetValue(ValueRemember); got != "remember" {
		t.Errorf("GetValue(remember) = %q, want remember", got)
	}
	if got := found.GetValue(valueIP); got != "127.0.0.1" {
		t.Errorf("GetValue(ip) = %q, want 127.0.0.1", got)
	}
	if got := found.GetValue(valueUA); got != "test-agent" {
		t.Errorf("GetValue(user_agent) = %q, want test-agent", got)
	}

	if err := adapter.SessionDeleteByKey(ctx, "test-key"); err != nil {
		t.Fatalf("SessionDeleteByKey: %v", err)
	}
	found, err = adapter.SessionFindByKey(ctx, "test-key")
	if err != nil || found != nil {
		t.Fatalf("find after delete = (%v, %v), want (nil, nil)", found, err)
	}
}

// TestSessionStoreAdapter_DeleteByUserID removes all sessions for a user.
func TestSessionStoreAdapter_DeleteByUserID(t *testing.T) {
	application := testutils.Setup(testutils.WithSessionStore(true))
	adapter := &SessionStoreAdapter{Store: application.GetSessionStore()}
	ctx := context.Background()

	for _, key := range []string{"k-a", "k-b"} {
		s := adapter.NewSession()
		s.SetUserID("u9")
		s.SetKey(key)
		s.SetExpiresAt(time.Now().Add(time.Hour))
		if err := adapter.SessionCreate(ctx, s); err != nil {
			t.Fatalf("SessionCreate %s: %v", key, err)
		}
	}
	// Another user's session must survive.
	other := adapter.NewSession()
	other.SetUserID("u2")
	other.SetKey("k-c")
	other.SetExpiresAt(time.Now().Add(time.Hour))
	if err := adapter.SessionCreate(ctx, other); err != nil {
		t.Fatalf("SessionCreate k-c: %v", err)
	}

	if err := adapter.SessionDeleteByUserID(ctx, "u9"); err != nil {
		t.Fatalf("SessionDeleteByUserID: %v", err)
	}
	for _, key := range []string{"k-a", "k-b"} {
		found, _ := adapter.SessionFindByKey(ctx, key)
		if found != nil {
			t.Errorf("key %s still findable after user delete", key)
		}
	}
	found, _ := adapter.SessionFindByKey(ctx, "k-c")
	if found == nil {
		t.Error("other user's session wrongly deleted")
	}
}

// TestBuildPasswordlessAuth verifies the builder returns a wired impl
// without panicking on a fully-stored test app.
func TestBuildPasswordlessAuth(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetTaskStoreUsed(true)
	application := testutils.Setup(testutils.WithCfg(cfg))

	a, err := BuildPasswordlessAuth(application)
	if err != nil {
		t.Fatalf("BuildPasswordlessAuth: %v", err)
	}
	if a == nil {
		t.Fatal("BuildPasswordlessAuth returned nil")
	}
	if !a.IsPasswordless() {
		t.Error("expected passwordless auth impl")
	}
}

// TestSecureFlag matches the environment-conditional Secure cookie flag.
func TestSecureFlag(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetAppEnv("development") // insecure HTTP → "no"
	application := testutils.Setup(testutils.WithCfg(cfg))
	if got := secureFlag(application); got != "no" {
		t.Errorf("secureFlag(dev) = %q, want no", got)
	}
}
