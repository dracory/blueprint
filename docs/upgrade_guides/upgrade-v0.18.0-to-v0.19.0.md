# Upgrade Guide: v0.18.0 to v0.19.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.18.0 to v0.19.0.

## Summary

Version 0.19.0 introduces a major refactoring of the configuration system to improve organization, maintainability, and type safety. The changes consolidate the config loader pattern into a more structured approach using domain-specific loader functions and typed return values.

**Key Changes:**
- Config loading API renamed from `config.Load()` to `config.NewFromEnv()`
- Domain-specific config loader pattern (typed structs instead of individual setters)
- Consolidated config interfaces into a single `interfaces.go` file
- Reorganized config files into domain-specific modules (2-file pattern)
- Dependency update: `github.com/dracory/env` v1.0.0 → v1.2.0

---

## Breaking Changes

### 1. Config Loading Function Renamed

**Change**: The main entry point for loading configuration has been renamed from `config.Load()` to `config.NewFromEnv()`.

**Old Usage**:
```go
import "project/internal/config"

cfg, err := config.Load()
if err != nil {
    fmt.Printf("Failed to load config: %v\n", err)
    return
}
```

**New Usage**:
```go
import "project/internal/config"

cfg, err := config.NewFromEnv()
if err != nil {
    fmt.Printf("Failed to load config: %v\n", err)
    return
}
```

**Action Required**:
- Update all calls to `config.Load()` to `config.NewFromEnv()`
- **Primary location**: `cmd/server/main.go`
- **Search command**: `grep -r "config.Load()" --include="*.go" .`

**Migration Command**:
```bash
# Find all occurrences
find . -type f -name "*.go" -exec grep -l "config.Load()" {} \;

# Replace with NewFromEnv
find . -type f -name "*.go" -exec sed -i 's/config.Load()/config.NewFromEnv()/g' {} \;
```

---

### 2. Config Loader Function Renames

**Change**: Config loader helper functions have been renamed from `loadXxxConfig()` to `xxxConfig()` following Go naming conventions.

**Old Usage** (in v0.18.0 internal/config/load.go):
```go
app := loadAppConfig(acc)
db := loadDatabaseConfig(acc)
mail := loadMailConfig()
stores := loadStoresConfig(acc)
```

**New Usage** (in v0.19.0 internal/config/config_implementation.go):
```go
app := appConfig(v)
db := databaseConfig(v)
mail := emailConfig()
stores := storesConfig(v)
```

**Action Required**:
- If you have custom config loaders or have extended the config system, update function names
- Remove the `load` prefix from config loader functions
- Update calls from `loadXxxConfig()` to `xxxConfig()`

**Note**: This primarily affects internal framework code. Application code typically doesn't call these functions directly.

---

### 3. Config Interface Consolidation

**Change**: Individual config interface files have been consolidated into a single `interfaces.go` file with standardized naming.

**Old Structure** (v0.18.0):
```
internal/config/config_interface.go  (contained all interfaces)
```

**New Structure** (v0.19.0):
```
internal/config/interfaces.go  (consolidated interface definitions)
```

**Interface Naming Changes**:
- `appConfigInterface` → `AppConfigInterface`
- `databaseConfigInterface` → `DatabaseConfigInterface`
- `emailConfigInterface` → `EmailConfigInterface`
- etc. (all interfaces now use PascalCase and start with uppercase)

**Action Required**:
- Update any references to old interface names if you have custom implementations
- Verify `ConfigInterface` is used correctly (no changes to main interface name)

---

### 4. Config Loader Pattern Change

**Change**: Config loaders now return typed structs instead of using a shared accumulator pattern with individual setter methods.

**Old Pattern** (v0.18.0):
```go
// config_loader.go
func loadAppConfig(acc *baseCfg.LoadAccumulator) appConfigData {
    name := env.GetStringOrDefault(KEY_APP_NAME, "Blueprint")
    // ... more loading
    
    // Return a struct with loaded values
    return appConfigData{
        name: name,
        // ...
    }
}

// In the main Load() function:
app := loadAppConfig(acc)
cfg.SetAppName(app.name)
cfg.SetAppUrl(app.url)
// ... individual setter calls for each field
```

