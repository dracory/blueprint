# Upgrade Guide: v0.10.0 to v0.11.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.10.0 to v0.11.0.

## ⚠️ Major Breaking Changes

### 1. Application Entry Point Moved
**Change**: `main.go` moved from project root to `cmd/server/main.go`

**Old Location**:
```bash
project/
├── main.go  # ❌ Removed
```

**New Location**:
```bash
project/
├── cmd/
│   └── server/
│       └── main.go  # ✅ New entry point
```

**Action Required**:
- Update any build scripts or deployment configurations that reference `main.go`
- Update import paths if they reference the old location

### 2. Registry Architecture Refactor
**Change**: Application renamed to Registry throughout codebase

**Old Usage**:
```go
// Old interface name
types.RegistryInterface
types.ConfigInterface

// Old struct names
Application
Config
```

**New Usage**:
```go
// New interface names
registry.RegistryInterface
config.ConfigInterface

// New struct names  
registryImplementation
configImplementation
```

**Action Required**:
- Update all imports from `types` package to specific packages:
  - `types.RegistryInterface` → `registry.RegistryInterface`
  - `types.ConfigInterface` → `config.ConfigInterface`
- Update type references in code
- Update any custom implementations

### 3. Package Reorganization
**Change**: Major packages moved and renamed for consistency

**Old Package Structure**:
```go
internal/app/           // ❌ Moved
internal/app_interface.go  // ❌ Renamed
internal/config_model.go   // ❌ Renamed
```

**New Package Structure**:
```go
internal/registry/        // ✅ New location
internal/registry_interface.go  // ✅ Renamed
internal/config_implementation.go  // ✅ Renamed
```

**Action Required**:
- Update import paths:
  - `project/internal/app` → `project/internal/registry`
  - Update any references to `app_interface.go` → `registry_interface.go`
  - Update any references to `config_model.go` → `config_implementation.go`

### 4. Task Store API Parameter Order Change
**Change**: `TaskDefinitionEnqueueByAlias` parameter order swapped

**Old Signature**:
```go
TaskDefinitionEnqueueByAlias(context, queueName, alias, params)
```

**New Signature**:
```go
TaskDefinitionEnqueueByAlias(context, alias, queueName, params)
```

**Action Required**:
- Update all calls to `TaskDefinitionEnqueueByAlias` to swap `queueName` and `alias` parameters
- This affects all task enqueue calls throughout the codebase

### 5. Cache Architecture Changes
**Change**: Caches moved from global singletons to registry instances

**Old Usage**:
```go
// Global cache access
cache.GetGlobalCache()
```

**New Usage**:
```go
// Registry-scoped cache access
registry.GetCacheStore()
registry.GetMemoryCache()
registry.GetFileCache()
```

**Action Required**:
- Replace any global cache access with registry-based access
- Update cache initialization code
- Remove any package-level cache variables

### 6. Database Method Rename
**Change**: `GetDB()` renamed to `GetDatabase()`

**Old Usage**:
```go
registry.GetDB()
```

**New Usage**:
```go
registry.GetDatabase()
```

**Action Required**:
- Update all calls from `GetDB()` to `GetDatabase()`

### 7. Store Configuration Changes
**Change**: Store configuration moved from environment variables to compile-time constants

**Old Configuration**:
```go
// Environment-based store configuration
if os.Getenv("BLOG_STORE_USED") == "yes" {
    // Initialize blog store
}
```

**New Configuration**:
```go
// Compile-time store configuration via config interface
if registry.GetConfig().GetBlogStoreUsed() {
    // Initialize blog store
}
```

**Action Required**:
- Update store initialization logic to use config interface methods
- Remove direct environment variable checks for store configuration
- Use the standardized `Get[StoreName]Used()` methods

### 8. Dependency Updates
**Change**: Major dependency updates and removals

**Removed Dependencies**:
- `github.com/gouniverse/validator`
- `github.com/gouniverse/responses`

