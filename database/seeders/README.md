# Seeds

Seed data — the code source of truth for rows that ship with the
application. Seeds are Go structs compiled into the binary and applied to
the database on boot via `seeders.SeedAll` (called from
`cmd/*/background_processes.go`).

Two kinds of seed data live here, distinguished by their **behavior
contract**:

| Directory  | Runs        | Existing row      | Source removed        | Source of truth |
|------------|-------------|-------------------|-----------------------|-----------------|
| `sync/`    | every boot  | updated to match  | deleted from DB       | the Go file     |
| `once/`    | first boot  | left untouched    | left untouched        | the DB          |

- **`sync/`** — canonical data. The Go definitions always win; the database
  is a materialization. Use for data the product owns (default settings,
  plans, vocabulary).
- **`once/`** — initial data. Inserted only when absent, never updated or
  removed afterwards. Use for bootstrap data (e.g. install timestamp, a
  default admin account) that is expected to drift in the database after
  creation.

Each subdirectory has its own README detailing the convention. See
`sync/README.md` and `once/README.md`.

## Layout

```
database/
  migrations/        schema changes
  seeders/
    sync/            canonical data, upserted every boot
      settings/      canonical app setting definitions
    once/            initial data, inserted only when absent
      settings/      records the install timestamp, once
```