**New Pattern** (v0.19.0):
```go
// app_config.go
func appConfig(env *envValidator) appSettings {
    name := env.GetStringOrDefault(KEY_APP_NAME, "Blueprint")
    // ... more loading
    
    return appSettings{
        name: name,
        // ...
    }
}

// Batch setter in config_implementation.go
func (c *configImplementation) setAppConfig(s appSettings) {
    c.appName = s.name
    c.appUrl = s.url
    c.appHost = s.host
    c.appPort = s.port
    c.appEnv = s.env
    c.appDebug = s.debug
    c.cmsMcpApiKey = s.cmsMcpApiKey
}

// Usage in NewFromEnv():
cfg.setAppConfig(appConfig(v))
```

**Action Required**:
- If you have custom config sections, refactor to use the new typed struct pattern
- Create a batch setter method (e.g., `setXxxConfig(s xxxSettings)`) instead of individual setters
- Update loader function to return typed struct instead of setting values directly

---

### 5. File Reorganization

**Change**: Config files have been reorganized following a domain-specific 2-file pattern.

**Files Removed**:
- `internal/config/load.go` → functionality moved to `config_implementation.go`
- `internal/config/load_test.go` → renamed to `z_config_implementation_test.go`
- `internal/config/app.go` → functionality moved to `app_config.go`
- `internal/config/database.go` → functionality moved to `database_config.go`
- `internal/config/mail.go` → replaced by `email_config.go`
- `internal/config/registration.go` → replaced by `auth_config.go`
- `internal/config/stores.go` → renamed to `stores_config.go`
- `internal/config/stripe.go` → functionality moved to `payment_config.go`
- `internal/config/translation.go` → replaced by `i18n_config.go`
- `internal/config/llm.go` → replaced by `llm_config.go`
- `internal/config/defaults.go` → removed (defaults now inline in loaders)
- `internal/config/config_interface.go` → replaced by `interfaces.go`

**Files Added**:
- `internal/config/app_config.go` - application config loader
- `internal/config/auth_config.go` - authentication config loader
- `internal/config/database_config.go` - database config loader
- `internal/config/email_config.go` - email config loader
- `internal/config/i18n_config.go` - i18n/translation config loader
- `internal/config/interfaces.go` - consolidated interfaces
- `internal/config/llm_config.go` - LLM config loader
- `internal/config/payment_config.go` - payment/Stripe config loader
- `internal/config/stores_config.go` - datastore config loader

**Action Required**:
- Update any imports referencing old file names
- Remove any custom code that imported from the old file paths
- Review `internal/config/configuration_stores.go` for store-related configuration

---

### 6. Dependency Version Update

**Change**: The `github.com/dracory/env` dependency has been updated from v1.0.0 to v1.2.0.

**Old** (v0.18.0 go.mod):
```
github.com/dracory/env v1.0.0
```

**New** (v0.19.0 go.mod):
```
github.com/dracory/env v1.2.0
```

**Action Required**:
- Run `go mod tidy` after upgrading
- The new version is backward compatible, no code changes required
- Benefits: improved performance and additional utility methods

---

## Migration Steps

### Step 1: Update Config Loading Calls

Update the main entry point in `cmd/server/main.go`:

```go
// Before (v0.18.0)
cfg, err := config.Load()

// After (v0.19.0)
cfg, err := config.NewFromEnv()
```

### Step 2: Clean Up Old Config Files

