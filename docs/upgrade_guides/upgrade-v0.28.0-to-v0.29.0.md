# Upgrade Guide: v0.28.0 to v0.29.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.28.0 to v0.29.0.

## Overview

This release renames the `internal/registry` package to `internal/app` to improve semantic clarity and developer intuition. The package was originally named `app`, then renamed to `registry`, but `registry` does not accurately reflect its purpose as an application instance holding infrastructure dependencies. This change aligns the naming with the actual pattern used in the codebase and avoids confusion with the classical registry pattern used in other Dracory modules.

**Key Changes:**
- Package renamed from `internal/registry` to `internal/app`
- Interface renamed from `RegistryInterface` to `AppInterface`
- Implementation struct renamed from `registryImplementation` to `appImplementation`
- Constructor call changed from `registry.New()` to `app.New()`
- External package AdminOptions structs now use exported `Registry` field instead of unexported `app` field
- Documentation updated to reflect the new naming

---

## ⚠️ Breaking Changes

### 1. Package Import Path Change

**Change**: The package import path has changed from `internal/registry` to `internal/app`.

**Old Usage**:
```go
import "project/internal/registry"

func myFunction(app registry.RegistryInterface) error {
    cfg := app.GetConfig()
    // ...
}
```

**New Usage**:
```go
import "project/internal/app"

func myFunction(app app.AppInterface) error {
    cfg := app.GetConfig()
    // ...
}
```

**Action Required**:
- Update all imports from `project/internal/registry` to `project/internal/app`
- Update all type references from `registry.RegistryInterface` to `app.AppInterface`
- Update all variable names if they use `registry` as a package alias

**Files to Check**:
- All `.go` files that import `project/internal/registry`
- All `.go` files that reference `registry.RegistryInterface`
- All `.go` files that call `registry.New()`

**Migration Command**:
```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|project/internal/registry|project/internal/app|g' {} \;

# Update type references
find . -type f -name "*.go" -exec sed -i 's|registry\.RegistryInterface|app.AppInterface|g' {} \;

# Update constructor calls
find . -type f -name "*.go" -exec sed -i 's|registry\.New(|app.New(|g' {} \;
```

---

### 2. Interface Name Change

**Change**: The main interface has been renamed from `RegistryInterface` to `AppInterface`.

**Old Usage**:
```go
import "project/internal/registry"

type MyService struct {
    app registry.RegistryInterface
}

func NewMyService(app registry.RegistryInterface) *MyService {
    return &MyService{app: app}
}
```

**New Usage**:
```go
import "project/internal/app"

type MyService struct {
    app app.AppInterface
}

func NewMyService(app app.AppInterface) *MyService {
    return &MyService{app: app}
}
```

**Action Required**:
- Update all type declarations using `registry.RegistryInterface` to `app.AppInterface`
- Update all function signatures using `registry.RegistryInterface` to `app.AppInterface`
- Update all struct fields using `registry.RegistryInterface` to `app.AppInterface`

**Migration Command**:
```bash
# Update interface type references
find . -type f -name "*.go" -exec sed -i 's|RegistryInterface|AppInterface|g' {} \;
```

---

### 3. Constructor Function Change

**Change**: The constructor function has been renamed from `registry.New()` to `app.New()`.

**Old Usage**:
```go
import "project/internal/registry"

func main() {
    cfg := config.NewFromEnv()
    app, err := registry.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer app.Close()
}
```

**New Usage**:
```go
import "project/internal/app"

func main() {
    cfg := config.NewFromEnv()
    app, err := app.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer app.Close()
}
```

**Action Required**:
- Update all calls to `registry.New()` to `app.New()`
- Update variable names if they conflict with the package name (e.g., `registry, err := registry.New()` should become `app, err := app.New()`)

**Migration Command**:
```bash
# Update constructor calls
find . -type f -name "*.go" -exec sed -i 's|registry\.New(|app.New(|g' {} \;
```

---

### 4. External Package AdminOptions Field Change

**Change**: AdminOptions structs in external packages (useradmin, shopadmin, blogadmin, fileadmin, logadmin) now use an exported `Registry` field instead of an unexported `app` field.

**Old Usage**:
```go
import "project/internal/registry"
import "project/pkg/useradmin"

admin, err := useradmin.New(useradmin.AdminOptions{
    app:          registryApp,
    AdminHomeURL: "/admin",
    UserAdminURL: "/admin/users",
    AuthUserID:   getAuthUserID,
})
```

**New Usage**:
```go
import "project/internal/app"
import "project/pkg/useradmin"

admin, err := useradmin.New(useradmin.AdminOptions{
    Registry:     app,
    AdminHomeURL: "/admin",
    UserAdminURL: "/admin/users",
    AuthUserID:   getAuthUserID,
})
```

**Action Required**:
- Update all AdminOptions struct literals to use `Registry:` instead of `app:`
- This affects useradmin, shopadmin, blogadmin, fileadmin, and logadmin packages
- The Blueprint rapid application development (RAD) starter template's controller files have been updated automatically

**Files to Check**:
- `internal/controllers/admin/users/users_controller.go` (already updated in template)
- `internal/controllers/admin/shop/*/` (already updated in template)
- `internal/controllers/admin/blog/blog_controller.go` (already updated in template)
- `internal/controllers/admin/files/file_manager_controller.go` (already updated in template)
- `internal/controllers/admin/logs/logs_controller.go` (already updated in template)
- Any custom controllers that use external admin packages

**Migration Command**:
```bash
# Update AdminOptions field names in struct literals
find . -type f -name "*.go" -exec sed -i 's|app:|Registry:|g' {} \;
```

---

### 5. Directory and File Renames

**Change**: The directory and files have been renamed:

