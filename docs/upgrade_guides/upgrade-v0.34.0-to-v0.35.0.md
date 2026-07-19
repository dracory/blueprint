# Upgrade Guide: v0.34.0 to v0.35.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.34.0 to v0.35.0.

## Overview

This release consolidates all datastore initialization into a single config-driven file, moves store builder functions into the `config` package, introduces CSRF secret configuration with automatic generation for non-production environments, replaces all individual admin shop controllers with a single delegation controller, removes the `testify` test dependency in favor of standard library assertions, adds Stripe payment support, refactors the file manager code organization, and adds the Stats middleware to the global middleware chain.

**Key Changes:**
- 20+ individual `internal/app/stores_*.go` files deleted — replaced by single `internal/app/datastores.go` with config-driven setup
- `internal/app/app_datastores_initialize.go` deleted — replaced by `dataStoresInitialize()` in `datastores.go`
- New `internal/config/store_builders.go` — all store construction functions moved from `internal/app/` to `internal/config/`
- New `AUTH_CSRF_SECRET` environment variable — required in production/staging, auto-generated in other environments
- Modified `internal/config/auth_config.go` — CSRF secret logic added with `authSettings` struct gaining `csrfSecret` field
- New `GetCsrfSecret()` / `SetCsrfSecret()` methods on `AuthConfigInterface`
- `middlewares.NewStatsMiddleware(app)` added to global middleware chain
- All individual admin shop controllers deleted (categories, discounts, products, productupdate, shared) — replaced by single `shopAdminController` delegating to `pkg/shopadmin`
- `github.com/stretchr/testify` dependency removed — all tests refactored to use standard library assertions
- `github.com/dracory/crud/v2` dependency removed
- `github.com/stripe/stripe-go/v73` dependency added for Stripe payment support
- New website shop cart controller, helpers (`cart_util.go`, `payment_code.go`, `stripe.go`), and shop cart link constants
- `pkg/blogai/database.go` deleted (was dead/commented code)
- File manager `handleLoadFilesAjax` extracted to separate file `load_files_ajax.go`
- Import order standardized: `project/internal/app` placed before other internal packages
- `pkg/shopadmin/categories/category_manager_controller.go` deleted (was placeholder)

---

## ⚠️ Breaking Changes

---

### 1. Datastore Initialization Consolidated into Single File

**Change**: All 20+ individual store initialization files (`internal/app/stores_*.go`), their corresponding test files (`internal/app/stores_*_test.go`), and `internal/app/app_datastores_initialize.go` have been deleted. Store initialization is now consolidated into a single `internal/app/datastores.go` file with a config-driven approach that filters stores by their `Used` config flags. All store initialization tests are consolidated into a single `internal/app/datastores_test.go` file.

**Deleted Files**:
- `internal/app/app_datastores_initialize.go`
- `internal/app/stores_audit.go` + `internal/app/stores_audit_test.go`
- `internal/app/stores_blindindex_email.go`
- `internal/app/stores_blindindex_first_name.go`
- `internal/app/stores_blindindex_last_name.go`
- `internal/app/stores_blog.go` + `internal/app/stores_blog_test.go`
- `internal/app/stores_cache.go` + `internal/app/stores_cache_test.go`
- `internal/app/stores_chat.go` + `internal/app/stores_chat_test.go`
- `internal/app/stores_cms.go` + `internal/app/stores_cms_test.go`
- `internal/app/stores_custom.go` + `internal/app/stores_custom_test.go`
- `internal/app/stores_entity.go` + `internal/app/stores_entity_test.go`
- `internal/app/stores_feed.go` + `internal/app/stores_feed_test.go`
- `internal/app/stores_geo.go` + `internal/app/stores_geo_test.go`
- `internal/app/stores_log.go` + `internal/app/stores_log_test.go`
- `internal/app/stores_meta.go` + `internal/app/stores_meta_test.go`
- `internal/app/stores_session.go` + `internal/app/stores_session_test.go`
- `internal/app/stores_setting.go` + `internal/app/stores_setting_test.go`
- `internal/app/stores_shop.go` + `internal/app/stores_shop_test.go`
- `internal/app/stores_sqlfile.go`
- `internal/app/stores_stats.go` + `internal/app/stores_stats_test.go`
- `internal/app/stores_subscription.go` + `internal/app/stores_subscription_test.go`
- `internal/app/stores_task.go` + `internal/app/stores_task_test.go`
- `internal/app/stores_user.go` + `internal/app/stores_user_test.go`
- `internal/app/stores_vault.go` + `internal/app/stores_vault_test.go`

