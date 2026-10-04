# Migrations

Schema changes, run by `migrations.MigrateAll(app)` on every boot (called
from `cmd/server/main.go`, `cmd/ai-browser/main.go`, and test setup before
the app serves requests).

All migrations implement `migrator.MigrationInterface` (`Signature`,
`Description`, `Up`, `Down` — embed `migrator.BaseMigration`) and run
through the neat migrator, which records applied signatures in the shared
`migration_tracker` table and skips them on subsequent runs.

## Two kinds of migrations

| Kind            | Registered in                    | Purpose                                  |
|-----------------|----------------------------------|------------------------------------------|
| Store migrations | `getStoreMigrations()`          | Create the tables a dracory store needs (`userstore`, `settingstore`, ...), gated by the matching `Get<Name>StoreUsed()` config flag |
| SQL migrations   | `getSQLMigrations()`            | Custom schema changes — create/alter your own tables via `m.GetSchema()` blueprint calls |

## Signature convention

`YYYY_MM_DD_NNNN_<name>` — e.g. `2026_03_22_0001_table_custom_create`.
Filenames mirror the signature. Signatures are stable identifiers recorded
in `migration_tracker`; never change one once deployed.

## Adding a migration

1. Copy `2026_03_22_0001_table_custom_create.go` as a template — it shows
   the `Blueprint` schema DSL (`ID`, `String`, `Unique`, `Timestamps`,
   `HasTable`/`DropIfExists` guards).
2. Name the file and signature per the convention above.
3. Register it in `getSQLMigrations()` in `migrate.go`.
4. Add a `*_test.go` alongside — the existing `*_migrate_test.go` files
   show the pattern (in-memory SQLite via `testutils`).

## See also

- `database/seeders/` — seed *data* (sync/once), applied after migrations.
  Once-seeds reuse the same `migration_tracker` mechanism with `seed_`-prefixed signatures.
