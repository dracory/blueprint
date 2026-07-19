# Proposal: Maintenance Mode

## Status

Implemented

## Context

Blueprint currently has no way to take the application offline for deployments, database migrations, or infrastructure work without stopping the web server entirely. When the server is stopped, users get connection errors (browser "can't reach" page) instead of a clear "we'll be back" message. There's also no way to allow administrators to bypass maintenance mode to verify fixes before bringing the site back online.

Laravel solves this with `php artisan down` / `php artisan up` — a file-based flag checked by middleware, with support for a secret bypass token, retry intervals, and excluded paths.

## Goal

Add a maintenance mode system to Blueprint that:
1. Can be toggled via CLI commands (`go run ./cmd/server maintenance down` / `up`) without restarting the server
2. Returns a 503 Service Unavailable response with a user-friendly maintenance page
3. Supports excluded IPs (so admins can access the site during maintenance by IP)
4. Supports excluded paths (so `/admin/*` or `/api/health` can remain accessible)
5. Sets a `Retry-After` HTTP header so browsers/crawlers know when to retry
6. Works even when the database is down (file-based, no DB dependency)

## Key Design Decisions

1. **File-based flag (not database, not env var)** — A JSON file at `maintenance_mode_state.json` (project root) indicates maintenance mode is active. This mirrors Django's `django-maintenance-mode` package approach. Rationale:
   - No database dependency — works even during DB migrations
   - No server restart required — file creation/removal takes effect on next request
   - Atomic and simple — file exists = maintenance, file absent = normal operation
   - Can be created/removed by CLI, deploy scripts, or manually

2. **Middleware implementation** — A `MaintenanceMiddleware` added to the global middleware chain, checked early (before auth, before route matching). This ensures all routes are blocked uniformly.

3. **CLI commands via existing dispatcher** — Uses the existing `cli.Dispatcher` pattern (`go run ./cmd/server maintenance down`, `maintenance up`, `maintenance status`). No new CLI infrastructure needed.

4. **Excluded IPs bypass** — When `maintenance down --ips="203.0.113.5,198.51.100.10"` is used, requests from those IPs pass through maintenance mode. This is Django's `MAINTENANCE_MODE_IGNORE_IP_ADDRESSES` approach — simpler than a secret token, no cookie needed, and the IP list is stored in the JSON file.

5. **Configurable excluded paths** — Paths like `/admin/*`, `/api/health`, `/liveflux/*` can be exempted so monitoring and admin access continue working during maintenance.

6. **Custom maintenance message** — The `maintenance_mode_state.json` file stores an optional message and retry interval, displayed on the 503 page.

7. **Project root placement** — File lives at the project root (same level as `.env`). No `storage/` directory needed.

## Proposed Implementation

### New Files

#### `internal/middlewares/maintenance_middleware.go`

Middleware that checks for the maintenance file and returns 503 if active:

```go
type MaintenanceConfig struct {
    FilePath      string        // Path to maintenance state file (default: "maintenance_mode_state.json")
    ExcludePaths  []string      // Paths exempt from maintenance (e.g., "/admin/*", "/api/health")
    ExcludeIPs    []string      // IPs exempt from maintenance (e.g., "203.0.113.5")
    TemplatePath  string        // Optional custom HTML template path
    CacheDuration time.Duration // File check cache duration (default: 30s)
}

func NewMaintenanceMiddleware(app app.AppInterface) rtr.MiddlewareInterface
```

Behavior:
- Reads `maintenance_mode_state.json` on each request (file stat + read, cached for 30 seconds to avoid disk I/O on every request)
- If file doesn't exist → pass through (normal operation)
- If file exists:
  - Extract client IP (X-Forwarded-For, X-Real-IP, or RemoteAddr)
  - Check if client IP matches any `ExcludeIPs` entry → pass through
  - Check if request path matches any `ExcludePaths` pattern → pass through
  - Otherwise → return 503 with `Retry-After` header and maintenance HTML page

#### `internal/middlewares/maintenance_middleware_test.go`

Tests:
- No maintenance file → request passes through
- Maintenance file exists → request gets 503
- Excluded path matches → request passes through despite maintenance
- Excluded IP matches → request passes through despite maintenance
- Excluded IP with X-Forwarded-For header → request passes through
- Non-excluded IP → request gets 503
- Retry-After header is set correctly (seconds, per RFC 7231)
- Custom message is displayed in 503 response
- File cache: multiple requests within 30s cache window don't re-read file

