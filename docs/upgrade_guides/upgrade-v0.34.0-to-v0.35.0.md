# Upgrade Guide: v0.34.0 to v0.35.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.34.0 to v0.35.0.

## Overview

This release consolidates all datastore initialization into a single config-driven file, moves store builder functions into the `config` package, introduces CSRF secret configuration with automatic generation for non-production environments, adds the Stats middleware to the global middleware chain, and removes the placeholder category manager controller.

**Key Changes:**
- 20+ individual `internal/app/stores_*.go` files deleted — replaced by single `internal/app/datastores.go` with config-driven setup
- `internal/app/app_datastores_initialize.go` deleted — replaced by `dataStoresInitialize()` in `datastores.go`
- New `internal/config/store_builders.go` — all store construction functions moved from `internal/app/` to `internal/config/`
- New `AUTH_CSRF_SECRET` environment variable — required in production/staging, auto-generated in other environments
- Modified `internal/config/auth_config.go` — CSRF secret logic added with `authSettings` struct gaining `csrfSecret` field
- New `GetCsrfSecret()` / `SetCsrfSecret()` methods on `AuthConfigInterface`
- `middlewares.NewStatsMiddleware(app)` added to global middleware chain
- `pkg/shopadmin/categories/category_manager_controller.go` deleted (was placeholder)

---

## ⚠️ Breaking Changes

---

### 1. Datastore Initialization Consolidated into Single File

**Change**: All 20+ individual store initialization files (`internal/app/stores_*.go`) and `internal/app/app_datastores_initialize.go` have been deleted. Store initialization is now consolidated into a single `internal/app/datastores.go` file with a config-driven approach that filters stores by their `Used` config flags.

**Deleted Files**:
- `internal/app/app_datastores_initialize.go`
- `internal/app/stores_audit.go`
- `internal/app/stores_blindindex_email.go`
- `internal/app/stores_blindindex_first_name.go`
- `internal/app/stores_blindindex_last_name.go`
- `internal/app/stores_blog.go`
- `internal/app/stores_cache.go`
- `internal/app/stores_chat.go`
- `internal/app/stores_cms.go`
- `internal/app/stores_custom.go`
- `internal/app/stores_entity.go`
- `internal/app/stores_feed.go`
- `internal/app/stores_geo.go`
- `internal/app/stores_log.go`
- `internal/app/stores_meta.go`
- `internal/app/stores_session.go`
- `internal/app/stores_setting.go`
- `internal/app/stores_shop.go`
- `internal/app/stores_sqlfile.go`
- `internal/app/stores_stats.go`
- `internal/app/stores_subscription.go`
- `internal/app/stores_task.go`
- `internal/app/stores_user.go`
- `internal/app/stores_vault.go`

**Old Usage**:
```go
// v0.34.0 — internal/app/app_datastores_initialize.go
func (r *appImplementation) dataStoresInitialize() error {
	initializers := []func(app AppInterface) error{
		auditStoreInitialize,
		blogStoreInitialize,
		// ... 20+ individual initializer functions
	}
	for _, initializer := range initializers {
		if err := initializer(r); err != nil {
			return err
		}
	}
	return nil
}

// v0.34.0 — internal/app/stores_audit.go
func auditStoreInitialize(app AppInterface) error {
	if !app.GetConfig().GetAuditStoreUsed() {
		return nil
	}
	store, err := newAuditStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetAuditStore(store)
	return nil
}

func newAuditStore(db *sql.DB) (auditstore.StoreInterface, error) {
	store, err := auditstore.NewStore(auditstore.NewStoreOptions{
		DB:             db,
		AuditTableName: "snv_audit_record",
	})
	// ...
	return store, nil
}
```

