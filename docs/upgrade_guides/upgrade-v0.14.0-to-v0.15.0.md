# Upgrade Guide: v0.14.0 to v0.15.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.14.0 to v0.15.0.

## ⚠️ Breaking Changes

### 1. Registry Variable Naming Refactor
**Change**: Variable names changed from `app` to `registry` across the codebase for consistency

**Old Usage**:
```go
app := registry.New()
app.Start()
```

**New Usage**:
```go
registry := registry.New()
registry.Start()
```

**Action Required**:
- Update any references from `app` to `registry` in your custom code
- Search for `var app` or `app :=` patterns and rename to `registry`
- Update function receivers that use `app` to use `registry` for consistency

### 2. MCP CMS Endpoint API Key Validation
**Change**: MCP CMS endpoint now requires API key authentication

**Old Behavior**: MCP CMS endpoint accessible without authentication

**New Behavior**: 
- Returns `401 Unauthorized` when no API key is provided
- Returns `401 Unauthorized` when wrong API key is provided
- Returns `200 OK` when correct API key is provided

**Action Required**:
- Ensure MCP API key is configured in your environment
- Add `MCP_API_KEY` to your configuration if using MCP features
- Update integration tests to include valid API key

### 3. Dependency Updates
**Change**: Several dependencies were updated to newer versions

**Updated Dependencies**:
- `github.com/dracory/blogstore`: v1.4.1 → v1.5.0
- `github.com/dracory/cmsstore`: v1.6.0 → v1.8.0
- `github.com/dracory/crud/v2`: v2.0.0-20251030193142-403ea1e5e710 → v2.0.0-20260211155951-8319e6826160
- `github.com/dracory/form`: v0.19.0 → v0.20.0
- `github.com/dracory/llm`: v1.2.0 → v1.3.0
- `github.com/dracory/rtr`: v1.3.0 → v1.4.0
- `github.com/dracory/vaultstore`: v0.27.0 → v0.32.0
- `github.com/dracory/versionstore`: Added v0.5.0 (new dependency)
- `github.com/dromara/carbon/v2`: v2.6.15 → v2.6.16
- `github.com/lmittmann/tint`: v1.1.2 → v1.1.3
- `modernc.org/sqlite`: v1.44.3 → v1.46.1

**New Indirect Dependencies**:
- `github.com/aws/smithy-go`: v1.24.1
- `atomicgo.dev/cursor`: v0.2.0
- `atomicgo.dev/keyboard`: v0.2.9
- `atomicgo.dev/schedule`: v0.1.0

**Action Required**:
- Run `go mod tidy` to update dependencies
- Review dependency release notes for any API changes
- Add new `versionstore` package if needed for versioning functionality

### 4. Blog Store API Context-Aware Methods
**Change**: Blog store methods now use context-aware API

**Old Usage**:
```go
posts, err := store.List(...)
```

**New Usage**:
```go
posts, err := store.ListWithContext(ctx, ...)
```

**Action Required**:
- Update blog store method calls to use context-aware variants
- Ensure context is properly passed to store operations
- Check for new `ListWithContext`, `GetWithContext`, etc. methods

### 5. Discount Controller Refactor
**Change**: Removed unused ReadFields from discount controller

**Old Usage**:
```go
type DiscountForm struct {
    ReadFields []string // removed
    // ... other fields
}
```

**New Usage**:
```go
type DiscountForm struct {
    // ReadFields removed
    // ... other fields
}
```

**Action Required**:
- Remove any `ReadFields` references from discount-related code
- Update discount form structures to match new schema

### 6. Content Security Policy (CSP) Update
**Change**: CSP headers now include additional trusted domains

**Old CSP**:
```
default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://unpkg.com; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net; img-src 'self' data: https://cdn.jsdelivr.net
```

**New CSP**:
- Added script domains: `'unsafe-inline'`, `https://code.jquery.com`, `https://cdnjs.cloudflare.com`, `https://www.googletagmanager.com`, `https://www.statcounter.com`
- Added style domains: `https://maxcdn.bootstrapcdn.com`, `https://cdnjs.cloudflare.com`
- Added font domains: `https://fonts.googleapis.com`, `https://cdnjs.cloudflare.com`, `https://maxcdn.bootstrapcdn.com`
- Added image domains: `https://sfs.ams3.digitaloceanspaces.com`, `https://lesichkov.ams3.digitaloceanspaces.com`

**Action Required**:
- No code changes required - this is an infrastructure change
- Verify your CDN domains are included if using external resources

