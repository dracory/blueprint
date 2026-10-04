# Once Seeds

Initial seed data: work that runs **exactly once**, like a migration.
After it is applied, the database owns the data — edits here do not
propagate.

## How "once" is enforced

Once seeds implement `migrator.MigrationInterface` and run through a
dedicated `Migrator` (in `database/seeders/seeders.go`) on the shared
`migration_tracker` table. Each seed is recorded under a `seed_` signature
and skipped on subsequent boots — same mechanism as schema migrations,
but they don't interfere with them: the main migrator only acts on its
own registered migration list.

## Contract

- Runs on boot via `seeders.SeedAll`, after the sync seeds — so canonical
  data a seed depends on already exists.
- Applied exactly once per `seed_` signature, tracked in
  `migration_tracker`.
- `Down()` is usually a no-op — seeded rows may have drifted since.

## When to use

- Bootstrap data that must exist on a fresh install (install timestamp,
  default admin account, required settings).
- One-time backfills of existing rows.

If the data should stay locked to the code definitions on every deploy,
it belongs in `sync/` instead — not here.

## Adding a once seed

1. Create `database/seeders/once/<name>/` with a struct implementing
   `migrator.MigrationInterface` (`Signature` prefixed `seed_`, `Up`,
   `Down` — embed `migrator.BaseMigration`, keep the app in a private
   field, and expose a `New<Name>Seed(app app.AppInterface)` constructor).
2. Make `Up` idempotent anyway — check before insert (find-or-create
   style), so a manually deleted tracker row or a partial failure can't
   produce duplicates.
3. Register it in `database/seeders/seeders.go` behind the appropriate
   store-enabled config checks.

## Current contents

- `settings/` — `InstalledAtSeed`: records the install timestamp in the
  setting store the first time the app boots.