**New Usage**:
```go
// v0.35.0 — internal/app/datastores.go
func (app *appImplementation) dataStoresInitialize() error {
	cfg := app.GetConfig()
	if cfg == nil {
		return errors.New("config is not initialized")
	}

	if app.GetDatabase() == nil {
		return errors.New("database is not initialized")
	}

	stores := []struct {
		enabled bool
		init    func(AppInterface) error
	}{
		{cfg.GetAuditStoreUsed(), setupAuditStore},
		{cfg.GetBlogStoreUsed(), setupBlogStore},
		// ... all stores with their config flags
	}

	enabledStores := lo.Filter(stores, func(s struct {
		enabled bool
		init    func(AppInterface) error
	}, _ int) bool {
		return s.enabled
	})

	for _, s := range enabledStores {
		if err := s.init(app); err != nil {
			return err
		}
	}
	return nil
}

func setupAuditStore(app AppInterface) error {
	st, err := config.NewAuditStore(app.GetDatabase())
	if err != nil {
		return err
	}
	app.SetAuditStore(st)
	return nil
}
```

**Action Required**:
- If you have custom store initializer files in `internal/app/stores_*.go`, migrate them to the new pattern: move store construction to `internal/config/store_builders.go` and add a `setup*Store` function in `internal/app/datastores.go`
- If you added new stores, register them in the `stores` slice in `dataStoresInitialize()` with their config flag
- The `blindIndexEnabled` flag is derived from `cfg.GetUserStoreUsed() && cfg.GetVaultStoreUsed()` — blind index stores are only initialized when both user and vault stores are enabled
- Remove any references to the old `*StoreInitialize` function names (e.g., `auditStoreInitialize`) — they are now named `setup*Store` (e.g., `setupAuditStore`)

---

### 2. Store Builder Functions Moved to Config Package

**Change**: All store construction functions (previously `newAuditStore`, `newBlogStore`, etc. in `internal/app/stores_*.go`) have been moved to `internal/config/store_builders.go` as exported functions (`NewAuditStore`, `NewBlogStore`, etc.). These are now pure functions that take a `*sql.DB` (and optionally a `debug bool`) and return a configured store instance.

**Old Usage**:
```go
// v0.34.0 — store construction was private inside internal/app/
func newAuditStore(db *sql.DB) (auditstore.StoreInterface, error) {
	store, err := auditstore.NewStore(auditstore.NewStoreOptions{
		DB:             db,
		AuditTableName: "snv_audit_record",
	})
	return store, nil
}

func newBlogStore(db *sql.DB) (blogstore.StoreInterface, error) {
	store, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                  db,
		PostTableName:       "snv_blogs_post",
		TaxonomyEnabled:     true,
		// ...
	})
	return store, nil
}
```

**New Usage**:
```go
// v0.35.0 — store construction is now in internal/config/store_builders.go
st, err := config.NewAuditStore(app.GetDatabase())

st, err := config.NewBlogStore(app.GetDatabase(), app.GetConfig().GetAppDebug())
```

**Action Required**:
- If you referenced `new*Store` functions from `internal/app/`, update imports to use `config.New*Store` instead
- The new builder functions are exported and can be used directly from the `config` package
- Some builders now accept a `debug bool` parameter (blog, cache, cms, custom, log, meta, session, shop, stats, task, vault) — pass `app.GetConfig().GetAppDebug()` when calling them
- The `NewSessionStore` function accepts an additional `isDev bool` parameter for timeout configuration

**Available Store Builders**:
| Function | Parameters |
|----------|-----------|
| `config.NewAuditStore` | `db *sql.DB` |
| `config.NewBlogStore` | `db *sql.DB, debug bool` |
| `config.NewBlindIndexEmailStore` | `db *sql.DB` |
| `config.NewBlindIndexFirstNameStore` | `db *sql.DB` |
| `config.NewBlindIndexLastNameStore` | `db *sql.DB` |
| `config.NewCacheStore` | `db *sql.DB, debug bool` |
| `config.NewChatStore` | `db *sql.DB` |
| `config.NewCmsStore` | `db *sql.DB, debug bool` |
| `config.NewCustomStore` | `db *sql.DB, debug bool` |
| `config.NewEntityStore` | `db *sql.DB` |
| `config.NewFeedStore` | `db *sql.DB` |
| `config.NewGeoStore` | `db *sql.DB` |
| `config.NewLogStore` | `db *sql.DB, debug bool` |
| `config.NewMetaStore` | `db *sql.DB, debug bool` |
| `config.NewSessionStore` | `db *sql.DB, debug bool, isDev bool` |
| `config.NewSettingStore` | `db *sql.DB` |
| `config.NewShopStore` | `db *sql.DB, debug bool` |
| `config.NewSqlFileStorage` | `db *sql.DB` |
| `config.NewStatsStore` | `db *sql.DB, debug bool` |
| `config.NewSubscriptionStore` | `db *sql.DB` |
| `config.NewTaskStore` | `db *sql.DB, debug bool` |
| `config.NewUserStore` | `db *sql.DB` |
| `config.NewVaultStore` | `db *sql.DB, debug bool` |

