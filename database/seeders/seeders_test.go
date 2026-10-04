package seeders

import (
	"context"
	settingsonce "project/database/seeders/once/settings"
	usersonce "project/database/seeders/once/users"
	settingsdata "project/database/seeders/sync/settings"
	"project/internal/testutils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSeedAll_NilApp(t *testing.T) {
	err := SeedAll(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestSeedAll_Success(t *testing.T) {
	a := testutils.Setup(
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)
	defer a.Close()

	if err := SeedAll(context.Background(), a); err != nil {
		t.Fatalf("SeedAll failed: %v", err)
	}
}

func TestSettingsSync_Idempotent(t *testing.T) {
	a := testutils.Setup(testutils.WithSettingStore(true))
	defer a.Close()

	ctx := context.Background()
	store := a.GetSettingStore()

	if err := settingsdata.Sync(ctx, store); err != nil {
		t.Fatalf("first Sync failed: %v", err)
	}
	if err := settingsdata.Sync(ctx, store); err != nil {
		t.Fatalf("second Sync should be idempotent: %v", err)
	}

	for key, want := range settingsdata.Defaults() {
		got, err := store.Get(ctx, key, "")
		if err != nil {
			t.Fatalf("Get(%q) failed: %v", key, err)
		}
		if got != want {
			t.Errorf("expected %q = %q, got %q", key, want, got)
		}
	}
}

func TestInstalledAtSeed_Signature(t *testing.T) {
	s := &settingsonce.InstalledAtSeed{}
	if s.Signature() != "seed_2026_10_04_0001_settings_installed_at" {
		t.Errorf("unexpected signature %q", s.Signature())
	}
}

func TestInstalledAtSeed_UpNilApp(t *testing.T) {
	s := &settingsonce.InstalledAtSeed{}
	if err := s.Up(); err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestInstalledAtSeed_UpSuccess(t *testing.T) {
	a := testutils.Setup(testutils.WithSettingStore(true))
	defer a.Close()

	s := settingsonce.NewInstalledAtSeed(a)
	if err := s.Up(); err != nil {
		t.Fatalf("InstalledAtSeed.Up failed: %v", err)
	}

	got, err := a.GetSettingStore().Get(context.Background(), settingsonce.KeyInstalledAt, "")
	if err != nil {
		t.Fatalf("Get(%q) failed: %v", settingsonce.KeyInstalledAt, err)
	}
	if got == "" {
		t.Error("expected install timestamp to be set")
	}
}

func TestUsersSeed_Signature(t *testing.T) {
	s := &usersonce.UsersSeed{}
	if s.Signature() != "seed_2026_10_04_0002_users_default" {
		t.Errorf("unexpected signature %q", s.Signature())
	}
}

func TestUsersSeed_UpNilApp(t *testing.T) {
	s := &usersonce.UsersSeed{}
	if err := s.Up(); err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestUsersSeed_UpIdempotent(t *testing.T) {
	a := testutils.Setup(testutils.WithUserStore(true))
	defer a.Close()

	s := usersonce.NewUsersSeed(a)
	if err := s.Up(); err != nil {
		t.Fatalf("first Up failed: %v", err)
	}
	if err := s.Up(); err != nil {
		t.Fatalf("second Up should be idempotent: %v", err)
	}

	user, err := a.GetUserStore().UserFindByEmail(context.Background(), "admin@example.com")
	if err != nil {
		t.Fatalf("UserFindByEmail failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected admin@example.com to be seeded")
	}
}
