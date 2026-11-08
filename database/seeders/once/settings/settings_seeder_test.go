package settingsonce

import (
	"context"
	"testing"

	"project/internal/testutils"
)

func TestInstalledAtSeedSignature(t *testing.T) {
	seed := NewInstalledAtSeed(nil)
	if seed.Signature() != "seed_2026_10_04_0001_settings_installed_at" {
		t.Errorf("Signature() = %q", seed.Signature())
	}
	if seed.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestInstalledAtSeedNilApp(t *testing.T) {
	seed := NewInstalledAtSeed(nil)
	if err := seed.Up(); err == nil {
		t.Fatal("expected error with nil app")
	}
}

func TestInstalledAtSeedDisabledStore(t *testing.T) {
	app := testutils.Setup() // setting store disabled
	defer app.Close()

	seed := NewInstalledAtSeed(app)
	if err := seed.Up(); err == nil {
		t.Fatal("expected error when setting store is disabled")
	}
}

func TestInstalledAtSeedSetsKey(t *testing.T) {
	app := testutils.Setup(testutils.WithSettingStore(true))
	defer app.Close()

	seed := NewInstalledAtSeed(app)
	if err := seed.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	value, err := app.GetSettingStore().Get(context.Background(), KeyInstalledAt, "")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value == "" {
		t.Error("expected installed_at value to be set")
	}
}

func TestInstalledAtSeedDoesNotOverwrite(t *testing.T) {
	app := testutils.Setup(testutils.WithSettingStore(true))
	defer app.Close()

	ctx := context.Background()
	if err := app.GetSettingStore().Set(ctx, KeyInstalledAt, "2020-01-01 00:00:00"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	seed := NewInstalledAtSeed(app)
	if err := seed.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	value, _ := app.GetSettingStore().Get(ctx, KeyInstalledAt, "")
	if value != "2020-01-01 00:00:00" {
		t.Errorf("installed_at = %q, expected existing value to be preserved", value)
	}
}
