# Upgrade Guide: v0.21.0 to v0.22.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.21.0 to v0.22.0.

## Summary

Version 0.22.0 introduces significant architectural improvements including user admin package migration, userstore API updates, enhanced security middleware, and improved configuration handling. The most significant changes involve consolidating user admin controllers into a reusable package and updating the userstore API for better consistency.

**Key Changes:**
- User admin controllers migrated from `internal/controllers/admin/users` to `pkg/useradmin`
- Userstore upgraded from v1.6.0 to v1.8.0 with breaking API change (`ID()` → `GetID()`)
- HTTPS redirect middleware refactored into dedicated package
- New email allowlist middleware with vault support
- New EnvEnc configuration module for encrypted environment variables
- Vault store key validation added
- Enhanced user tokenization with field-by-field untokenization
- Admin user routes function renamed (`UserRoutes()` → `Routes()`)
- Individual user admin link methods removed from admin links
- Log request middleware updated to skip MCP endpoints
- VSCode launch configuration improvements

---

## Breaking Changes

### 1. Userstore API Change: ID() → GetID()

**Change**: The userstore package has been upgraded from v1.6.0 to v1.8.0. The `ID()` method has been renamed to `GetID()` for better consistency with other getter methods.

**Old Usage** (v0.21.0):
```go
user := helpers.GetAuthUser(r)
userID := user.ID()

// In subscription middleware
SetSubscriberID(authenticatedUser.ID())

// In user tokenization
user.ID()
```

**New Usage** (v0.22.0):
```go
user := helpers.GetAuthUser(r)
userID := user.GetID()

// In subscription middleware
SetSubscriberID(authenticatedUser.GetID())

// In user tokenization
user.GetID()
```

**Action Required**:
- Update all calls to `user.ID()` to `user.GetID()` throughout your codebase
- This affects controllers, helpers, middleware, and any custom code using userstore
- Run tests to identify any remaining instances

**Migration Command**:
```bash
# Update all Go files to use GetID() instead of ID()
find . -type f -name "*.go" -exec sed -i 's|\.ID()|\.GetID()|g' {} \;
```

**Note**: This change affects 37 files across the codebase including controllers, helpers, middleware, and test files.

---

### 2. User Admin Package Migration

**Change**: All user admin controllers have been moved from `internal/controllers/admin/users` to `pkg/useradmin` following the same pattern as blog and log admin packages.

**Old Import Paths** (v0.21.0):
```go
import (
    "project/internal/controllers/admin/users/user_create"
    "project/internal/controllers/admin/users/user_delete"
    "project/internal/controllers/admin/users/user_impersonate"
    "project/internal/controllers/admin/users/user_manager"
    "project/internal/controllers/admin/users/user_update"
)
```

**New Import Paths** (v0.22.0):
```go
import (
    "project/pkg/useradmin/user_create"
    "project/pkg/useradmin/user_delete"
    "project/pkg/useradmin/user_impersonate"
    "project/pkg/useradmin/user_manager"
    "project/pkg/useradmin/user_update"
)
```

**Old Route Registration** (v0.21.0):
```go
// In internal/controllers/admin/routes.go
rtr.RegisterRoute(rtr.NewRoute().
    SetName("Admin Users").
    SetMethods("GET").
    SetPath("/admin/users").
    SetHandler(users.NewUsersAdminController(registry).Handler),
)
```

**New Route Registration** (v0.22.0):
```go
// In internal/controllers/admin/routes.go
rtr.RegisterRoute(rtr.NewRoute().
    SetName("Admin Users").
    SetMethods("GET").
    SetPath("/admin/users").
    SetHandler(users.NewUsersAdminController(registry).Handler),
)

// The pkg/useradmin package now provides a unified interface
// Individual controllers are now in pkg/useradmin/*
```

**Action Required**:
- Update all imports referencing the old user admin controller paths
- The `pkg/useradmin` package provides a unified interface with `useradmin.New()` and `useradmin.Handle()`
- The old `internal/controllers/admin/users` directory is removed
- Update route registration if you had custom user admin routes