**Updated Dependencies**:
- `github.com/dracory/auth` v0.29.0
- `github.com/dracory/base` v0.26.0
- All dracory store packages updated to latest versions
- AWS SDK v2 and Google Cloud libraries updated

**Action Required**:
- Remove any usage of removed `gouniverse` packages
- Update any custom code that relied on specific dependency APIs
- Test all integrations with updated dependencies

### 9. Routes Renamed
**Change**: `Routes` renamed to `Router`

**Old Usage**:
```go
routes.Routes(registry)
```

**New Usage**:
```go
routes.Router(registry)
```

**Action Required**:
- Update all calls from `Routes()` to `Router()`

### 10. Environment Variable Standardization
**Change**: Environment variable naming standardized

**Old Variables**:
```bash
# Various inconsistent naming patterns
```

**New Variables**:
```bash
# Standardized naming (see .env.example)
APP_ENV="development"
APP_NAME="YOUR APP NAME"
# ... etc
```

**Action Required**:
- Review `.env.example` for standardized variable names
- Update any environment variable references to match new naming
- Update deployment configurations

## 🔄 Migration Steps

### Step 1: Update Entry Point
1. Move any custom logic from old `main.go` to `cmd/server/main.go`
2. Update build scripts and deployment configurations

### Step 2: Update Imports
```bash
# Find and replace old imports
find . -name "*.go" -exec sed -i 's|project/internal/app|project/internal/registry|g' {} +
find . -name "*.go" -exec sed -i 's|types\.RegistryInterface|registry.RegistryInterface|g' {} +
find . -name "*.go" -exec sed -i 's|types\.ConfigInterface|config.ConfigInterface|g' {} +
```

### Step 3: Update Task Store Calls
```bash
# Update TaskDefinitionEnqueueByAlias calls
# Manual review required for each call site
```

### Step 4: Update Cache Access
```bash
# Replace global cache access with registry-based access
# Manual review required for each cache usage
```

### Step 5: Update Database Access
```bash
# Replace GetDB() calls with GetDatabase()
find . -name "*.go" -exec sed -i 's|\.GetDB()|\.GetDatabase()|g' {} +
```

### Step 6: Update Store Configuration
```bash
# Replace environment variable checks with config interface calls
# Manual review required for each store initialization
```

### Step 7: Update Routes
```bash
# Replace Routes() calls with Router()
find . -name "*.go" -exec sed -i 's|routes\.Routes(|routes.Router(|g' {} +
```

### Step 8: Test Everything
```bash
# Run full test suite
go test ./...

# Run application
go run ./cmd/server
```

## 🧪 Testing After Migration

1. **Unit Tests**: Ensure all existing tests pass
2. **Integration Tests**: Test all store integrations
3. **Task Queue**: Verify task enqueue and execution
4. **Cache**: Test cache functionality
5. **Database**: Verify database connections and operations
6. **Web Server**: Test HTTP endpoints and middleware

## 📝 Additional Notes

### New Features Added
- Enhanced log manager with UI components
- EasyMDE markdown editor support
- Improved component app initialization
- Better error handling and logging

### Removed Features
- Deprecated `post_update_v1` controller
- Commented-out replace directives for cachestore and sessionstore
- Global cache singletons

### Configuration Changes
- Store configuration now centralized in config package
- Environment variable naming standardized
- Better validation and error handling

## 🆘 Common Issues

### Issue: Import Path Not Found
**Solution**: Update import paths from `internal/app` to `internal/registry`

### Issue: Task Enqueue Fails
**Solution**: Check parameter order in `TaskDefinitionEnqueueByAlias` calls

### Issue: Cache Not Working
**Solution**: Replace global cache access with registry-based access

### Issue: Store Not Initializing
**Solution**: Update store configuration to use config interface methods instead of environment variables

## 📞 Support

For issues during migration:
1. Check this guide first
2. Review the updated documentation in `docs/overview.md`
3. Check the architectural review in `docs/review.md`
4. Run tests to identify specific issues

---

**Version**: v0.11.0  
**Previous Version**: v0.10.0  
**Release Date**: January 2026
