# Upgrade Guide: v0.33.0 to v0.34.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.33.0 to v0.34.0.

## Overview

This release introduces unified debug mode across all stores, the database layer, and the console logger. When `APP_DEBUG=true`, all stores now automatically enable debug logging, the database ORM outputs debug SQL logs, and the console logger switches to `slog.LevelDebug`. Additionally, multiple Dracory store dependencies and third-party packages are upgraded to their latest versions.

**Key Changes:**
- All 11 store initializers now call `store.EnableDebug(app.GetConfig().GetAppDebug())` during initialization
- Console logger now respects `APP_DEBUG` — uses `slog.LevelDebug` when debug is enabled, `slog.LevelInfo` otherwise
- `DatabaseNeatConfig` now passes `Debug: cfg.GetAppDebug()` to the neat ORM config
- `github.com/samber/lo` added as a direct import in `app_implementation.go` for ternary debug level selection
- `neat` upgraded from v0.27.0 to v0.31.0 (significant jump)
- `auth` upgraded from v0.29.0 to v0.34.0 (significant jump)
- `shopstore` upgraded from v1.18.0 to v1.24.0
- `statsstore` upgraded from v1.3.0 to v1.9.0
- `taskstore` upgraded from v1.26.0 to v1.29.0
- `base` upgraded from v0.37.0 to v0.39.0
- Multiple other Dracory store dependency version bumps (25 Dracory packages total)
- Third-party package bumps: `tint` v1.2.0, `goldmark` v1.8.4, `ttlcache` v3.4.1
- New indirect dependency: `github.com/LumenResearch/uasurfer`

---

## ⚠️ Breaking Changes

---

### 1. All Stores Now Enable Debug Mode Automatically

**Change**: All 11 store initializers in `internal/app/stores_*.go` now call `store.EnableDebug(app.GetConfig().GetAppDebug())` after creating the store instance. This means when `APP_DEBUG=true`, all stores will output debug-level logs (SQL queries, cache operations, etc.).

**Affected Files**:
- `internal/app/stores_blog.go`
- `internal/app/stores_cache.go`
- `internal/app/stores_cms.go`
- `internal/app/stores_custom.go`
- `internal/app/stores_log.go`
- `internal/app/stores_meta.go`
- `internal/app/stores_session.go`
- `internal/app/stores_shop.go`
- `internal/app/stores_stats.go`
- `internal/app/stores_task.go`
- `internal/app/stores_vault.go`

**Old Usage**:
```go
// v0.33.0 — stores_blog.go
func blogStoreInitialize(app AppInterface) error {
	if store, err := newBlogStore(app.GetDatabase()); err != nil {
		return err
	} else {
		app.SetBlogStore(store)
	}
	return nil
}
```

**New Usage**:
```go
// v0.34.0 — stores_blog.go
func blogStoreInitialize(app AppInterface) error {
	if store, err := newBlogStore(app.GetDatabase()); err != nil {
		return err
	} else {
		store.EnableDebug(app.GetConfig().GetAppDebug())
		app.SetBlogStore(store)
	}
	return nil
}
```

**Action Required**:
- If you have custom store initializers, add `store.EnableDebug(app.GetConfig().GetAppDebug())` before calling `app.Set*Store(store)`
- If you have custom stores that don't implement `EnableDebug(bool)`, you may need to add this method to match the store interface
- Be aware that enabling `APP_DEBUG=true` in production will now produce significantly more log output from all stores

---

### 2. Console Logger Now Respects APP_DEBUG Setting

**Change**: The console logger in `app_implementation.go` now uses `tint.NewHandler` with a `Level` option that switches between `slog.LevelDebug` and `slog.LevelInfo` based on `cfg.GetAppDebug()`. Previously, the handler used `nil` options (defaulting to `slog.LevelInfo`).

**Old Usage**:
```go
// v0.33.0 — app_implementation.go
consoleLogger := slog.New(tint.NewHandler(os.Stdout, nil))
```