### 7. Post Versioning Feature
**Change**: Blog posts now support versioning through versionstore

**New Feature**: Automatic versioning of blog posts when content changes

**Action Required**:
- Ensure `versionstore` is properly configured
- Blog store `VersioningEnabled()` check added - verify versioning is enabled in config if using this feature
- New versioning methods available: `VersioningList`, `VersioningCreate`

### 8. New Routes Structure
**Change**: Added new dedicated route files for blog and CMS

**New Files**:
- `internal/controllers/website/blog/routes.go`
- `internal/controllers/website/cms/routes.go`

**Action Required**:
- No code changes required for existing functionality
- Routes are automatically registered via router setup

### 9. Database Driver Registration
**Change**: Added explicit database driver registration

**New File**: `internal/registry/dbdrivers.go`

**Content**:
```go
import (
    _ "github.com/go-sql-driver/mysql"
    // _ "github.com/lib/pq"
    // _ "modernc.org/sqlite"
)
```

**Action Required**:
- No code changes required
- MySQL driver now explicitly registered
- Additional drivers (PostgreSQL, SQLite) available via blank imports if needed

### 10. Authentication Controller Updates
**Change**: Updated authentication controller with better error handling and context support

**Changes**:
- `callAuthKnight` now accepts context for proper cancellation/timeout support
- Error logging changed from `Error` to `Warn` for invalid authentication response status
- Improved documentation comments

**Action Required**:
- No code changes required for existing functionality
- Tests updated to include context-aware calls

### 11. LLM Temperature Parameter Type Change
**Change**: LLM Temperature parameter now requires pointer type

**Old Usage**:
```go
resp, err := engine.Generate(prompt, llm.LlmOptions{
    MaxTokens:   512,
    Temperature: 0.7,  // float64
})
```

**New Usage**:
```go
import "github.com/aws/smithy-go/ptr"

resp, err := engine.Generate(prompt, llm.LlmOptions{
    MaxTokens:   512,
    Temperature: ptr.Float64(0.7),  // *float64
})
```

**Action Required**:
- Update LLM calls to use `ptr.Float64()` for Temperature parameter
- Add import: `github.com/aws/smithy-go/ptr`

### 12. SSH Command Validation
**Change**: SSH deploy commands now include security validation

**New Feature**: Added `validateCommand` function to prevent command injection attacks

**Validation**:
- Only allowed commands: `ls`, `pwd`, `cat`, `grep`, `find`, `ps`, `df`, `du`, `whoami`, `id`, `date`, `uptime`, `top`, `free`, `uname`
- Blocks dangerous characters: `;`, `&`, `|`, `` ` ``, `$`, `(`, `)`, `<`, `>`, `"`, `'`

**Action Required**:
- If using custom SSH commands in deployment, ensure they are in the allowed list
- Update any scripts that use disallowed commands

### 13. New Response Utils Package
**Change**: Added new utility function for HTTP response body handling

**New File**: `internal/utils/response_utils.go`

**Content**:
```go
func SafeCloseResponseBody(body io.Closer)
```

**Action Required**:
- No action required - this is a new utility
- Can be used for consistent HTTP response cleanup

### 14. Authentication Rate Limiting
**Change**: Added rate limiting middleware to authentication routes

**Changes**:
- Login route: 5 requests per minute (per IP)
- Registration route: 3 requests per minute (per IP)
- Logout route: No rate limiting

**Action Required**:
- No code changes required
- Ensure rate limiting is compatible with your authentication flow

### 15. Thumbnail Controller Refactor
**Change**: Thumbnail controller updated with context support and URL handling

**Changes**:
- Added `context` import for context-aware operations
- Added `neturl` import for URL parsing
- Possible breaking changes if using thumbnail controller directly

**Action Required**:
- Test thumbnail generation and serving functionality

### 16. Vault Store Configuration Update
**Change**: Vault store now includes additional configuration options

**New Options**:
- `VaultMetaTableName`: "snv_vault_meta"
- `PasswordMinLength`: 6

**Action Required**:
- No action required if using default configuration
- If using custom vault store configuration, ensure new options are included

## 🔄 Migration Steps

### Step 1: Update Variable Naming
Search and replace `app` with `registry`:

```bash
# Find all occurrences
grep -r "var app" --include="*.go" .
grep -r "app :=" --include="*.go" .

# Replace in your custom code
```

**Example Changes**:
```go
// Before
app := registry.New()
app.GetConfig().SetDatabaseDriver("sqlite")

// After
registry := registry.New()
registry.GetConfig().SetDatabaseDriver("sqlite")
```

