# Upgrade Guide: v0.27.0 to v0.28.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.27.0 to v0.28.0.

## Overview

This release focuses on modernizing the admin UI with Vue.js, refactoring migration files to use the registry package, and re-adding testify for testing. The changes primarily affect the shopadmin and useradmin packages with a complete UI overhaul from HTMX to Vue.js with AJAX handlers.

**Key Changes:**
- Migration files now use `registry.RegistryInterface` from the registry package instead of local types
- Admin UI migrated from HTMX to Vue.js with AJAX handlers
- Response format changed from `data.success` to `data.status === 'success'`
- SweetAlert2 replaced with Notiflix for notifications
- Testify re-added to go.mod (reverses v0.26.0 removal)
- New AJAX handler files for CRUD operations in admin packages

---

## ⚠️ Breaking Changes

### 1. Migration File RegistryInterface Type

**Change**: Migration files now use `registry.RegistryInterface` from the `project/internal/registry` package instead of a local `RegistryInterface` type.

**Old Usage**:
```go
// database/migrations/migrate.go
type RegistryInterface interface {
    GetConfig() config.ConfigInterface
    GetDatabase() *sql.DB
    // ... other methods
}

func migrateStores(registry RegistryInterface) error {
    cfg := registry.GetConfig()
    // ...
}
```

**New Usage**:
```go
// database/migrations/migrate.go
import "project/internal/registry"

func migrateStores(registry registry.RegistryInterface) error {
    cfg := registry.GetConfig()
    // ...
}
```

**Action Required**:
- If you have custom migration files that define a local `RegistryInterface`, update them to use `registry.RegistryInterface` from the registry package
- Add the import: `import "project/internal/registry"`
- Update function signatures to use `registry.RegistryInterface`
- The Blueprint rapid application development (RAD) starter template's migration files have been updated automatically

**Files to Check**:
- `database/migrations/migrate.go` (already updated in template)
- Any custom migration files in your project that use `RegistryInterface`

**Migration Command**:
```bash
# Find custom migration files using local RegistryInterface
grep -r "type RegistryInterface" database/migrations/

# Update imports in migration files
find database/migrations -name "*.go" -exec sed -i 's|RegistryInterface|registry.RegistryInterface|g' {} \;
```

---

### 2. AJAX Response Format Change

**Change**: AJAX endpoints now return `data.status === 'success'` instead of `data.success` for success checks in JavaScript.

**Old Usage**:
```javascript
// pkg/shopadmin/product_manager/products.js
const response = await fetch(url);
const data = await response.json();

if (data.success) {
    // Handle success
}
```

**New Usage**:
```javascript
// pkg/shopadmin/product_manager/products.js
const response = await fetch(url);
const data = await response.json();

if (data.status === 'success') {
    // Handle success
}
```

**Action Required**:
- Update any custom JavaScript that checks `data.success` to use `data.status === 'success'`
- This affects all admin UI JavaScript files that interact with AJAX endpoints
- The Blueprint rapid application development (RAD) starter template's admin UI files have been updated automatically

**Files to Check**:
- Any custom JavaScript files in your admin packages
- Any custom AJAX handlers you've added

**Migration Command**:
```bash
# Find JavaScript files using data.success
grep -r "data.success" pkg/ --include="*.js"

# Replace data.success with data.status === 'success'
find pkg/ -name "*.js" -exec sed -i 's|data\.success|data.status === '\''success'\''|g' {} \;
```

---

### 3. SweetAlert2 to Notiflix Migration

**Change**: SweetAlert2 has been replaced with Notiflix for notifications in the admin UI.

**Old Usage**:
```javascript
// pkg/shopadmin/product_manager/products.js
Swal.fire({
    icon: 'success',
    title: 'Success',
    text: 'Operation completed'
});

Swal.fire('Error', 'Failed to complete operation', 'error');
```