**New Files**:
- `internal/app/datastores.go` — consolidated `dataStoresInitialize()` with config-driven store list + `setup*Store` functions
- `internal/app/datastores_test.go` — 19 test functions covering each store's `_Success` and `_NotUsed` initialization cases using `testutils.Setup()`

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

### 6. Admin Shop Controllers Replaced by Single Delegation Controller

**Change**: All individual admin shop controllers in `internal/controllers/admin/shop/` have been deleted and replaced by a single `shopAdminController` that delegates all shop admin requests to the `pkg/shopadmin` package. The old route handler that dispatched to individual controllers based on a `controller` query parameter has been replaced with a simple two-route setup (base + catchall) that routes everything through `shopAdminController`.

**Deleted Directories and Files**:
- `internal/controllers/admin/shop/categories/` (entire directory)
  - `category_create_controller.go` and test
  - `categorymanager/` (entire subdirectory)
  - `categoryupdate/` (entire subdirectory)
- `internal/controllers/admin/shop/discounts/` (entire directory)
  - `discount_controller.go`, `discount_form_validation_rule.go` and tests
- `internal/controllers/admin/shop/home_controller.go`
- `internal/controllers/admin/shop/products/` (entire directory)
  - `product_create_controller.go`, `product_delete_controller.go`, `product_manager_controller.go` and tests
  - `productupdate/` (entire subdirectory with all components: details, media, metadata, tags)
- `internal/controllers/admin/shop/shared/` (entire directory)
  - `constants.go`, `header.go`, `links.go` and tests

**New File**:
- `internal/controllers/admin/shop/shop_controller.go` — `shopAdminController` struct with `NewShopAdminController(app)` constructor and `Handler(w, r)` method

**Old Usage**:
```go
// v0.34.0 — internal/controllers/admin/shop/routes.go
handler := func(w http.ResponseWriter, r *http.Request) string {
	controller := req.GetStringTrimmed(r, "controller")

	if controller == shared.CONTROLLER_CATEGORIES {
		return categorymanager.NewCategoryManagerController(app).Handler(w, r)
	}
	if controller == shared.CONTROLLER_CATEGORY_CREATE {
		return categories.NewCategoryCreateController(app).Handler(w, r)
	}
	if controller == shared.CONTROLLER_PRODUCT_UPDATE {
		return productupdate.NewProductUpdateController(app).Handler(w, r)
	}
	// ... many more controller dispatches
	return NewHomeController(app).Handler(w, r)
}

shopOrders := rtr.NewRoute().
	SetName("Admin > Shop > Orders").
	SetPath(links.ADMIN_SHOP).
	SetHTMLHandler(handler)
```

**New Usage**:
```go
// v0.35.0 — internal/controllers/admin/shop/routes.go
shop := rtr.NewRoute().
	SetName("Admin > Shop").
	SetPath(links.ADMIN_SHOP).
	SetHandler(NewShopAdminController(app).Handler)

shopCatchAll := rtr.NewRoute().
	SetName("Admin > Shop > Catchall").
	SetPath(links.ADMIN_SHOP + links.CATCHALL).
	SetHandler(NewShopAdminController(app).Handler)
```

**Action Required**:
- If you imported any of the deleted controller packages (e.g., `project/internal/controllers/admin/shop/categories`, `project/internal/controllers/admin/shop/products`, `project/internal/controllers/admin/shop/discounts`, `project/internal/controllers/admin/shop/shared`), remove those imports and use `NewShopAdminController(app)` instead
- If you had custom routes pointing to individual shop controllers, route them through `NewShopAdminController(app).Handler` instead
- The `pkg/shopadmin` package now handles all shop admin routing internally — the `shopAdminController` is a thin wrapper that initializes `shopadmin.New()` with `AdminOptions` and delegates to `admin.Handle(w, r)`
- The route handler type changed from `SetHTMLHandler` (returning string) to `SetHandler` (standard `http.HandlerFunc`) — update any custom route registrations accordingly

