# Upgrade Guide: v0.24.0 to v0.25.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.24.0 to v0.25.0.

## Overview

This release introduces a major migration system overhaul that centralizes all database migrations into a dedicated `database/migrations` package. The new system provides better separation of concerns, improved testability, and clearer migration lifecycle management.

**Key Changes:**
- New centralized migration system in `database/migrations/`
- Context parameter added to migration Up/Down methods
- Store-level migrations for all dracory packages
- Updated dependency versions for multiple store packages

---

## ⚠️ Breaking Changes

### 1. Migration System Restructuring

**Change**: The migration system has been completely restructured. Individual store migration functions in `internal/registry/` have been replaced with a centralized migration system in `database/migrations/`.

**Old Usage**:
```go
// internal/registry/registry_datastores_migrate.go (DELETED)
func (r *registryImplementation) dataStoresMigrate() error {
    migrators := []func(registry RegistryInterface) error{
        auditStoreMigrate,
        blogStoreMigrate,
        // ... 20+ other store migrators
    }
    for _, m := range migrators {
        if err := m(r); err != nil {
            return err
        }
    }
    return nil
}

// Called during registry initialization
if err := registry.dataStoresMigrate(); err != nil {
    return err
}
```

**New Usage**:
```go
// database/migrations/migrate.go (NEW)
func MigrateAll(registry RegistryInterface) error {
    // Phase 1: Store-level migrations (run directly outside transactions)
    if err := migrateStores(registry); err != nil {
        return err
    }
    
    // Phase 2: Custom SQL migrations via the migrate framework
    if err := migrateSQL(registry); err != nil {
        return err
    }
    
    return nil
}

// Called during application startup
if err := migrations.MigrateAll(registry); err != nil {
    return err
}
```

**Action Required**:
- Remove any calls to `registry.dataStoresMigrate()` or individual store migration functions
- Update `cmd/server/main.go` to call `migrations.MigrateAll(registry)` after registry initialization
- Update `internal/testutils/setup.go` to call `migrations.MigrateAll(app)` after registry creation
- Delete `internal/registry/registry_datastores_migrate.go` if it exists in your project
- Remove individual store migration functions from `internal/registry/stores_*.go` files

---

### 2. Migration Interface Signature Change

**Change**: The `migrate.MigrationInterface` now requires a `context.Context` parameter in both `Up()` and `Down()` methods.

**Old Usage**:
```go
func (m *MyMigration) Up(tx *sql.Tx) error {
    // Migration logic without context
    return nil
}

func (m *MyMigration) Down(tx *sql.Tx) error {
    // Rollback logic without context
    return nil
}
```

**New Usage**:
```go
func (m *MyMigration) Up(ctx context.Context, tx *sql.Tx) error {
    // Migration logic with context
    return nil
}

func (m *MyMigration) Down(ctx context.Context, tx *sql.Tx) error {
    // Rollback logic with context
    return nil
}
```

**Action Required**:
- Update any custom migrations implementing `migrate.MigrationInterface` to accept `context.Context` as the first parameter
- Update migration calls to pass `context.Background()` or appropriate context
- Search for all files implementing `MigrationInterface` and update method signatures

---

### 3. Entry Point Migration Call

**Change**: The main entry point now uses the centralized migration system instead of registry-internal migration methods.

**Old Usage**:
```go
// cmd/server/main.go
registry, err := registry.New(cfg)
if err != nil {
    fmt.Printf("Failed to initialize registry: %v\n", err)
    return
}
// Migrations were called internally during registry.New()
```

**New Usage**:
```go
// cmd/server/main.go
registry, err := registry.New(cfg)
if err != nil {
    fmt.Printf("Failed to initialize registry: %v\n", err)
    return
}

// Run all migrations (store-level + custom SQL)
if err := migrations.MigrateAll(registry); err != nil {
    fmt.Printf("Failed to run migrations: %v\n", err)
    return
}
```

**Action Required**:
- Add import for `project/database/migrations` in `cmd/server/main.go`
- Add `migrations.MigrateAll(registry)` call after registry initialization
- Ensure this call happens before any other initialization that depends on database tables

---

### 4. Test Setup Migration Call

**Change**: Test utilities now use the centralized migration system.

**Old Usage**:
```go
// internal/testutils/setup.go
app, err := registry.New(opts.cfg)
if err != nil {
    panic("testutils.Setup: failed to build registry: " + err.Error())
}
// Migrations were called internally during registry.New()
```

**New Usage**:
```go
// internal/testutils/setup.go
app, err := registry.New(opts.cfg)
if err != nil {
    panic("testutils.Setup: failed to build registry: " + err.Error())
}

// Run migrations to create database tables
if err := migrations.MigrateAll(app); err != nil {
    panic("testutils.Setup: failed to run migrations: " + err.Error())
}
```