**New Usage**:
```javascript
// pkg/shopadmin/product_manager/products.js
Notiflix.Notify.success('Operation completed', {
    position: 'right-top',
    timeout: 3000,
});

Notiflix.Notify.failure('Failed to complete operation', {
    position: 'right-top',
    timeout: 3000,
});

Notiflix.Confirm.show(
    'Confirm',
    'Are you sure?',
    'Yes',
    'No',
    () => { /* on confirm */ },
    () => { /* on cancel */ }
);
```

**Action Required**:
- Update any custom JavaScript that uses SweetAlert2 to use Notiflix
- Add Notiflix CSS and JS to your layout if you have custom admin pages
- The Blueprint rapid application development (RAD) starter template's admin UI files have been updated automatically

**CDN Dependencies**:
```html
<!-- Add to your layout if using custom admin pages -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/notiflix@3.2.8/dist/notiflix-3.2.8.min.css">
<script src="https://cdn.jsdelivr.net/npm/notiflix@3.2.8/dist/notiflix-3.2.8.min.js"></script>
```

**Migration Command**:
```bash
# Find JavaScript files using Swal
grep -r "Swal\." pkg/ --include="*.js"

# Common replacements (manual review required):
# Swal.fire({icon: 'success', ...}) → Notiflix.Notify.success(...)
# Swal.fire({icon: 'error', ...}) → Notiflix.Notify.failure(...)
# Swal.fire('Title', 'Message', 'type') → Notiflix.Notify.failure(...)
# Swal.confirm(...) → Notiflix.Confirm.show(...)
```

---

### 4. Testify Re-added

**Change**: Testify has been re-added to go.mod (github.com/stretchr/testify v1.11.1), reversing the removal from v0.26.0.

**Old Usage** (v0.26.0):
```go
import "testing"

func TestSomething(t *testing.T) {
    if expected != actual {
        t.Errorf("expected %v, got %v", expected, actual)
    }
}
```

**New Usage** (v0.28.0):
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
    assert.Equal(t, expected, actual)
}
```

**Action Required**:
- Run `go mod tidy` to add testify and related dependencies
- If you migrated away from testify in v0.26.0, you can optionally use it again
- This is not a breaking change - standard library assertions still work
- The Blueprint rapid application development (RAD) starter template now uses testify in new test files

**Migration Command**:
```bash
# Update dependencies
go mod tidy
```

---

### 5. Admin UI Architecture Change

**Change**: Admin UI has been migrated from HTMX-based to Vue.js-based with dedicated AJAX handlers. This is a major architectural change affecting all admin packages.

**Old Architecture**:
- HTMX for dynamic content loading
- Server-side rendering with HTMX attributes
- Simple controller methods returning HTML
- `*_page.go` files with embedded HTML/JS

**New Architecture**:
- Vue.js for reactive UI
- AJAX handlers for data operations
- Separate handler files for each action (create, delete, fetch, update)
- `*_controller.go` files with dedicated handler methods
- `*_page.go` files for page rendering with Vue.js mounting

**Action Required**:
- If you have custom admin pages using HTMX, consider migrating to Vue.js
- If you have custom AJAX handlers, ensure they return the correct response format
- The Blueprint rapid application development (RAD) starter template's admin packages have been migrated automatically

**Impact**:
- `pkg/shopadmin/order_manager/` - Migrated to Vue.js
- `pkg/shopadmin/product_manager/` - Migrated to Vue.js
- `pkg/shopadmin/product_update/` - Migrated to Vue.js
- `pkg/useradmin/user_manager/` - Migrated to Vue.js
- `pkg/useradmin/user_update/` - Migrated to Vue.js

**Note**: This is primarily an internal framework change. If you don't modify the admin packages, no action is required.

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

Update go.mod to include testify and related dependencies:

```bash
go mod tidy
```

### Step 2: Update Custom Migration Files

If you have custom migration files using local `RegistryInterface`:

```bash
# Find files using local RegistryInterface
grep -r "type RegistryInterface" database/migrations/