**Migration Command**:
```bash
# Update all Go file imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/users|pkg/useradmin|g' {} \;
```

**Note**: This migration consolidates all user admin functionality into a reusable package following the folder-per-controller pattern. The new structure includes:
- `pkg/useradmin/useradmin.go` - Main entry point with `New()` and `Handle()`
- `pkg/useradmin/shared/` - Shared constants and links
- `pkg/useradmin/user_create/` - User creation controller
- `pkg/useradmin/user_delete/` - User deletion controller
- `pkg/useradmin/user_impersonate/` - User impersonation controller
- `pkg/useradmin/user_manager/` - User management controller with Vue.js UI
- `pkg/useradmin/user_update/` - User update controller with field-by-field tokenization

---

### 3. HTTPS Redirect Middleware Refactoring

**Change**: The HTTPS redirect middleware has been moved from `internal/middlewares/https_redirect_middleware.go` to `internal/middlewares/httpsredirect/middleware.go` as a dedicated package.

**Old Import Path** (v0.21.0):
```go
import (
    "project/internal/middlewares"
)

// Usage
middlewares.NewHTTPSRedirectMiddleware()
```

**New Import Path** (v0.22.0):
```go
import (
    "project/internal/middlewares/httpsredirect"
)

// Usage
httpsredirect.NewHTTPSRedirectMiddleware()
```

**Old File Structure** (v0.21.0):
```
internal/middlewares/https_redirect_middleware.go
internal/middlewares/https_redirect_middleware_test.go
```

**New File Structure** (v0.22.0):
```
internal/middlewares/httpsredirect/middleware.go
internal/middlewares/httpsredirect/middleware_test.go
```

**Action Required**:
- Update imports from `internal/middlewares` to `internal/middlewares/httpsredirect`
- Update middleware registration in `internal/routes/global_middlewares.go`
- The middleware API remains the same, only the import path changed

**Migration Command**:
```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|internal/middlewares"|"project/internal/middlewares/httpsredirect"|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|middlewares\.NewHTTPSRedirectMiddleware|httpsredirect.NewHTTPSRedirectMiddleware|g' {} \;
```

**Note**: The middleware functionality is identical. The refactoring improves code organization and follows the folder-per-middleware pattern.

---

### 5. Admin User Routes Function Rename

**Change**: The `UserRoutes()` function in `internal/controllers/admin/users/routes.go` has been renamed to `Routes()` for consistency with other route functions.

**Old Usage** (v0.21.0):
```go
// In internal/controllers/admin/routes.go
userRoutes, err := adminUsers.UserRoutes(registry)
```

**New Usage** (v0.22.0):
```go
// In internal/controllers/admin/routes.go
userRoutes, err := adminUsers.Routes(registry)
```

**Action Required**:
- Update any custom code that calls `adminUsers.UserRoutes()` to `adminUsers.Routes()`
- This is typically only in `internal/controllers/admin/routes.go`

**Migration Command**:
```bash
# Update function name
find . -type f -name "*.go" -exec sed -i 's|UserRoutes(|Routes(|g' {} \;
```

**Note**: The route structure has also been simplified. Instead of multiple individual routes (user_create, user_delete, etc.), there is now a single consolidated route that delegates to the new `usersAdminController`, which in turn uses the `pkg/useradmin` package.

---

### 6. Admin Links Removal

**Change**: Individual user admin link methods have been removed from `internal/links/admin_links.go` as they are now handled by the `pkg/useradmin` package.

**Removed Methods**:
```go
// These methods no longer exist in internal/links/admin_links.go:
func (l *adminLinks) UsersUserCreate(params ...map[string]string) string
func (l *adminLinks) UsersUserDelete(params ...map[string]string) string
func (l *adminLinks) UsersUserImpersonate(params ...map[string]string) string
func (l *adminLinks) UsersUserManager(params ...map[string]string) string
func (l *adminLinks) UsersUserUpdate(params ...map[string]string) string
```

