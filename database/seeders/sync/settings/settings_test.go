package settingsdata

import (
	"context"
	"testing"

	"project/internal/testutils"
)

func TestSyncSeedsDefaults(t *testing.T) {
	app := testutils.Setup(testutils.WithSettingStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetSettingStore()

	if err := Sync(ctx, store); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	for key, want := range Defaults() {
		got, err := store.Get(ctx, key, "")
		if err != nil {
			t.Fatalf("Get(%s) failed: %v", key, err)
		}
		if got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestSyncUpdatesDriftedValues(t *testing.T) {
	app := testutils.Setup(testutils.WithSettingStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetSettingStore()

	if err := store.Set(ctx, KeySiteName, "Drifted Name"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := Sync(ctx, store); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	got, _ := store.Get(ctx, KeySiteName, "")
	if got != Defaults()[KeySiteName] {
		t.Errorf("%s = %q, want canonical %q", KeySiteName, got, Defaults()[KeySiteName])
	}
}

func TestSyncIsIdempotent(t *testing.T) {
	app := testutils.Setup(testutils.WithSettingStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetSettingStore()

	if err := Sync(ctx, store); err != nil {
		t.Fatalf("first Sync failed: %v", err)
	}
	if err := Sync(ctx, store); err != nil {
		t.Fatalf("second Sync failed: %v", err)
	}
}
