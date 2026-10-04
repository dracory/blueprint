# Upgrade Guide: v0.49.0 to v0.50.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.49.0 to v0.50.0.

## Overview

This release adds a **database seeders framework** that mirrors the existing migrations pattern. A new `database/seeders` package provides a single `SeedAll(ctx, app)` entrypoint, called from `cmd/server/background_processes.go` and `cmd/ai-browser/background_processes.go` on boot. It distinguishes two seed types:

- **Sync seeds** (`database/seeders/sync/`) — canonical data whose Go definitions are the source of truth. Upserted on every boot; edits to the Go file propagate to the database on deploy.
- **Once seeds** (`database/seeders/once/`) — insert-if-absent bootstrap data. Each seed is a `migrator.MigrationInterface` with a `seed_YYYY_MM_DD_NNNN_<name>` signature, applied exactly once via a dedicated `migrator.Migrator` on the shared `migration_tracker` table. Afterwards the database owns the rows.

All seeders shipped in the template are **commented-out examples** — `SeedAll` is a no-op until you opt in. The `once/users` seeder creates well-known credentials (`admin@example.com` / `password`) and must only be enabled for local development.

The release also adds documentation: `database/migrations/README.md`, `database/models/README.md`, an example `database/models/note.go` ORM model for prototyping, and `AGENTS.md` files for `cmd/deploy`, `database`, `internal/config`, and `internal/testutils`.

**Key Changes:**
- New `database/seeders` package with `SeedAll(ctx, app) error` as the single entrypoint
- `SeedAll` is called from `startBackgroundProcesses` in both `cmd/server` and `cmd/ai-browser`, after store checks and before the task store startup
- Sync seed example: `database/seeders/sync/settings` (`settingsdata.Defaults()` / `settingsdata.Sync()` upserting `app.site_name`, `app.site_description`, `app.theme`)
- Once seed examples: `once/settings` (`InstalledAtSeed`, signature `seed_2026_10_04_0001_settings_installed_at`) and `once/users` (`UsersSeed`, signature `seed_2026_10_04_0002_users_default`), both registered through a `migrator.Migrator` with transactions disabled
- Once seeds gate on the `Get<Name>StoreUsed()` config flags and check `app.IsDisabled<Name>Store()` before running
- `database/migrations/README.md` documents migration signature conventions and the store-migration vs SQL-migration registration points
- `database/models/note.go` added as an example neat ORM model (prototype-to-store workflow, not for production use)
- No dependency changes in `go.mod`

---

## ⚠️ Breaking Changes

None. All changes are additive: the new `database/seeders` package is self-contained and `SeedAll` is a no-op unless seeds are enabled. The only integration point is the new `SeedAll` call in `startBackgroundProcesses`.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.50.0"
```

### Step 2: Copy the seeders package

Copy the new directory into your project:

```
database/seeders/
  seeders.go          entrypoint (SeedAll) — all shipped seeds commented out
  seeders_test.go
  README.md
  sync/README.md
  sync/settings/settings.go      example sync seed
  once/README.md
  once/settings/doc.go
  once/settings/settings_seeder.go   example once seed (installed_at)
  once/users/doc.go
  once/users/users.go            example account definitions
  once/users/users_seeder.go     example once seed (default users)
```

The package depends only on `project/internal/app`, `github.com/dracory/neat/database/migrator`, `github.com/dracory/settingstore`, `github.com/dracory/userstore`, `github.com/dracory/base/cfmt`, and `github.com/dromara/carbon/v2` — all already in `go.mod`. The commented imports in `seeders.go` require no action until you enable a seed.

### Step 3: Call `SeedAll` on boot

In `cmd/server/background_processes.go` (and `cmd/ai-browser/background_processes.go` if you use the AI browser sandbox), add inside `startBackgroundProcesses`, after the session-store check and before the task-store block:

```go
// Run all seeds (database/seeders.SeedAll): sync seeds upsert canonical
// data, once seeds insert-if-absent.
if err := seeders.SeedAll(ctx, app); err != nil {
    return fmt.Errorf("seed data: %w", err)
}
```

Add the import `"project/database/seeders"` and ensure `"fmt"` is imported.

### Step 4: Copy documentation files (optional)

- `database/migrations/README.md` — migration conventions
- `database/models/README.md`, `database/models/note.go` — example ORM model
- `AGENTS.md` files in `cmd/deploy/`, `database/`, `internal/config/`, `internal/testutils/` — agent-facing rules for those packages

### Step 5: Verify the build

```bash
go mod tidy
go build ./...
go test ./database/...
```

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...` (see `database/seeders/seeders_test.go` for coverage of the framework)
3. **No-op boot**: with all seeds commented out, `SeedAll` returns nil and startup is unchanged
4. **Once seed**: enable `settingsonce.NewInstalledAtSeed(app)` in `seeders.go`, boot, and confirm `app.installed_at` exists in the setting store and a `seed_2026_10_04_0001_settings_installed_at` row is recorded in `migration_tracker`; reboot and confirm it is not re-applied
5. **Sync seed**: enable `settingsdata.Sync`, change a value in `Defaults()`, reboot, and confirm the database reflects the new value

---

## 📝 Additional Notes

- **Signatures are immutable.** Once-seed signatures are recorded in `migration_tracker`; renaming a signature makes the migrator treat it as a new seed. Follow the `seed_YYYY_MM_DD_NNNN_<name>` convention.
- **Once seeds must be idempotent.** The migrator guards re-runs, but seeds are additionally written as find-or-create (see `findOrCreateUser`) so they are safe even if the tracker row is lost.
- **`Down()` is a no-op** on once seeds — seeded rows are left in place on rollback because they may have drifted or been referenced since.
- **Slow sync seeds** (many upserts) should run in a background goroutine inside `SeedAll` so they don't delay the server from accepting requests — a commented example is included in `seeders.go`.
- **Store flags matter**: seeds check `cfg.Get<Name>StoreUsed()` and `app.IsDisabled<Name>Store()`; a disabled store is an error at seed time, not a silent skip, once the seed is enabled.
- Do not enable `once/users` outside local development — it seeds `admin@example.com`/`password` and `user@example.com`/`password`. Change the credentials in `once/users/users.go` first if you use it.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: seeders.SeedAll`

**Cause**: The `database/seeders` directory was not copied, or `startBackgroundProcesses` was updated without the package.

**Solution**: Copy `database/seeders/` (Step 2) and add the `"project/database/seeders"` import.

### Issue: `seed data: user store is not initialized`

**Cause**: `usersonce.NewUsersSeed` was enabled while `USER_STORE_USED` is off or the store failed to initialize.

**Solution**: Enable and initialize the user store via the `GetUserStoreUsed()` config path, or keep the seed commented out.

### Issue: A once seed ran twice

**Cause**: The seed's `Signature()` was renamed between boots — the tracker treats it as a different seed.

**Solution**: Never change a deployed signature; write a new seed with a new `seed_` signature instead.

### Issue: Sync seed overwrote admin-edited settings

**Cause**: Sync seeds are upserts — the Go file always wins on every boot.

**Solution**: That is by design; if a row should drift in the DB, model it as a `once/` seed instead of `sync/`.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