- `internal/registry/` → `internal/app/`
- `internal/registry/registry_interface.go` → `internal/app/app_interface.go`
- `internal/registry/registry_implementation.go` → `internal/app/app_implementation.go`
- `internal/registry/registry_datastores_initialize.go` → `internal/app/app_datastores_initialize.go`
- `internal/registry/registry_datastores_test.go` → `internal/app/app_datastores_test.go`
- `internal/registry/registry_logger_test.go` → `internal/app/app_logger_test.go`
- `internal/registry/registry_close_test.go` → `internal/app/app_close_test.go`

**Action Required**:
- No action required for applications using Blueprint - this is an internal template change
- If you have custom code that directly references these file paths, update them
- If you have test files that reference the old directory structure, update them

---

## 🔄 Migration Steps

### Step 1: Update Package Imports

Update all imports from `internal/registry` to `internal/app`:

```bash
# Update imports in all Go files
find . -type f -name "*.go" -exec sed -i 's|project/internal/registry|project/internal/app|g' {} \;
```

### Step 2: Update Interface Type References

Update all type references from `RegistryInterface` to `AppInterface`:

```bash
# Update interface type references
find . -type f -name "*.go" -exec sed -i 's|RegistryInterface|AppInterface|g' {} \;
```

### Step 3: Update Constructor Calls

Update all constructor calls from `registry.New()` to `app.New()`:

```bash
# Update constructor calls
find . -type f -name "*.go" -exec sed -i 's|registry\.New(|app.New(|g' {} \;
```

### Step 4: Update AdminOptions Field Names

Update AdminOptions struct literals to use `Registry:` instead of `app:`:

```bash
# Update AdminOptions field names
find . -type f -name "*.go" -exec sed -i 's|app:|Registry:|g' {} \;
```

### Step 5: Update Variable Names (if needed)

If you have variable names that conflict with the package name:

```bash
# Manual review required for variable name changes
# Example: registry, err := app.New() → app, err := app.New()
```

### Step 6: Update Version Constant

Update the version constant in `internal/config/version.go`:

```go
const Version = "0.29.0"
```

### Step 7: Build and Test

Build the application to ensure all changes are compatible:

```bash
go build ./...
```

Run tests:

```bash
go test ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run all unit tests to ensure no regressions:

```bash
go test ./...
```

### 2. Integration Tests

Run integration tests if applicable:

```bash
go test -tags=integration ./...
```

### 3. Build Verification

Verify the application builds successfully:

```bash
go build -o ./bin/server ./cmd/server
```

### 4. Manual Testing

Test the application manually:

```bash
# Start the server
go run ./cmd/server

# Test key functionality:
# - Application starts successfully
# - Database connections work
# - Admin interfaces load
# - Controllers function correctly
```

---

## 📝 Additional Notes

### Rationale for the Change

The package was originally named `app`, then renamed to `registry`. However, `registry` does not accurately reflect its purpose:

1. **Semantic Mismatch**: Other Dracory modules use `registry` for lookup tables (ork, liveflux, neat), but Blueprint's struct is an application instance holding infrastructure dependencies.

2. **Poor Developer Intuition**: In tests, developers naturally write `app := registry.New()` which is confusing. The name `app` is more intuitive.

3. **Alignment with Purpose**: The struct is booted once at startup and holds infrastructure dependencies (database, cache, session, logger, etc.), making "app" the more accurate name.

### Benefits

- Better semantic clarity in the codebase
- Improved developer intuition, especially in tests
- Alignment with the actual pattern used (application instance, not lookup table)
- Consistency with the original naming before the rename to `registry`

### No Functional Changes

This is purely a naming change with no functional changes to the API behavior. All methods and interfaces remain the same, only the names have changed.

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: registry" after migration

**Symptom**: Compilation errors about undefined `registry` package or types.

**Solution**: Ensure all imports have been updated from `project/internal/registry` to `project/internal/app`:
```bash
find . -type f -name "*.go" -exec sed -i 's|project/internal/registry|project/internal/app|g' {} \;
```

### Issue 2: "undefined: RegistryInterface" after migration

**Symptom**: Compilation errors about undefined `RegistryInterface` type.

**Solution**: Update all type references to `AppInterface`:
```bash
find . -type f -name "*.go" -exec sed -i 's|RegistryInterface|AppInterface|g' {} \;
```

### Issue 3: "undefined: registry.New" after migration

**Symptom**: Compilation errors about undefined `registry.New` function.

**Solution**: Update constructor calls to `app.New()`:
```bash
find . -type f -name "*.go" -exec sed -i 's|registry\.New(|app.New(|g' {} \;
```

### Issue 4: "unknown field app" in AdminOptions

**Symptom**: Compilation errors about unknown field `app` in AdminOptions struct literals.

**Solution**: Update field names from `app:` to `Registry:`:
```bash
find . -type f -name "*.go" -exec sed -i 's|app:|Registry:|g' {} \;
```

### Issue 5: Variable name conflicts with package name

**Symptom**: Confusion between variable names and package names (e.g., `app := app.New()`).

**Solution**: Rename variables to avoid conflicts:
```go
// Old
registry, err := registry.New(cfg)

// New
app, err := app.New(cfg)
```

### Issue 6: External package references not updated

**Symptom**: External packages (useradmin, shopadmin, etc.) still reference old field names.

**Solution**: Ensure AdminOptions struct literals use `Registry:` instead of `app:`. The Blueprint starter template's controllers have been updated automatically.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Review the [rename proposal](docs/proposals/rename-registry-to-app.md) for detailed rationale
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
- [x] Git tag verified for previous version (v0.28.0)
- [x] Previous guides reviewed for consistency
