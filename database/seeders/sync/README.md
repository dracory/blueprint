# Sync Seeds

Canonical seed data: the Go files here are the **source of truth** and the
database is a downstream materialization.

## Contract

- Runs on **every boot** via an applier (e.g. `settingsdata.Sync`).
- Existing rows are **updated** to match the definitions (upsert).
- Rows whose source definition was removed may be **deleted**, depending on
  the applier.

## Adding sync seed data

1. Create `database/seeders/sync/<name>/` with the definitions as Go structs
   (or a function returning them, e.g. `Defaults()`).
2. Write an applier that upserts them into the appropriate store — it can
   live alongside the definitions or in `pkg/<name>/`.
3. Call the applier from `database/seeders/seeders.go` (`SeedAll`), behind the
   appropriate store-enabled config check.

## Current contents

- `settings/` — canonical app setting definitions (`settingsdata` package),
  applied by `settingsdata.Sync()` into the setting store.
