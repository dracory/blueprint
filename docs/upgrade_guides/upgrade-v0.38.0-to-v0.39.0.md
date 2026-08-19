# Upgrade Guide: v0.38.0 to v0.39.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.38.0 to v0.39.0.

## Overview

This release **extracts the four in-repo admin packages** (`pkg/blogadmin`, `pkg/fileadmin`, `pkg/logadmin`, `pkg/shopadmin`) **into standalone external modules** (`github.com/dracory/blogadmin`, `github.com/dracory/fileadmin`, `github.com/dracory/logadmin`, `github.com/dracory/shopadmin`). A new **adapters layer** (`internal/controllers/admin/adapters/`) bridges the external packages to the blueprint's internal services.

The external packages are now decoupled from the host project: they no longer accept a `Registry app.AppInterface` field. Instead, each package takes its dependencies directly (stores, logger, layout function, LLM factory, customer resolver). The `Handle()` method no longer returns an HTML string — it writes directly to the `http.ResponseWriter`. The `Routes()` helper functions have been removed from the admin packages; the host project now registers a single **catch-all route** per admin section that delegates to the controller's `Handler`.

**External Admin Packages** — `pkg/blogadmin`, `pkg/fileadmin`, `pkg/logadmin`, and `pkg/shopadmin` have been deleted from the repository. Their code now lives in standalone modules published to `github.com/dracory/*admin`. The host project wires them up via the new `internal/controllers/admin/adapters/` package, which provides `NewLayoutFunc`, `NewLlmFactory`, and `NewUserStoreCustomerResolver`.

**`Registry` Field Removed** — `AdminOptions.Registry app.AppInterface` has been removed from all four admin packages. Each package now takes only the stores/services it actually needs (e.g. `Store`, `Logger`, `FuncLayout`, `LlmFactory`, `CustomerResolver`, `Storage`, `RootDirPath`). This keeps the external packages free of any host-project dependency.

**`Handle()` Returns Nothing** — `AdminInterface.Handle(w, r)` previously returned an HTML `string` that the caller wrote to the response. It now writes directly to the `http.ResponseWriter` and returns nothing. Callers no longer need to set the `Content-Type` header or write the returned string.

**Catch-All Routing** — The admin packages no longer expose a `Routes()` function. The host project registers a base route and a `/*` catch-all route per admin section, both delegating to the same controller `Handler`. The external package's `Handle()` inspects the request path and dispatches internally.

**`shopadmin` Customer Resolution** — `shopadmin` no longer depends on `userstore`. Customer name/email resolution for order views is now done via a `CustomerResolverInterface` that the host project implements. The blueprint provides `adapters.NewUserStoreCustomerResolver(userStore)`.

**Key Changes:**
- `pkg/blogadmin`, `pkg/fileadmin`, `pkg/logadmin`, `pkg/shopadmin` deleted — migrated to external modules
- New `internal/controllers/admin/adapters/` package (`adapters.go`)
- `AdminOptions.Registry` removed in all four admin packages; replaced with explicit stores/services
- `Handle(w, r) string` → `Handle(w, r)` (no return value) in all four admin packages
- `blogadmin.Routes()`, `shopadmin.Routes()` removed — host registers catch-all routes
- `shopadmin`: `userstore` dependency replaced with `CustomerResolverInterface`
- `blogadmin`: new `LlmFactory`, `CustomStore`, `SettingStore`, `FuncLayout`, `FileManagerURL` fields
- `logadmin`: new `Store`, `Logger`, `FuncLayout`, `FileManagerURL` fields
- `fileadmin`: new `Storage`, `RootDirPath`, `FuncLayout` fields; root-dir derivation moved to controller
- Import path `project/pkg/blogadmin/post_update` → `github.com/dracory/blogadmin/post_update`
- `go.mod`: added `blogadmin v0.2.0`, `fileadmin v0.1.0`, `logadmin v0.1.0`, `shopadmin v0.3.0`
- `go.mod`: bumped `blogstore` v1.34.0 → v1.35.0, `shopstore` v1.25.0 → v1.26.0
- `go.mod`: `versionstore`, `wf`, `flosch/pongo2` moved from direct to indirect
- `AGENTS.md` updated with adapters layer and external admin packages documentation

---

## ⚠️ Breaking Changes

---

### 1. In-Repo Admin Packages Removed — Migrated to External Modules

