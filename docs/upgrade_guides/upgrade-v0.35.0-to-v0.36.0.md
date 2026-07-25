# Upgrade Guide: v0.35.0 to v0.36.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.35.0 to v0.36.0.

## Overview

This release introduces two major features: **Maintenance Mode** and **Turso Database Support**. It also includes updates to 15 dracory store dependencies.

**Maintenance Mode** — A file-based, CLI-toggleable maintenance mode that returns `503 Service Unavailable` for all requests except excluded IPs and paths. The middleware is added as the **first** item in the global middleware chain so it short-circuits all processing when active.

**Turso Database Support** — Turso (libSQL) is now a supported database driver alongside SQLite, MySQL, and PostgreSQL. Turso is treated as SQLite-compatible — same connection pool constraints, same GORM dialector, same SSL behavior. Connections use `DB_DSN` rather than host/port/user/password.

**Key Changes:**
- New `maintenance` CLI command (`go run ./cmd/server maintenance enable|disable|status`) with aliases `on`/`down`/`off`/`up`
- New `MaintenanceMiddleware` added as the first global middleware
- New `AppConfigInterface` methods: `SetAppMaintenanceEnabled`, `GetAppMaintenanceEnabled`, `SetAppMaintenanceFilePath`, `GetAppMaintenanceFilePath`
- New environment variables: `APP_MAINTENANCE_ENABLED`, `APP_MAINTENANCE_FILE_PATH`
- New files: `internal/middlewares/maintenance_middleware.go`, `internal/cli/maintenance_handler.go` (+ test files)
- `maintenance_mode_state.json` added to `.gitignore`
- README and `docs/environment-variables.md` updated with maintenance mode documentation
- New `driverTurso` constant in `internal/config/constants.go`
- New libSQL driver import in `internal/app/dbdrivers.go`
- `database_config.go` updated to treat `turso` identically to `sqlite` (pool settings, SSL mode, required fields)
- New dependency: `github.com/tursodatabase/libsql-client-go`
- 15 dracory store dependencies updated to latest versions

---

## ⚠️ Breaking Changes

---

### 1. New Methods on `AppConfigInterface`

**Change**: The `AppConfigInterface` in `internal/config/config_interfaces.go` gained four new methods for maintenance mode configuration. Any custom implementation of this interface (e.g., a mock or alternative config implementation) must now implement these methods.

**Old Interface**:
```go
// v0.35.0 — internal/config/config_interfaces.go
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppType(string)
	GetAppType() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}
```

**New Interface**:
```go
// v0.36.0 — internal/config/config_interfaces.go
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppType(string)
	GetAppType() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool

	// Maintenance mode
	SetAppMaintenanceEnabled(bool)
	GetAppMaintenanceEnabled() bool

	SetAppMaintenanceFilePath(string)
	GetAppMaintenanceFilePath() string
}
```

**Action Required**:
- If you have a custom type implementing `AppConfigInterface`, add the four new methods:

```go
func (c *myConfig) SetAppMaintenanceEnabled(v bool) {
	c.maintenanceEnabled = v
}

func (c *myConfig) GetAppMaintenanceEnabled() bool {
	return c.maintenanceEnabled
}

func (c *myConfig) SetAppMaintenanceFilePath(v string) {
	c.maintenanceFilePath = v
}

func (c *myConfig) GetAppMaintenanceFilePath() string {
	if c.maintenanceFilePath == "" {
		return "maintenance_mode_state.json"
	}
	return c.maintenanceFilePath
}
```

- If you use Blueprint's built-in `config.NewFromEnv()`, no changes are needed — the implementation is already updated.

**Migration Command**:
```bash
# Find custom implementations of AppConfigInterface
grep -rn "AppConfigInterface" --include="*.go" .
```

---

### 2. New `maintenance` Command Registered in CLI Dispatcher

**Change**: A new `maintenance` command is now registered in the CLI dispatcher (`internal/cli/cli.go`). If your application has a custom CLI dispatcher setup or overrides `NewDispatcher()`, you may want to register the maintenance command to enable the feature.