### Step 2: Configure MCP API Key
Add MCP API key to your environment:

```bash
# Add to .env file
MCP_API_KEY=your-secure-api-key-here
```

### Step 3: Update Dependencies
Update Go modules:

```bash
go mod tidy
go mod download
```

### Step 4: Update Blog Store Calls
Update blog store method calls:

```bash
# Find all blog store usages
grep -r "store\." --include="*.go" internal/controllers/admin/blog/ internal/controllers/website/blog/
```

**Example Changes**:
```go
// Before
posts, err := r.store.List(&blogstore.ListOptions{})

// After
posts, err := r.store.ListWithContext(ctx, &blogstore.ListOptions{})
```

### Step 5: Update Discount Controller
Remove ReadFields from discount code:

```bash
# Find discount controller usages
grep -r "ReadFields" --include="*.go" internal/controllers/admin/shop/discounts/
```

### Step 6: Verify Version Store
Check if versionstore is needed for your application:

```bash
# Review dependency changes
go list -m all | grep versionstore
```

### Step 7: Review CSP Changes
Verify your external CDN resources are still allowed:

```bash
# Check if your domains are in the allowed list
# New domains added: code.jquery.com, cdnjs.cloudflare.com, 
# googletagmanager.com, statcounter.com, maxcdn.bootstrapcdn.com,
# sfs.ams3.digitaloceanspaces.com, lesichkov.ams3.digitaloceanspaces.com
```

### Step 8: Test Post Versioning
If using blog versioning:

```bash
# Verify versioning is enabled
go test ./internal/controllers/admin/blog/post_update/...
```

### Step 9: Update LLM Temperature Parameter
Update LLM calls to use pointer type:

```bash
# Find all LLM Temperature usages
grep -r "Temperature:" --include="*.go" .
```

**Example Changes**:
```go
// Before
resp, err := engine.Generate(prompt, llm.LlmOptions{
    Temperature: 0.7,
})

// After
resp, err := engine.Generate(prompt, llm.LlmOptions{
    Temperature: ptr.Float64(0.7),
})
```

### Step 10: Update SSH Commands
If using deploy commands, verify allowed commands:

```bash
# Check cmd/deploy/functions.go for allowed commands list
# Allowed: ls, pwd, cat, grep, find, ps, df, du, whoami, id, date, uptime, top, free, uname
```

### Step 11: Review New Utility Functions
New utility functions available:

```go
// internal/utils/response_utils.go
utils.SafeCloseResponseBody(body)
```

### Step 12: Test Rate Limiting
Verify authentication rate limiting works as expected:

```bash
# Test login rate limiting
go test ./internal/controllers/auth/... -run RateLimit

# Verify registration rate limiting
go test ./internal/controllers/auth/... -run Register
```

### Step 13: Test Thumbnail Controller
Verify thumbnail generation still works:

```bash
go test ./internal/controllers/shared/thumb/...
```

### Step 14: Review Vault Store Configuration
Verify vault store configuration includes new options:

```go
// Check internal/registry/stores_vault.go for new options
// New options added:
// - VaultMetaTableName: "snv_vault_meta"
// - PasswordMinLength: 6
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

# Run blog controller tests
go test ./internal/controllers/website/blog/...

# Run MCP tests
go test ./internal/controllers/website/blog/... -run MCP

# Run discount controller tests
go test ./internal/controllers/admin/shop/discounts/...
```

### 2. Integration Tests
```bash
# Run integration tests with tags
go test -tags=integration ./...

# Test MCP functionality specifically
go test -tags=integration ./internal/controllers/website/blog/... -run MCP
```

### 3. Manual Testing
- Verify MCP API key validation works correctly
- Test registry variable naming throughout the application
- Confirm blog store operations work with context-aware methods
- Check discount controller functionality
- Verify all dependencies are properly loaded

## 📝 Additional Notes

### New Features
- **MCP API Key Security**: Added API key validation for MCP CMS endpoint for enhanced security
- **Version Store Support**: New versionstore dependency for versioning functionality
- **Registry Consistency**: Improved variable naming consistency across the codebase
- **Context-Aware Stores**: Blog and other stores now support context-aware operations
- **Post Versioning**: Automatic versioning of blog posts when content changes
- **Enhanced CSP**: Expanded Content Security Policy with additional trusted CDN domains
- **New Route Structure**: Dedicated route files for blog and CMS controllers
- **Database Driver Support**: Explicit MySQL driver registration with optional PostgreSQL/SQLite
- **LLM Temperature Pointer**: LLM temperature parameter now uses pointer type for better nil handling
- **SSH Command Validation**: Security enhancement to prevent command injection in deploy commands
- **Response Utils**: New utility function for safe HTTP response body handling
- **Auth Rate Limiting**: Added rate limiting to login (5/min) and registration (3/min) endpoints
- **Thumbnail Controller**: Updated with context support
- **Vault Store Config**: New VaultMetaTableName and PasswordMinLength options