**Change**: The four admin packages that lived inside the blueprint repository under `pkg/` have been deleted and published as standalone external modules:

| Old in-repo path | New external module |
|---|---|
| `project/pkg/blogadmin` | `github.com/dracory/blogadmin` v0.2.0 |
| `project/pkg/fileadmin` | `github.com/dracory/fileadmin` v0.1.0 |
| `project/pkg/logadmin` | `github.com/dracory/logadmin` v0.1.0 |
| `project/pkg/shopadmin` | `github.com/dracory/shopadmin` v0.3.0 |

All sub-packages keep the same names (e.g. `post_update`, `file_manager`, `log_manager`, `product_manager`), so only the import **root** changes.

**Old Usage**:
```go
import (
	blogadmin "project/pkg/blogadmin"
	fileadmin "project/pkg/fileadmin"
	logadmin "project/pkg/logadmin"
	shopadmin "project/pkg/shopadmin"

	// sub-packages
	"project/pkg/blogadmin/post_update"
)
```

**New Usage**:
```go
import (
	blogadmin "github.com/dracory/blogadmin"
	fileadmin "github.com/dracory/fileadmin"
	logadmin "github.com/dracory/logadmin"
	shopadmin "github.com/dracory/shopadmin"

	// sub-packages
	"github.com/dracory/blogadmin/post_update"
)
```

**Action Required**:
- Find all imports of the old in-repo packages:
  ```bash
  grep -rln "project/pkg/blogadmin\|project/pkg/fileadmin\|project/pkg/logadmin\|project/pkg/shopadmin" --include="*.go" .
  ```
- Replace `project/pkg/blogadmin` with `github.com/dracory/blogadmin` (and likewise for the other three).
- Delete the old `pkg/blogadmin`, `pkg/fileadmin`, `pkg/logadmin`, `pkg/shopadmin` directories from your project if you had not already done so. **They are no longer part of the blueprint repository.**
- Add the new external dependencies:
  ```bash
  go get github.com/dracory/blogadmin@v0.2.0
  go get github.com/dracory/fileadmin@v0.1.0
  go get github.com/dracory/logadmin@v0.1.0
  go get github.com/dracory/shopadmin@v0.3.0
  go mod tidy
  ```
- If you had **custom modifications** in your `pkg/*admin` copies, port them to the external module (preferably via the new callback/interface extension points: `FuncLayout`, `LlmFactory`, `CustomerResolver`) or maintain a fork.

**Migration Command**:
```bash
# Bulk-rename import paths (review the diff afterwards)
grep -rl "project/pkg/blogadmin" --include="*.go" . | xargs sed -i 's|project/pkg/blogadmin|github.com/dracory/blogadmin|g'
grep -rl "project/pkg/fileadmin" --include="*.go" . | xargs sed -i 's|project/pkg/fileadmin|github.com/dracory/fileadmin|g'
grep -rl "project/pkg/logadmin" --include="*.go" . | xargs sed -i 's|project/pkg/logadmin|github.com/dracory/logadmin|g'
grep -rl "project/pkg/shopadmin" --include="*.go" . | xargs sed -i 's|project/pkg/shopadmin|github.com/dracory/shopadmin|g'
```

---

### 2. `AdminOptions.Registry` Removed — Explicit Dependencies Instead

**Change**: All four admin packages no longer accept a `Registry app.AppInterface` field on `AdminOptions`. Each package now takes only the stores and services it actually needs. This decouples the external packages from any host-project type.

**Old Usage** (blogadmin):
```go
admin, err := blogadmin.New(blogadmin.AdminOptions{
	Store:        controller.app.GetBlogStore(),
	AdminHomeURL: links.Admin().Home(),
	BlogAdminURL: links.Admin().Blog(),
	Registry:     controller.app, // <-- removed
})
```

**New Usage** (blogadmin):
```go
admin, err := blogadmin.New(blogadmin.AdminOptions{
	Store:          controller.app.GetBlogStore(),
	Logger:         controller.app.GetLogger(),
	CustomStore:    controller.app.GetCustomStore(),
	SettingStore:   controller.app.GetSettingStore(),
	LlmFactory:     adapters.NewLlmFactory(controller.app),
	FuncLayout:     adapters.NewLayoutFunc(controller.app),
	AdminHomeURL:   links.Admin().Home(),
	BlogAdminURL:   links.Admin().Blog(),
	FileManagerURL: links.Admin().FileManager(),
	AuthUserID: func(r *http.Request) string {
		user := helpers.GetAuthUser(r)
		if user == nil {
			return ""
		}
		return user.GetID()
	},
})
```