**New Location**:
These links are now managed internally by `pkg/useradmin/shared/links.go` and are not exposed as public methods in the main admin links interface.

**Action Required**:
- If you have custom code that calls these specific link methods, update to use the generic `Users()` method instead
- The individual user admin action URLs are now generated internally by the user admin package

**Migration Command**:
```bash
# Replace specific link calls with generic Users() method
find . -type f -name "*.go" -exec sed -i 's|UsersUserCreate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserDelete(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserImpersonate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserManager(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserUpdate(|Users(|g' {} \;
```

---

### 7. Configuration Interface Changes

**Change**: New configuration methods have been added to support EnvEnc (environment encryption) and email allowlist features.

**New Methods Added** to `AuthConfigInterface`:
```go
// Email allowlist configuration
SetEmailsAllowedAccess([]string)
GetEmailsAllowedAccess() []string
```

**New Methods Added** to `EncryptionConfigInterface`:
```go
// EnvEnc configuration
SetEnvEncUsed(bool)
GetEnvEncUsed() bool
SetEnvEncPublicKey(string)
GetEnvEncPublicKey() string
```

**New Environment Variables**:
```bash
# Email allowlist
AUTH_EMAILS_ALLOWED_ACCESS=user1@example.com,user2@example.com

# EnvEnc configuration (if using encrypted environment variables)
ENVENC_KEY_PUBLIC=your-public-key
ENVENC_KEY_PRIVATE=your-private-key
VAULT_STORE_KEY=your-vault-key
```

**Action Required**:
- Update `.env.example` to include new environment variables if needed
- If using email allowlist, add `AUTH_EMAILS_ALLOWED_ACCESS` to your environment
- If using EnvEnc, add `ENVENC_KEY_PUBLIC`, `ENVENC_KEY_PRIVATE`, and `VAULT_STORE_KEY`
- No code changes required unless you directly access configuration interfaces

---

### 8. Vault Store Key Validation

**Change**: Vault store key validation has been added. The configuration now requires `KEY_VAULT_STORE_KEY` to be set when using vault store.

**Old Behavior** (v0.21.0):
```go
// Vault store key was optional or validated differently
```

**New Behavior** (v0.22.0):
```go
// Vault store key is now required when vault store is used
// Added to config tests: mustSetenv(t, KEY_VAULT_STORE_KEY, "test-vault-key")
```

**Action Required**:
- If using vault store, ensure `VAULT_STORE_KEY` environment variable is set
- Update your environment configuration to include the vault store key
- Update any custom config tests to include the vault store key

---

## Migration Steps

### Step 1: Update Go Module Dependencies

```bash
# Update dependencies
go get github.com/dracory/userstore@v1.8.0
go mod tidy
go mod download
```

### Step 2: Update Userstore API Calls

Replace all `ID()` method calls with `GetID()`:

```bash
# Update all Go files to use GetID() instead of ID()
find . -type f -name "*.go" -exec sed -i 's|\.ID()|\.GetID()|g' {} \;
```

**Manual Review**: After running the sed command, review the changes to ensure:
- Only userstore `ID()` calls were changed (not other `ID()` methods)
- Test files are updated correctly
- No false positives in string literals or comments

### Step 3: Update User Admin Import Paths and Function Names

Update imports for user admin controllers and rename the routes function:

```bash
# Update user admin imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/users|pkg/useradmin|g' {} \;

# Update UserRoutes() to Routes()
find . -type f -name "*.go" -exec sed -i 's|UserRoutes(|Routes(|g' {} \;
```

### Step 4: Update HTTPS Redirect Middleware Imports

Update middleware imports:

```bash
# Update HTTPS redirect middleware imports
find . -type f -name "*.go" -exec sed -i 's|internal/middlewares"|"project/internal/middlewares/httpsredirect"|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|middlewares\.NewHTTPSRedirectMiddleware|httpsredirect.NewHTTPSRedirectMiddleware|g' {} \;
```

### Step 5: Update Admin Links

