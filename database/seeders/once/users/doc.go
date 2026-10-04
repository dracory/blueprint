// Package usersonce is the once-seed for users:
//
//	Default accounts (an administrator and a regular user) are created
//	via find-or-create on first boot.
//
// Once-seeds implement migrator.MigrationInterface and run through a
// dedicated Migrator on the shared migration_tracker table — each seed
// runs exactly once, recorded under its seed_ signature (see
// database/seeders/once/README.md). Up is additionally written to be
// idempotent (find-or-create), so it is safe to re-run by hand.
package usersonce