---

### 3. New CSRF Secret Configuration

**Change**: A new `AUTH_CSRF_SECRET` environment variable has been introduced for CSRF token generation and validation. In production and staging environments, this variable is **required** — the application will panic on startup if it is not set. In development, testing, and local environments, a random secret is generated automatically and a warning is logged.

**Modified Files**:
- `internal/config/auth_config.go` — CSRF secret logic added to `authConfig()` function; new imports: `crypto/rand`, `encoding/hex`, `fmt`, `log/slog`
- `internal/config/config_implementation.go` — `csrfSecret` field added to `configImplementation` struct; `SetCsrfSecret()`/`GetCsrfSecret()` methods added; duplicate `LLM Config Implementation` comment block removed
- `internal/config/config_interfaces.go` — `SetCsrfSecret(string)` and `GetCsrfSecret() string` added to `AuthConfigInterface`
- `internal/config/constants.go` — `KEY_AUTH_CSRF_SECRET` constant added

**New Config Interface Methods**:
```go
// AuthConfigInterface now includes:
SetCsrfSecret(string)
GetCsrfSecret() string
```

**New Environment Variable**:
```
AUTH_CSRF_SECRET=your-secret-key-here
```

**Old Usage**:
```go
// v0.34.0 — No CSRF secret configuration existed
// Auth config only had registration, emails, and password auth
```

**New Usage**:
```go
// v0.35.0 — CSRF secret is now part of auth config
secret := app.GetConfig().GetCsrfSecret()
```

**Action Required**:
- **Production/Staging**: Add `AUTH_CSRF_SECRET` to your `.env` file or environment variables with a strong random value (e.g., 32-byte hex string)
- **Development/Testing**: No action required — a random secret is auto-generated. However, for persistent CSRF protection across restarts, set the variable explicitly
- If you implement `AuthConfigInterface` in a custom config, add `GetCsrfSecret()` and `SetCsrfSecret()` methods
- If you have a custom `authConfig()` function, add the CSRF secret reading logic

**Generating a CSRF Secret**:
```bash
# Generate a 32-byte hex secret
openssl rand -hex 32
```

---

### 4. Stats Middleware Added to Global Middleware Chain

**Change**: `middlewares.NewStatsMiddleware(app)` has been added to the global middleware chain in `internal/routes/global_middlewares.go`. The Stats middleware was already implemented but was not previously included in the global middleware list.

**Old Usage**:
```go
// v0.34.0 — global_middlewares.go
globalMiddlewares := append(globalMiddlewares,
	middlewares.NewSecurityHeadersMiddleware(app),
	middlewares.ThemeMiddleware(),
	middlewares.AuthMiddleware(app),
)
```

**New Usage**:
```go
// v0.35.0 — global_middlewares.go
globalMiddlewares := append(globalMiddlewares,
	middlewares.NewSecurityHeadersMiddleware(app),
	middlewares.ThemeMiddleware(),
	middlewares.AuthMiddleware(app),
	middlewares.NewStatsMiddleware(app),
)
```

**Action Required**:
- If you have a custom `globalMiddlewares()` function, add `middlewares.NewStatsMiddleware(app)` to your middleware chain
- The Stats middleware is a no-op when `GetStatsStoreUsed()` returns false or when the stats store is nil, so it is safe to include unconditionally
- If you previously added Stats middleware conditionally in routes, you can remove those conditional additions since it is now global

---

### 5. Category Manager Controller Deleted

**Change**: `pkg/shopadmin/categories/category_manager_controller.go` has been deleted. This was a placeholder implementation with a TODO comment — it rendered a simple "Category Manager - TODO: Implement full migration" message.

**Deleted File**:
- `pkg/shopadmin/categories/category_manager_controller.go`