**New Usage**:
```go
// v0.34.0 — app_implementation.go
consoleLogger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
	Level: lo.Ternary(cfg.GetAppDebug(), slog.LevelDebug, slog.LevelInfo),
}))
```

**Action Required**:
- If you customize the console logger in your application, update it to respect `APP_DEBUG` for log level
- Add `github.com/samber/lo` to your imports if you use the same ternary pattern (already a dependency)
- Be aware that `APP_DEBUG=true` now produces debug-level console output, which is much more verbose

---

### 3. DatabaseNeatConfig Now Passes Debug Flag

**Change**: `DatabaseNeatConfig` in `internal/config/database_config.go` now sets `Debug: cfg.GetAppDebug()` in the returned `db.DBConfig`. This means the neat ORM will output debug SQL logs when `APP_DEBUG=true`.

**Old Usage**:
```go
// v0.33.0 — database_config.go
func DatabaseNeatConfig(cfg ConfigInterface) db.DBConfig {
	return db.DBConfig{
		Default:     defaultConnection,
		Connections: connections,
		Pool:        pool,
	}
}
```

**New Usage**:
```go
// v0.34.0 — database_config.go
func DatabaseNeatConfig(cfg ConfigInterface) db.DBConfig {
	return db.DBConfig{
		Default:     defaultConnection,
		Connections: connections,
		Pool:        pool,
		Debug:       cfg.GetAppDebug(),
	}
}
```

**Action Required**:
- If you call `DatabaseNeatConfig` directly, be aware that the `Debug` field is now populated
- If you construct `db.DBConfig` manually, add `Debug: cfg.GetAppDebug()` to enable ORM debug logging
- The `neat` package (v0.31.0) must support the `Debug` field in `DBConfig` — ensure your `neat` version is v0.31.0 or later

---

### 4. Dependency Upgrades — neat v0.27.0 to v0.31.0

**Change**: `github.com/dracory/neat` is upgraded from v0.27.0 to v0.31.0, a significant jump of 4 minor versions. This version adds the `Debug` field support in `DBConfig` and may include other API changes.