### Removed Features
- **Discount ReadFields**: Removed unused ReadFields from discount controller
- **Legacy App Variable**: Replaced `app` variable naming with `registry` for consistency

### Dependency Improvements
- Updated SQLite to latest stable version for better performance and security
- Multiple dracory package updates include performance improvements and bug fixes
- Added support for new database drivers (mysql, postgres, sqlite)

## 🆘 Common Issues and Solutions

### Issue 1: MCP API Key 401 Errors
**Problem**: MCP endpoint returns 401 Unauthorized

**Solution**: 
- Verify `MCP_API_KEY` is set in environment
- Check API key matches configuration
- Ensure header is properly formatted in requests

### Issue 2: Registry Variable Not Found
**Problem**: Compilation errors due to missing registry variable

**Solution**: 
- Search for remaining `app` references
- Update all `var app` to `var registry`
- Update constructor calls to assign to `registry` variable

### Issue 3: Blog Store Context Errors
**Problem**: Compilation errors with blog store method calls

**Solution**: 
- Update to context-aware methods (e.g., `ListWithContext`)
- Ensure context is imported and passed correctly
- Check store interface for available methods

### Issue 4: Dependency Conflicts
**Problem**: `go mod tidy` reports version conflicts

**Solution**: Clean go.mod and go.sum files, then run:
```bash
rm go.sum
go mod tidy
```

### Issue 5: Discount Controller Compilation Errors
**Problem**: ReadFields not found in discount form

**Solution**: Remove ReadFields references from discount-related structs and code

### Issue 6: CSP Blocking Resources
**Problem**: External resources blocked by Content Security Policy

**Solution**: 
- Verify your CDN domains are in the allowed list
- Add required domains to CSP configuration in `internal/middlewares/security_headers.go`
- Common allowed domains now include: jquery.com, cdnjs.cloudflare.com, googletagmanager.com

### Issue 7: Post Versioning Not Working
**Problem**: Post versioning not creating versions

**Solution**:
- Verify `versionstore` dependency is installed: `go list -m all | grep versionstore`
- Check if blog store versioning is enabled via `store.VersioningEnabled()`
- Ensure context is properly passed to versioning operations

### Issue 8: LLM Temperature Compilation Errors
**Problem**: Compilation errors related to Temperature parameter in LLM options

**Solution**:
- Import `github.com/aws/smithy-go/ptr`
- Change `Temperature: 0.7` to `Temperature: ptr.Float64(0.7)`

### Issue 9: SSH Command Blocked
**Problem**: SSH deploy command fails with "command not allowed" error

**Solution**:
- Check if command is in allowed list: `ls`, `pwd`, `cat`, `grep`, `find`, `ps`, `df`, `du`, `whoami`, `id`, `date`, `uptime`, `top`, `free`, `uname`
- If command is not allowed, it must be added to the allowed list in `cmd/deploy/functions.go`
- Never use dangerous characters in commands: `;`, `&`, `|`, `` ` ``, `$`, `(`, `)`, `<`, `>`, `"`, `'`

### Issue 10: Rate Limiting Blocking Legitimate Requests
**Problem**: Authentication requests blocked by rate limiting

**Solution**:
- Login: 5 requests per minute per IP
- Registration: 3 requests per minute per IP
- If you need higher limits, modify rate limiting in `internal/controllers/auth/routes.go`

### Issue 11: Vault Store Configuration
**Problem**: Vault store may need new configuration options

**Solution**:
- Ensure new options are included in vault store initialization:
```go
vaultstore.NewStoreOptions{
    DB:                 db,
    VaultTableName:     "snv_vault_vault",
    VaultMetaTableName: "snv_vault_meta",  // NEW
    PasswordMinLength:  6,                  // NEW
}
```

## 📞 Support

For additional support:
- Check the [Blueprint GitHub repository](https://github.com/dracory/blueprint)
- Review existing upgrade guides in `docs/upgrade_guides/`
- Consult the project documentation for detailed API references
- Open an issue on GitHub for specific problems not covered in this guide

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