#### `internal/cli/maintenance_handler.go`

CLI command handler for `maintenance` command:

```go
func handleMaintenanceCommand(app app.AppInterface, args []string) error
```

Subcommands:
- `maintenance enable [--message="Back soon"] [--retry=60] [--ips="203.0.113.5,198.51.100.10"] [--exclude="/admin/*,/api/health"]`
  - `--retry` is in **seconds** (matches the HTTP `Retry-After` header spec, RFC 7231)
  - Creates `maintenance_mode_state.json` with the specified options
  - Aliases: `on`, `down`
- `maintenance disable`
  - Removes `maintenance_mode_state.json`
  - Aliases: `off`, `up`
- `maintenance status`
  - Prints whether maintenance mode is active and the current configuration

#### `internal/cli/maintenance_handler_test.go`

Tests:
- `maintenance enable` creates the file with correct JSON
- `maintenance disable` removes the file
- `maintenance status` reports correctly
- `maintenance enable` with options (message, retry_seconds, ips, exclude) writes correct config
- `maintenance disable` when already off is a no-op
- `maintenance enable` when already on overwrites with new config

#### `maintenance_mode_state.json` (runtime file, gitignored)

JSON structure created by `maintenance enable`:

```json
{
  "message": "We'll be right back.",
  "retry_after_seconds": 60,
  "exclude_ips": ["203.0.113.5", "198.51.100.10"],
  "exclude_paths": ["/admin/*", "/api/health"],
  "created_at": "2026-07-19T07:49:00Z"
}
```

### Modified Files

#### `internal/routes/global_middlewares.go`

Add `MaintenanceMiddleware` as the **first** middleware in the global chain (before JailBots, before everything):

```go
globalMiddlewares := []rtr.MiddlewareInterface{
    middlewares.NewMaintenanceMiddleware(app), // Maintenance check first
    rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
        // ... existing config
    }),
    // ... rest of existing middlewares
}
```

Rationale for first position: maintenance mode should block bots, rate limiting, and all other processing. No point in running JailBots or rate limiting on a 503 response.

#### `internal/cli/cli.go`

Register the `maintenance` command in `NewDispatcher()`:

```go
dispatcher.RegisterCommand(CommandMaintenance, "Manage maintenance mode", handleMaintenanceCommand)
```

Add constant:
```go
const CommandMaintenance = "maintenance"
```

#### `internal/config/config_interfaces.go`

Add maintenance config methods to `AppConfigInterface`:

```go
type AppConfigInterface interface {
    // ... existing methods ...

    SetAppMaintenanceEnabled(bool)
    GetAppMaintenanceEnabled() bool

    SetAppMaintenanceFilePath(string)
    GetAppMaintenanceFilePath() string
}
```

#### `internal/config/config_implementation.go`

Implement the new methods. Default file path: `maintenance_mode_state.json`.

#### `internal/config/app_config.go`

