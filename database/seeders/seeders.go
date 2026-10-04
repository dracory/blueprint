// Package seeders is the single entrypoint for all seed data, mirroring
// database/migrations.MigrateAll. SeedAll runs every registered seed in a
// defined order:
//
//  1. Sync seeds (database/seeders/sync/) — canonical data, upserted every
//     boot; the Go definitions are the source of truth.
//  2. Once seeds (database/seeders/once/) — insert-if-absent data; the DB
//     owns it afterwards.
//
// Fast seeds run synchronously so their data exists before requests are
// served. Slow seeds should run in a background goroutine so they don't
// delay startup.
package seeders

import (
	"context"
	"errors"
	"fmt"
	settingsonce "project/database/seeders/once/settings"

	// usersonce "project/database/seeders/once/users" // enable with UsersSeed below
	settingsdata "project/database/seeders/sync/settings"
	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

// SeedAll runs all seeds for the given app. Safe to call on every boot.
func SeedAll(ctx context.Context, app app.AppInterface) error {
	if app == nil {
		return errors.New("app is nil")
	}
	cfg := app.GetConfig()
	if cfg == nil {
		return errors.New("config is nil")
	}

	// --- Sync seeds (fast, blocking) ---

	// Canonical app settings: a handful of upserts.
	if cfg.GetSettingStoreUsed() {
		if err := settingsdata.Sync(ctx, app.GetSettingStore()); err != nil {
			return fmt.Errorf("sync settings: %w", err)
		}
	}

	// --- Once seeds (fast, blocking) ---

	// Once seeds run through their own migrator on the shared
	// migration_tracker table — each is applied exactly once, recorded
	// under its seed_ signature. They run after the sync seeds so
	// canonical data they depend on already exists.
	onceSeeds := []migrator.MigrationInterface{}

	// Record the install timestamp on first boot.
	if cfg.GetSettingStoreUsed() {
		onceSeeds = append(onceSeeds, settingsonce.NewInstalledAtSeed(app))
	}

	// Default user accounts — DISABLED by default: seeding well-known
	// credentials is a security risk if forgotten in production. Enable
	// for local development only, and change the passwords in
	// once/users/users.go first.
	// if cfg.GetUserStoreUsed() {
	// 	onceSeeds = append(onceSeeds, usersonce.NewUsersSeed(app))
	// }

	if len(onceSeeds) > 0 {
		m := migrator.NewMigrator(app.GetNeatDatabase())
		m.SetTransactionsEnabled(false)
		if err := m.AddMigrations(onceSeeds); err != nil {
			return fmt.Errorf("add once seeds: %w", err)
		}
		if err := m.Up(ctx); err != nil {
			return fmt.Errorf("run once seeds: %w", err)
		}
	}

	// --- Sync seeds (slow, background) ---

	// Slow seeds (many upserts) should run in a goroutine so they don't
	// delay the web server from accepting requests, e.g.:
	//
	//	go func() {
	//		if err := slowdata.Sync(ctx); err != nil {
	//			cfmt.Errorln("Failed to sync slow data:", err.Error())
	//		}
	//	}()

	return nil
}