**Action Required**:
- Update `go.mod` to use `github.com/dracory/neat v0.31.0`
- Run `go mod tidy` to resolve transitive dependencies
- Test all database operations thoroughly after upgrading
- If you use `neat` types directly in your code, check the [neat changelog](https://github.com/dracory/neat) for any API changes between v0.27.0 and v0.31.0

**Migration Command**:
```bash
go get github.com/dracory/neat@v0.31.0
go mod tidy
```

---

### 5. Dependency Upgrades — Dracory Store Packages

**Change**: Multiple Dracory store packages are upgraded to their latest versions. Several of these are significant jumps.

**Upgraded Dependencies**:

| Package | Old Version | New Version |
|---------|-------------|-------------|
| `github.com/dracory/auditstore` | v1.8.0 | v1.9.0 |
| `github.com/dracory/auth` | v0.29.0 | v0.34.0 |
| `github.com/dracory/base` | v0.37.0 | v0.39.0 |
| `github.com/dracory/blindindexstore` | v1.14.0 | v1.15.0 |
| `github.com/dracory/cachestore` | v1.7.0 | v1.8.0 |
| `github.com/dracory/chatstore` | v1.2.0 | v1.3.0 |
| `github.com/dracory/cmsstore` | v1.34.0 | v1.35.0 |
| `github.com/dracory/customstore` | v1.11.0 | v1.12.0 |
| `github.com/dracory/entitystore` | v1.11.0 | v1.12.0 |
| `github.com/dracory/feedstore` | v1.2.0 | v1.3.0 |
| `github.com/dracory/filesystem` | v1.3.0 | v1.4.0 |
| `github.com/dracory/geostore` | v1.6.0 | v1.8.0 |
| `github.com/dracory/logstore` | v1.19.0 | v1.20.0 |
| `github.com/dracory/metastore` | v1.8.0 | v1.9.0 |
| `github.com/dracory/neat` | v0.27.0 | v0.31.0 |
| `github.com/dracory/sessionstore` | v1.16.0 | v1.17.0 |
| `github.com/dracory/settingstore` | v1.9.0 | v1.11.0 |
| `github.com/dracory/shopstore` | v1.18.0 | v1.24.0 |
| `github.com/dracory/statsstore` | v1.3.0 | v1.9.0 |
| `github.com/dracory/subscriptionstore` | v1.3.0 | v1.4.0 |
| `github.com/dracory/taskstore` | v1.26.0 | v1.29.0 |
| `github.com/dracory/userstore` | v1.16.0 | v1.17.0 |
| `github.com/dracory/vaultstore` | v1.3.0 | v1.4.0 |
| `github.com/dracory/versionstore` | v1.6.0 | v1.7.0 |
| `github.com/dracory/sqlfilestore` (indirect) | v1.6.0 | v1.8.0 |

**Action Required**:
- Update all dependencies in `go.mod` to the versions listed above
- Run `go mod tidy` to resolve transitive dependencies
- Test all store operations after upgrading
- Pay special attention to `auth` (v0.29.0 → v0.34.0), `shopstore` (v1.18.0 → v1.24.0), `statsstore` (v1.3.0 → v1.9.0), and `neat` (v0.27.0 → v0.31.0) as these are the largest jumps
- If you use store interfaces directly, check for any interface method signature changes

**Migration Command**:
```bash
go get github.com/dracory/auditstore@v1.9.0
go get github.com/dracory/auth@v0.34.0
go get github.com/dracory/base@v0.39.0
go get github.com/dracory/blindindexstore@v1.15.0
go get github.com/dracory/cachestore@v1.8.0
go get github.com/dracory/chatstore@v1.3.0
go get github.com/dracory/cmsstore@v1.35.0
go get github.com/dracory/customstore@v1.12.0
go get github.com/dracory/entitystore@v1.12.0
go get github.com/dracory/feedstore@v1.3.0
go get github.com/dracory/filesystem@v1.4.0
go get github.com/dracory/geostore@v1.8.0
go get github.com/dracory/logstore@v1.20.0
go get github.com/dracory/metastore@v1.9.0
go get github.com/dracory/neat@v0.31.0
go get github.com/dracory/sessionstore@v1.17.0
go get github.com/dracory/settingstore@v1.11.0
go get github.com/dracory/shopstore@v1.24.0
go get github.com/dracory/statsstore@v1.9.0
go get github.com/dracory/subscriptionstore@v1.4.0
go get github.com/dracory/taskstore@v1.29.0
go get github.com/dracory/userstore@v1.17.0
go get github.com/dracory/vaultstore@v1.4.0
go get github.com/dracory/versionstore@v1.7.0
go mod tidy
```

---

### 6. Third-Party Package Upgrades

**Change**: Several third-party packages are upgraded to their latest versions.

**Upgraded Packages**:

| Package | Old Version | New Version |
|---------|-------------|-------------|
| `github.com/jellydator/ttlcache/v3` | v3.4.0 | v3.4.1 |
| `github.com/lmittmann/tint` | v1.1.3 | v1.2.0 |
| `github.com/yuin/goldmark` | v1.8.2 | v1.8.4 |
| `cloud.google.com/go/aiplatform` | v1.125.0 | v1.126.0 |
| `cloud.google.com/go/longrunning` | v1.0.0 | v1.2.0 |
| `github.com/googleapis/enterprise-certificate-proxy` | v0.3.16 | v0.3.17 |
| `github.com/googleapis/gax-go/v2` | v2.22.0 | v2.23.0 |
| `golang.org/x/crypto` | v0.53.0 | v0.54.0 |
| `golang.org/x/net` | v0.56.0 | v0.57.0 |
| `golang.org/x/sync` | v0.21.0 | v0.22.0 |
| `golang.org/x/sys` | v0.46.0 | v0.47.0 |
| `golang.org/x/term` | v0.44.0 | v0.45.0 |
| `golang.org/x/text` | v0.38.0 | v0.40.0 |
| `google.golang.org/api` | v0.285.0 | v0.287.1 |
| `google.golang.org/grpc` | v1.81.1 | v1.82.0 |
| `modernc.org/libc` | v1.73.5 | v1.74.1 |
| `github.com/aws/smithy-go` (indirect) | v1.27.2 | v1.27.3 |
| `golang.org/x/exp` (indirect) | v0.0.0-20260611... | v0.0.0-20260709... |
| `google.golang.org/genproto/googleapis/api` (indirect) | v0.0.0-20260618... | v0.0.0-20260630... |
| `google.golang.org/genproto/googleapis/rpc` (indirect) | v0.0.0-20260618... | v0.0.0-20260630... |

**New Indirect Dependency**:
- `github.com/LumenResearch/uasurfer v0.0.0-20260126094926-dace53404a8d`

**Action Required**:
- Run `go mod tidy` to pick up all transitive dependency updates
- The `tint` upgrade to v1.2.0 is required for the new `tint.Options{Level: ...}` usage in the console logger
- No code changes required for third-party packages — these are backward compatible updates

---

## 🔄 Migration Steps

### Step 1: Update go.mod Dependencies

Update all Dracory store dependencies and third-party packages:

```bash
# Update Dracory store packages
go get github.com/dracory/auditstore@v1.9.0
go get github.com/dracory/auth@v0.34.0
go get github.com/dracory/base@v0.39.0
go get github.com/dracory/blindindexstore@v1.15.0
go get github.com/dracory/cachestore@v1.8.0
go get github.com/dracory/chatstore@v1.3.0
go get github.com/dracory/cmsstore@v1.35.0
go get github.com/dracory/customstore@v1.12.0
go get github.com/dracory/entitystore@v1.12.0
go get github.com/dracory/feedstore@v1.3.0
go get github.com/dracory/filesystem@v1.4.0
go get github.com/dracory/geostore@v1.8.0
go get github.com/dracory/logstore@v1.20.0
go get github.com/dracory/metastore@v1.9.0
go get github.com/dracory/neat@v0.31.0
go get github.com/dracory/sessionstore@v1.17.0
go get github.com/dracory/settingstore@v1.11.0
go get github.com/dracory/shopstore@v1.24.0
go get github.com/dracory/statsstore@v1.9.0
go get github.com/dracory/subscriptionstore@v1.4.0
go get github.com/dracory/taskstore@v1.29.0
go get github.com/dracory/userstore@v1.17.0
go get github.com/dracory/vaultstore@v1.4.0
go get github.com/dracory/versionstore@v1.7.0

# Update third-party packages
go get github.com/lmittmann/tint@v1.2.0
go get github.com/yuin/goldmark@v1.8.4
go get github.com/jellydator/ttlcache/v3@v3.4.1

# Tidy up
go mod tidy
```

### Step 2: Add EnableDebug Calls to Custom Store Initializers

If you have custom store initializers, add `EnableDebug` calls:

```go
// Example for a custom store
func myCustomStoreInitialize(app AppInterface) error {
	if store, err := newMyCustomStore(app.GetDatabase()); err != nil {
		return err
	} else {
		store.EnableDebug(app.GetConfig().GetAppDebug())
		app.SetCustomStore(store)
	}
	return nil
}
```

### Step 3: Update Console Logger (If Customized)

If you have customized the console logger, update it to respect `APP_DEBUG`:

```go
import (
	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	"log/slog"
	"os"
)

consoleLogger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
	Level: lo.Ternary(cfg.GetAppDebug(), slog.LevelDebug, slog.LevelInfo),
}))
```

### Step 4: Update DatabaseNeatConfig (If Customized)

If you construct `db.DBConfig` manually, add the `Debug` field:

```go
config := db.DBConfig{
	Default:     defaultConnection,
	Connections: connections,
	Pool:        pool,
	Debug:       cfg.GetAppDebug(),
}
```

### Step 5: Verify Build

```bash
go build ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
go test ./...
```

### 2. Verify Debug Mode Behavior

Test with `APP_DEBUG=true` to verify debug logging is enabled:

```bash
# Set debug mode
export APP_DEBUG=true

# Run the application and verify:
# - Console logger outputs DEBUG-level messages
# - Stores output debug SQL queries
# - Database ORM outputs debug SQL logs

# Run tests
go test ./internal/app/... -v
```

### 3. Verify Production Mode Behavior

Test with `APP_DEBUG=false` to verify debug logging is disabled:

```bash
# Set production mode
export APP_DEBUG=false

# Run the application and verify:
# - Console logger only outputs INFO-level messages
# - Stores do not output debug logs
# - Database ORM does not output debug SQL logs
```

### 4. Store-Specific Tests

Run tests for each store to verify the `EnableDebug` calls don't cause issues:

```bash
go test ./internal/app/... -v -run TestStore
```

---

## 📝 Additional Notes

### New Features
- **Unified Debug Mode**: All stores, database ORM, and console logger now respect `APP_DEBUG` setting, providing consistent debug output across the entire application stack
- **Debug-level SQL Logging**: When `APP_DEBUG=true`, the neat ORM now outputs SQL queries for debugging purposes

### Behavioral Changes
- **Increased Log Verbosity**: Enabling `APP_DEBUG=true` now produces significantly more log output from stores and database layer. Use only in development/debugging environments
- **Console Log Level**: The console logger now defaults to `slog.LevelInfo` (same as before) but switches to `slog.LevelDebug` when `APP_DEBUG=true`

### Removed Features
- None

---

## 🆘 Common Issues and Solutions

### Issue 1: Store Does Not Implement EnableDebug
**Symptom**: Build error: `store.EnableDebug undefined (type X has no field or method EnableDebug)`
**Solution**: Ensure your custom store type implements the `EnableDebug(bool)` method. If using a Dracory store package, update to the latest version which includes this method.

### Issue 2: tint.Options Level Field Not Found
**Symptom**: Build error: `unknown field 'Level' in struct literal of type tint.Options`
**Solution**: Update `github.com/lmittmann/tint` to v1.2.0 or later:
```bash
go get github.com/lmittmann/tint@v1.2.0
```

### Issue 3: db.DBConfig Does Not Have Debug Field
**Symptom**: Build error: `unknown field 'Debug' in struct literal of type db.DBConfig`
**Solution**: Update `github.com/dracory/neat` to v0.31.0 or later:
```bash
go get github.com/dracory/neat@v0.31.0
```

### Issue 4: Excessive Log Output in Production
**Symptom**: Application produces too many log lines in production
**Solution**: Ensure `APP_DEBUG=false` (or unset) in production environment variables. The debug mode is controlled by the `APP_DEBUG` environment variable.

### Issue 5: go mod tidy Fails After Dependency Updates
**Symptom**: `go mod tidy` reports errors about incompatible versions
**Solution**: Update all dependencies together rather than one at a time. Run `go get` for all packages first, then `go mod tidy`:
```bash
go get github.com/dracory/neat@v0.31.0 github.com/dracory/shopstore@v1.24.0 github.com/dracory/statsstore@v1.9.0
go mod tidy
```

---

## Support

- [Blueprint Repository](https://github.com/dracory/blueprint)
- [Upgrade Guides Directory](docs/upgrade_guides/)
- [Version Workflow Documentation](docs/version_workflow.md)

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Git tag verified for previous version (v0.33.0)
- [x] Previous guides reviewed for consistency
- [x] Quality checklist included in generated guide
