// Package settingsdata defines the canonical app settings — a sync seed.
// These definitions are the source of truth; Sync upserts them into the
// setting store on every boot, so edits here propagate on deploy.
// See database/seeders/sync/README.md for the contract.
package settingsdata

import (
	"context"

	"github.com/dracory/settingstore"
)

// Setting keys are stable identifiers stored in the database. They may be
// referenced from code, so they must never change once deployed.
const (
	KeySiteName        = "app.site_name"
	KeySiteDescription = "app.site_description"
	KeyTheme           = "app.theme"
)

// Defaults returns the canonical setting definitions — the code source of
// truth, as key/value pairs.
func Defaults() map[string]string {
	return map[string]string{
		KeySiteName:        "Blueprint",
		KeySiteDescription: "A Blueprint application",
		KeyTheme:           "default",
	}
}

// Sync upserts the canonical settings into the setting store. Existing
// keys keep their values only if they already match — definitions here
// always win.
func Sync(ctx context.Context, store settingstore.StoreInterface) error {
	for key, value := range Defaults() {
		current, err := store.Get(ctx, key, "")
		if err != nil {
			return err
		}
		if current == value {
			continue
		}
		if err := store.Set(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}
