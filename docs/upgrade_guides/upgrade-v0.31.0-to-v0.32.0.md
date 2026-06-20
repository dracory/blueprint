# Upgrade Guide: v0.31.0 to v0.32.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.31.0 to v0.32.0.

## Overview

This release consolidates the database configuration files into a single file. The `DatabaseNeatConfig` function and its helpers (`connectionNeatConfig`, `portToInt`) previously in `internal/config/database_neat_config.go` have been merged into `internal/config/database_config.go`. The test file was renamed accordingly. No public API changes were made — the exported function name remains `DatabaseNeatConfig`.

**Key Changes:**
- `internal/config/database_neat_config.go` deleted — contents merged into `internal/config/database_config.go`
- `internal/config/database_neat_config_test.go` renamed to `internal/config/database_config_test.go`
- No breaking changes to exported APIs

---

## ⚠️ Breaking Changes

There are no breaking changes in this release. The exported function `DatabaseNeatConfig` retains its signature and behavior. All internal types (`databaseSettings`, `databaseConnectionSettings`) and unexported helpers (`databaseConfig`, `connectionNeatConfig`, `portToInt`) remain in the same `config` package.

---

## 🔄 Migration Steps

### Step 1: Remove `database_neat_config.go` (If Present)

If your project has a copy of `internal/config/database_neat_config.go`, delete it. Its contents have been merged into `database_config.go`.

```bash
rm internal/config/database_neat_config.go
```

### Step 2: Rename the Test File (If Present)

Rename `database_neat_config_test.go` to `database_config_test.go` to match the new file structure:

```bash
git mv internal/config/database_neat_config_test.go internal/config/database_config_test.go
```

### Step 3: Verify `DatabaseNeatConfig` Is Still Accessible

The function `DatabaseNeatConfig` is now defined in `database_config.go` instead of `database_neat_config.go`. Since both files are in the same `config` package, all imports and call sites remain valid. No changes to calling code are needed.

For example, this call in `internal/app/database_open.go` continues to work unchanged:

```go
neatCfg := config.DatabaseNeatConfig(cfg)
return neatdatabase.New(neatCfg)
```

### Step 4: Update Version Constant

Update the version constant in `internal/config/version.go`:

```go
// Old (v0.31.0)
const Version = "0.31.0"

// New (v0.32.0)
const Version = "0.32.0"
```

### Step 5: Run `go mod tidy` and Tests

```bash
go mod tidy
go test ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run the full test suite to verify everything compiles and passes:

```bash
go test ./...
```

### 2. Database Config Tests

Verify the database configuration tests pass under the new file name:

```bash
go test ./internal/config/... -run TestDatabase
```

### 3. Verify Application Startup

Start the application to confirm the database connection still initializes correctly:

```bash
go run ./cmd/server
```

---

## 📝 Additional Notes

### What Changed

- **File consolidation**: `database_neat_config.go` (114 lines) was merged into `database_config.go`. The resulting file contains all database configuration logic: environment variable loading (`databaseConfig`), struct definitions (`databaseConnectionSettings`, `databaseSettings`), getter methods, and the neat config conversion (`DatabaseNeatConfig`, `connectionNeatConfig`, `portToInt`).
- **Test file rename**: `database_neat_config_test.go` → `database_config_test.go` (content unchanged).
- **No new dependencies**: The `github.com/dracory/neat/database/db` import was already present; it simply moved from `database_neat_config.go` to `database_config.go`.

### What Did NOT Change

- `DatabaseNeatConfig` function name and signature
- `connectionNeatConfig` helper function
- `portToInt` helper function
- All getter methods on `databaseConnectionSettings`
- Any environment variable keys or defaults
- Any calling code in `internal/app/database_open.go` or elsewhere

---

## 🆘 Common Issues and Solutions

### Issue 1: Duplicate symbol errors

**Symptom**: Compile error like `DatabaseNeatConfig redeclared in this block` or `portToInt redeclared in this block`.

**Solution**: This happens if you merged the new `database_config.go` but forgot to delete `database_neat_config.go`. Remove the old file:

```bash
rm internal/config/database_neat_config.go
```

### Issue 2: Missing `database_neat_config_test.go`

**Symptom**: `go test` reports `no files to test` or a missing test file.

**Solution**: The test file was renamed to `database_config_test.go`. Ensure the rename was applied:

```bash
git mv internal/config/database_neat_config_test.go internal/config/database_config_test.go
```

If the old test file was already deleted, simply ensure `database_config_test.go` exists with the test contents.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Open an issue on GitHub for upgrade-specific problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Previous guides reviewed for consistency