**Removed Functions**:
- `NewCategoryManagerController(app AppInterface) *categoryManagerController`
- `NewCategoryCreateController(app AppInterface) *categoryManagerController`
- `NewCategoryUpdateController(app AppInterface) *categoryManagerController`

**Action Required**:
- If you imported `project/pkg/shopadmin/categories` and used `NewCategoryManagerController`, `NewCategoryCreateController`, or `NewCategoryUpdateController`, remove those references
- If you have routes pointing to these controllers, remove or replace them with a proper implementation

---

## 🔄 Migration Steps

### Step 1: Update Version Constant
Update the version constant in `internal/config/version.go`:
```go
const Version = "0.35.0"
```

### Step 2: Add AUTH_CSRF_SECRET to Environment
Add the new environment variable to your `.env` file:

**For production/staging**:
```bash
# Generate a strong secret
openssl rand -hex 32

# Add to .env
AUTH_CSRF_SECRET=<generated-secret>
```

**For development** (optional — auto-generated if missing):
```bash
AUTH_CSRF_SECRET=dev-only-secret-key
```

### Step 3: Remove Old Store Initializer Files
Delete the old individual store initializer files if you have local modifications:
```bash
rm internal/app/app_datastores_initialize.go
rm internal/app/stores_audit.go
rm internal/app/stores_blindindex_email.go
rm internal/app/stores_blindindex_first_name.go
rm internal/app/stores_blindindex_last_name.go
rm internal/app/stores_blog.go
rm internal/app/stores_cache.go
rm internal/app/stores_chat.go
rm internal/app/stores_cms.go
rm internal/app/stores_custom.go
rm internal/app/stores_entity.go
rm internal/app/stores_feed.go
rm internal/app/stores_geo.go
rm internal/app/stores_log.go
rm internal/app/stores_meta.go
rm internal/app/stores_session.go
rm internal/app/stores_setting.go
rm internal/app/stores_shop.go
rm internal/app/stores_sqlfile.go
rm internal/app/stores_stats.go
rm internal/app/stores_subscription.go
rm internal/app/stores_task.go
rm internal/app/stores_user.go
rm internal/app/stores_vault.go
```

### Step 4: Verify Datastore Initialization
Ensure `internal/app/datastores.go` and `internal/config/store_builders.go` are present. If you have custom stores, add them to the `stores` slice in `dataStoresInitialize()` and create a builder function in `store_builders.go`.

### Step 5: Remove Category Manager Controller References
Search for and remove any references to the deleted category manager controller:
```bash
# Find references
grep -r "NewCategoryManagerController\|NewCategoryCreateController\|NewCategoryUpdateController" .
```

### Step 6: Verify Build
```bash
go build ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
go test ./...
```

### 2. Verify CSRF Secret Behavior
- **Production**: Verify the application panics if `AUTH_CSRF_SECRET` is not set
- **Development**: Verify a warning is logged and a random secret is generated
- **Testing**: Verify the application starts without `AUTH_CSRF_SECRET` set

### 3. Verify Store Initialization
- Verify all enabled stores are still initialized correctly
- Verify disabled stores (config flag set to false) are not initialized
- Verify blind index stores only initialize when both user and vault stores are enabled

### 4. Verify Stats Middleware
- Verify visitor tracking works when stats store is enabled
- Verify no errors when stats store is disabled or nil

---

## 📝 Additional Notes

### New Features
- **CSRF Secret Configuration**: Centralized CSRF secret management with environment-aware behavior (required in production, auto-generated in development)
- **Config-Driven Store Initialization**: Store initialization is now driven by a single config-driven list, making it easier to add, remove, or reorder stores
- **Exported Store Builders**: Store construction functions are now exported from the `config` package, enabling reuse in tests and other contexts

### Removed Features
- `pkg/shopadmin/categories/category_manager_controller.go` — placeholder controller deleted
- All private `new*Store` functions in `internal/app/` — replaced by exported `config.New*Store` functions
- All individual `*StoreInitialize` functions in `internal/app/` — replaced by `setup*Store` functions in `datastores.go`