---

### 7. testify Dependency Removed

**Change**: `github.com/stretchr/testify` has been removed from `go.mod`. All tests that previously used `testify/assert` or `testify/require` have been refactored to use standard library assertions and error checking patterns.

**Affected Files**:
- `go.mod` — `github.com/stretchr/testify` removed from require block; `github.com/davecgh/go-spew`, `github.com/pmezard/go-difflib`, `gopkg.in/yaml.v3` removed from indirect dependencies
- `pkg/useradmin/user_update/handle_timezones_fetch_ajax_test.go`
- `pkg/useradmin/user_update/handle_user_fetch_ajax_test.go`
- `pkg/useradmin/user_update/handle_user_update_ajax_test.go`
- `pkg/useradmin/user_update/user_update_page_test.go`

**Old Usage**:
```go
// v0.34.0 — using testify
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
	result := doSomething()
	assert.True(t, result)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}
```

**New Usage**:
```go
// v0.35.0 — using standard library
func TestSomething(t *testing.T) {
	result := doSomething()
	if !result {
		t.Errorf("expected true, got false")
	}
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

**Action Required**:
- If your custom tests use `testify`, either keep the dependency in your `go.mod` or refactor to standard library assertions
- If you relied on testify being a transitive dependency from Blueprint, add it explicitly to your `go.mod`:
```bash
go get github.com/stretchr/testify@v1.11.1
```
- Run `go mod tidy` after updating

---

### 8. New Stripe Dependency Added

**Change**: `github.com/stripe/stripe-go/v73` has been added as a direct dependency for Stripe payment integration support.

**New Files**:
- `internal/helpers/stripe.go` — `GenerateStripePaymentsCheckoutURL()` function and `GenerateStripePaymentsCheckoutURLOptions` / `LineItem` types
- `internal/helpers/payment_code.go` — `PaymentCodeData` struct, `GetPaymentCodeData()`, `SetPaymentCodeData()`, `DeletePaymentCodeData()` functions for cache-based payment code management
- `internal/helpers/cart_util.go` — `GenerateCartCacheKey()` for guest cart cache key generation

**New go.mod Entry**:
```
github.com/stripe/stripe-go/v73 v73.16.0
```

**Action Required**:
- Run `go mod tidy` to pick up the new Stripe dependency
- If you don't use Stripe payments, the dependency is still required for compilation but the helpers are not called unless you explicitly use them
- If you implement Stripe payment flows, use `helpers.GenerateStripePaymentsCheckoutURL()` to create checkout sessions

---

### 9. New Website Shop Cart Routes and Link Constants

**Change**: New shop cart routes, link constants, and a cart controller have been added to the website controllers. The cart routes are **commented out by default** in the website routes file.

**New Files**:
- `internal/controllers/website/shop/cart/cart_controller.go` — `CartController` with `NewCartController(app)` and `Handler(w, r)`
- `internal/controllers/website/shop/cart/routes.go` — `Routes(app)` returning cart API route

**New Link Constants** (`internal/links/constants.go`):
```go
const SHOP_CART = SHOP + "/cart"
const SHOP_CART_ADD = SHOP_CART + "/add"
const SHOP_CART_REMOVE = SHOP_CART + "/remove"
const SHOP_CART_UPDATE = SHOP_CART + "/update"
const SHOP_CART_API = SHOP_CART + "/api"
const SHOP_CHECKOUT = SHOP + "/checkout"
const SHOP_PRODUCT = SHOP + "/product"
```

**New Website Link Methods** (`internal/links/website_links.go`):
- `ShopCart()`, `ShopCartAdd()`, `ShopCartRemove()`, `ShopCartUpdate()`, `ShopCartAPI()`
- `ShopCheckout()`, `ShopProduct()`

**Old Usage**:
```go
// v0.34.0 — No shop cart routes or link constants existed
```

**New Usage**:
```go
// v0.35.0 — internal/controllers/website/routes.go (commented out by default)
// websiteRoutes = append(websiteRoutes, cart.Routes(app)...)