Read `APP_MAINTENANCE_ENABLED` from env (optional, defaults to `false`). This allows enabling maintenance mode via env var as an alternative to the CLI command (useful for containerized deployments where file access isn't possible).

#### `internal/config/constants.go`

Add:
```go
const KEY_APP_MAINTENANCE_ENABLED = "APP_MAINTENANCE_ENABLED"
const KEY_APP_MAINTENANCE_FILE_PATH = "APP_MAINTENANCE_FILE_PATH"
```

#### `.env.example`

Add maintenance mode section:

```env
# ============================================================================
# Maintenance Mode Configuration
# ============================================================================

# Maintenance Mode Enabled
# When enabled, the app returns 503 Service Unavailable for all requests
# except excluded paths and excluded IPs.
# Valid values: true, false
# Default: false
# APP_MAINTENANCE_ENABLED=false

# Maintenance File Path
# Path to the maintenance flag file (when using CLI commands).
# Default: maintenance_mode_state.json
# APP_MAINTENANCE_FILE_PATH="maintenance_mode_state.json"
```

#### `.gitignore`

Add:
```
maintenance_mode_state.json
```

## Maintenance Page Template

The 503 response uses an inline HTML template (no external file dependency). The default template:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Maintenance</title>
    <style>
        body { font-family: system-ui, sans-serif; text-align: center; padding: 50px; }
        h1 { font-size: 2em; color: #333; }
        p { color: #666; }
    </style>
</head>
<body>
    <h1>Undergoing Maintenance</h1>
    <p>{{Message}}</p>
</body>
</html>
```

Variables:
- `{{Message}}` — from `maintenance_mode_state.json` message field, or default

A custom template path can be specified in `MaintenanceConfig.TemplatePath` for projects that want a branded maintenance page.

## CLI Usage

```bash
# Enable maintenance mode with defaults
go run ./cmd/server maintenance enable

# Enable with custom message and retry interval
go run ./cmd/server maintenance enable --message="Database migration in progress" --retry=120

# Enable with excluded IPs (admin bypass)
go run ./cmd/server maintenance enable --ips="203.0.113.5,198.51.100.10"

# Enable with excluded paths
go run ./cmd/server maintenance enable --exclude="/admin/*,/api/health,/liveflux/*"

# Disable maintenance mode
go run ./cmd/server maintenance disable

# Check maintenance status
go run ./cmd/server maintenance status
```

Aliases: `on`/`down` for enable, `off`/`up` for disable.

## File Caching Strategy

To avoid reading the maintenance file on every request (disk I/O on every HTTP request would be wasteful):

- The middleware caches the file check result for **30 seconds** using a simple `time.Time` last-checked timestamp
- On cache expiry, it re-stats the file (a stat syscall is negligible cost)
- If the file was created/removed, it re-reads the content
- This means maintenance mode takes effect within 30 seconds of the file being created — acceptable for a feature that's toggled infrequently
- The cache duration is configurable via `MaintenanceConfig.CacheDuration` (default: 30s)

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/middlewares/maintenance_middleware.go` | Create | Maintenance mode middleware |
| `internal/middlewares/maintenance_middleware_test.go` | Create | Middleware tests |
| `internal/cli/maintenance_handler.go` | Create | CLI command handler for `maintenance` command |
| `internal/cli/maintenance_handler_test.go` | Create | CLI handler tests |
| `internal/routes/global_middlewares.go` | Modify | Add MaintenanceMiddleware as first middleware |
| `internal/cli/cli.go` | Modify | Register `maintenance` command |
| `internal/config/config_interfaces.go` | Modify | Add maintenance config methods |
| `internal/config/config_implementation.go` | Modify | Implement maintenance config methods |
| `internal/config/app_config.go` | Modify | Read maintenance env vars |
| `internal/config/constants.go` | Modify | Add maintenance env var constants |
| `.env.example` | Modify | Add maintenance mode section |
| `.gitignore` | Modify | Add `maintenance_mode_state.json` |

## Testing

- `go test ./internal/middlewares/...` — middleware tests pass
- `go test ./internal/cli/...` — CLI handler tests pass
- `go test ./internal/config/...` — config tests pass with new maintenance fields
- `go test ./...` — all existing tests still pass

### Manual Verification

1. Start the server: `task dev`
2. Enable maintenance: `go run ./cmd/server maintenance down --message="Testing" --ips="127.0.0.1"`
3. Visit `http://localhost:32322` → see 503 maintenance page (from non-excluded IP)
4. Visit `http://127.0.0.1:32322` → site works (IP is excluded)
5. Disable maintenance: `go run ./cmd/server maintenance up`
6. Visit `http://localhost:32322` → site works normally

## Benefits

- **Zero-downtime deployments** — Take the app offline gracefully without killing the server process
- **User-friendly** — Users see a clear maintenance message instead of connection errors
- **Admin bypass** — Excluded IPs allow admins to verify the site before public re-enablement
- **SEO-friendly** — 503 status code with `Retry-After` tells search engines to come back later
- **Database-independent** — File-based flag works even during DB migrations
- **Flexible** — Env var support for containerized deployments, CLI for traditional deployments
- **Familiar** — Mirrors Laravel's `artisan down`/`up` pattern, easy for Laravel developers to understand

## Future Extensions

- Admin UI toggle in the admin panel (uses the same file-based mechanism)
- Scheduled maintenance windows (auto-enable/disable at configured times)
- Webhook notification when maintenance mode is toggled
- Multiple maintenance templates (branded, minimal, custom)
- API-specific maintenance response (JSON instead of HTML for `/api/*` paths)