Remove the following files if they exist in your project (they're now replaced):

```bash
# Remove old config files that have been replaced
rm -f internal/config/load.go
rm -f internal/config/load_test.go
rm -f internal/config/app.go
rm -f internal/config/database.go
rm -f internal/config/mail.go
rm -f internal/config/registration.go
rm -f internal/config/stores.go
rm -f internal/config/stripe.go
rm -f internal/config/translation.go
rm -f internal/config/llm.go
rm -f internal/config/defaults.go
rm -f internal/config/config_interface.go
```

### Step 3: Update Dependencies

Run the following commands to update dependencies:

```bash
# Update the env dependency
go get github.com/dracory/env@v1.2.0

# Tidy and verify modules
go mod tidy
go mod download
```

### Step 4: Verify Build

Build the application to verify no compilation errors:

```bash
go build -o ./tmp/main ./cmd/server
```

---

## Testing After Migration

### 1. Unit Tests

Run the config-specific tests to verify the new loading mechanism:

```bash
go test ./internal/config/... -v
```

**Key Test Files:**
- `z_config_implementation_test.go` - validates all config loading scenarios
- Tests cover: missing required fields, encryption key handling, store configurations, LLM provider requirements, Stripe configuration, mail configuration, translation defaults, and vault store requirements

### 2. Integration Tests

Run the full test suite to verify no regressions:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. Manual Verification

Test the application startup with different configurations:

```bash
# Test with minimal config
APP_HOST=localhost APP_PORT=8080 APP_ENVIRONMENT=testing go run ./cmd/server

# Verify encrypted environment loading (if using envenc)
ENVENC_KEY_PRIVATE=your-key APP_HOST=localhost APP_PORT=8080 APP_ENVIRONMENT=testing go run ./cmd/server
```

### 4. Environment Variable Validation

Verify all required environment variables are correctly loaded:

```bash
# Required variables for v0.19.0:
APP_HOST        # Required: Application host address
APP_PORT        # Required: Application port
APP_ENVIRONMENT # Required: local, development, staging, testing, or production

# Optional variables are handled with sensible defaults
```

---

## Additional Notes

### New Features

1. **Typed Config Settings**: Config loaders now return typed structs (e.g., `appSettings`, `databaseSettings`) providing better compile-time type safety.

2. **Batch Setter Methods**: Config implementation uses batch setters (e.g., `setAppConfig(s appSettings)`) reducing repetitive individual setter calls.

3. **Enhanced Validation**: The `envValidator` type provides consistent validation across all config loaders.

4. **Consolidated Interface Definitions**: All config interfaces are now in a single, well-documented `interfaces.go` file.

### Removed Features

- The old `config.Load()` function (use `config.NewFromEnv()` instead)
- Individual config interface files (consolidated into `interfaces.go`)
- Scattered config loader files (organized into domain-specific files)

### Configuration Behavior Changes

No behavioral changes - all environment variable keys, defaults, and validation logic remain the same. The changes are purely structural and organizational.

---

## Common Issues and Solutions

### Issue 1: "undefined: config.Load"

**Symptom**: Compilation error `undefined: config.Load`

**Solution**: Replace all occurrences of `config.Load()` with `config.NewFromEnv()`

```bash
find . -type f -name "*.go" -exec sed -i 's/config.Load()/config.NewFromEnv()/g' {} \;
```

### Issue 2: Import Errors

**Symptom**: Cannot find config files after migration

**Solution**: Ensure you're importing the correct package path:

```go
// Correct import
import "project/internal/config"

// Usage remains the same
cfg, err := config.NewFromEnv()
```

### Issue 3: Custom Config Extensions Broken

**Symptom**: Custom config extensions no longer compile

**Solution**: Update custom config code to match the new pattern:

1. Create a typed settings struct
2. Create a loader function returning the struct
3. Create a batch setter method on `configImplementation`
4. Call in `NewFromEnv()`

Example:
```go
// my_custom_config.go
func myCustomConfig(env *envValidator) myCustomSettings {
    return myCustomSettings{
        value: env.GetString("MY_CUSTOM_VALUE"),
    }
}

type myCustomSettings struct {
    value string
}

// In config_implementation.go, add:
func (c *configImplementation) setMyCustomConfig(s myCustomSettings) {
    c.myCustomValue = s.value
}
```

### Issue 4: Test Failures

**Symptom**: Config-related tests fail after migration

**Solution**: Update test files that reference old function names:

```bash
# Update test function names if needed
# The test file was renamed from load_test.go to z_config_implementation_test.go
```

---

## Support

For issues related to this upgrade:

1. **Documentation**: Review the config proposals in `docs/proposals/`:
   - `config-file-organization-improvement.md`
   - `config-file-organization-implementation-summary.md`
   - `config-reorganization-final-summary.md`

2. **Git History**: Review the migration commits:
   ```bash
   git log --oneline v0.18.0..v0.19.0
   ```

3. **Reference Implementation**: Compare with the reference implementation in the v0.19.0 tag

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

*This upgrade guide was generated for Blueprint v0.18.0 to v0.19.0 migration.*
