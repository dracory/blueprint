# database/

## Layout

- `migrations/` — schema changes, run on every boot via `migrations.MigrateAll`
- `seeders/` — seed data (`sync/` upserted every boot, `once/` run once via the migration tracker)
- `models/` — example neat ORM models for prototyping only

## Rules

- **Never change a deployed signature.** Migration and once-seed signatures are recorded in `migration_tracker`; renaming skips/re-runs them.
- **Add migrations** by copying `migrations/2026_03_22_0001_table_custom_create.go`, naming per `YYYY_MM_DD_NNNN_<name>`, and registering in `getSQLMigrations()`.
- **Add seeders** in `seeders/sync/` (canonical data, upserted) or `seeders/once/` (bootstrap data, `seed_`-prefixed signature, idempotent find-or-create), registered in `seeders.go`.
- **All seeders shipped here are commented-out examples** — do not enable them in `seeders.go` without being asked. `once/users` seeds well-known credentials: dev-only, never in production.
- Prefer the dracory store APIs over `models/`; `models/` is for throwaway prototyping, not production code.
- Run `go test ./database/...` after changes; mirror the existing `*_test.go` patterns.

See `migrations/README.md`, `seeders/README.md`, `models/README.md` for full conventions.