**Old Usage**:
```go
// v0.35.0 — internal/cli/cli.go
func NewDispatcher() *cli.Dispatcher[app.AppInterface] {
	dispatcher := cli.NewDispatcher[app.AppInterface]()
	dispatcher.RegisterCommand(CommandTask, "Execute a task by alias", handleTaskCommand)
	dispatcher.RegisterCommand(CommandJob, "Execute a job with arguments", handleJobCommand)
	dispatcher.RegisterCommand(CommandRoutes, "List all registered routes", handleRoutesCommand)
	return dispatcher
}
```

**New Usage**:
```go
// v0.36.0 — internal/cli/cli.go
func NewDispatcher() *cli.Dispatcher[app.AppInterface] {
	dispatcher := cli.NewDispatcher[app.AppInterface]()
	dispatcher.RegisterCommand(CommandTask, "Execute a task by alias", handleTaskCommand)
	dispatcher.RegisterCommand(CommandJob, "Execute a job with arguments", handleJobCommand)
	dispatcher.RegisterCommand(CommandRoutes, "List all registered routes", handleRoutesCommand)
	dispatcher.RegisterCommand(CommandMaintenance, "Manage maintenance mode", handleMaintenanceCommand)
	return dispatcher
}
```

**Action Required**:
- If you use the default `NewDispatcher()`, no changes needed — the maintenance command is already registered.
- If you have a custom dispatcher, add the maintenance command registration if you want the feature.

---

### 3. `MaintenanceMiddleware` Added to Global Middleware Chain

**Change**: `middlewares.NewMaintenanceMiddleware(app)` is now the **first** middleware in the global middleware chain in `internal/routes/global_middlewares.go`. When maintenance mode is active, it short-circuits all subsequent middleware and route processing by returning `503 Service Unavailable`.

**Old Usage**:
```go
// v0.35.0 — internal/routes/global_middlewares.go
func globalMiddlewares(app app.AppInterface) []rtr.MiddlewareInterface {
	// ...
	globalMiddlewares := []rtr.MiddlewareInterface{
		rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
			Exclude: []string{"/new"},
			// ...
		}),
		// ... other middlewares
	}
	return globalMiddlewares
}
```

**New Usage**:
```go
// v0.36.0 — internal/routes/global_middlewares.go
func globalMiddlewares(app app.AppInterface) []rtr.MiddlewareInterface {
	// ...
	globalMiddlewares := []rtr.MiddlewareInterface{
		// Maintenance mode check first — blocks all processing if active
		middlewares.NewMaintenanceMiddleware(app),
		rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
			Exclude: []string{"/new"},
			// ...
		}),
		// ... other middlewares
	}
	return globalMiddlewares
}
```

**Action Required**:
- If you use the default `globalMiddlewares()`, no changes needed.
- If you have a custom global middleware chain, consider adding `middlewares.NewMaintenanceMiddleware(app)` as the first entry if you want maintenance mode support.

---

### 4. Turso Database Driver Support

**Change**: A new `turso` database driver is now supported. This required adding a new `driverTurso` constant, a new libSQL driver import, and updating all SQLite-specific checks in `database_config.go` to also handle `turso`. The `DB_DRIVER` environment variable now accepts `turso` as a valid value.

**Old Code**:
```go
// v0.35.0 — internal/config/constants.go
const driverSQLite = "sqlite"

// v0.35.0 — internal/app/dbdrivers.go
import (
	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	// _ "modernc.org/sqlite"
)

// v0.35.0 — internal/config/database_config.go
if driver == driverSQLite {
	maxOpenConns = 1
}
```

**New Code**:
```go
// v0.36.0 — internal/config/constants.go
const driverSQLite = "sqlite"
const driverTurso = "turso"

// v0.36.0 — internal/app/dbdrivers.go
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	// _ "github.com/lib/pq"
	// _ "modernc.org/sqlite"
)

// v0.36.0 — internal/config/database_config.go
if driver == driverSQLite || driver == driverTurso {
	maxOpenConns = 1
}
```

