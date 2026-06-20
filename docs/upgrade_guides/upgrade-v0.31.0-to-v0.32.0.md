# Upgrade Guide: v0.31.0 to v0.32.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.31.0 to v0.32.0.

## Overview

This release consolidates the database configuration files into a single file and moves the task constants package under `internal/tasks`. The `DatabaseNeatConfig` function and its helpers (`connectionNeatConfig`, `portToInt`) previously in `internal/config/database_neat_config.go` have been merged into `internal/config/database_config.go`. The test file was renamed accordingly. The `internal/taskconstants` package was moved to `internal/tasks/constants` and the package was renamed from `taskconstants` to `constants`.

**Key Changes:**
- `internal/config/database_neat_config.go` deleted — contents merged into `internal/config/database_config.go`
- `internal/config/database_neat_config_test.go` renamed to `internal/config/database_config_test.go`
- `internal/taskconstants` moved to `internal/tasks/constants` (package renamed from `taskconstants` to `constants`)
- `internal/taskconstants/taskconstants.go` renamed to `internal/tasks/constants/constants.go`
- All import paths updated from `project/internal/taskconstants` to `project/internal/tasks/constants`
- All package references updated from `taskconstants.X` to `constants.X`

---

## ⚠️ Breaking Changes

### 1. Task Constants Package Renamed

**Change**: The `internal/taskconstants` package was moved to `internal/tasks/constants` and the Go package name was renamed from `taskconstants` to `constants`. All import paths and qualified references must be updated.

**Old Usage**:
```go
// v0.31.0
import "project/internal/taskconstants"

func (h *helloWorldTask) Alias() string {
    return taskconstants.HelloWorldTaskAlias
}
```

**New Usage**:
```go
// v0.32.0
import "project/internal/tasks/constants"

func (h *helloWorldTask) Alias() string {
    return constants.HelloWorldTaskAlias
}
```

**Action Required**:
- Replace all imports of `project/internal/taskconstants` with `project/internal/tasks/constants`.
- Replace all qualified references from `taskconstants.X` to `constants.X`.
- If you had a local copy of `internal/taskconstants/`, move it to `internal/tasks/constants/` and rename the package declaration from `package taskconstants` to `package constants`.

**Files to Check**:
- All task handler files in `internal/tasks/*/` (already updated in template)
- `pkg/useradmin/user_update/handle_user_update_ajax.go` (already updated in template)
- Any custom code importing `project/internal/taskconstants`

---

## 🔄 Migration Steps

### Step 1: Move and Rename Task Constants Package

Move `internal/taskconstants/` to `internal/tasks/constants/` and rename the package:

```bash
git mv internal/taskconstants internal/tasks/constants
git mv internal/tasks/constants/taskconstants.go internal/tasks/constants/constants.go
```

Update the package declaration in `constants.go`:

```go
// Old
package taskconstants

// New
package constants
```

Update all import paths and references:

```bash
# Find files with old import
grep -rn 'project/internal/taskconstants' --include='*.go' .
```

Replace `project/internal/taskconstants` with `project/internal/tasks/constants` and `taskconstants.` with `constants.` in all matching files.

### Step 2: Remove `database_neat_config.go` (If Present)

If your project has a copy of `internal/config/database_neat_config.go`, delete it. Its contents have been merged into `database_config.go`.

```bash
rm internal/config/database_neat_config.go
```

### Step 3: Rename the Test File (If Present)

Rename `database_neat_config_test.go` to `database_config_test.go` to match the new file structure:

```bash
git mv internal/config/database_neat_config_test.go internal/config/database_config_test.go
```

### Step 4: Verify `DatabaseNeatConfig` Is Still Accessible

The function `DatabaseNeatConfig` is now defined in `database_config.go` instead of `database_neat_config.go`. Since both files are in the same `config` package, all imports and call sites remain valid. No changes to calling code are needed.

For example, this call in `internal/app/database_open.go` continues to work unchanged:

```go
neatCfg := config.DatabaseNeatConfig(cfg)
return neatdatabase.New(neatCfg)
```

### Step 5: Update Version Constant

Update the version constant in `internal/config/version.go`:

```go
// Old (v0.31.0)
const Version = "0.31.0"

// New (v0.32.0)
const Version = "0.32.0"
```

### Step 6: Run `go mod tidy` and Tests

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

### 2. Task Constants Tests

Verify the task constants package compiles and all task alias references are valid:

```bash
go test ./internal/tasks/...
```

### 3. Database Config Tests

Verify the database configuration tests pass under the new file name:

```bash
go test ./internal/config/... -run TestDatabase
```

### 4. Verify Application Startup

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
- **Task constants relocation**: `internal/taskconstants` moved to `internal/tasks/constants` to colocate with the task implementations. Package renamed from `taskconstants` to `constants` for brevity. File renamed from `taskconstants.go` to `constants.go`.

### What Did NOT Change

- `DatabaseNeatConfig` function name and signature
- `connectionNeatConfig` helper function
- `portToInt` helper function
- All getter methods on `databaseConnectionSettings`
- Any environment variable keys or defaults
- Any calling code in `internal/app/database_open.go` or elsewhere
- Task alias constant names and values (e.g., `HelloWorldTaskAlias = "HelloWorldTask"`)

---

## 🆘 Common Issues and Solutions

### Issue 1: Duplicate symbol errors

**Symptom**: Compile error like `DatabaseNeatConfig redeclared in this block` or `portToInt redeclared in this block`.

**Solution**: This happens if you merged the new `database_config.go` but forgot to delete `database_neat_config.go`. Remove the old file:

```bash
rm internal/config/database_neat_config.go
```

### Issue 2: `cannot find package "project/internal/taskconstants"`

**Symptom**: Compile error like `cannot find package "project/internal/taskconstants"`.

**Solution**: The package was moved to `internal/tasks/constants`. Update all imports from `project/internal/taskconstants` to `project/internal/tasks/constants` and rename all `taskconstants.X` references to `constants.X`.

### Issue 3: Missing `database_neat_config_test.go`

**Symptom**: `go test` reports `no files to test` or a missing test file.

**Solution**: The test file was renamed to `database_config_test.go`. Ensure the rename was applied:

```bash
git mv internal/config/database_neat_config_test.go internal/config/database_config_test.go
```

If the old test file was already deleted, simply ensure `database_config_test.go` exists with the test contents.

### Issue 4: `constants` package name conflicts

**Symptom**: Compile error like `constants redeclared in this block` or ambiguous import.

**Solution**: If your file already imports another package named `constants`, use a package alias:

```go
import taskconstants "project/internal/tasks/constants"
```

Then reference as `taskconstants.HelloWorldTaskAlias`.

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