### Behavior Changes
- **Nil store checks removed from builders**: The old `new*Store` functions in `internal/app/` included `if st == nil { return nil, errors.New("xxx.NewStore returned a nil store") }` checks. The new `config.New*Store` functions do NOT include these checks — they rely on the underlying store package's `NewStore()` to return an error on failure.
- **Database nil checks consolidated**: The old `new*Store` functions each had individual `if db == nil { return nil, errors.New("database is not initialized") }` checks. The new `dataStoresInitialize()` checks `if app.GetDatabase() == nil` once at the top level instead.
- **Config nil checks consolidated**: The old `*StoreInitialize` functions each had `if app.GetConfig() == nil` checks. The new `dataStoresInitialize()` checks once at the top.
- **EnableDebug now called inside store builders**: In v0.34.0, `EnableDebug(app.GetConfig().GetAppDebug())` was called in the `*StoreInitialize` function after store creation. In v0.35.0, `EnableDebug(debug)` is called inside `config.New*Store` builder functions. This means if you call `config.New*Store` directly (e.g., in tests), debug mode is already enabled based on the `debug` parameter you pass.
- **Import order fixed in `global_middlewares.go`**: `project/internal/app` import moved before `project/internal/middlewares` to fix import ordering.
- **`AUTH_CSRF_SECRET` not added to `.env.example`**: The new environment variable was added to constants and config logic but was not added to the `.env.example` template file. You should add it manually to your `.env` files.

### Architecture Improvements
- **Single source of truth for store construction**: All table names and store options are now defined in `internal/config/store_builders.go`
- **Cleaner separation of concerns**: Store construction (config package) vs. store wiring (app package)
- **Reduced file count**: 24 files consolidated into 2 (`datastores.go` + `store_builders.go`)

---

## 🆘 Common Issues and Solutions

### Issue 1: Application Panics on Startup in Production
**Symptom**: `FATAL: AUTH_CSRF_SECRET must be set in the production environment`
**Solution**: Set the `AUTH_CSRF_SECRET` environment variable in your production/staging environment:
```bash
# Generate a strong secret
openssl rand -hex 32
# Add to your .env or environment configuration
AUTH_CSRF_SECRET=<generated-secret>
```

### Issue 2: Build Error — Undefined new*Store Functions
**Symptom**: Build error: `undefined: newAuditStore` or similar
**Solution**: Update references from `new*Store` (private, in `internal/app`) to `config.New*Store` (exported, in `internal/config`):
```go
// Old
store, err := newAuditStore(db)
// New
store, err := config.NewAuditStore(db)
```

### Issue 3: Build Error — Undefined *StoreInitialize Functions
**Symptom**: Build error: `undefined: auditStoreInitialize` or similar
**Solution**: These functions have been renamed to `setup*Store` and are now in `datastores.go`. If you referenced them externally, update to the new names or use the `dataStoresInitialize()` method instead.

### Issue 4: Build Error — Missing GetCsrfSecret/SetCsrfSecret
**Symptom**: Build error if you have a custom config implementation that doesn't implement the new methods
**Solution**: Add the following methods to your custom config implementation:
```go
func (c *yourConfig) SetCsrfSecret(v string) {
	c.csrfSecret = v
}

func (c *yourConfig) GetCsrfSecret() string {
	return c.csrfSecret
}
```

### Issue 5: Category Manager Controller Not Found
**Symptom**: Build error: `undefined: categories.NewCategoryManagerController`
**Solution**: Remove all references to the deleted controller. The package `pkg/shopadmin/categories` no longer has these controller constructors.

### Issue 6: Nil Store Returned Without Error
**Symptom**: Store is nil but no error is returned (previously would return `"xxx.NewStore returned a nil store"` error)
**Solution**: The new `config.New*Store` builder functions do not include nil store checks. If a store package's `NewStore()` returns nil without an error, it will be set on the app without detection. If you rely on nil store checks, add them manually after calling `config.New*Store`:
```go
st, err := config.NewAuditStore(db)
if err != nil {
    return err
}
if st == nil {
    return errors.New("audit store is nil")
}
```

---

## Support

- [Blueprint Repository](https://github.com/dracory/blueprint)
- [Upgrade Guides Directory](docs/upgrade_guides/)
- [Version Workflow Documentation](docs/version_workflow.md)

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Git tag verified for previous version (v0.34.0)
- [x] Previous guides reviewed for consistency
- [x] Quality checklist included in generated guide