**Action Required**:
- If you use SQLite or MySQL, no changes needed — existing drivers continue to work.
- To use Turso, set `DB_DRIVER="turso"` and provide a libSQL DSN via `DB_DSN`.
- If you have custom code that checks `driver == "sqlite"`, consider also checking for `"turso"` since Turso is SQLite-compatible.
- If you have GORM-based stores with driver switches, add `"turso"` to the `"sqlite"` case:

```go
// Old
switch driver {
case "sqlite":
	gdb, err = gorm.Open(sqlite.Dialector{Conn: opts.DB}, gcfg)

// New
switch driver {
case "sqlite", "turso":
	gdb, err = gorm.Open(sqlite.Dialector{Conn: opts.DB}, gcfg)
```

**Migration Command**:
```bash
# Find all SQLite-only driver checks that may need turso support
grep -rn 'case "sqlite"' --include="*.go" . | grep -v _test.go
```

---

### 5. Updated Dracory Store Dependencies

**Change**: 15 dracory store dependencies were updated to their latest versions. This is a non-breaking change for most applications, but if you use replace directives for local development, you may need to update your local module versions to match.

**Updated Dependencies**:

| Package | Old Version | New Version |
|---------|------------|-------------|
| `github.com/dracory/auditstore` | v1.9.0 | v1.10.0 |
| `github.com/dracory/auth` | v0.34.0 | v0.35.0 |
| `github.com/dracory/blindindexstore` | v1.15.0 | v1.16.0 |
| `github.com/dracory/blogstore` | v1.29.0 | v1.33.0 |
| `github.com/dracory/cmsstore` | v1.35.0 | v1.40.0 |
| `github.com/dracory/customstore` | v1.12.0 | v1.13.0 |
| `github.com/dracory/entitystore` | v1.12.0 | v1.13.0 |
| `github.com/dracory/feedstore` | v1.3.0 | v1.4.0 |
| `github.com/dracory/geostore` | v1.8.0 | v1.9.0 |
| `github.com/dracory/logstore` | v1.20.0 | v1.22.0 |
| `github.com/dracory/metastore` | v1.9.0 | v1.10.0 |
| `github.com/dracory/neat` | v0.31.0 | v0.33.0 |
| `github.com/dracory/sessionstore` | v1.17.0 | v1.18.0 |
| `github.com/dracory/sqlfilestore` | v1.8.0 | v1.9.0 |
| `github.com/dracory/statsstore` | v1.9.0 | v1.14.0 |

**Action Required**:
- Run `go mod tidy` to sync indirect dependencies.
- Run `go build ./...` to verify compatibility.
- If you have replace directives pointing to local copies, ensure your local modules are compatible with the new versions.

---

## 🔄 Migration Steps

### Step 1: Update Custom `AppConfigInterface` Implementations

If you have a custom config type implementing `AppConfigInterface`, add the four new maintenance methods. See the code example in Breaking Change #1 above.

```bash
# Search for custom implementations
grep -rn "AppConfigInterface" --include="*.go" . | grep -v "_test.go" | grep -v "config_interfaces.go"
```

### Step 2: Add New Environment Variables (Optional)

Add the following to your `.env` file if you want to configure maintenance mode via environment variables:

```bash
# Maintenance Mode Configuration
APP_MAINTENANCE_ENABLED=false
APP_MAINTENANCE_FILE_PATH="maintenance_mode_state.json"
```

Both variables are optional with sensible defaults (`false` and `maintenance_mode_state.json` respectively).

### Step 3: Verify `.gitignore` Entry

Ensure `maintenance_mode_state.json` is in your `.gitignore` so the runtime state file is not committed:

```bash
grep "maintenance_mode_state.json" .gitignore
```

If missing, add:
```
# Maintenance mode state file
maintenance_mode_state.json
```

### Step 4: Add Turso Driver Dependency (If Using Turso)

If you plan to use Turso, ensure the libSQL driver is available:

```bash
go get github.com/tursodatabase/libsql-client-go/libsql
```