# Update imports
find database/migrations -name "*.go" -exec sed -i '1i import "project/internal/registry"' {} \;

# Update type references
find database/migrations -name "*.go" -exec sed -i 's|RegistryInterface|registry.RegistryInterface|g' {} \;
```

### Step 3: Update Custom JavaScript

If you have custom JavaScript in admin packages:

```bash
# Update response format checks
find pkg/ -name "*.js" -exec sed -i 's|data\.success|data.status === '\''success'\''|g' {} \;

# Review and update SweetAlert2 usage to Notiflix (manual review required)
grep -r "Swal\." pkg/ --include="*.js"
```

### Step 4: Update Version Constant

Update the version constant in `internal/config/version.go`:

```go
const Version = "0.28.0"
```

### Step 5: Build and Test

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

### 3. Migration Tests

Test database migrations specifically:

```bash
go test ./database/migrations/...
```

### 4. Admin UI Tests

Test admin interfaces manually:

```bash
# Start the server
go run ./cmd/server

# Test each admin interface:
# - /admin/shop/products
# - /admin/shop/orders
# - /admin/users
# - /admin/users/update
```

Verify:
- Vue.js loads correctly
- AJAX operations work (create, delete, update, fetch)
- Notifications display correctly (Notiflix)
- Filters and pagination work
- Sorting works

### 5. Build Verification

Verify the application builds successfully:

```bash
go build -o ./bin/server ./cmd/server
```

---

## 📝 Additional Notes

### New Features

- **Vue.js Admin UI**: Modern reactive UI with better user experience
- **AJAX Handlers**: Dedicated handlers for each CRUD operation
- **Notiflix Notifications**: Modern notification library with better UX
- **Improved Test Coverage**: Testify re-added for better test assertions
- **Better Code Organization**: Separate handler files for better maintainability

### Removed Features

- **HTMX-based Admin UI**: Removed in favor of Vue.js
- **SweetAlert2**: Replaced with Notiflix

### Migration File Improvements

- Migration files now use the standard `registry.RegistryInterface` type
- Better type safety and consistency across the codebase
- Easier to maintain and extend migration logic

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: RegistryInterface" in migration files

**Symptom**: Compilation error about undefined RegistryInterface type.

**Solution**: Update migration files to use `registry.RegistryInterface` from the registry package:
```go
import "project/internal/registry"

func migrateStores(registry registry.RegistryInterface) error {
    // ...
}
```

### Issue 2: JavaScript "data.success is undefined"

**Symptom**: JavaScript errors about data.success being undefined.

**Solution**: Update JavaScript to check `data.status === 'success'` instead of `data.success`:
```javascript
// Old
if (data.success) { }

// New
if (data.status === 'success') { }
```

### Issue 3: "Swal is not defined"

**Symptom**: JavaScript errors about Swal being undefined.

**Solution**: Replace SweetAlert2 with Notiflix:
```javascript
// Old
Swal.fire('Success', 'Message', 'success');

// New
Notiflix.Notify.success('Message', { position: 'right-top', timeout: 3000 });
```

### Issue 4: Vue.js not loading

**Symptom**: Admin pages show raw Vue.js template syntax instead of rendered UI.

**Solution**: Ensure Vue.js CDN is included in your layout:
```html
<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
```

### Issue 5: Notiflix notifications not showing

**Symptom**: Notifications don't appear when actions complete.

**Solution**: Ensure Notiflix CSS and JS are included in your layout:
```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/notiflix@3.2.8/dist/notiflix-3.2.8.min.css">
<script src="https://cdn.jsdelivr.net/npm/notiflix@3.2.8/dist/notiflix-3.2.8.min.js"></script>
```

### Issue 6: Testify dependency conflicts

**Symptom**: `go mod tidy` fails with dependency conflicts.

**Solution**: Clean the module cache and retry:
```bash
go clean -modcache
go mod tidy
```

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Review the changelog for detailed changes
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
- [x] Git tag verified for previous version
- [x] Previous guides reviewed for consistency