The new `AdminOptions` fields per package:

| Package | New required fields | New optional fields | Removed |
|---|---|---|---|
| `blogadmin` | `Store`, `Logger` | `CustomStore`, `SettingStore`, `LlmFactory`, `FuncLayout`, `FileManagerURL` | `Registry` |
| `shopadmin` | `Store`, `Logger` | `CustomerResolver`, `FuncLayout`, `FileManagerURL` | `Registry` |
| `fileadmin` | `Storage`, `RootDirPath` | `FuncLayout` | `Registry` |
| `logadmin` | `Store`, `Logger` | `FuncLayout`, `FileManagerURL` | `Registry` |

**Action Required**:
- Remove the `Registry:` field from every `AdminOptions{...}` literal.
- Add the new explicit dependency fields as shown above. Use the `adapters` package (see Breaking Change #5) to supply `FuncLayout`, `LlmFactory`, and `CustomerResolver`.
- `Logger` is `*slog.Logger` — pass `app.GetLogger()`.
- `FuncLayout` is optional but recommended so the admin UI renders inside your project layout. If nil, a bare-bones default HTML page is used.
- For `blogadmin` AI controllers: `CustomStore`, `SettingStore`, and `LlmFactory` are required for AI features. If nil, AI controllers return an error to the user instead of panicking.
- For `shopadmin` order views: `CustomerResolver` is optional. If nil, customer fields stay empty and customer filtering is disabled.

**Migration Command**:
```bash
# Find all AdminOptions literals that still pass Registry
grep -rn "Registry:" --include="*.go" internal/controllers/admin/
```

---

### 3. `Handle()` No Longer Returns a String

**Change**: `AdminInterface.Handle(w http.ResponseWriter, r *http.Request)` previously returned an HTML `string`. It now writes directly to the `http.ResponseWriter` and returns nothing. The caller no longer sets the `Content-Type` header or writes the returned string.

**Old Usage** (fileadmin):
```go
html := admin.Handle(w, r)

if html != "" {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil { // #nosec G705 -- html is built by trusted admin handler
		if logger := c.app.GetLogger(); logger != nil {
			logger.Error("At FileManagerController > Handler", "write_error", err.Error())
		}
	}
}
```

**New Usage** (fileadmin):
```go
admin.Handle(w, r)
```

**Action Required**:
- Find all call sites that capture the return value of `Handle()`:
  ```bash
  grep -rn ":= .*\.Handle(w, r)" --include="*.go" .
  grep -rn "html := .*Handle(" --include="*.go" .
  ```
- Replace `html := admin.Handle(w, r)` with `admin.Handle(w, r)`.
- Remove the subsequent `w.Header().Set("Content-Type", ...)` and `w.Write([]byte(html))` block — the external package now handles writing the response itself.
- This applies to all four admin packages (`blogadmin`, `shopadmin`, `fileadmin`, `logadmin`).

---

### 4. `Routes()` Functions Removed — Catch-All Route Registration

**Change**: The admin packages no longer expose a `Routes()` function that returned `[]rtr.RouteInterface`. The host project now registers a base route and a `/*` catch-all route per admin section, both delegating to the same controller `Handler`. The external package's `Handle()` inspects the request path and dispatches to the correct sub-controller internally.

**Old Usage** (admin routes.go):
```go
import (
	"project/pkg/blogadmin"
	"project/pkg/logadmin"
	"project/pkg/shopadmin"
)

blogRoutes, err := blogadmin.Routes(app, blogadmin.AdminOptions{
	Store:        app.GetBlogStore(),
	AdminHomeURL: links.Admin().Home(),
	BlogAdminURL: links.Admin().Blog(),
	Registry:     app,
})
if err == nil {
	adminRoutes = append(adminRoutes, blogRoutes...)
}

shopRoutes, err := shopadmin.Routes(app, shopadmin.AdminOptions{
	Registry:       app,
	AdminHomeURL:   links.Admin().Home(),
	ShopAdminURL:   links.Admin().Shop(),
	FileManagerURL: links.Admin().FileManager(),
})
if err == nil {
	adminRoutes = append(adminRoutes, shopRoutes...)
}

logRoutes, err := logadmin.Routes(app)
if err == nil {
	adminRoutes = append(adminRoutes, logRoutes...)
}
```

**New Usage** (admin routes.go):
```go
import (
	"net/http"
	"project/internal/app"
	adminBlog "project/internal/controllers/admin/blog"
	adminLogs "project/internal/controllers/admin/logs"
	adminShop "project/internal/controllers/admin/shop"
	"project/internal/links"
	"github.com/dracory/rtr"
)

// Blog: base + catch-all
blogController := adminBlog.NewBlogAdminController(app)
blog := rtr.NewRoute().
	SetName("Admin > Blog").
	SetPath(links.ADMIN_BLOG).
	SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
		blogController.Handler(w, r)
		return ""
	})
blogCatchAll := rtr.NewRoute().
	SetName("Admin > Blog > Catchall").
	SetPath(links.ADMIN_BLOG + links.CATCHALL).
	SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
		blogController.Handler(w, r)
		return ""
	})
adminRoutes = append(adminRoutes, blog, blogCatchAll)

// Shop: base + catch-all
shopController := adminShop.NewShopAdminController(app)
shop := rtr.NewRoute().
	SetName("Admin > Shop").
	SetPath(links.ADMIN_SHOP).
	SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
		shopController.Handler(w, r)
		return ""
	})
shopCatchAll := rtr.NewRoute().
	SetName("Admin > Shop > Catchall").
	SetPath(links.ADMIN_SHOP + links.CATCHALL).
	SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
		shopController.Handler(w, r)
		return ""
	})
adminRoutes = append(adminRoutes, shop, shopCatchAll)

// Logs: base + catch-all (registered via adminLogs.Routes(app))
logRoutes, err := adminLogs.Routes(app)
if err == nil {
	adminRoutes = append(adminRoutes, logRoutes...)
}
```

For `fileadmin` and `logadmin`, the catch-all route is registered inside their own `routes.go` files:

```go
// internal/controllers/admin/files/routes.go
fileManagerCatchAll := rtr.NewRoute().
	SetName("Admin > File Manager > Catchall").
	SetPath(links.ADMIN_FILE_MANAGER + links.CATCHALL).
	SetHandler(NewFileManagerController(app).Handler)

return []rtr.RouteInterface{fileManager, fileManagerCatchAll}, nil
```

**Action Required**:
- Remove all calls to `blogadmin.Routes(...)`, `shopadmin.Routes(...)`, `logadmin.Routes(...)` — these functions no longer exist in the external packages.
- For `blogadmin` and `shopadmin`: register a base route and a catch-all route (`links.ADMIN_BLOG` + `links.CATCHALL`, `links.ADMIN_SHOP` + `links.CATCHALL`) that both call the controller's `Handler`. Use `SetHTMLHandler` with a wrapper that returns `""` (since `Handler` now writes directly to `w`).
- For `fileadmin` and `logadmin`: add a catch-all route alongside the existing base route in their `routes.go` files.
- `links.CATCHALL` is the constant `"/*"` — append it to the base admin path to match sub-paths.
- Update route-count assertions in tests (e.g. `routes_test.go`) to expect the extra catch-all route.

**Migration Command**:
```bash
# Find all calls to the removed Routes() functions
grep -rn "blogadmin\.Routes(\|shopadmin\.Routes(\|logadmin\.Routes(" --include="*.go" .
```

---

### 5. New `internal/controllers/admin/adapters/` Package

**Change**: A new package `internal/controllers/admin/adapters/` provides bridge implementations that wire the external admin packages to the blueprint's internal services. It exposes three constructors:

- `NewLayoutFunc(app)` — returns a `FuncLayout` adapter that bridges the external packages' layout callback signature to `layouts.NewAdminLayout`. The same function can be passed to both `blogadmin` and `shopadmin` (and `fileadmin`/`logadmin`) since they share the same anonymous-struct `options` shape.
- `NewLlmFactory(app)` — returns an LLM factory function for `blogadmin`'s AI controllers. Uses the blueprint config to select the provider (mock in testing, OpenRouter otherwise).
- `NewUserStoreCustomerResolver(userStore)` — returns a `*UserStoreCustomerResolver` implementing `shopadmin.CustomerResolverInterface`, backed by the blueprint's `userstore`.

**New Usage**:
```go
import (
	"project/internal/controllers/admin/adapters"
)

// blogadmin
admin, err := blogadmin.New(blogadmin.AdminOptions{
	// ...
	LlmFactory: adapters.NewLlmFactory(controller.app),
	FuncLayout: adapters.NewLayoutFunc(controller.app),
	// ...
})

// shopadmin
admin, err := shopadmin.New(shopadmin.AdminOptions{
	// ...
	CustomerResolver: adapters.NewUserStoreCustomerResolver(controller.app.GetUserStore()),
	FuncLayout:       adapters.NewLayoutFunc(controller.app),
	// ...
})

// fileadmin / logadmin
admin, err := fileadmin.New(fileadmin.AdminOptions{
	// ...
	FuncLayout: adapters.NewLayoutFunc(c.app),
	// ...
})
```

**Action Required**:
- Create the `internal/controllers/admin/adapters/` package in your project (copy `adapters.go` from the blueprint).
- Pass `adapters.NewLayoutFunc(app)` as the `FuncLayout` field to every admin package's `AdminOptions`.
- Pass `adapters.NewLlmFactory(app)` as the `LlmFactory` field to `blogadmin.AdminOptions` (required for AI controllers).
- Pass `adapters.NewUserStoreCustomerResolver(app.GetUserStore())` as the `CustomerResolver` field to `shopadmin.AdminOptions` (optional but needed for order customer resolution).
- The `FuncLayout` signature uses an anonymous struct for `options` that is identical across all four admin packages, so a single `NewLayoutFunc` works for all of them.

---

### 6. `shopadmin` Customer Resolution via `CustomerResolverInterface`

**Change**: `shopadmin` no longer imports `github.com/dracory/userstore` or accepts a `UserStore` field. Customer name/email resolution for order views is now done via a `CustomerResolverInterface` that the host project implements:

```go
type CustomerResolverInterface interface {
	FindByID(ctx context.Context, customerID string) (name, email string)
	SearchIDs(ctx context.Context, name, email string) ([]string, error)
}
```

**Old Usage**:
```go
// shopadmin internally called registry.GetUserStore().UserFindByID(...)
admin, err := shopadmin.New(shopadmin.AdminOptions{
	Registry:       controller.app, // provided GetUserStore() internally
	AdminHomeURL:   links.Admin().Home(),
	ShopAdminURL:   links.Admin().Shop(),
	FileManagerURL: links.Admin().FileManager(),
})
```

**New Usage**:
```go
admin, err := shopadmin.New(shopadmin.AdminOptions{
	Store:            controller.app.GetShopStore(),
	Logger:           controller.app.GetLogger(),
	CustomerResolver: adapters.NewUserStoreCustomerResolver(controller.app.GetUserStore()),
	FuncLayout:       adapters.NewLayoutFunc(controller.app),
	AdminHomeURL:     links.Admin().Home(),
	ShopAdminURL:     links.Admin().Shop(),
	FileManagerURL:   links.Admin().FileManager(),
	AuthUserID: func(r *http.Request) string {
		user := helpers.GetAuthUser(r)
		if user == nil {
			return ""
		}
		return user.GetID()
	},
})
```

**Action Required**:
- Implement `shopadmin.CustomerResolverInterface` in your project, or use `adapters.NewUserStoreCustomerResolver(app.GetUserStore())` which is the blueprint's built-in implementation backed by `userstore`.
- `CustomerResolver` is optional (nil-safe). If nil, customer fields in order views stay empty and customer filtering is disabled.
- If you had a custom user store that is not `userstore`, implement the interface directly.

---

### 7. `fileadmin` Root-Dir Derivation Moved to Controller

**Change**: The old in-repo `pkg/fileadmin` derived the root directory path internally from `Registry.GetConfig().GetMediaRoot()`. The external `github.com/dracory/fileadmin` package is decoupled from config and takes `Storage filesystem.StorageInterface` and `RootDirPath string` directly. The root-dir derivation logic now lives in the blueprint's `FileManagerController`.

**Old Usage**:
```go
admin, err := fileadmin.New(fileadmin.AdminOptions{
	Registry:     c.app, // internally derived root dir from config
	AdminHomeURL: links.Admin().Home(),
	FileAdminURL: links.Admin().FileManager(),
	AuthUserID:   func(r *http.Request) string { ... },
})
```

**New Usage**:
```go
cfg := c.app.GetConfig()

// Derive root dir path from config (same logic as the old pkg/fileadmin)
rootDirPath := strings.TrimSpace(cfg.GetMediaRoot())
rootDirPath = strings.Trim(rootDirPath, "/")
rootDirPath = strings.Trim(rootDirPath, ".")
rootDirPath = "/" + rootDirPath

admin, err := fileadmin.New(fileadmin.AdminOptions{
	Storage:      c.app.GetSqlFileStorage(),
	RootDirPath:  rootDirPath,
	FuncLayout:   adapters.NewLayoutFunc(c.app),
	AdminHomeURL: links.Admin().Home(),
	FileAdminURL: links.Admin().FileManager(),
	AuthUserID:   func(r *http.Request) string { ... },
})
```

**Action Required**:
- In your `FileManagerController.Handler`, derive `rootDirPath` from your config's media root (trim surrounding `/` and `.`, then prefix with `/`).
- Pass `app.GetSqlFileStorage()` as `Storage` (type `filesystem.StorageInterface`).
- Add the `strings` import to the controller file.
- Remove the old `Registry:` field.

---

### 8. `post_update` Import Path Changed in Website Blog Controller

**Change**: The website blog post controller (`internal/controllers/website/blog/post/post_controller.go`) imported `project/pkg/blogadmin/post_update`. This import path is now `github.com/dracory/blogadmin/post_update`.

**Old Usage**:
```go
import (
	"project/pkg/blogadmin/post_update"
)
```

**New Usage**:
```go
import (
	"github.com/dracory/blogadmin/post_update"
)
```

**Action Required**:
- Update the import in `internal/controllers/website/blog/post/post_controller.go` (and any other file importing `project/pkg/blogadmin/post_update`).
- This is covered by the bulk rename in Breaking Change #1.

---

### 9. Dependency Updates in `go.mod`

**Change**: Four new external admin modules were added as direct dependencies. `blogstore` and `shopstore` were bumped. `versionstore`, `wf`, and `flosch/pongo2` were moved from direct to indirect (they are still in the module graph but no longer imported directly by the blueprint).

**Old `go.mod`** (key direct dependencies):
```go
github.com/dracory/blogstore v1.34.0
github.com/dracory/shopstore v1.25.0
github.com/dracory/versionstore v1.7.1   // direct
github.com/dracory/wf v0.6.0             // direct
github.com/flosch/pongo2/v6 v6.1.0       // direct
```

**New `go.mod`** (key direct dependencies):
```go
github.com/dracory/blogadmin v0.2.0      // NEW
github.com/dracory/blogstore v1.35.0
github.com/dracory/fileadmin v0.1.0      // NEW
github.com/dracory/logadmin v0.1.0       // NEW
github.com/dracory/shopadmin v0.3.0      // NEW
github.com/dracory/shopstore v1.26.0
// github.com/dracory/versionstore — moved to indirect
// github.com/dracory/wf — moved to indirect
// github.com/flosch/pongo2/v6 — moved to indirect
```

**Action Required**:
- Add the new direct dependencies and bump the existing ones:
  ```bash
  go get github.com/dracory/blogadmin@v0.2.0
  go get github.com/dracory/fileadmin@v0.1.0
  go get github.com/dracory/logadmin@v0.1.0
  go get github.com/dracory/shopadmin@v0.3.0
  go get github.com/dracory/blogstore@v1.35.0
  go get github.com/dracory/shopstore@v1.26.0
  go mod tidy
  ```
- After `go mod tidy`, `versionstore`, `wf`, and `flosch/pongo2` will move to the indirect block automatically if your code no longer imports them directly.
- If your application code directly imports `versionstore`, `wf`, or `pongo2`, they will remain direct dependencies — that is fine.
- Review the resulting `go.mod` / `go.sum` diff for any unexpected version jumps.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.39.0"
```

### Step 2: Add the external admin dependencies

```bash
go get github.com/dracory/blogadmin@v0.2.0
go get github.com/dracory/fileadmin@v0.1.0
go get github.com/dracory/logadmin@v0.1.0
go get github.com/dracory/shopadmin@v0.3.0
go get github.com/dracory/blogstore@v1.35.0
go get github.com/dracory/shopstore@v1.26.0
```

### Step 3: Create the `adapters` package

Create `internal/controllers/admin/adapters/adapters.go` with the three adapters (`NewLayoutFunc`, `NewLlmFactory`, `NewUserStoreCustomerResolver`). Copy the file from the blueprint repository. Ensure the imports match your module path (`project/internal/app`, `project/internal/layouts`).

### Step 4: Rename import paths

Replace all `project/pkg/*admin` imports with the external module paths:

```bash
grep -rl "project/pkg/blogadmin" --include="*.go" . | xargs sed -i 's|project/pkg/blogadmin|github.com/dracory/blogadmin|g'
grep -rl "project/pkg/fileadmin" --include="*.go" . | xargs sed -i 's|project/pkg/fileadmin|github.com/dracory/fileadmin|g'
grep -rl "project/pkg/logadmin" --include="*.go" . | xargs sed -i 's|project/pkg/logadmin|github.com/dracory/logadmin|g'
grep -rl "project/pkg/shopadmin" --include="*.go" . | xargs sed -i 's|project/pkg/shopadmin|github.com/dracory/shopadmin|g'
```

### Step 5: Update each admin controller's `AdminOptions`

For each of the four admin controllers (`blog_controller.go`, `shop_controller.go`, `file_manager_controller.go`, `logs_controller.go`):

1. Remove the `Registry:` field.
2. Add the new explicit dependency fields (see Breaking Change #2 for the per-package field list).
3. Pass `adapters.NewLayoutFunc(app)` as `FuncLayout`.
4. For `blogadmin`: also pass `adapters.NewLlmFactory(app)`, `app.GetCustomStore()`, `app.GetSettingStore()`, `app.GetLogger()`, and `links.Admin().FileManager()`.
5. For `shopadmin`: also pass `adapters.NewUserStoreCustomerResolver(app.GetUserStore())`, `app.GetLogger()`.
6. For `fileadmin`: derive `rootDirPath` from config and pass `app.GetSqlFileStorage()` as `Storage` (see Breaking Change #7).
7. For `logadmin`: also pass `app.GetLogStore()`, `app.GetLogger()`, and `links.Admin().FileManager()`.

### Step 6: Update `Handle()` call sites

Remove the `html := admin.Handle(w, r)` pattern and the subsequent `w.Write` block. Replace with a plain `admin.Handle(w, r)` call (see Breaking Change #3).

```bash
grep -rn ":= .*\.Handle(w, r)" --include="*.go" internal/controllers/admin/
```

### Step 7: Replace `Routes()` calls with catch-all route registration

In `internal/controllers/admin/routes.go`:

1. Remove the `blogadmin.Routes(...)`, `shopadmin.Routes(...)` calls and their imports.
2. Register a base route + catch-all route for blog and shop, delegating to the controller `Handler` (see Breaking Change #4).
3. Keep `adminLogs.Routes(app)` and `adminFiles.Routes(app)` — these now internally register the catch-all route.
4. Add the `"net/http"` import and the `adminBlog` / `adminShop` / `adminLogs` controller imports.

In `internal/controllers/admin/files/routes.go` and `internal/controllers/admin/logs/routes.go`: add the catch-all route (see Breaking Change #4).

### Step 8: Update route-count tests

Update the route-count assertions in `routes_test.go` files to expect the extra catch-all route:

- `internal/controllers/admin/files/routes_test.go`: expect `2` routes (file manager + catchall).
- `internal/controllers/admin/logs/routes_test.go`: expect `2` routes (logs + catchall).

### Step 9: Delete the old `pkg/*admin` directories

Once all imports are updated and the build passes, delete the old in-repo admin packages:

```bash
rm -rf pkg/blogadmin pkg/fileadmin pkg/logadmin pkg/shopadmin
go mod tidy
```

### Step 10: Run `go mod tidy` and verify

```bash
go mod tidy
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...` — confirms all import paths resolve and `AdminOptions` literals compile.
2. **Unit Tests**: `go test ./...` — confirms route registration, controller wiring, and adapter behavior.
3. **Route Tests**: Verify that `internal/controllers/admin/files/routes_test.go` and `internal/controllers/admin/logs/routes_test.go` expect 2 routes (base + catchall).
4. **Integration Tests**: `go test -tags=integration ./...` — confirms the admin endpoints respond correctly end-to-end.
5. **Manual Smoke Test**: Start the server (`go run ./cmd/server`) and verify:
   - `/admin/blog` and `/admin/blog/*` render the blog admin (including AI controllers if configured).
   - `/admin/shop` and `/admin/shop/*` render the shop admin (including order customer resolution).
   - `/admin/file-manager` and `/admin/file-manager/*` render the file manager.
   - `/admin/logs` and `/admin/logs/*` render the log manager.
   - The admin layout (menus, branding) renders correctly via `FuncLayout`.
6. **AI Controllers** (blogadmin): Set a valid `OPENROUTER_API_KEY` and verify the AI title generator / post editor work. In testing mode (`APP_ENV=testing`), the mock LLM provider is used automatically.
7. **Customer Resolution** (shopadmin): Open an order with a known customer ID and verify the customer name/email appear in the order details view.

---

## 📝 Additional Notes

- **Why the extraction?** The admin packages were large (the four `pkg/*admin` directories accounted for ~32,000 deleted lines) and tightly coupled to the blueprint's `app.AppInterface`. Extracting them to standalone modules makes them reusable across projects, testable in isolation, and independently versioned.
- **`FuncLayout` shared signature**: All four external admin packages use the same anonymous-struct shape for the `options` parameter of `FuncLayout` (`Styles`, `StyleURLs`, `Scripts`, `ScriptURLs`). This means a single `adapters.NewLayoutFunc(app)` can be passed to all four packages.
- **`Handle()` writes directly**: The external packages set their own `Content-Type` headers and write the response body. Do not wrap the call with additional response writing.
- **Catch-all routing**: The external packages parse the request path internally to dispatch to sub-controllers (e.g. `/admin/blog/post-manager`, `/admin/blog/ai-title-generator`). The host project only needs to register the base path and the `/*` catch-all.
- **`AGENTS.md` updated**: The project structure section now documents the `internal/controllers/admin/adapters/` directory and the external admin packages.
- **No `.env` changes**: This release does not introduce or remove any environment variables.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: blogadmin.Routes` (or `shopadmin.Routes`, `logadmin.Routes`)
**Cause**: The `Routes()` functions were removed from the external packages.
**Solution**: Replace with catch-all route registration (see Breaking Change #4).

### Issue: `cannot use controller.app (variable of type app.AppInterface) as ... value in struct literal`
**Cause**: The `Registry` field was removed from `AdminOptions`.
**Solution**: Remove the `Registry:` field and add the explicit dependency fields (see Breaking Change #2).

### Issue: `admin.Handle(w, r) (string) used as value` or `too many return values`
**Cause**: `Handle()` no longer returns a string.
**Solution**: Change `html := admin.Handle(w, r)` to `admin.Handle(w, r)` and remove the `w.Write` block (see Breaking Change #3).

### Issue: `undefined: adapters` in admin controllers
**Cause**: The `adapters` package has not been created or imported.
**Solution**: Create `internal/controllers/admin/adapters/adapters.go` (copy from blueprint) and add `"project/internal/controllers/admin/adapters"` to the controller imports.

### Issue: `fileadmin.AdminOptions` missing `Storage` / `RootDirPath`
**Cause**: The external `fileadmin` package takes `Storage` and `RootDirPath` instead of `Registry`.
**Solution**: Derive `rootDirPath` from config and pass `app.GetSqlFileStorage()` (see Breaking Change #7).

### Issue: Route tests fail with "expected 1 route, got 2"
**Cause**: The catch-all route was added but the test assertion was not updated.
**Solution**: Update the test to expect 2 routes (base + catchall).

### Issue: AI controllers in blogadmin return "LlmFactory is nil"
**Cause**: `LlmFactory` was not passed to `blogadmin.AdminOptions`.
**Solution**: Pass `adapters.NewLlmFactory(app)` as the `LlmFactory` field.

### Issue: Customer fields empty in shopadmin order views
**Cause**: `CustomerResolver` was not passed (or is nil).
**Solution**: Pass `adapters.NewUserStoreCustomerResolver(app.GetUserStore())` as the `CustomerResolver` field.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- External admin modules: `github.com/dracory/blogadmin`, `github.com/dracory/fileadmin`, `github.com/dracory/logadmin`, `github.com/dracory/shopadmin`
- For questions about this upgrade guide, open an issue on the blueprint repository.