Then set in your `.env`:
```bash
DB_DRIVER="turso"
DB_DSN="libsql://your-database-url"
DB_DATABASE="your-database-name"
```

Note: `DB_HOST`, `DB_PORT`, `DB_USERNAME`, and `DB_PASSWORD` are **not required** when using Turso.

### Step 5: Update Store Dependencies

Run `go mod tidy` to sync all updated store dependencies:

```bash
go mod tidy
```

### Step 6: Verify Build

```bash
go build ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run the full test suite to verify the new maintenance code doesn't break existing functionality:

```bash
go test ./...
```

### 2. Test Maintenance Mode CLI

Test the new CLI commands:

```bash
# Enable maintenance mode
go run ./cmd/server maintenance enable

# Check status
go run ./cmd/server maintenance status

# Disable maintenance mode
go run ./cmd/server maintenance disable

# Test aliases
go run ./cmd/server maintenance down
go run ./cmd/server maintenance up
```

### 3. Test Maintenance Middleware Behavior

1. Enable maintenance mode: `go run ./cmd/server maintenance enable`
2. Start the server: `go run ./cmd/server`
3. Verify all requests return `503 Service Unavailable` with the maintenance HTML page
4. Test IP exclusions: `go run ./cmd/server maintenance enable --ips="127.0.0.1"`
5. Test path exclusions: `go run ./cmd/server maintenance enable --exclude="/api/health,/admin/*"`
6. Verify excluded IPs/paths bypass the 503 response
7. Disable maintenance mode: `go run ./cmd/server maintenance disable`
8. Verify normal traffic resumes

### 4. Test Environment Variable Override

1. Set `APP_MAINTENANCE_ENABLED=true` in your `.env`
2. Start the server (without a state file present)
3. Verify the server returns 503 with the default message "We'll be right back."
4. Remove the env var or set to `false`
5. Verify normal traffic resumes

### 5. Test Turso Database Connection (If Using Turso)

1. Set `DB_DRIVER="turso"` and `DB_DSN="libsql://your-database-url"` in your `.env`
2. Run `go build ./...` to verify compilation
3. Start the server: `go run ./cmd/server`
4. Verify database operations work (migrations, queries, etc.)
5. Verify connection pool settings (maxOpenConns=1, maxIdleConns=1)

---

## 📝 Additional Notes

### New Features

- **Maintenance Mode Middleware** (`internal/middlewares/maintenance_middleware.go`):
  - File-based flag — presence of `maintenance_mode_state.json` activates maintenance mode
  - Returns `503 Service Unavailable` with a styled HTML page
  - Supports `Retry-After` header via `retry_after_seconds` in the state file
  - IP exclusions via `exclude_ips` field (exact match)
  - Path exclusions via `exclude_paths` field (supports wildcards: `/admin/*`, `/api/*`)
  - 30-second file stat caching for performance
  - Environment variable override: `APP_MAINTENANCE_ENABLED=true` forces maintenance mode even without the state file (useful for containerized deployments)

- **Maintenance CLI Command** (`internal/cli/maintenance_handler.go`):
  - Subcommands: `enable`, `disable`, `status` (aliases: `on`/`down`, `off`/`up`)
  - Enable options: `--message="..."`, `--retry=120`, `--ips="..."`, `--exclude="..."`
  - Writes/reads/removes the `maintenance_mode_state.json` file

- **Maintenance State File Format** (`maintenance_mode_state.json`):
  ```json
  {
    "message": "We'll be right back.",
    "retry_after_seconds": 60,
    "exclude_ips": ["203.0.113.5"],
    "exclude_paths": ["/admin/*", "/api/health"],
    "created_at": "2026-07-19T14:30:00Z"
  }
  ```

### Removed Features

- None

### Configuration Changes

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_MAINTENANCE_ENABLED` | No | `false` | Enable maintenance mode via env var (returns 503 for all requests except excluded IPs/paths) |
| `APP_MAINTENANCE_FILE_PATH` | No | `maintenance_mode_state.json` | Path to the maintenance state JSON file |
| `DB_DRIVER` | Yes | `sqlite` | Now accepts `turso` as a valid value (in addition to `sqlite`, `mysql`, `postgres`) |

### Turso Driver Details

- **Driver name**: `turso`
- **Go SQL driver**: `github.com/tursodatabase/libsql-client-go/libsql` (registered via blank import in `internal/app/dbdrivers.go`)
- **Connection**: Uses `DB_DSN` with `libsql://` scheme (e.g., `libsql://your-db.turso.io`)
- **Connection pool**: Same as SQLite (maxOpenConns=1, maxIdleConns=1, connMaxLifetime=30s)
- **SSL mode**: Empty (not applicable)
- **Required env vars**: `DB_DRIVER`, `DB_DATABASE` (or `DB_DSN`)
- **Not required**: `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`
- **GORM compatibility**: Turso uses the SQLite GORM dialector (libSQL is SQLite-compatible)

---

## 🆘 Common Issues and Solutions

### Issue 1: Compile Error — Missing Methods on Custom Config Type

**Symptom**: After upgrading, you get a compile error like:
```
cannot use &myConfig{} as config.AppConfigInterface in assignment:
	*myConfig does not implement config.AppConfigInterface (missing SetAppMaintenanceEnabled method)
```

**Solution**: Add the four new methods to your custom config type. See the code example in Breaking Change #1.

### Issue 2: Maintenance Mode Stuck On

**Symptom**: The server keeps returning 503 even after running `maintenance disable`.

**Solution**:
1. Check if `APP_MAINTENANCE_ENABLED=true` is set in your `.env` — this env var override forces maintenance mode regardless of the state file.
2. Verify the state file was removed: `go run ./cmd/server maintenance status`
3. Manually delete the state file if needed: `rm maintenance_mode_state.json`

### Issue 3: State File Not Found After Enable

**Symptom**: `maintenance status` reports OFF after running `maintenance enable`.

**Solution**: Check the file path. If `APP_MAINTENANCE_FILE_PATH` is set to a custom path, the CLI writes there. Verify the path is writable and accessible from the directory where you run the server.

### Issue 4: Excluded Paths Not Working

**Symptom**: Requests to excluded paths still return 503.

**Solution**: Verify the wildcard syntax in your exclude paths:
- `/admin/*` matches `/admin/anything` and `/admin/anything/deep`
- `/api/health` matches only the exact path `/api/health`
- `/api/*` matches all paths under `/api/`

### Issue 5: Turso Driver Not Found

**Symptom**: After setting `DB_DRIVER="turso"`, you get a connection error like `driver: bad connection` or `unknown driver`.

**Solution**:
1. Verify the libSQL driver is imported in `internal/app/dbdrivers.go`:
   ```go
   import (
       _ "github.com/tursodatabase/libsql-client-go/libsql"
   )
   ```
2. Run `go mod tidy` to ensure the dependency is downloaded.
3. Verify `DB_DSN` is set to a valid `libsql://` URL.

### Issue 6: GORM Store Fails with Turso Driver

**Symptom**: A GORM-based store returns an error like `unsupported driver: turso`.

**Solution**: Update the driver switch in the store to include `"turso"` in the `"sqlite"` case:
```go
// Old
case "sqlite":
	gdb, err = gorm.Open(sqlite.Dialector{Conn: opts.DB}, gcfg)

// New
case "sqlite", "turso":
	gdb, err = gorm.Open(sqlite.Dialector{Conn: opts.DB}, gcfg)
```

---

## Support

For issues or questions about this upgrade:
- Review the maintenance mode proposal: `docs/proposals/maintenance-mode.md`
- Check the maintenance middleware source: `internal/middlewares/maintenance_middleware.go`
- Check the CLI handler source: `internal/cli/maintenance_handler.go`
- Check the Turso driver registration: `internal/app/dbdrivers.go`
- Check the database config: `internal/config/database_config.go`
- Open an issue on the [Blueprint repository](https://github.com/dracory/blueprint)
