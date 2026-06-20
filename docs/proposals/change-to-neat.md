# Proposal: Migrate Database Connection Management to `dracory/neat`

## Status

Completed

## Context

The project currently uses two separate packages for database concerns:

1. `github.com/dracory/database` — for opening and configuring `*sql.DB` connections
2. `github.com/dracory/neat` — for database migrations and ORM-style operations

The `neat` package already contains a full database connection layer with support for multiple named connections, read/write replicas, connection pooling, and a Laravel-style configuration model. Continuing to depend on `dracory/database` adds an unnecessary package and prevents us from using `neat`'s more advanced features.

## Goal

Replace `github.com/dracory/database` with `github.com/dracory/neat` for database connection management while:

- Keeping `*sql.DB` compatibility for all existing stores
- Adding Laravel-like multi-connection configuration support
- Unifying connection, migration, and ORM concerns under one package
- Maintaining backward compatibility with existing `.env` variables

## Key Design Decisions

1. **Keep `*sql.DB` compatibility** — All existing stores continue to receive `*sql.DB`, so no store constructors need to change.
2. **Expose the `neat.Database` instance** — Add it to `AppInterface` so migrations and future ORM code can use it directly.
3. **Laravel-style multi-connection configuration** — Add optional connection configuration to the config layer and `.env`.
4. **Single connection by default** — Current behavior is preserved; multi-connection is opt-in.

## Proposed Configuration Model

### `.env` (existing keys remain)

```env
DB_DRIVER=sqlite
DB_DATABASE=database.db
DB_HOST=
DB_PORT=
DB_USERNAME=
DB_PASSWORD=
DB_SSL_MODE=require
DB_CHARSET=utf8mb4
DB_TIMEZONE=UTC
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME_SECONDS=300
DB_CONN_MAX_IDLE_TIME_SECONDS=5
```

### New optional `.env` keys

```env
# Default connection name (default: "default")
DB_DEFAULT_CONNECTION=default

# Direct DSN override (optional)
DB_DSN=

# Table prefix (optional)
DB_PREFIX=

# Future multi-connection support
# DB_CONNECTIONS='{"default":{"driver":"sqlite","database":"database.db"}}'
```

### Internal config representation

Add `DatabaseConnectionConfigInterface` and extend `DatabaseConfigInterface`:

```go
type DatabaseConnectionConfigInterface interface {
    GetName() string
    GetDriver() string
    GetHost() string
    GetPort() string
    GetDatabase() string
    GetUsername() string
    GetPassword() string
    GetSSLMode() string
    GetCharset() string
    GetTimezone() string
    GetDSN() string
    GetPrefix() string
}

type DatabaseConfigInterface interface {
    // existing single-db getters/setters...

    SetDatabaseDefaultConnection(string)
    GetDatabaseDefaultConnection() string

    GetDatabaseConnections() []DatabaseConnectionConfigInterface
    GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface
}
```

The existing single-db getters (`GetDatabaseDriver`, `GetDatabaseHost`, etc.) will delegate to the default connection for backward compatibility.

### Mapping to `neat.DBConfig`

A new helper `internal/config/database_neat_config.go` will map the blueprint config to `neat/database/db.DBConfig`:

```go
func databaseNeatConfig(cfg ConfigInterface) neatdb.DBConfig {
    // Build default connection from cfg
    // Set Pool config from cfg
    // Return neat DBConfig
}
```

## Proposed App Changes

### `internal/app/database_open.go`

Replace `dracory/database.Open()` with `neat.New()`:

```go
func databaseOpen(cfg config.ConfigInterface) (*neatdatabase.Database, error) {
    if cfg == nil {
        return nil, errors.New("databaseOpen: cfg is nil")
    }

    neatCfg := databaseNeatConfig(cfg)
    return neatdatabase.New(neatCfg)
}
```

### `internal/app/app_interface.go`

Add neat accessors while keeping existing `*sql.DB` methods:

```go
type AppInterface interface {
    // ... existing methods ...

    GetDatabase() *sql.DB
    SetDatabase(db *sql.DB)

    GetNeatDatabase() *neatdatabase.Database
    SetNeatDatabase(db *neatdatabase.Database)

    GetDatabaseConnection(name string) *sql.DB
}
```

### `internal/app/app_implementation.go`

- Add `neatDB *neatdatabase.Database` field
- In `New()`: open neat database, derive `*sql.DB` for stores
- Update `Close()` to close `neatDB`

### `database/migrations/migrate.go`

Replace `neatdb.NewFromSQLDB(app.GetDatabase())` with `app.GetNeatDatabase()`.

## Files to Modify

| File | Change |
|------|--------|
| `internal/config/z_config_interfaces.go` | Add `DatabaseConnectionConfigInterface`, extend `DatabaseConfigInterface` |
| `internal/config/z_config_constants.go` | Add new env keys |
| `internal/config/database_config.go` | Parse connection settings into default connection |
| `internal/config/z_config_implementation.go` | Store connection configs, implement new getters |
| `internal/config/database_neat_config.go` | New: map config to `neat.DBConfig` |
| `internal/app/app_interface.go` | Add neat database accessors |
| `internal/app/app_implementation.go` | Store neat DB, update `New()`/`Close()` |
| `internal/app/database_open.go` | Replace with neat database opener |
| `database/migrations/migrate.go` | Use neat DB directly |
| `go.mod` | Remove `github.com/dracory/database` |
| `internal/app/app_datastores_test.go` | Update test setup |
| `internal/app/app_close_test.go` | Add neat DB nil tests |
| `internal/config/database_neat_config_test.go` | New test file |
| `.env.example` | New env vars |
| `docs/upgrade_guides/*` | Migration notes |
| `docs/comparisons/goravel-vs-blueprint-comparison.md` | Update references |

## Migration Plan

### Phase 1: Configuration Layer
1. Add new interfaces and env keys
2. Parse existing single DB config into a default connection
3. Implement `databaseNeatConfig` mapper

### Phase 2: App Layer
4. Update `AppInterface` to expose neat database
5. Refactor `database_open.go`
6. Update `app_implementation.go` to own neat database lifecycle

### Phase 3: Migrations
7. Update `migrate.go` to use `app.GetNeatDatabase()`

### Phase 4: Dependencies
8. Remove `dracory/database` from `go.mod`
9. Run `go mod tidy`

### Phase 5: Tests
10. Update existing tests
11. Add new tests for `databaseNeatConfig`

### Phase 6: Documentation
12. Update `.env.example`
13. Add upgrade guide notes

## Risks and Considerations

- **Driver naming**: Ensure driver names map correctly (e.g., `sqlite` stays `sqlite`).
- **SQLite pooling**: Current config forces `MaxOpenConns=1` for SQLite; must preserve this.
- **SSL default**: Current code defaults to `require` for non-SQLite; preserve in neat config.
- **Port type**: Current config stores port as string; neat expects int.
- **Close semantics**: `neat.Database.Close()` closes underlying `*sql.DB`; avoid double-close.
- **Backward compatibility**: Existing `.env` keys remain unchanged.

## Verification

- `go test ./...` passes
- `go run ./cmd/server` starts successfully
- `go mod tidy` produces no unexpected changes
- SQLite in-memory and file modes work
- PostgreSQL/MySQL connections work
- Migrations run successfully via `neat`
- All stores initialize without errors

## Decision

Approve this proposal to proceed with the migration.