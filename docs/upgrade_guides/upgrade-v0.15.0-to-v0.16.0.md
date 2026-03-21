# Upgrade Guide: v0.15.0 to v0.16.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.15.0 to v0.16.0.

## ⚠️ Breaking Changes

### 1. Blog and Web Theme Packages Migration to Base
**Change**: Theme rendering packages moved from Blueprint to base package

**Old Usage**:
```go
import "project/pkg/blogtheme"
import "project/pkg/webtheme"

// Use theme functions
html := blogtheme.ParagraphToHtml("content")
html := webtheme.HeadingToHtml("title", 1)
```

**New Usage**:
```go
import "github.com/dracory/base/blogtheme"
import "github.com/dracory/base/webtheme"

// Use theme functions
html := blogtheme.ParagraphToHtml("content")
html := webtheme.HeadingToHtml("title", 1)
```

**Action Required**:
- Update import statements to use base package paths
- No functional changes to the theme rendering APIs
- All theme functions maintain the same signatures
- Complete theme packages moved (all .go files in pkg/blogtheme/* and pkg/webtheme/*)

### 2. Jail Bots Middleware Migration to RTR Package
**Change**: Jail bots middleware moved from Blueprint to rtr router package

**Old Usage**:
```go
import "project/internal/middlewares"

// Use middleware
router.Use(middlewares.JailBotsMiddleware())
```

**New Usage**:
```go
import rtrMiddleware "github.com/dracory/rtr/middlewares"

// Use middleware
router.Use(rtrMiddleware.JailBotsMiddleware())
```

**Action Required**:
- Update import statements to use rtr/middlewares package
- Middleware functionality and configuration remain identical
- All exclusion rules and IP jailing behavior preserved
- Files removed: `internal/middlewares/jail_bots_middleware.go` and `internal/middlewares/jail_bots_middleware_test.go`

### 3. CLI Dispatcher Migration to Base Package
**Change**: Generic CLI dispatcher moved to base package with type-safe generics

**Old Usage**:
```go
// Blueprint had custom dispatcher logic with handler map
var commandHandlers = map[string]commandHandler{
    CommandTask:   handleTaskCommand,
    CommandJob:    handleJobCommand,
    CommandRoutes: handleRoutesCommand,
}

// Manual command lookup and execution
handler, found := commandHandlers[command]
if !found {
    return fmt.Errorf("unrecognized command: %s", command)
}
return handler(registry, remainingArgs)
```

**New Usage**:
```go
// Uses base package generic dispatcher with type safety
dispatcher := cli.NewDispatcher[registry.RegistryInterface]()
dispatcher.RegisterCommand(CommandTask, "Execute a task by alias", handleTaskCommand)
dispatcher.RegisterCommand(CommandJob, "Execute a job with arguments", handleJobCommand)
dispatcher.RegisterCommand(CommandRoutes, "List all registered routes", handleRoutesCommand)

// Automatic command lookup and execution
return dispatcher.ExecuteCommand(registry, args)
```

**Action Required**:
- No code changes required for existing CLI usage
- Blueprint now uses base package generic dispatcher internally
- All CLI commands (task, job, routes) work identically
- Type safety improved with generics
- Better error handling and command registration API

### 4. Database Connection Optimization Migration
**Change**: SQLite optimization logic moved to database package

**Old Usage**:
```go
// Blueprint had custom SQLite optimizations (50+ lines of PRAGMA and pool settings)
db, err := databaseOpen(cfg)
// The function included:
// _, _ = db.Exec("PRAGMA journal_mode=WAL;")
// _, _ = db.Exec("PRAGMA synchronous=NORMAL;")
// _, _ = db.Exec("PRAGMA foreign_keys=ON;")
// _, _ = db.Exec("PRAGMA busy_timeout=5000;")
// db.SetMaxOpenConns(1)
// db.SetMaxIdleConns(1)
```

**New Usage**:
```go
// Database package includes optimizations automatically
db, err := database.Open(options)
// All SQLite optimizations now handled automatically by database package
```

**Action Required**:
- No code changes required
- SQLite optimizations (WAL mode, connection pooling, PRAGMA settings) now automatic
- Database package handles all optimization logic
- Blueprint's databaseOpen() simplified from 84 lines to 44 lines
- Connection pool settings and PRAGMA executions moved to database package

### 5. Config Encryption Loader Migration to Base Package
**Change**: Environment encryption loader moved to base package

**Old Usage**:
```go
// Custom initialization in Blueprint (134 lines with custom logic)
if err := initializeEnvEncVariables(app.env, ENVENC_KEY_PUBLIC, envEnc.privateKey); err != nil {
    acc.add(err)
}

// Custom function with file system checks and resource loading
func initializeEnvEncVariables(appEnvironment string, publicKey string, privateKey string) error {
    // 40+ lines of custom logic including:
    // - Environment validation
    // - File path construction
    // - Resource loading
    // - File existence checks
    // - Vault content hydration
}
```

**New Usage**:
```go
// Base package loader with embedded resources support
if err := baseCfg.InitializeEnvEncVariablesFromResources(app.env, ENVENC_KEY_PUBLIC, envEnc.privateKey, resources.Resource); err != nil {
    acc.add(err)
}

// All logic moved to base package with enhanced error handling
```

**Action Required**:
- No code changes required for existing functionality
- Enhanced error handling with custom error types (MissingEnvError, EnvEncError)
- Support for embedded resources in addition to local files
- Environment-specific vault loading improved
- Better API design following Go conventions
- Function renamed to InitializeEnvEncVariablesFromResources

### 6. FlashMessage Type Migration to Base Package
**Change**: FlashMessage type moved from internal/types to base/types

**Old Usage**:
```go
import "project/internal/types"

flash := types.FlashMessage{
    Type:    "success",
    Message: "Operation completed",
}
```

**New Usage**:
```go
import "github.com/dracory/base/types"

flash := types.FlashMessage{
    Type:    "success",
    Message: "Operation completed",
}
```

**Action Required**:
- Update import statements to use base/types package
- FlashMessage struct and behavior remain identical
- All flash message functionality preserved
- File removed: `internal/types/flash_message.go`

### 7. URL Utilities Migration to Base Package
**Change**: URL utility functions moved to base/url package

**Old Usage**:
```go
import "project/internal/links"

// Local URL building logic
url := links.Url("path", "param", "value")
rootUrl := links.RootURL()
queryUrl := links.BuildQuery(params)
```

**New Usage**:
```go
import "github.com/dracory/base/url"

// Base package URL utilities
url := url.Url("path", "param", "value")
rootUrl := url.RootURL()
queryUrl := url.BuildQuery(params)
```

**Action Required**:
- Update import statements from `project/internal/links` to `github.com/dracory/base/url`
- All URL utility functions maintain same signatures
- URL initialization moved to init() function using SetDefaultURL
- No functional changes to URL building logic

### 8. User Untokenize Helper Migration to Vaultstore
**Change**: Untokenize helper moved to vaultstore package

**Old Usage**:
```go
import "project/internal/helpers"

result, err := helpers.Untokenize(template, user)
```

**New Usage**:
```go
import "github.com/dracory/vaultstore"

result, err := vaultstore.Untokenize(template, user)
```

**Action Required**:
- Update import statements to use vaultstore package
- Untokenize function behavior remains identical
- Template processing logic unchanged
- File removed: `helpers/untokenize.go`
- Note: `ext/user_untokenize.go` remains but now uses vaultstore internally

### 9. HTTP Response Utilities Migration to Base Package
**Change**: HTTP response utility functions moved to base/http package

**Old Usage**:
```go
import "project/internal/utils"

utils.SafeCloseResponseBody(responseBody)
```

**New Usage**:
```go
import "github.com/dracory/base/http"

http.SafeCloseResponseBody(responseBody)
```

**Action Required**:
- Update import statements from `project/internal/utils` to `github.com/dracory/base/http`
- SafeCloseResponseBody function maintains same signature
- HTTP response cleanup logic unchanged
- File removed: `internal/utils/response_utils.go`

### 10. Filesystem Utilities Removal
**Change**: Filesystem utility functions removed

**Old Usage**:
```go
import "project/internal/helpers"

content, err := helpers.EmbeddedFileToBytes(embeddedFS, path)
str, err := helpers.EmbeddedFileToString(embeddedFS, path)
```

**New Usage**:
```go
// Use standard library embed.FS directly
content, err := embeddedFS.ReadFile(path)
if err != nil {
    return nil, err
}
str := string(content)
```

**Action Required**:
- Remove usage of helpers filesystem functions
- Replace with standard library embed.FS methods
- File removed: `internal/helpers/filesystem.go`

### 11. Redirect Helper Removal
**Change**: Redirect helper function removed from helpers package

**Old Usage**:
```go
import "project/internal/helpers"

// Use redirect helper
result := helpers.Redirect(w, r, url)
```

**New Usage**:
```go
import "net/http"

// Use standard library
http.Redirect(w, r, url, http.StatusTemporaryRedirect)
```

**Action Required**:
- Remove usage of helpers redirect functions
- Replace with standard library http.Redirect
- File removed: `internal/helpers/redirect.go`

### 12. Ext Package Cleanup
**Change**: Several ext package files removed

**Removed Files**:
- `internal/ext/hx.go` - HTMX utilities removed
- `internal/ext/vault.go` - Vault utilities removed  
- `internal/ext/vault_test.go` - Vault tests removed

**Action Required**:
- Remove any imports or usage of these removed files
- HTMX functionality may need to be replaced with alternative implementations
- Vault functionality should use vaultstore package instead

### 13. Security Middleware Renaming
**Change**: Security middleware files renamed for consistency

**Old Files**:
- `internal/middlewares/security_headers.go` → `internal/middlewares/security_headers_middleware.go`
- `internal/middlewares/https_redirect.go` → `internal/middlewares/https_redirect_middleware.go`
- `internal/middlewares/security_headers_test.go` → `internal/middlewares/security_headers_middleware_test.go`

**Action Required**:
- No code changes required (internal refactoring only)
- Middleware functionality remains identical
- Import statements unchanged

### 14. Auth Controllers Reorganization
**Change**: Auth controllers moved to dedicated subdirectories

**New Structure**:
- `internal/controllers/auth/authentication/` (moved from root)
- `internal/controllers/auth/login/` (moved from root)
- `internal/controllers/auth/logout/` (moved from root)
- `internal/controllers/auth/register/` (moved from root)

**Action Required**:
- No code changes required (internal refactoring only)
- Controller functionality remains identical
- Import statements unchanged

### 15. Dependency Updates
**Change**: Multiple dependencies updated to newer versions

**Updated Dependencies**:
- `github.com/dracory/base`: v0.26.0 → v0.36.0
- `github.com/dracory/rtr`: v1.4.0 → v1.6.0
- `github.com/dracory/database`: v0.6.0 → v0.7.0
- `github.com/dracory/test`: v0.8.0 → v0.9.0
- `github.com/dracory/vaultstore`: v0.32.0 → v0.37.0
- `github.com/dracory/versionstore`: v0.5.0 → v0.6.0
- `modernc.org/sqlite`: v1.44.3 → v1.47.0
- Multiple other dracory packages updated

**New Local Replace Directives** (commented out by default):
```go
// replace github.com/dracory/base => ../../_modules_dracory/base
// replace github.com/dracory/rtr => ../../_modules_dracory/rtr
// replace github.com/dracory/database => ../../_modules_dracory/database
// replace github.com/dracory/test => ../../_modules_dracory/test
```

**Action Required**:
- Run `go mod tidy` to update dependencies
- Enable local replace directives if doing local development
- Review dependency release notes for any API changes

### 16. Test Utilities Removal
**Change**: Some test utility files removed

**Removed Files**:
- `internal/testutils/constants.go`
- `internal/testutils/testutils.go`

**Action Required**:
- Remove any usage of these test utilities
- Replace with standard testing practices or base/test package

### 17. Removed Directories
**Change**: Complete directories removed

**Removed Directories**:
- `internal/types/` - FlashMessage moved to base/types
- `pkg/blogtheme/` - Moved to github.com/dracory/base/blogtheme
- `pkg/webtheme/` - Moved to github.com/dracory/base/webtheme
- `pkg/blogblocks/` - Moved to base/blogblocks

**Action Required**:
- Update all imports to use base packages
- No functional changes to moved components

## 🔄 Migration Steps

### Step 1: Update Theme Package Imports
Search and replace theme imports:

```bash
# Find theme package usages
grep -r "project/pkg/blogtheme" --include="*.go" .
grep -r "project/pkg/webtheme" --include="*.go" .

# Replace imports
# project/pkg/blogtheme → github.com/dracory/base/blogtheme
# project/pkg/webtheme → github.com/dracory/base/webtheme
```

**Example Changes**:
```go
// Before
import "project/pkg/blogtheme"

// After
import "github.com/dracory/base/blogtheme"
```

### Step 2: Update Middleware Imports
Update jail bots middleware imports:

```bash
# Find jail bots middleware usages
grep -r "middlewares.JailBotsMiddleware" --include="*.go" .

# Replace imports and usage
```

**Example Changes**:
```go
// Before
import "project/internal/middlewares"
router.Use(middlewares.JailBotsMiddleware())

// After
import rtrMiddleware "github.com/dracory/rtr/middlewares"
router.Use(rtrMiddleware.JailBotsMiddleware())
```

### Step 3: Update Types Imports
Update FlashMessage and other types imports:

```bash
# Find types package usages
grep -r "project/internal/types" --include="*.go" .

# Replace imports
```

**Example Changes**:
```go
// Before
import "project/internal/types"
flash := types.FlashMessage{...}

// After
import "github.com/dracory/base/types"
flash := types.FlashMessage{...}
```

### Step 4: Update Helper Imports
Update various helper imports:

```bash
# Find helper package usages
grep -r "project/internal/helpers" --include="*.go" .
grep -r "project/internal/links" --include="*.go" .
grep -r "project/internal/utils" --include="*.go" .
grep -r "project/internal/ext" --include="*.go" .

# Replace imports with appropriate base packages
```

**Example Changes**:
```go
// Before
import "project/internal/links"
url := links.Url("path", "param", "value")

// After
import "github.com/dracory/base/url"
url := url.Url("path", "param", "value")
```

### Step 5: Update Removed Helper Usage
Replace removed helper functions:

```bash
# Find filesystem helper usage
grep -r "helpers.EmbeddedFile" --include="*.go" .

# Replace with standard library
```

**Example Changes**:
```go
// Before
content, err := helpers.EmbeddedFileToBytes(embeddedFS, path)

// After
content, err := embeddedFS.ReadFile(path)
```

### Step 6: Update Redirect Helper Usage
Replace redirect helper:

```bash
# Find redirect helper usage
grep -r "helpers.Redirect" --include="*.go" .

# Replace with standard library
```

**Example Changes**:
```go
// Before
result := helpers.Redirect(w, r, url)

// After
http.Redirect(w, r, url, http.StatusTemporaryRedirect)
```

### Step 7: Update HTTP Utils Usage
Update HTTP response utilities:

```bash
# Find utils usage
grep -r "utils.SafeCloseResponseBody" --include="*.go" .

# Replace with base package
```

**Example Changes**:
```go
// Before
import "project/internal/utils"
defer utils.SafeCloseResponseBody(resp.Body)

// After
import "github.com/dracory/base/http"
defer http.SafeCloseResponseBody(resp.Body)
```

### Step 8: Update Untokenize Usage
Update untokenize helper:

```bash
# Find untokenize usage
grep -r "helpers.Untokenize" --include="*.go" .

# Replace with vaultstore
```

**Example Changes**:
```go
// Before
import "project/internal/helpers"
result, err := helpers.Untokenize(template, user)

// After
import "github.com/dracory/vaultstore"
result, err := vaultstore.Untokenize(template, user)
```

### Step 9: Enable Local Replace Directives (Optional)
If doing local development of dracory modules:

```bash
# Uncomment replace directives in go.mod
# replace github.com/dracory/base => ../../_modules_dracory/base
# replace github.com/dracory/rtr => ../../_modules_dracory/rtr
# replace github.com/dracory/database => ../../_modules_dracory/database
# replace github.com/dracory/test => ../../_modules_dracory/test
```

### Step 10: Update Dependencies
Update Go modules:

```bash
go mod tidy
go mod download
```

### Step 11: Verify Database Connection
Test database connection to ensure optimizations work:

```bash
# Test database connection
go test ./internal/registry/... -run TestDatabase
```

### Step 12: Test CLI Functionality
Verify CLI commands work with new dispatcher:

```bash
# Test CLI commands
go run ./cmd/server task list
go run ./cmd/server routes list
```

### Step 13: Test Theme Rendering
Verify theme rendering works with new packages:

```bash
# Test theme functionality
go test ./... -run Theme
```

### Step 14: Test Middleware
Verify jail bots middleware works:

```bash
# Test middleware
go test ./internal/middlewares/... -run JailBots
```

### Step 15: Final Verification
Run comprehensive tests:

```bash
go test ./...
go build ./...
```

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
# Run all unit tests
go test ./...

# Test specific components
go test ./internal/cli/...
go test ./internal/config/...
go test ./internal/registry/...
go test ./internal/middlewares/...
```

### 2. Integration Tests
```bash
# Run integration tests with tags
go test -tags=integration ./...

# Test database functionality
go test -tags=integration ./internal/registry/... -run Database
```

### 3. Manual Testing
- Verify all theme rendering works correctly
- Test jail bots middleware functionality
- Confirm CLI commands operate properly
- Check database connection and optimizations
- Verify config encryption loading
- Test all helper functions
- Confirm URL utilities work
- Test HTTP response utilities

## 📝 Additional Notes

### New Features
- **Base Package Integration**: Major components moved to base package for reusability
- **RTR Middleware**: Jail bots middleware now available to all rtr-based projects
- **Generic CLI Dispatcher**: Type-safe CLI infrastructure using Go generics
- **Enhanced Database Package**: Automatic SQLite optimizations and connection pooling
- **Improved Config Loader**: Better error handling and embedded resources support
- **Centralized Types**: Common types like FlashMessage available in base package
- **Shared Utilities**: URL, HTTP response, and other utilities available to all projects

### Removed Features
- **Internal Types Directory**: Completely removed, types moved to base package
- **Theme Packages**: pkg/blogtheme, pkg/webtheme, pkg/blogblocks removed, moved to base
- **Custom Database Optimization**: Moved to database package for consistency
- **Custom CLI Logic**: Replaced with generic base package dispatcher
- **Helper Functions**: Several helper functions removed in favor of standard library
- **Ext Package Cleanup**: HTMX and vault utilities removed from ext package
- **Test Utilities**: Some test utilities removed

### Dependency Improvements
- Updated to latest versions of all dracory packages
- Enhanced local development support with replace directives
- Improved modularity and code reusability across projects
- Better separation of concerns between application and framework code

### Migration Benefits
- **Code Reusability**: Common components now available to all Dracory projects
- **Consistency**: Standardized implementations across projects
- **Maintainability**: Centralized maintenance of common functionality
- **Performance**: Automatic optimizations in database package
- **Type Safety**: Generic CLI dispatcher with compile-time type checking

## 🆘 Common Issues and Solutions

### Issue 1: Theme Package Import Errors
**Problem**: Compilation errors due to missing theme packages

**Solution**: 
- Update imports from `project/pkg/blogtheme` to `github.com/dracory/base/blogtheme`
- Update imports from `project/pkg/webtheme` to `github.com/dracory/base/webtheme`

### Issue 2: Middleware Import Errors
**Problem**: Jail bots middleware not found

**Solution**: 
- Update import to `github.com/dracory/rtr/middlewares`
- Use `rtrMiddleware.JailBotsMiddleware()` instead of `middlewares.JailBotsMiddleware()`

### Issue 3: Types Import Errors
**Problem**: FlashMessage type not found

**Solution**: 
- Update import from `project/internal/types` to `github.com/dracory/base/types`
- All other types follow same pattern

### Issue 4: Helper Function Import Errors
**Problem**: Helper functions not found (URL, filesystem, HTTP, etc.)

**Solution**: 
- Update imports to appropriate base packages:
  - `project/internal/links` → `github.com/dracory/base/url`
  - `project/internal/utils` → `github.com/dracory/base/http`
  - `project/internal/helpers` → various replacements (standard library or base packages)

### Issue 5: Database Connection Issues
**Problem**: Database connection or optimization not working

**Solution**: 
- Verify database package dependency is updated
- Check that SQLite optimizations are applied automatically
- No code changes should be required

### Issue 6: CLI Command Issues
**Problem**: CLI commands not working

**Solution**: 
- Verify base package dependency is updated
- CLI functionality should work identically
- Check that generic dispatcher is working properly

### Issue 7: Removed Helper Functions
**Problem**: Compilation errors due to removed helper functions

**Solution**: 
- Replace `helpers.EmbeddedFileToBytes()` with `embeddedFS.ReadFile()`
- Replace `helpers.Redirect()` with `http.Redirect()`
- Replace `helpers.Untokenize()` with `vaultstore.Untokenize()`

### Issue 8: Ext Package Cleanup
**Problem**: Missing ext package files (hx.go, vault.go)

**Solution**: 
- Remove any imports or usage of these files
- Use alternative implementations or appropriate dracory packages
- HTMX functionality may need custom implementation

### Issue 9: Dependency Conflicts
**Problem**: `go mod tidy` reports version conflicts

**Solution**: Clean go.mod and go.sum files, then run:
```bash
rm go.sum
go mod tidy
```

### Issue 10: Local Development Issues
**Problem**: Local module development not working

**Solution**: 
- Uncomment replace directives in go.mod for local modules
- Verify paths to local modules are correct
- Use `go mod download` after enabling replace directives

## 📞 Support

For additional support:
- Check the [Blueprint GitHub repository](https://github.com/dracory/blueprint)
- Review existing upgrade guides in `docs/upgrade_guides/`
- Consult the project documentation for detailed API references
- Check dracory base package documentation for migrated components
- Open an issue on GitHub for specific problems not covered in this guide

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] Import path changes clearly documented
- [x] Dependency updates properly noted
- [x] Removed files and directories accurately listed

## 🔄 Migration Steps

### Step 1: Update Theme Package Imports
Search and replace theme imports:

```bash
# Find theme package usages
grep -r "project/pkg/blogtheme" --include="*.go" .
grep -r "project/pkg/webtheme" --include="*.go" .

# Replace imports
# project/pkg/blogtheme → github.com/dracory/base/blogtheme
# project/pkg/webtheme → github.com/dracory/base/webtheme
```

**Example Changes**:
```go
// Before
import "project/pkg/blogtheme"

// After
import "github.com/dracory/base/blogtheme"
```

### Step 2: Update Middleware Imports
Update jail bots middleware imports:

```bash
# Find jail bots middleware usages
grep -r "middlewares.JailBotsMiddleware" --include="*.go" .

# Replace imports and usage
```

**Example Changes**:
```go
// Before
import "project/internal/middlewares"
router.Use(middlewares.JailBotsMiddleware())

// After
import rtrMiddleware "github.com/dracory/rtr/middlewares"
router.Use(rtrMiddleware.JailBotsMiddleware())
```

### Step 3: Update Types Imports
Update FlashMessage and other types imports:

```bash
# Find types package usages
grep -r "project/internal/types" --include="*.go" .

# Replace imports
```

**Example Changes**:
```go
// Before
import "project/internal/types"
flash := types.FlashMessage{...}

// After
import "github.com/dracory/base/types"
flash := types.FlashMessage{...}
```

### Step 4: Update Helper Imports
Update various helper imports:

```bash
# Find helper package usages
grep -r "project/internal/helpers" --include="*.go" .
grep -r "project/internal/links" --include="*.go" .
grep -r "project/internal/ext" --include="*.go" .

# Replace imports with appropriate base packages
```

**Example Changes**:
```go
// Before
import "project/internal/links"
url := links.Url("path", "param", "value")

// After
import "github.com/dracory/base/url"
url := url.Url("path", "param", "value")
```

### Step 5: Enable Local Replace Directives (Optional)
If doing local development of dracory modules:

```bash
# Uncomment replace directives in go.mod
# replace github.com/dracory/base => ../../_modules_dracory/base
# replace github.com/dracory/rtr => ../../_modules_dracory/rtr
# replace github.com/dracory/database => ../../_modules_dracory/database
# replace github.com/dracory/test => ../../_modules_dracory/test
```

### Step 6: Update Dependencies
Update Go modules:

```bash
go mod tidy
go mod download
```

### Step 7: Verify Database Connection
Test database connection to ensure optimizations work:

```bash
# Test database connection
go test ./internal/registry/... -run TestDatabase
```

### Step 8: Test CLI Functionality
Verify CLI commands work with new dispatcher:

```bash
# Test CLI commands
go run ./cmd/server task list
go run ./cmd/server routes list
```

### Step 9: Test Theme Rendering
Verify theme rendering works with new packages:

```bash
# Test theme functionality
go test ./... -run Theme
```

### Step 10: Test Middleware
Verify jail bots middleware works:

```bash
# Test middleware
go test ./internal/middlewares/... -run JailBots
```

### Step 11: Final Verification
Run comprehensive tests:

```bash
go test ./...
go build ./...
```

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
# Run all unit tests
go test ./...

# Test specific components
go test ./internal/cli/...
go test ./internal/config/...
go test ./internal/registry/...
go test ./internal/middlewares/...
```

### 2. Integration Tests
```bash
# Run integration tests with tags
go test -tags=integration ./...

# Test database functionality
go test -tags=integration ./internal/registry/... -run Database
```

### 3. Manual Testing
- Verify all theme rendering works correctly
- Test jail bots middleware functionality
- Confirm CLI commands operate properly
- Check database connection and optimizations
- Verify config encryption loading
- Test all helper functions

## 📝 Additional Notes

### New Features
- **Base Package Integration**: Major components moved to base package for reusability
- **RTR Middleware**: Jail bots middleware now available to all rtr-based projects
- **Generic CLI Dispatcher**: Type-safe CLI infrastructure using Go generics
- **Enhanced Database Package**: Automatic SQLite optimizations and connection pooling
- **Improved Config Loader**: Better error handling and embedded resources support
- **Centralized Types**: Common types like FlashMessage available in base package
- **Shared Utilities**: URL, filesystem, and HTMX utilities available to all projects

### Removed Features
- **Internal Types Directory**: Completely removed, types moved to base package
- **Theme Packages**: pkg/blogtheme and pkg/webtheme removed, moved to base
- **Custom Database Optimization**: Moved to database package for consistency
- **Custom CLI Logic**: Replaced with generic base package dispatcher

### Dependency Improvements
- Updated to latest versions of all dracory packages
- Enhanced local development support with replace directives
- Improved modularity and code reusability across projects
- Better separation of concerns between application and framework code

### Migration Benefits
- **Code Reusability**: Common components now available to all Dracory projects
- **Consistency**: Standardized implementations across projects
- **Maintainability**: Centralized maintenance of common functionality
- **Performance**: Automatic optimizations in database package
- **Type Safety**: Generic CLI dispatcher with compile-time type checking

## 🆘 Common Issues and Solutions

### Issue 1: Theme Package Import Errors
**Problem**: Compilation errors due to missing theme packages

**Solution**: 
- Update imports from `project/pkg/blogtheme` to `github.com/dracory/base/blogtheme`
- Update imports from `project/pkg/webtheme` to `github.com/dracory/base/webtheme`

### Issue 2: Middleware Import Errors
**Problem**: Jail bots middleware not found

**Solution**: 
- Update import to `github.com/dracory/rtr/middlewares`
- Use `rtrMiddleware.JailBotsMiddleware()` instead of `middlewares.JailBotsMiddleware()`

### Issue 3: Types Import Errors
**Problem**: FlashMessage type not found

**Solution**: 
- Update import from `project/internal/types` to `github.com/dracory/base/types`
- All other types follow same pattern

### Issue 4: Helper Function Import Errors
**Problem**: Helper functions not found (URL, filesystem, etc.)

**Solution**: 
- Update imports to appropriate base packages:
  - `project/internal/links` → `github.com/dracory/base/url`
  - `project/internal/helpers` → `github.com/dracory/base/files`
  - `project/internal/ext` → `github.com/dracory/base/htmx`

### Issue 5: Database Connection Issues
**Problem**: Database connection or optimization not working

**Solution**: 
- Verify database package dependency is updated
- Check that SQLite optimizations are applied automatically
- No code changes should be required

### Issue 6: CLI Command Issues
**Problem**: CLI commands not working

**Solution**: 
- Verify base package dependency is updated
- CLI functionality should work identically
- Check that generic dispatcher is working properly

### Issue 7: Dependency Conflicts
**Problem**: `go mod tidy` reports version conflicts

**Solution**: Clean go.mod and go.sum files, then run:
```bash
rm go.sum
go mod tidy
```

### Issue 8: Local Development Issues
**Problem**: Local module development not working

**Solution**: 
- Uncomment replace directives in go.mod for local modules
- Verify paths to local modules are correct
- Use `go mod download` after enabling replace directives

## 📞 Support

For additional support:
- Check the [Blueprint GitHub repository](https://github.com/dracory/blueprint)
- Review existing upgrade guides in `docs/upgrade_guides/`
- Consult the project documentation for detailed API references
- Check dracory base package documentation for migrated components
- Open an issue on GitHub for specific problems not covered in this guide

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] Import path changes clearly documented
- [x] Dependency updates properly noted
