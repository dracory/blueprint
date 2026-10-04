// Package settingsonce is the once-seed for settings:
//
//	The install timestamp (app.installed_at) is recorded in the setting
//	store the first time the app boots.
//
// Once-seeds implement migrator.MigrationInterface and run through a
// dedicated Migrator on the shared migration_tracker table — each seed
// runs exactly once, recorded under its seed_ signature (see
// database/seeders/once/README.md). SeedAll (database/seeders) wires them up
// after the sync seeds.
package settingsonce