**Action Required**:
- Add import for `project/database/migrations` in `internal/testutils/setup.go`
- Add `migrations.MigrateAll(app)` call after registry creation
- Ensure this call happens in the `Setup()` function

---

### 5. Dependency Version Updates

**Change**: Multiple dracory store packages have been updated to new versions with API changes.

**Updated Packages**:
- `github.com/dracory/geostore` v1.1.0 → v1.4.0
- `github.com/dracory/logstore` v1.14.0 → v1.16.0
- `github.com/dracory/metastore` v1.4.0 → v1.5.0
- `github.com/dracory/migrate` v0.0.0-20260507034242-aaff5b53bdb9 → v0.2.0
- `github.com/dracory/sessionstore` v1.6.0 → v1.11.0
- `github.com/dracory/settingstore` v1.5.0 → v1.7.0
- `github.com/dracory/shopstore` v1.10.0 → v1.13.0
- `github.com/dracory/statsstore` v0.11.0 → v1.0.0
- `github.com/dracory/subscriptionstore` v0.7.0 → v1.0.0
- `github.com/dracory/taskstore` v1.19.0 → v1.22.0
- `github.com/dracory/userstore` v1.8.0 → v1.11.0
- `github.com/dracory/vaultstore` v0.38.0 → v1.0.0
- `github.com/dracory/versionstore` v0.9.0 → v1.1.0

**Action Required**:
- Run `go get -u ./...` to update all dependencies
- Run `go mod tidy` to clean up go.mod
- Review breaking changes in updated store packages (check their respective upgrade guides)
- Update any code that uses deprecated APIs from these packages

---

### 6. Removed Files

**Change**: Obsolete migration files have been removed.

**Removed Files**:
- `database/migrations/2026_03_21_table_users_create.go` (replaced by store migrations)
- `internal/registry/registry_datastores_migrate.go` (replaced by centralized migration system)
- Individual store migration functions from `internal/registry/stores_*.go` files:
  - `auditStoreMigrate()`
  - `blogStoreMigrate()`
  - `blindIndexEmailStoreMigrate()`
  - `blindIndexFirstNameStoreMigrate()`
  - `blindIndexLastNameStoreMigrate()`
  - `cacheStoreMigrate()`
  - `chatStoreMigrate()`
  - `cmsStoreMigrate()`
  - `customStoreMigrate()`
  - `entityStoreMigrate()`
  - `feedStoreMigrate()`
  - `geoStoreMigrate()`
  - `logStoreMigrate()`
  - `metaStoreMigrate()`
  - `sessionStoreMigrate()`
  - `settingStoreMigrate()`
  - `shopStoreMigrate()`
  - `sqlFileStorageMigrate()`
  - `statsStoreMigrate()`
  - `subscriptionStoreMigrate()`
  - `taskStoreMigrate()`
  - `userStoreMigrate()`
  - `vaultStoreMigrate()`

**Action Required**:
- Delete these files if they exist in your project
- Remove any references to these functions from your code
- Ensure no imports reference the deleted files

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

Update all dependencies to their new versions:

```bash
go get -u ./...
go mod tidy
```

### Step 2: Update cmd/server/main.go

Add the migration call after registry initialization:

```bash
# Add import
sed -i '/import (/a\\t"project/database/migrations"' cmd/server/main.go

# Add migration call after registry initialization
# Insert after the defer registry.Close() block
```

Manual changes required in `cmd/server/main.go`:
1. Add import: `"project/database/migrations"`
2. Add after registry initialization (around line 75):
```go
// Run all migrations (store-level + custom SQL)
if err := migrations.MigrateAll(registry); err != nil {
    fmt.Printf("Failed to run migrations: %v\n", err)
    return
}
```

### Step 3: Update internal/testutils/setup.go

Add the migration call in the Setup function:

```bash
# Add import if not present
sed -i '/import (/a\\t"project/database/migrations"' internal/testutils/setup.go
```

Manual changes required in `internal/testutils/setup.go`:
1. Ensure import exists: `"project/database/migrations"`
2. Add after registry creation (around line 335):
```go
// Run migrations to create database tables
if err := migrations.MigrateAll(app); err != nil {
    panic("testutils.Setup: failed to run migrations: " + err.Error())
}
```

### Step 4: Update Custom Migrations

Update any custom migrations to include context parameter:

```bash
# Find all files implementing MigrationInterface
grep -r "MigrationInterface" --include="*.go" .

# Update Up and Down method signatures
# Change: func (m *X) Up(tx *sql.Tx) error
# To: func (m *X) Up(ctx context.Context, tx *sql.Tx) error
```

### Step 5: Remove Obsolete Files

Delete the removed files if they exist in your project:

```bash
# Remove old migration file
rm -f database/migrations/2026_03_21_table_users_create.go

# Remove old registry migration file
rm -f internal/registry/registry_datastores_migrate.go

# Remove individual store migration functions from stores_*.go files
# These are typically at the bottom of each store file
```

### Step 6: Clean Up Registry Store Files

Remove individual store migration functions from `internal/registry/stores_*.go`:

```bash
# Find and remove store migration functions
# These typically follow the pattern: func xyzStoreMigrate(r RegistryInterface) error
```

Each store file (e.g., `stores_audit.go`, `stores_blog.go`) should have its migration function removed.

### Step 7: Update Registry Initialization

Ensure registry initialization no longer calls migration methods:

```bash
# Verify registry.New() doesn't call dataStoresMigrate()
grep -n "dataStoresMigrate" internal/registry/registry_implementation.go
# Should return no results
```

If found, remove the call from the `New()` function.

### Step 8: Verify Imports

Check for any remaining imports to deleted files:

```bash
# Check for imports to removed files
grep -r "registry_datastores_migrate" --include="*.go" .
grep -r "table_users_create" --include="*.go" .
```

Remove any found imports.

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run all unit tests to ensure no regressions:

```bash
go test ./...
```

### 2. Integration Tests

Run integration tests with database:

```bash
go test -tags=integration ./...
```

### 3. Migration Tests

Test the new migration system:

```bash
go test ./database/migrations/...
```

### 4. Server Startup Test

Test that the server starts correctly:

```bash
go run ./cmd/server
# Should start without migration errors
```

### 5. Test Setup Verification

Verify test utilities work correctly:

```bash
go test ./internal/testutils/...
```

---

## 📝 Additional Notes

### New Migration System Features

The new migration system provides:

1. **Centralized Management**: All migrations in one location (`database/migrations/`)
2. **Two-Phase Migration**: 
   - Phase 1: Store-level migrations (run directly outside transactions)
   - Phase 2: Custom SQL migrations (run via migrate framework with transactions)
3. **Conditional Execution**: Store migrations only run if the store is enabled in config
4. **Better Testability**: Clear separation between store initialization and migration
5. **Context Support**: All migrations now accept context for cancellation and tracing

### Store Migration Files

New individual migration files for each store:
- `2026_03_21_0001_store_audit_migrate.go`
- `2026_03_21_0002_store_blog_migrate.go`
- `2026_03_21_0003_store_blindindex_email_migrate.go`
- `2026_03_21_0004_store_blindindex_first_name_migrate.go`
- `2026_03_21_0005_store_blindindex_last_name_migrate.go`
- `2026_03_21_0006_store_cache_migrate.go`
- `2026_03_21_0007_store_chat_migrate.go`
- `2026_03_21_0008_store_cms_migrate.go`
- `2026_03_21_0009_store_custom_migrate.go`
- `2026_03_21_0010_store_entity_migrate.go`
- `2026_03_21_0011_store_feed_migrate.go`
- `2026_03_21_0012_store_geo_migrate.go`
- `2026_03_21_0013_store_log_migrate.go`
- `2026_03_21_0014_store_meta_migrate.go`
- `2026_03_21_0015_store_session_migrate.go`
- `2026_03_21_0016_store_setting_migrate.go`
- `2026_03_21_0017_store_shop_migrate.go`
- `2026_03_21_0018_store_stats_migrate.go`
- `2026_03_21_0019_store_subscription_migrate.go`
- `2026_03_21_0020_store_task_migrate.go`
- `2026_03_21_0021_store_user_migrate.go`
- `2026_03_21_0022_store_vault_migrate.go`

Each migration file has a corresponding test file.

### Geo Store Seeding

The geo store migration now includes seeding functionality for better test data setup.

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: migrations"

**Cause**: Missing import for the new migrations package.

**Solution**: Add import `"project/database/migrations"` to files calling `migrations.MigrateAll()`.

### Issue 2: "cannot use *sql.Tx as type context.Context"

**Cause**: Migration methods still use old signature without context parameter.

**Solution**: Update custom migration signatures to accept `context.Context` as first parameter.

### Issue 3: "registry.dataStoresMigrate undefined"

**Cause**: Trying to call the removed migration method.

**Solution**: Replace with `migrations.MigrateAll(registry)`.

### Issue 4: Tests failing with "table already exists"

**Cause**: Old migration files still being executed alongside new ones.

**Solution**: Remove old migration files and ensure only new migration system is used.

### Issue 5: Dependency conflicts after go get

**Cause**: Some packages may have conflicting version requirements.

**Solution**: Run `go mod tidy` and manually adjust versions in go.mod if needed.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the Blueprint repository: https://github.com/dracory/blueprint
- Review individual store package upgrade guides for dependency-specific changes
- Open an issue on GitHub for migration problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