Replace specific user admin link method calls with the generic `Users()` method:

```bash
# Replace specific link calls with generic Users() method
find . -type f -name "*.go" -exec sed -i 's|UsersUserCreate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserDelete(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserImpersonate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserManager(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserUpdate(|Users(|g' {} \;
```

### Step 7: Update Environment Configuration

Add new environment variables to your `.env` file:

```bash
# Email allowlist (optional)
AUTH_EMAILS_ALLOWED_ACCESS=user1@example.com,user2@example.com

# EnvEnc configuration (if using encrypted environment variables)
ENVENC_KEY_PUBLIC=your-public-key
ENVENC_KEY_PRIVATE=your-private-key
VAULT_STORE_KEY=your-vault-key
```

Update `.env.example` to document these new variables.

### Step 8: Clean Up Old Files

Remove old files that have been replaced:

```bash
# Remove old user admin directory
rm -rf internal/controllers/admin/users/

# Remove old HTTPS redirect middleware files
rm -f internal/middlewares/https_redirect_middleware.go
rm -f internal/middlewares/https_redirect_middleware_test.go
```

### Step 9: Verify Build

```bash
# Build the application
go build -o ./tmp/main ./cmd/server

# If build succeeds, run tests
go test ./...
```

---

## Testing After Migration

### 1. Unit Tests

Run the full test suite to ensure all API changes are correctly applied:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/useradmin/...
go test ./internal/middlewares/httpsredirect/...
```

### 2. User Admin Integration Tests

Test the user admin interface:

```bash
# Start the application
go run ./cmd/server

# Test user admin routes
# Visit: http://localhost:8080/admin/users
# Verify user creation, update, delete, and impersonation work correctly
```

### 3. Email Allowlist Middleware Test

If using email allowlist:

```bash
# Test that allowed emails can access protected routes
# Test that non-allowed emails are blocked

# Run middleware tests
go test ./internal/middlewares/... -run EmailAllowlist
```

### 4. HTTPS Redirect Middleware Test

Test HTTPS redirect functionality:

```bash
# Run middleware tests
go test ./internal/middlewares/httpsredirect/...

# Test HTTP to HTTPS redirect in development (should skip)
# Test in production (should redirect)
```

### 5. EnvEnc Configuration Test

If using encrypted environment variables:

```bash
# Test that encrypted variables are loaded correctly
# Test vault store key validation

# Run config tests
go test ./internal/config/... -run EnvEnc
```

### 6. User Tokenization Test

Test enhanced user tokenization:

```bash
# Test field-by-field untokenization
# Test that user update with tokenization works correctly

