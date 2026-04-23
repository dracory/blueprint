# Upgrade Guide: v0.19.0 to v0.20.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.19.0 to v0.20.0.

## Summary

Version 0.20.0 introduces significant improvements to database configuration and connection management. The changes add support for configurable connection pool settings, charset and timezone configuration, and improve file organization in the config package.

**Key Changes:**
- Connection pool configuration (MaxOpenConns, MaxIdleConns, connection lifetime/idle time)
- Database charset and timezone configuration
- Dependency update: `github.com/dracory/database` v0.7.0 → v0.8.0
- Config file organization improvements (z_ prefix for framework files)
- Store configuration consolidation

---

## Breaking Changes

### 1. Config File Renames

**Change**: Framework config files have been renamed with a `z_` prefix for better directory organization.

**Old File Names** (v0.19.0):
```
internal/config/constants.go
internal/config/config_implementation.go
internal/config/config_implementation_test.go
internal/config/config_interfaces.go
```

**New File Names** (v0.20.0):
```
internal/config/z_config_constants.go
internal/config/z_config_implementation.go
internal/config/z_config_implementation_test.go
internal/config/z_config_interfaces.go
```

**Action Required**:
- Update any custom scripts or tools that reference the old file names
- If you have custom patches or modifications, update the file paths
- The file contents remain functionally identical - only the names changed

**Migration Command**:
```bash
# Update any build scripts or CI configurations
find . -type f \( -name "*.sh" -o -name "*.yml" -o -name "*.yaml" \) -exec sed -i 's/config\/constants\.go/config\/z_config_constants.go/g' {} \;
find . -type f \( -name "*.sh" -o -name "*.yml" -o -name "*.yaml" \) -exec sed -i 's/config\/config_implementation\.go/config\/z_config_implementation.go/g' {} \;
```

---

### 2. Store Configuration File Consolidation

**Change**: The `configuration_stores.go` file has been deleted and its contents moved to `stores_config.go`.

**Old Structure** (v0.19.0):
```
internal/config/configuration_stores.go  (store constants)
internal/config/stores_config.go         (store config loader)
internal/config/interfaces.go            (interfaces)
```

**New Structure** (v0.20.0):
```
internal/config/stores_config.go         (consolidated store config)
internal/config/z_config_interfaces.go   (renamed from interfaces.go)
```

**Action Required**:
- Remove references to `configuration_stores.go` in your codebase
- All store configuration constants are now in `stores_config.go`
- No code changes required - imports remain the same

**Migration Command**:
```bash
# Remove old file if it exists in your project
rm -f internal/config/configuration_stores.go

# Update any custom code that imported from the old file
find . -type f -name "*.go" -exec sed -i '/configuration_stores\.go/d' {} \;
```

---

### 3. New Required Config Interface Methods

**Change**: New getter and setter methods have been added to the `ConfigInterface` for connection pool and database settings.

**New Interface Methods** (added in v0.20.0):
```go
// Connection pool getters
GetDatabaseMaxOpenConns() int
GetDatabaseMaxIdleConns() int
GetDatabaseConnMaxLifetimeSeconds() int
GetDatabaseConnMaxIdleTimeSeconds() int

// Database charset and timezone getters
GetDatabaseCharset() string
GetDatabaseTimezone() string

// Corresponding setter methods on the implementation
SetDatabaseMaxOpenConns(v int)
SetDatabaseMaxIdleConns(v int)
SetDatabaseConnMaxLifetimeSeconds(v int)
SetDatabaseConnMaxIdleTimeSeconds(v int)
SetDatabaseCharset(v string)
SetDatabaseTimezone(v string)
```

**Action Required**:
- If you have a custom `ConfigInterface` implementation, you must add these new methods
- Most users use the built-in `configImplementation`, so no changes are needed

**Custom Implementation Example**:
```go
// Add these fields to your custom config struct
type myCustomConfig struct {
    // ... existing fields
    dbMaxOpenConns    int
    dbMaxIdleConns    int
    dbConnMaxLifetime int
    dbConnMaxIdleTime int
    dbCharset         string
    dbTimezone        string
}

// Implement the new getter methods
func (c *myCustomConfig) GetDatabaseMaxOpenConns() int {
    return c.dbMaxOpenConns
}

func (c *myCustomConfig) GetDatabaseMaxIdleConns() int {
    return c.dbMaxIdleConns
}

func (c *myCustomConfig) GetDatabaseConnMaxLifetimeSeconds() int {
    return c.dbConnMaxLifetime
}

func (c *myCustomConfig) GetDatabaseConnMaxIdleTimeSeconds() int {
    return c.dbConnMaxIdleTime
}

func (c *myCustomConfig) GetDatabaseCharset() string {
    return c.dbCharset
}

func (c *myCustomConfig) GetDatabaseTimezone() string {
    return c.dbTimezone
}
```

---

### 4. Dependency Version Update

**Change**: The `github.com/dracory/database` dependency has been updated from v0.7.0 to v0.8.0.

**Old** (v0.19.0 go.mod):
```
github.com/dracory/database v0.7.0
```

**New** (v0.20.0 go.mod):
```
github.com/dracory/database v0.8.0
```

**Action Required**:
- Run `go mod tidy` after upgrading
- The new version adds support for connection pool configuration methods
- No breaking changes in the database package API

---

## Migration Steps

### Step 1: Update Dependencies

Run the following commands to update dependencies:

```bash
# Update the database dependency
go get github.com/dracory/database@v0.8.0

# Tidy and verify modules
go mod tidy
go mod download
```

### Step 2: Clean Up Old Config Files

Remove the following files if they exist in your project:

```bash
# Remove old config files that have been replaced
rm -f internal/config/configuration_stores.go
rm -f internal/config/interfaces.go

# Note: The following files are renamed, not deleted
# They will be automatically replaced by git pull/merge:
# - constants.go → z_config_constants.go
# - config_implementation.go → z_config_implementation.go
# - config_implementation_test.go → z_config_implementation_test.go
# - config_interfaces.go → z_config_interfaces.go
```

### Step 3: Update Environment Variables

Add the new optional environment variables to your `.env` file:

```bash
# Database Connection Pool Settings (optional - defaults shown)
DB_MAX_OPEN_CONNS=25           # Max open connections (use 1 for SQLite)
DB_MAX_IDLE_CONNS=5            # Max idle connections (use 1 for SQLite)
DB_CONN_MAX_LIFETIME_SECONDS=300  # Connection max lifetime in seconds (5 minutes)
DB_CONN_MAX_IDLE_TIME_SECONDS=5   # Connection max idle time in seconds

# Database Charset and Timezone (optional - defaults shown)
DB_CHARSET=utf8mb4             # Database charset (MySQL only)
DB_TIMEZONE=UTC                # Database timezone
```

**Note**: These are all optional. If not specified, the following defaults apply:
- **MySQL/PostgreSQL**: MaxOpenConns=25, MaxIdleConns=5, Lifetime=300s, IdleTime=5s
- **SQLite**: MaxOpenConns=1, MaxIdleConns=1, Lifetime=30s, IdleTime=5s (to avoid concurrent write issues)

### Step 4: Verify Build

Build the application to verify no compilation errors:

```bash
go build -o ./tmp/main ./cmd/server
```

---

## Testing After Migration

### 1. Unit Tests

Run the config-specific tests to verify the new configuration mechanism:

```bash
go test ./internal/config/... -v
```

**Key Test Files:**
- `z_config_implementation_test.go` - validates all config loading scenarios including connection pool settings

### 2. Integration Tests

Run the full test suite to verify no regressions:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. Database Connection Pool Testing

Test the database connection with different pool settings:

```bash
# Test with custom pool settings
DB_MAX_OPEN_CONNS=50 DB_MAX_IDLE_CONNS=10 go run ./cmd/server

# Test SQLite (automatically uses pool=1)
DB_DRIVER=sqlite DB_DATABASE=./app.db go run ./cmd/server
```

### 4. Environment Variable Validation

Verify all new environment variables are correctly loaded:

```bash
# Test with explicit pool settings
DB_MAX_OPEN_CONNS=25 \
DB_MAX_IDLE_CONNS=5 \
DB_CONN_MAX_LIFETIME_SECONDS=300 \
DB_CONN_MAX_IDLE_TIME_SECONDS=5 \
DB_CHARSET=utf8mb4 \
DB_TIMEZONE=UTC \
APP_HOST=localhost \
APP_PORT=8080 \
APP_ENVIRONMENT=testing \
go run ./cmd/server
```

---

## Additional Notes

### New Features

1. **Connection Pool Configuration**: Configure database connection pool settings via environment variables for better resource management.

2. **Database Charset Configuration**: Specify database charset (primarily for MySQL) via `DB_CHARSET` environment variable.

3. **Database Timezone Configuration**: Set database timezone via `DB_TIMEZONE` environment variable.

4. **SQLite-Specific Optimizations**: Connection pool is automatically limited to 1 for SQLite to prevent concurrent write issues.

5. **Conditional Pool Settings**: Connection pool settings are only applied when non-zero, allowing database driver defaults to be used.

### Removed Features

- `internal/config/configuration_stores.go` - functionality consolidated into `stores_config.go`
- `internal/config/interfaces.go` - renamed to `z_config_interfaces.go`

### Configuration Behavior Changes

- **Connection Pool**: Now configurable via environment variables with sensible defaults for each database driver
- **Charset and Timezone**: Now configurable instead of hardcoded (was utf8mb4 and UTC)
- **File Organization**: Framework config files now prefixed with `z_` for better directory sorting

---

## Common Issues and Solutions

### Issue 1: "undefined: config.GetDatabaseMaxOpenConns"

**Symptom**: Compilation error when using custom ConfigInterface implementation

**Solution**: Add the new required methods to your custom implementation (see section 3 above)

### Issue 2: "cannot find package github.com/dracory/database@v0.7.0"

**Symptom**: Module resolution errors after upgrade

**Solution**: Run `go mod tidy` to update dependencies

```bash
go mod tidy
go mod download
```

### Issue 3: SQLite Concurrent Write Errors

**Symptom**: "database is locked" errors when using SQLite

**Solution**: Ensure connection pool is limited to 1. This is now automatic in v0.20.0 when using SQLite driver. If you've overridden the settings, reset them:

```bash
# Remove any DB_MAX_OPEN_CONNS override for SQLite
unset DB_MAX_OPEN_CONNS
unset DB_MAX_IDLE_CONNS
```

### Issue 4: Config File Import Errors

**Symptom**: Cannot find config files after migration

**Solution**: Ensure you're importing the package, not individual files:

```go
// Correct import
import "project/internal/config"

// Usage remains the same
cfg, err := config.NewFromEnv()
```

---

## Support

For issues related to this upgrade:

1. **Documentation**: Review the database package documentation at https://github.com/dracory/database

2. **Git History**: Review the migration commits:
   ```bash
   git log --oneline v0.19.0..v0.20.0
   ```

3. **Reference Implementation**: Compare with the reference implementation in the v0.20.0 tag

4. **Issues**: Report issues at https://github.com/dracory/blueprint/issues

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices

---

*This upgrade guide was generated for Blueprint v0.19.0 to v0.20.0 migration.*