// Using link helpers
url := links.Website().ShopCart()
url := links.Website().ShopCheckout()
url := links.Website().ShopProduct(productID, productSlug, nil)
```

**Action Required**:
- No action required if you don't use shop cart functionality — routes are commented out by default
- To enable cart routes, uncomment `cart.Routes(app)` in `internal/controllers/website/routes.go`
- If you have custom shop link constants, verify they don't conflict with the new `SHOP_CART*`, `SHOP_CHECKOUT`, and `SHOP_PRODUCT` constants

---

### 10. blogai/database.go Deleted

**Change**: `pkg/blogai/database.go` has been deleted. This file contained only commented-out code (a dead `Initialize()` function that was never active).

**Action Required**:
- If you imported `pkg/blogai` and referenced `blogai.Initialize()` or `blogai.customStore`, remove those references — they were never functional
- No functional impact — the file was entirely commented out

---

### 11. File Manager Code Organization Refactored

**Change**: The `handleLoadFilesAjax` method has been extracted from `pkg/fileadmin/file_manager/file_manager_controller.go` into a separate file `pkg/fileadmin/file_manager/load_files_ajax.go`. A new test file `load_files_ajax_test.go` was also added.

**Old Usage**:
```go
// v0.34.0 — handleLoadFilesAjax was a method inside file_manager_controller.go
func (controller *FileManagerController) handleLoadFilesAjax(r *http.Request) string {
	// ... implementation inline in controller file
}
```

**New Usage**:
```go
// v0.35.0 — handleLoadFilesAjax is now in load_files_ajax.go (same package, same method)
// No API change — the method signature and behavior are identical
// The controller file is now ~86 lines smaller
```

**Action Required**:
- No code changes required — the method is in the same package and has the same signature
- If you have local modifications to `file_manager_controller.go`, be aware that `handleLoadFilesAjax` is now in a separate file

---

### 12. Import Order Standardized

**Change**: Import order has been standardized across the codebase to place `project/internal/app` before other `project/internal/*` packages. This is a style-only change with no functional impact.

**Affected Files** (28 files across controllers, middlewares, widgets, emails, testutils, and pkg packages):
- `internal/emails/admin_email_contact_form_submitted.go`
- `internal/emails/admin_email_new_user_registered.go`
- `internal/emails/user_email_invite_friend.go`
- `internal/middlewares/api_auth_middleware.go`
- `internal/middlewares/cms_layout_middleware.go`
- `internal/middlewares/email_allowlist_middleware.go`
- `internal/middlewares/log_request_mIddleware.go`
- `internal/middlewares/subscription_middleware.go`
- `internal/testutils/login_as.go`
- `internal/testutils/setup.go`
- `internal/widgets/authenticated_widget.go` (and 5 other widget files)
- `pkg/blogadmin/` (8 files)
- `pkg/fileadmin/routes.go`
- `pkg/logadmin/routes.go`
- `pkg/shopadmin/home/home_controller.go`
- `pkg/useradmin/` (6 files)

**Old Usage**:
```go
// v0.34.0 — internal/app was mixed with other internal packages
import (
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/app"
)
```

**New Usage**:
```go
// v0.35.0 — internal/app comes before other internal packages
import (
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"
)
```

**Action Required**:
- Run `gofmt -w .` to automatically standardize import order in your code
- No functional impact — this is purely a style change

---

### 13. crud/v2 Dependency Removed

**Change**: `github.com/dracory/crud/v2` has been removed from `go.mod` as it was no longer used by any code in the project.

**Action Required**:
- If your custom code depends on `github.com/dracory/crud/v2`, add it explicitly to your `go.mod`:
```bash
go get github.com/dracory/crud/v2
```
- Run `go mod tidy` to clean up unused dependencies

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
rm internal/app/stores_audit_test.go
rm internal/app/stores_blindindex_email.go
rm internal/app/stores_blindindex_first_name.go
rm internal/app/stores_blindindex_last_name.go
rm internal/app/stores_blog.go
rm internal/app/stores_blog_test.go
rm internal/app/stores_cache.go
rm internal/app/stores_cache_test.go
rm internal/app/stores_chat.go
rm internal/app/stores_chat_test.go
rm internal/app/stores_cms.go
rm internal/app/stores_cms_test.go
rm internal/app/stores_custom.go
rm internal/app/stores_custom_test.go
rm internal/app/stores_entity.go
rm internal/app/stores_entity_test.go
rm internal/app/stores_feed.go
rm internal/app/stores_feed_test.go
rm internal/app/stores_geo.go
rm internal/app/stores_geo_test.go
rm internal/app/stores_log.go
rm internal/app/stores_log_test.go
rm internal/app/stores_meta.go
rm internal/app/stores_meta_test.go
rm internal/app/stores_session.go
rm internal/app/stores_session_test.go
rm internal/app/stores_setting.go
rm internal/app/stores_setting_test.go
rm internal/app/stores_shop.go
rm internal/app/stores_shop_test.go
rm internal/app/stores_sqlfile.go
rm internal/app/stores_stats.go
rm internal/app/stores_stats_test.go
rm internal/app/stores_subscription.go
rm internal/app/stores_subscription_test.go
rm internal/app/stores_task.go
rm internal/app/stores_task_test.go
rm internal/app/stores_user.go
rm internal/app/stores_user_test.go
rm internal/app/stores_vault.go
rm internal/app/stores_vault_test.go
```

### Step 4: Verify Datastore Initialization
Ensure `internal/app/datastores.go`, `internal/app/datastores_test.go`, and `internal/config/store_builders.go` are present. If you have custom stores, add them to the `stores` slice in `dataStoresInitialize()` and create a builder function in `store_builders.go`. Add corresponding `_Success` and `_NotUsed` test cases in `datastores_test.go` following the existing pattern.

### Step 5: Remove Deleted Admin Shop Controller References
Search for and remove any references to the deleted admin shop controllers:
```bash
# Find references to deleted controllers
grep -r "NewCategoryManagerController\|NewCategoryCreateController\|NewCategoryUpdateController\|NewProductCreateController\|NewProductDeleteController\|NewProductManagerController\|NewProductUpdateController\|NewDiscountController\|shop/shared" .
```
Replace all admin shop routes with `NewShopAdminController(app).Handler`.

### Step 6: Update Dependencies
```bash
# Remove unused dependencies
go mod tidy

# If you use testify in your own tests, add it explicitly
go get github.com/stretchr/testify@v1.11.1

# If you use crud/v2 in your own code, add it explicitly
go get github.com/dracory/crud/v2
```

### Step 7: Standardize Import Order
Run gofmt to standardize import order:
```bash
gofmt -w .
```

### Step 8: Verify Build
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

### 5. Verify Admin Shop Routes
- Verify `/admin/shop` routes through `NewShopAdminController`
- Verify all sub-paths (categories, products, discounts, orders) are handled by `pkg/shopadmin`
- Verify no 404s for previously working shop admin pages

### 6. Verify No testify Dependencies
```bash
# Check if testify is still referenced
grep -r "stretchr/testify" . --include="*.go"
# If found in your own tests, either keep the dependency or refactor
```

---

## 📝 Additional Notes

### New Features
- **CSRF Secret Configuration**: Centralized CSRF secret management with environment-aware behavior (required in production, auto-generated in development)
- **Config-Driven Store Initialization**: Store initialization is now driven by a single config-driven list, making it easier to add, remove, or reorder stores
- **Exported Store Builders**: Store construction functions are now exported from the `config` package, enabling reuse in tests and other contexts
- **Unified Admin Shop Controller**: All admin shop routes now delegate to `pkg/shopadmin` through a single `shopAdminController`, eliminating the need for individual controller packages
- **Stripe Payment Support**: New `helpers.GenerateStripePaymentsCheckoutURL()` and payment code cache helpers for Stripe checkout integration
- **Shop Cart Link Constants**: New `SHOP_CART*`, `SHOP_CHECKOUT`, and `SHOP_PRODUCT` link constants with corresponding `websiteLinks` methods
- **Standard Library Tests**: Tests no longer depend on `testify` — all assertions use standard library patterns

### Removed Features
- `pkg/shopadmin/categories/category_manager_controller.go` — placeholder controller deleted
- All individual admin shop controllers in `internal/controllers/admin/shop/` (categories, discounts, products, productupdate, shared) — replaced by single `shopAdminController`
- `pkg/blogai/database.go` — dead/commented code removed
- All private `new*Store` functions in `internal/app/` — replaced by exported `config.New*Store` functions
- All individual `*StoreInitialize` functions in `internal/app/` — replaced by `setup*Store` functions in `datastores.go`
- `github.com/stretchr/testify` dependency removed
- `github.com/dracory/crud/v2` dependency removed

### Behavior Changes
- **Nil store checks removed from builders**: The old `new*Store` functions in `internal/app/` included `if st == nil { return nil, errors.New("xxx.NewStore returned a nil store") }` checks. The new `config.New*Store` functions do NOT include these checks — they rely on the underlying store package's `NewStore()` to return an error on failure.
- **Database nil checks consolidated**: The old `new*Store` functions each had individual `if db == nil { return nil, errors.New("database is not initialized") }` checks. The new `dataStoresInitialize()` checks `if app.GetDatabase() == nil` once at the top level instead.
- **Config nil checks consolidated**: The old `*StoreInitialize` functions each had `if app.GetConfig() == nil` checks. The new `dataStoresInitialize()` checks once at the top.
- **EnableDebug now called inside store builders**: In v0.34.0, `EnableDebug(app.GetConfig().GetAppDebug())` was called in the `*StoreInitialize` function after store creation. In v0.35.0, `EnableDebug(debug)` is called inside `config.New*Store` builder functions. This means if you call `config.New*Store` directly (e.g., in tests), debug mode is already enabled based on the `debug` parameter you pass.
- **Import order fixed in `global_middlewares.go`**: `project/internal/app` import moved before `project/internal/middlewares` to fix import ordering.
- **`AUTH_CSRF_SECRET` not added to `.env.example`**: The new environment variable was added to constants and config logic but was not added to the `.env.example` template file. You should add it manually to your `.env` files.
- **Shop cart routes commented out by default**: The new `cart.Routes(app)` call in `internal/controllers/website/routes.go` is commented out. Uncomment to enable.
- **Route handler type changed for admin shop**: Changed from `SetHTMLHandler` (returning string) to `SetHandler` (standard `http.HandlerFunc`) for admin shop routes.
- **Stripe currency hardcoded to GBP**: The `GenerateStripePaymentsCheckoutURL` helper currently hardcodes `priceCurrency = "GBP"`. Modify the helper if you need a different currency.

### Architecture Improvements
- **Single source of truth for store construction**: All table names and store options are now defined in `internal/config/store_builders.go`
- **Cleaner separation of concerns**: Store construction (config package) vs. store wiring (app package)
- **Reduced file count**: 24 files consolidated into 2 (`datastores.go` + `store_builders.go`)
- **Admin shop controller consolidation**: 70+ files in `internal/controllers/admin/shop/` replaced by single `shop_controller.go` delegating to `pkg/shopadmin`
- **Standardized import order**: `project/internal/app` consistently placed before other internal packages across 28+ files

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

### Issue 5: Admin Shop Controller Build Errors
**Symptom**: Build errors like `undefined: categories.NewCategoryManagerController`, `undefined: shopProducts.NewProductCreateController`, `undefined: shopDiscounts.NewDiscountController`, `undefined: shared.CONTROLLER_CATEGORIES`
**Solution**: All individual admin shop controllers have been deleted. Replace all references with `NewShopAdminController(app).Handler`:
```go
// Old
import "project/internal/controllers/admin/shop/products"
handler := shopProducts.NewProductCreateController(app).Handler

// New
import admin "project/internal/controllers/admin/shop"
handler := admin.NewShopAdminController(app).Handler
```

### Issue 6: testify Not Found
**Symptom**: Build error: `cannot find package "github.com/stretchr/testify"` in your tests
**Solution**: Add testify explicitly to your `go.mod` if your own tests use it:
```bash
go get github.com/stretchr/testify@v1.11.1
```
Or refactor your tests to use standard library assertions.

### Issue 7: Nil Store Returned Without Error
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