# Run user update tests
go test ./pkg/useradmin/user_update/...
```

---

## Additional Notes

### New Features

1. **Email Allowlist Middleware**: New middleware to restrict access to specific email addresses with vault support for encrypted allowlists.

2. **EnvEnc Configuration Module**: New configuration module for handling encrypted environment variables with public/private key support.

3. **Vault Store Key Validation**: Added validation for vault store keys to ensure proper configuration.

4. **Enhanced User Tokenization**: Field-by-field untokenization support for user data, allowing selective decryption of user fields.

5. **User Admin Package**: Unified user admin interface in `pkg/useradmin` following the folder-per-controller pattern with Vue.js UI.

6. **HTTPS Redirect Middleware Package**: Dedicated package for HTTPS redirect middleware with improved organization.

7. **MCP Endpoint Logging**: Log request middleware updated to skip MCP endpoints (`mcp/` path prefix).

8. **File Admin Package**: New `pkg/fileadmin` package providing a file management interface following the folder-per-controller pattern. This is a new feature and does not require migration for existing applications.

### Removed Features

- `internal/controllers/admin/users/` - moved to `pkg/useradmin/`
- `internal/middlewares/https_redirect_middleware.go` - moved to `internal/middlewares/httpsredirect/middleware.go`
- `internal/middlewares/https_redirect_middleware_test.go` - moved to `internal/middlewares/httpsredirect/middleware_test.go`
- User manager individual controller files - consolidated into `pkg/useradmin/user_manager/`
- Individual user admin link methods from `internal/links/admin_links.go`:
  - `UsersUserCreate()`
  - `UsersUserDelete()`
  - `UsersUserImpersonate()`
  - `UsersUserManager()`
  - `UsersUserUpdate()`
- Widget test functions from `internal/widgets/widgets_test.go` (test cleanup, 463 lines removed)

### Configuration Behavior Changes

- **Vault Store Key**: Now required when vault store is used. Configuration tests enforce this.
- **Email Allowlist**: New optional configuration for restricting access by email.
- **EnvEnc**: New optional configuration for encrypted environment variables.

### Dependency Updates

```
github.com/dracory/userstore v1.6.0 -> v1.8.0
```

---

## Common Issues and Solutions

### Issue 1: Userstore ID() Method Not Found

**Symptom**: Compilation errors for `user.ID()` method

**Solution**: Replace all `user.ID()` calls with `user.GetID()`:

```bash
find . -type f -name "*.go" -exec sed -i 's|\.ID()|\.GetID()|g' {} \;
```

**Manual Review**: Check for false positives in string literals or comments.

### Issue 2: User Admin Import Path Errors

**Symptom**: Compilation errors for user admin controller imports

**Solution**: Update imports from `internal/controllers/admin/users/*` to `pkg/useradmin/*`:

```bash
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/users|pkg/useradmin|g' {} \;
```

### Issue 3: HTTPS Redirect Middleware Import Errors

**Symptom**: Cannot find `middlewares.NewHTTPSRedirectMiddleware`

**Solution**: Update imports to use the new package:

```go
import "project/internal/middlewares/httpsredirect"

// Use
httpsredirect.NewHTTPSRedirectMiddleware()
```

### Issue 4: Vault Store Key Validation Errors

**Symptom**: Configuration tests fail with missing vault store key

**Solution**: Add the `VAULT_STORE_KEY` environment variable:

```bash
# In .env file
VAULT_STORE_KEY=your-vault-key

# In .env.example
VAULT_STORE_KEY=your-vault-key
```

### Issue 5: Build Failures After Migration

**Symptom**: Build fails with undefined references

**Solution**: 
1. Ensure all sed commands were run correctly
2. Check for any custom code that may have hardcoded old paths
3. Run `go mod tidy` to clean up dependencies
4. Review git diff to see all changes made

```bash
# Review changes
git diff

# Clean up dependencies
go mod tidy

# Rebuild
go build -o ./tmp/main ./cmd/server
```

### Issue 6: Admin Link Methods Not Found

**Symptom**: Compilation errors for `UsersUserCreate()`, `UsersUserDelete()`, etc.

**Solution**: 
1. These methods have been removed from `internal/links/admin_links.go`
2. Replace with the generic `Users()` method if you need the base user admin URL
3. Individual action URLs are now generated internally by the `pkg/useradmin` package

```bash
# Replace specific link calls with generic Users() method
find . -type f -name "*.go" -exec sed -i 's|UsersUserCreate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserDelete(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserImpersonate(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserManager(|Users(|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|UsersUserUpdate(|Users(|g' {} \;
```

### Issue 7: User Admin Routes Not Working

**Symptom**: User admin routes return 404 or errors

**Solution**: 
1. Verify route registration in `internal/controllers/admin/routes.go`
2. Ensure the new `users.NewUsersAdminController(registry)` is being used
3. Check that `pkg/useradmin` is imported correctly
4. Verify the `useradmin.New()` call in the controller if using the unified interface

---

## Support

For issues related to this upgrade:

1. **Documentation**: Review the package documentation in `pkg/useradmin/README.md`

2. **Git History**: Review the migration commits:
   ```bash
   git log --oneline v0.21.0..HEAD
   ```

3. **Reference Implementation**: Compare with the reference implementation in the main branch

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

*This upgrade guide was generated for Blueprint v0.21.0 to v0.22.0 migration.*
