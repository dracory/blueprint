# Upgrade Guide: v0.39.0 to v0.40.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.39.0 to v0.40.0.

## Overview

This release continues the **externalization** trend started in v0.39.0. Three more in-repo packages are removed and replaced with standalone external modules, and a large number of internal helpers are migrated into the shared `github.com/dracory/base` module. The goal is the same as v0.39.0: keep the blueprint repository thin and push reusable code into `base` so it can be shared across projects without copy/paste.

**`pkg/useradmin` → external `github.com/dracory/useradmin` v0.2.0** — The in-repo user admin package is deleted and replaced with the external module. As with the v0.39.0 admin packages, `Registry app.AppInterface` is removed and replaced with explicit callback adapters. The `Handle()` method no longer returns an HTML string. The callback signatures are simplified (userID passed directly instead of event structs), `OnUserUpdate` is renamed to `OnUserUpdated`, and `OnUserSearch` now takes a `[]SearchCondition` slice with OR/AND combine logic. The `VaultTokenizer`, `SessionStore`, `BlindIndexStore`, and `GeoStore` fields are replaced with `UserPiiSeal`/`UserPiiUnseal`/`UsersPiiUnseal`, `OnUserImpersonate`, `OnUserSearch`, and `GeoResolver` callback adapters respectively. The `SecureCookie` field is removed — the Secure flag is now derived from the request TLS state inside the `OnUserImpersonate` adapter.

**`pkg/social` → external `github.com/dracory/social` v0.1.0** — The in-repo social sharing package is deleted and replaced with the external module. Only the import root changes.

**`pkg/blogai` removed** — The in-repo blog AI package is deleted. It is no longer used by the blueprint; AI blog features are handled by `blogadmin`'s `LlmFactory` callback (introduced in v0.39.0).

**Internal helpers migrated to `github.com/dracory/base`** — `helpers.GetAuthUser`, `helpers.GetAuthSession`, `helpers.ExtendSession`, `helpers.BlogPostBlocksToString`, `helpers.GenerateCartCacheKey`, `helpers.TimezoneFromRequest`, `helpers.UserSettings`, and the stripe helper are removed from `internal/helpers/`. Their replacements live in `github.com/dracory/base/session`, `github.com/dracory/base/blogblocks`, `github.com/dracory/base/payment`, and `github.com/dracory/base/tz`.

**Internal test utilities migrated to `github.com/dracory/test` and `github.com/dracory/base/testutils`** — `testutils.NewRequest`/`NewRequestOptions` and `testutils.FlashMessageFind*` are removed from `internal/testutils/`. Use `github.com/dracory/test.NewRequest` and the flash helpers from `github.com/dracory/base/testutils` instead.

**Internal layout types migrated to `github.com/dracory/base/layouts`** — `layouts.Options`, `layouts.Breadcrumb`, `layouts.Breadcrumbs()`, and `layouts.LayoutInterface` are removed from `internal/layouts/`. Use `baselayouts.Options`, `baselayouts.Breadcrumb`, `baselayouts.Breadcrumbs()`, and `baselayouts.LayoutInterface` from `github.com/dracory/base/layouts` instead.

**Context key types aligned with `base/session`** — `config.AuthenticatedUserContextKey`, `config.AuthenticatedSessionContextKey`, `config.APIAuthenticatedUserContextKey`, and `config.APIAuthenticatedSessionContextKey` are now type aliases of the corresponding `github.com/dracory/base/session` types. This fixes a latent bug where values stored under one type could not be retrieved under the other.

**Key Changes:**
- `pkg/useradmin`, `pkg/social`, `pkg/blogai` deleted — migrated to external modules (or removed)
- New external deps: `github.com/dracory/useradmin` v0.2.0, `github.com/dracory/social` v0.1.0
- `useradmin.AdminOptions.Registry` removed; replaced with explicit callback adapters
- `useradmin.Handle(w, r) string` → `Handle(w, r)` (no return value)
- `useradmin` callbacks simplified: `OnUserUpdated(ctx, userID)`, `OnUserSearch(ctx, []SearchCondition)`, `OnUserImpersonate(w, r, userID) error`
- `useradmin`: `VaultTokenizer` → `UserPiiSeal`/`UserPiiUnseal`/`UsersPiiUnseal`; `SessionStore` → `OnUserImpersonate`; `BlindIndexStore` → `OnUserSearch`; `GeoStore` → `GeoResolver`
- `useradmin`: `SecureCookie` field removed — Secure flag derived from `r.TLS != nil`
- `helpers.GetAuthUser` → `basesession.GetAuthUser` (`github.com/dracory/base/session`)
- `helpers.BlogPostBlocksToString` → `blogblocks.BlocksToString` (`github.com/dracory/base/blogblocks`)
- `helpers.GenerateCartCacheKey` → `basepayment.GenerateCartCacheKey` (`github.com/dracory/base/payment`)
- `helpers.TimezoneFromRequest` → `basetz.FromRequest` / `basetz.FromUser` (`github.com/dracory/base/tz`)
- `helpers.ExtendSession`, `helpers.GetAuthSession`, `helpers.UserSettings`, `helpers/stripe.go` removed
- `testutils.NewRequest` / `NewRequestOptions` → `test.NewRequest` / `test.NewRequestOptions` (`github.com/dracory/test`)
- `testutils.FlashMessageFind*` → use `github.com/dracory/base/testutils`
- `layouts.Options` / `layouts.Breadcrumb` / `layouts.Breadcrumbs()` / `layouts.LayoutInterface` → `baselayouts.*` (`github.com/dracory/base/layouts`)
- `config.*ContextKey` types are now aliases of `basesession.*ContextKey`
- `database_config.go`: `Driver` field type changed to `neatcontracts.Driver`
- `go.mod`: `base` v0.39.0 → v0.42.3, `neat` v0.39.0 → v0.41.0, `blogstore` v1.35.0 → v1.35.1, `sqlite` v1.56.0 → v1.57.0
- `go.mod`: `stripe-go/v73` removed (replaced by indirect `stripe-go/v86`); `govalidator`, `smithy-go`, `gjson` moved to indirect
- `AGENTS.md` already documents the external admin packages and adapters layer

---

## ⚠️ Breaking Changes

---

### 1. In-Repo `pkg/useradmin` Removed — Migrated to External Module

**Change**: The in-repo user admin package `pkg/useradmin` has been deleted and published as a standalone external module:

| Old in-repo path | New external module |
|---|---|
| `project/pkg/useradmin` | `github.com/dracory/useradmin` v0.2.0 |

All sub-packages keep the same names, so only the import **root** changes.

**Old Usage**:
```go
import (
	useradmin "project/pkg/useradmin"
)
```

**New Usage**:
```go
import (
	useradmin "github.com/dracory/useradmin"
)
```

**Action Required**:
- Find all imports of the old in-repo package:
  ```bash
  grep -rln "project/pkg/useradmin" --include="*.go" .
  ```
- Replace `project/pkg/useradmin` with `github.com/dracory/useradmin`.
- Delete the old `pkg/useradmin` directory from your project if you had not already done so. **It is no longer part of the blueprint repository.**
- Add the new external dependency:
  ```bash
  go get github.com/dracory/useradmin@v0.2.0
  go mod tidy
  ```
- If you had **custom modifications** in your `pkg/useradmin` copy, port them to the external module via the new callback/interface extension points (see Breaking Change #3) or maintain a fork.

**Migration Command**:
```bash
grep -rl "project/pkg/useradmin" --include="*.go" . | xargs sed -i 's|project/pkg/useradmin|github.com/dracory/useradmin|g'
```

---

### 2. `useradmin.AdminOptions.Registry` Removed — Explicit Callback Adapters

**Change**: `useradmin` no longer accepts a `Registry app.AppInterface` field on `AdminOptions`. The package now takes only the stores and callback adapters it actually needs. This decouples the external package from any host-project type.

**Old Usage**:
```go
admin, err := useradmin.New(useradmin.AdminOptions{
	Registry:     controller.app,
	AdminHomeURL: links.Admin().Home(),
	UserAdminURL: links.Admin().Users(),
	AuthUserID: func(r *http.Request) string {
		user := helpers.GetAuthUser(r)
		if user == nil {
			return ""
		}
		return user.GetID()
	},
})
```

**New Usage**:
```go
admin, err := useradmin.New(useradmin.AdminOptions{
	UserStore:         controller.app.GetUserStore(),
	GeoResolver:       adapters.NewGeoResolver(controller.app.GetGeoStore()),
	Logger:            controller.app.GetLogger(),
	OnUserImpersonate: adapters.NewOnUserImpersonateFunc(controller.app),
	OnUserSearch:      adapters.NewOnUserSearchFunc(controller.app),
	OnUserUpdated:     adapters.NewOnUserUpdatedFunc(controller.app, constants.BlindIndexRebuildTaskAlias),
	UserPiiSeal:       adapters.NewUserPiiSealFunc(controller.app),
	UserPiiUnseal:     adapters.NewUserPiiUnsealFunc(controller.app),
	UsersPiiUnseal:    adapters.NewUsersPiiUnsealFunc(controller.app),
	FuncLayout:        adapters.NewLayoutFunc(controller.app),
	FlashRedirect:     adapters.NewFlashRedirectFunc(controller.app),
	AdminHomeURL:      links.Admin().Home(),
	UserAdminURL:      links.Admin().Users(),
	UserHomeURL:       links.User().Home(),
})
```

The new `AdminOptions` fields:

| Field | Type | Purpose |
|---|---|---|
| `UserStore` | `userstore.StoreInterface` | Required — user storage |
| `Logger` | `*slog.Logger` | Required — logging |
| `GeoResolver` | `useradmin.GeoResolverInterface` | Country/timezone lookups (use `adapters.NewGeoResolver`) |
| `OnUserImpersonate` | `useradmin.OnUserImpersonateFunc` | Creates session + sets auth cookie on impersonation |
| `OnUserSearch` | `useradmin.OnUserSearchFunc` | Blind-index search with `[]SearchCondition` (OR/AND combine) |
| `OnUserUpdated` | `useradmin.OnUserUpdatedFunc` | Called after a user is updated (e.g. enqueue blind-index rebuild) |
| `UserPiiSeal` | `useradmin.UserPiiSealFunc` | Tokenize/encrypt user fields before storage |
| `UserPiiUnseal` | `useradmin.UserPiiUnsealFunc` | Detokenize/decrypt a single user for display |
| `UsersPiiUnseal` | `useradmin.UsersPiiUnsealFunc` | Detokenize/decrypt a batch of users |
| `FuncLayout` | `func(...)` | Admin layout callback (use `adapters.NewLayoutFunc`) |
| `FlashRedirect` | `useradmin.FlashRedirectFunc` | Flash message + redirect helper (use `adapters.NewFlashRedirectFunc`) |
| `UserHomeURL` | `string` | URL to the user home page (new) |

**Removed fields**: `Registry`, `AuthUserID`, `AuthUser`, `SecureCookie`, `VaultTokenizer`, `SessionStore`, `BlindIndexStore`, `GeoStore`.

**Action Required**:
- Remove the `Registry:` field from every `AdminOptions{...}` literal.
- Remove the `AuthUserID:` / `AuthUser:` callback fields — the controller now resolves the auth user itself via `basesession.GetAuthUser(r)` before constructing the admin instance.
- Remove the `SecureCookie:` field — the Secure flag is now derived from `r.TLS != nil` inside the `OnUserImpersonate` adapter (see Breaking Change #6).
- Add the new explicit dependency fields as shown above. Use the `adapters` package (see Breaking Change #3) to supply the callbacks.
- `UserHomeURL` is a new required field — pass `links.User().Home()`.

**Migration Command**:
```bash
grep -rn "Registry:" --include="*.go" internal/controllers/admin/users/
```

---

### 3. New `adapters` Constructors for `useradmin`

**Change**: The `internal/controllers/admin/adapters/` package (introduced in v0.39.0 for `blogadmin`/`shopadmin`) is extended with new constructors that bridge the external `useradmin` package to the blueprint's internal services. The package doc comment is updated to mention `useradmin`.

New constructors:

- `NewGeoResolver(geoStore)` — returns a `*GeoResolver` implementing `useradmin.GeoResolverInterface`, backed by the blueprint's `geostore`. Adapts the geostore query-option API to the simpler `Countries(ctx)` / `Timezones(ctx, countryCode...)` shape.
- `NewOnUserSearchFunc(app)` — returns an `OnUserSearchFunc` that maps each `useradmin.SearchCondition` to the corresponding blind index store (`FirstName`/`LastName`/`Email`) and combines results: AND conditions intersect, OR conditions union. Includes `unionIDSets` and `intersectIDSets` helpers.
- `NewOnUserImpersonateFunc(app)` — returns an `OnUserImpersonateFunc` that creates a session in the blueprint's `sessionstore` and sets the auth cookie. The cookie Secure flag is set from `httpReq.TLS != nil`.
- `NewOnUserUpdatedFunc(app, taskAlias)` — returns an `OnUserUpdatedFunc` that enqueues a blind index rebuild task after a user is updated.
- `NewUserPiiSealFunc(app)` — returns a `UserPiiSealFunc` that tokenizes/encrypts user fields via `ext.UserTokenize`. When vault is disabled, the user is returned unchanged.
- `NewUserPiiUnsealFunc(app)` — returns a `UserPiiUnsealFunc` that detokenizes/decrypts a single user via `ext.UserUntokenize`. When vault is disabled, the user is returned unchanged.
- `NewUsersPiiUnsealFunc(app)` — returns a `UsersPiiUnsealFunc` that loops over the single-user unseal for a batch.
- `NewFlashRedirectFunc(app)` — returns a `FlashRedirectFunc` that bridges to `helpers.ToFlashError/Success/Info/Warning`.

**New Usage**:
```go
import (
	"project/internal/controllers/admin/adapters"
	"project/internal/tasks/constants"
)

admin, err := useradmin.New(useradmin.AdminOptions{
	UserStore:         controller.app.GetUserStore(),
	GeoResolver:       adapters.NewGeoResolver(controller.app.GetGeoStore()),
	Logger:            controller.app.GetLogger(),
	OnUserImpersonate: adapters.NewOnUserImpersonateFunc(controller.app),
	OnUserSearch:      adapters.NewOnUserSearchFunc(controller.app),
	OnUserUpdated:     adapters.NewOnUserUpdatedFunc(controller.app, constants.BlindIndexRebuildTaskAlias),
	UserPiiSeal:       adapters.NewUserPiiSealFunc(controller.app),
	UserPiiUnseal:     adapters.NewUserPiiUnsealFunc(controller.app),
	UsersPiiUnseal:    adapters.NewUsersPiiUnsealFunc(controller.app),
	FuncLayout:        adapters.NewLayoutFunc(controller.app),
	FlashRedirect:     adapters.NewFlashRedirectFunc(controller.app),
	AdminHomeURL:      links.Admin().Home(),
	UserAdminURL:      links.Admin().Users(),
	UserHomeURL:       links.User().Home(),
})
```

**Action Required**:
- Update your `internal/controllers/admin/adapters/adapters.go` to include the new constructors (copy them from the blueprint).
- Pass the new callback adapters to `useradmin.AdminOptions` as shown above.
- `NewOnUserUpdatedFunc` takes a task alias (e.g. `constants.BlindIndexRebuildTaskAlias`) — make sure your project has the corresponding task definition registered.

---

### 4. `useradmin.Handle()` No Longer Returns a String

**Change**: `AdminInterface.Handle(w http.ResponseWriter, r *http.Request)` previously returned an HTML `string`. It now writes directly to the `http.ResponseWriter` and returns nothing. The caller no longer sets the `Content-Type` header or writes the returned string.

**Old Usage**:
```go
html := admin.Handle(w, r)

if html != "" {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil { // #nosec G705 -- html is built by trusted admin handler
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At usersAdminController > Handler", "write_error", err.Error())
		}
	}
}
```

**New Usage**:
```go
admin.Handle(w, r)
```

**Action Required**:
- Find all call sites that capture the return value of `Handle()`:
  ```bash
  grep -rn ":= .*\.Handle(w, r)" --include="*.go" internal/controllers/admin/users/
  grep -rn "html := .*Handle(" --include="*.go" internal/controllers/admin/users/
  ```
- Replace `html := admin.Handle(w, r)` with `admin.Handle(w, r)`.
- Remove the subsequent `w.Header().Set("Content-Type", ...)` and `w.Write([]byte(html))` block.

---

### 5. `useradmin` Callback Signatures Simplified

**Change**: Several `useradmin` callbacks are simplified to pass `userID` directly instead of event structs, and the search callback now takes a `[]SearchCondition` slice with OR/AND combine logic.

| Old signature | New signature |
|---|---|
| `OnUserUpdate(ctx, event OnUserUpdateEvent)` | `OnUserUpdated(ctx, userID string)` |
| `OnUserSearch(ctx, event OnUserSearchEvent)` | `OnUserSearch(ctx, conditions []SearchCondition) ([]string, error)` |
| `OnUserImpersonate(ctx, event)` / `SessionStore` field | `OnUserImpersonate(w, r, userID) error` |

`SearchCondition` fields:
- `Field` — one of `useradmin.SearchFieldFirstName`, `SearchFieldLastName`, `SearchFieldEmail`
- `Value` — the search term
- `Op` — `useradmin.SearchOpEquals` or `useradmin.SearchOpContains`
- `CombineWith` — `useradmin.SearchOr` or `useradmin.SearchAnd` (how to combine with the previous condition)

**Action Required**:
- If you implemented custom `useradmin` callbacks, update their signatures to match the new types.
- `OnUserUpdate` is renamed to `OnUserUpdated` — update the field name in `AdminOptions`.
- For `OnUserSearch`: replace any single-term search logic with the new `[]SearchCondition` loop. AND conditions intersect results, OR conditions union them. See `adapters.NewOnUserSearchFunc` for a reference implementation.

---

### 6. `useradmin.SecureCookie` Removed — Secure Flag from Request TLS

**Change**: The `SecureCookie bool` field on `useradmin.AdminOptions` is removed. The auth cookie Secure flag is now determined from the request TLS state (`httpReq.TLS != nil`) inside the `OnUserImpersonate` adapter, rather than from environment config.

**Old Usage**:
```go
secure := true
if controller.app.GetConfig() != nil && controller.app.GetConfig().IsEnvDevelopment() {
	secure = false
}

admin, err := useradmin.New(useradmin.AdminOptions{
	// ...
	SecureCookie: secure,
	// ...
})
```

**New Usage**:
```go
// No SecureCookie field — the OnUserImpersonate adapter sets it from r.TLS
admin, err := useradmin.New(useradmin.AdminOptions{
	// ...
	OnUserImpersonate: adapters.NewOnUserImpersonateFunc(controller.app),
	// ...
})
```

**Action Required**:
- Remove the `SecureCookie:` field from every `useradmin.AdminOptions{...}` literal.
- Remove the `secure := true; if ... IsEnvDevelopment() { secure = false }` block.
- Ensure you pass `adapters.NewOnUserImpersonateFunc(controller.app)` — the adapter calls `auth.AuthCookieSet(w, r, key, types.WithSecure(httpReq.TLS != nil))`.

---

### 7. `helpers.GetAuthUser` → `basesession.GetAuthUser`

**Change**: `internal/helpers/get_auth_user.go` is deleted. The function is replaced by `basesession.GetAuthUser` from `github.com/dracory/base/session`. This applies across all controllers, middlewares, layouts, and widgets.

**Old Usage**:
```go
import (
	"project/internal/helpers"
)

authUser := helpers.GetAuthUser(r)
```

**New Usage**:
```go
import (
	basesession "github.com/dracory/base/session"
)

authUser := basesession.GetAuthUser(r)
```

**Action Required**:
- Find all usages:
  ```bash
  grep -rn "helpers\.GetAuthUser" --include="*.go" .
  ```
- Replace `helpers.GetAuthUser(r)` with `basesession.GetAuthUser(r)`.
- Add the import `basesession "github.com/dracory/base/session"`.
- Remove the now-unused `"project/internal/helpers"` import where applicable.

**Migration Command**:
```bash
grep -rl "helpers\.GetAuthUser" --include="*.go" . | xargs sed -i 's|helpers\.GetAuthUser|basesession.GetAuthUser|g'
```

---

### 8. Other Removed `internal/helpers` Functions

**Change**: Several more helper functions are removed from `internal/helpers/` and replaced by equivalents in `github.com/dracory/base`:

| Old helper | New replacement | Import |
|---|---|---|
| `helpers.GetAuthSession` | `basesession.GetAuthSession` | `github.com/dracory/base/session` |
| `helpers.ExtendSession` | removed (use `basesession` + `sessionstore` directly) | — |
| `helpers.BlogPostBlocksToString` | `blogblocks.BlocksToString` | `github.com/dracory/base/blogblocks` |
| `helpers.GenerateCartCacheKey` | `basepayment.GenerateCartCacheKey` | `github.com/dracory/base/payment` |
| `helpers.TimezoneFromRequest` | `basetz.FromRequest` / `basetz.FromUser` | `github.com/dracory/base/tz` |
| `helpers.UserSettings` | removed | — |
| `helpers/stripe.go` | removed (stripe-go moved to indirect dep) | — |

**Old Usage** (examples):
```go
helpers.BlogPostBlocksToString(content)
helpers.GenerateCartCacheKey(r)
helpers.TimezoneFromRequest(r)
```

**New Usage**:
```go
blogblocks.BlocksToString(content)
basepayment.GenerateCartCacheKey(r)
basetz.FromRequest(r)
// or basetz.FromUser(user)
```

**Action Required**:
- Find all usages and replace as above:
  ```bash
  grep -rn "helpers\.BlogPostBlocksToString\|helpers\.GenerateCartCacheKey\|helpers\.TimezoneFromRequest\|helpers\.GetAuthSession\|helpers\.ExtendSession\|helpers\.UserSettings" --include="*.go" .
  ```
- Add the corresponding `github.com/dracory/base/*` imports.
- If you used `helpers.ExtendSession` or `helpers.UserSettings`, re-implement the logic inline or in your own helper using `basesession` + `sessionstore` / `userstore` directly. The blueprint no longer ships these.
- The stripe helper is removed entirely. If you had custom Stripe integration code depending on it, use `github.com/stripe/stripe-go/v86` directly (it is now an indirect dependency).

---

### 9. `internal/testutils` Helpers Removed — Use `github.com/dracory/test`

**Change**: `testutils.NewRequest` / `testutils.NewRequestOptions` and `testutils.FlashMessageFind*` are removed from `internal/testutils/`. Their replacements live in `github.com/dracory/test` and `github.com/dracory/base/testutils`.

**Old Usage**:
```go
import (
	"project/internal/testutils"
)

r, _ := testutils.NewRequest("GET", "/", testutils.NewRequestOptions{})
```

**New Usage**:
```go
import (
	"github.com/dracory/test"
)

r, _ := test.NewRequest("GET", "/", test.NewRequestOptions{})
```

**Action Required**:
- Find all usages:
  ```bash
  grep -rn "testutils\.NewRequest\|testutils\.NewRequestOptions\|testutils\.FlashMessageFind" --include="*.go" .
  ```
- Replace `testutils.NewRequest` with `test.NewRequest` and `testutils.NewRequestOptions` with `test.NewRequestOptions`.
- For flash message helpers, use the equivalents from `github.com/dracory/base/testutils`.
- Add `github.com/dracory/test` to your `go.mod`:
  ```bash
  go get github.com/dracory/test
  ```

**Migration Command**:
```bash
grep -rl "testutils\.NewRequest" --include="*.go" . | xargs sed -i 's|testutils\.NewRequest|test.NewRequest|g'
grep -rl "testutils\.NewRequestOptions" --include="*.go" . | xargs sed -i 's|testutils\.NewRequestOptions|test.NewRequestOptions|g'
```

---

### 10. `internal/layouts` Types Removed — Use `github.com/dracory/base/layouts`

**Change**: The `Options`, `Breadcrumb`, `Breadcrumbs()`, and `LayoutInterface` types are removed from `internal/layouts/`. They are replaced by equivalents in `github.com/dracory/base/layouts`. The layout constructor functions (`NewAdminLayout`, `NewUserLayout`, `NewBlankLayout`, `NewCmsLayout`, `NewPageLayout`) now accept `baselayouts.Options` instead of `layouts.Options`.

**Old Usage**:
```go
import (
	"project/internal/layouts"
)

layouts.NewUserLayout(app, r, layouts.Options{
	Title:   "Home",
	Content: page,
})

breadcrumbs := []layouts.Breadcrumb{
	{Name: "Dashboard", URL: links.User().Home()},
}
layouts.Breadcrumbs(breadcrumbs)

var l layouts.LayoutInterface
```

**New Usage**:
```go
import (
	baselayouts "github.com/dracory/base/layouts"
	"project/internal/layouts" // constructors still live here
)

layouts.NewUserLayout(app, r, baselayouts.Options{
	Title:   "Home",
	Content: page,
})

breadcrumbs := []baselayouts.Breadcrumb{
	{Name: "Dashboard", URL: links.User().Home()},
}
baselayouts.Breadcrumbs(breadcrumbs)

var l baselayouts.LayoutInterface
```

**Action Required**:
- Find all usages:
  ```bash
  grep -rn "layouts\.Options\|layouts\.Breadcrumb\|layouts\.Breadcrumbs\|layouts\.LayoutInterface" --include="*.go" .
  ```
- Replace `layouts.Options` with `baselayouts.Options`, `layouts.Breadcrumb` with `baselayouts.Breadcrumb`, `layouts.Breadcrumbs(` with `baselayouts.Breadcrumbs(`, and `layouts.LayoutInterface` with `baselayouts.LayoutInterface`.
- Add the import `baselayouts "github.com/dracory/base/layouts"`.
- The layout **constructor functions** (`NewAdminLayout`, `NewUserLayout`, `NewBlankLayout`, `NewCmsLayout`, `NewPageLayout`, `NewUserBreadcrumbsSection`, etc.) still live in `project/internal/layouts` — only the **types** moved.
- `internal/layouts/breadcrumb.go`, `breadcrumbs.go`, `breadcrumbs_test.go`, `options.go`, and `layout_interface.go` are deleted — remove them from your project if present.

**Migration Command**:
```bash
grep -rl "layouts\.Options\b" --include="*.go" . | xargs sed -i 's|layouts\.Options\b|baselayouts.Options|g'
grep -rl "layouts\.Breadcrumb\b" --include="*.go" . | xargs sed -i 's|layouts\.Breadcrumb\b|baselayouts.Breadcrumb|g'
grep -rl "layouts\.Breadcrumbs(" --include="*.go" . | xargs sed -i 's|layouts\.Breadcrumbs(|baselayouts.Breadcrumbs(|g'
grep -rl "layouts\.LayoutInterface" --include="*.go" . | xargs sed -i 's|layouts\.LayoutInterface|baselayouts.LayoutInterface|g'
```

---

### 11. `config.*ContextKey` Types Are Now Aliases of `basesession.*ContextKey`

**Change**: `config.AuthenticatedUserContextKey`, `config.AuthenticatedSessionContextKey`, `config.APIAuthenticatedUserContextKey`, and `config.APIAuthenticatedSessionContextKey` are now type aliases of the corresponding `github.com/dracory/base/session` types. Previously they were independently defined `struct{}` types, which meant a value stored under one type could not be retrieved under the other — a latent bug.

**Old Usage** (still works, but now resolves to the same type):
```go
import (
	"project/internal/config"
)

ctx = context.WithValue(ctx, config.AuthenticatedUserContextKey{}, user)
```

**New Usage** (identical, but now interoperable with `basesession`):
```go
import (
	basesession "github.com/dracory/base/session"
	"project/internal/config"
)

// These two are now the SAME type:
ctx = context.WithValue(ctx, config.AuthenticatedUserContextKey{}, user)
ctx = context.WithValue(ctx, basesession.AuthenticatedUserContextKey{}, user)
```

**Action Required**:
- No code changes are required — existing code using `config.*ContextKey` continues to compile and now correctly interoperates with `basesession.*ContextKey`.
- If you had workarounds that stored/retrieved the same value under both key types, you can simplify them.
- If you defined your own context key types that need to interoperate, alias them to `basesession.*ContextKey` as well.

---

### 12. `pkg/social` Removed — Migrated to External Module

**Change**: The in-repo social sharing package `pkg/social` is deleted and replaced with the external module `github.com/dracory/social` v0.1.0. Only the import root changes.

**Old Usage**:
```go
import (
	"project/pkg/social"
)
```

**New Usage**:
```go
import (
	"github.com/dracory/social"
)
```

**Action Required**:
- Find all imports:
  ```bash
  grep -rln "project/pkg/social" --include="*.go" .
  ```
- Replace `project/pkg/social` with `github.com/dracory/social`.
- Delete the old `pkg/social` directory from your project.
- Add the new external dependency:
  ```bash
  go get github.com/dracory/social@v0.1.0
  go mod tidy
  ```

**Migration Command**:
```bash
grep -rl "project/pkg/social" --include="*.go" . | xargs sed -i 's|project/pkg/social|github.com/dracory/social|g'
```

---

### 13. `pkg/blogai` Removed

**Change**: The in-repo blog AI package `pkg/blogai` is deleted. It is no longer used by the blueprint — AI blog features are handled by `blogadmin`'s `LlmFactory` callback (introduced in v0.39.0).

**Action Required**:
- Find any imports:
  ```bash
  grep -rln "project/pkg/blogai" --include="*.go" .
  ```
- If you used `pkg/blogai` directly, port your AI blog logic to use the `LlmFactory` callback on `blogadmin.AdminOptions` (see the v0.39.0 upgrade guide, Breaking Change #5).
- Delete the old `pkg/blogai` directory from your project.

---

### 14. `database_config.go`: `Driver` Field Type Changed

**Change**: In `internal/config/database_config.go`, the `db.ConnectionConfig.Driver` field type changed from `string` to `neatcontracts.Driver` (from `github.com/dracory/neat/contracts/database`).

**Old Usage**:
```go
nc := db.ConnectionConfig{
	Driver:   driver, // string
	// ...
}
```

**New Usage**:
```go
import (
	neatcontracts "github.com/dracory/neat/contracts/database"
)

nc := db.ConnectionConfig{
	Driver:   neatcontracts.Driver(driver),
	// ...
}
```

**Action Required**:
- If you construct `db.ConnectionConfig` anywhere in your project, cast the driver string to `neatcontracts.Driver`.
- This is a type-safety improvement from the `neat` v0.41.0 bump; no runtime behavior change.

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

```bash
go get github.com/dracory/base@v0.42.3
go get github.com/dracory/neat@v0.41.0
go get github.com/dracory/blogstore@v1.35.1
go get github.com/dracory/useradmin@v0.2.0
go get github.com/dracory/social@v0.1.0
go get github.com/dracory/test
go mod tidy
```

### Step 2: Replace Import Paths

```bash
# useradmin
grep -rl "project/pkg/useradmin" --include="*.go" . | xargs sed -i 's|project/pkg/useradmin|github.com/dracory/useradmin|g'

# social
grep -rl "project/pkg/social" --include="*.go" . | xargs sed -i 's|project/pkg/social|github.com/dracory/social|g'
```

### Step 3: Replace `helpers.GetAuthUser` with `basesession.GetAuthUser`

```bash
grep -rl "helpers\.GetAuthUser" --include="*.go" . | xargs sed -i 's|helpers\.GetAuthUser|basesession.GetAuthUser|g'
```

Then add `basesession "github.com/dracory/base/session"` imports where needed and remove unused `"project/internal/helpers"` imports.

### Step 4: Replace Other Removed Helpers

```bash
grep -rl "helpers\.BlogPostBlocksToString" --include="*.go" . | xargs sed -i 's|helpers\.BlogPostBlocksToString|blogblocks.BlocksToString|g'
grep -rl "helpers\.GenerateCartCacheKey" --include="*.go" . | xargs sed -i 's|helpers\.GenerateCartCacheKey|basepayment.GenerateCartCacheKey|g'
grep -rl "helpers\.TimezoneFromRequest" --include="*.go" . | xargs sed -i 's|helpers\.TimezoneFromRequest|basetz.FromRequest|g'
```

Add imports: `"github.com/dracory/base/blogblocks"`, `"github.com/dracory/base/payment"`, `"github.com/dracory/base/tz"`.

### Step 5: Replace `layouts` Types with `baselayouts` Types

```bash
grep -rl "layouts\.Options\b" --include="*.go" . | xargs sed -i 's|layouts\.Options\b|baselayouts.Options|g'
grep -rl "layouts\.Breadcrumb\b" --include="*.go" . | xargs sed -i 's|layouts\.Breadcrumb\b|baselayouts.Breadcrumb|g'
grep -rl "layouts\.Breadcrumbs(" --include="*.go" . | xargs sed -i 's|layouts\.Breadcrumbs(|baselayouts.Breadcrumbs(|g'
grep -rl "layouts\.LayoutInterface" --include="*.go" . | xargs sed -i 's|layouts\.LayoutInterface|baselayouts.LayoutInterface|g'
```

Add import: `baselayouts "github.com/dracory/base/layouts"`.

### Step 6: Replace `testutils.NewRequest` with `test.NewRequest`

```bash
grep -rl "testutils\.NewRequest" --include="*.go" . | xargs sed -i 's|testutils\.NewRequest|test.NewRequest|g'
grep -rl "testutils\.NewRequestOptions" --include="*.go" . | xargs sed -i 's|testutils\.NewRequestOptions|test.NewRequestOptions|g'
```

Add import: `"github.com/dracory/test"`.

### Step 7: Update `useradmin` Controller

Update `internal/controllers/admin/users/users_controller.go` to:
- Resolve the auth user via `basesession.GetAuthUser(r)` before constructing the admin instance.
- Remove the `Registry`, `AuthUserID`, and `SecureCookie` fields.
- Add the new callback adapter fields (see Breaking Change #2).
- Replace `html := admin.Handle(w, r)` with `admin.Handle(w, r)` and remove the write block.

### Step 8: Delete Removed Files

Delete the following files/directories from your project if they still exist:
- `pkg/useradmin/`
- `pkg/social/`
- `pkg/blogai/`
- `internal/helpers/get_auth_user.go`
- `internal/helpers/get_auth_sesson.go`
- `internal/helpers/extend_session.go`
- `internal/helpers/blog_post_blocks_to_string.go`
- `internal/helpers/cart_util.go`
- `internal/helpers/timezone_from_request.go`
- `internal/helpers/user_settings.go`
- `internal/helpers/stripe.go`
- `internal/testutils/new_request.go`
- `internal/testutils/flash_message.go`
- `internal/layouts/breadcrumb.go`
- `internal/layouts/breadcrumbs.go`
- `internal/layouts/options.go`
- `internal/layouts/layout_interface.go`

### Step 9: Update `database_config.go`

If you construct `db.ConnectionConfig`, cast the driver to `neatcontracts.Driver`:
```go
import neatcontracts "github.com/dracory/neat/contracts/database"

nc := db.ConnectionConfig{
	Driver: neatcontracts.Driver(driver),
	// ...
}
```

### Step 10: Run `go mod tidy` and Build

```bash
go mod tidy
go build ./...
```

---

## 🧪 Testing After Migration

1. **Unit Tests**: Run `go test ./...` and fix any remaining compile errors (mostly import-related).
2. **User Admin**: Test user create/update/delete, user search (with AND/OR conditions), user impersonation (verify the auth cookie is set and Secure flag matches TLS state), and PII seal/unseal (with vault enabled and disabled).
3. **Blog Post Rendering**: Verify block-area blog posts still render correctly via `blogblocks.BlocksToString`.
4. **Cart**: Verify guest cart cache key generation still works via `basepayment.GenerateCartCacheKey`.
5. **Auth**: Verify login/logout, session extension, and that `basesession.GetAuthUser` returns the correct user across all controllers/middlewares/widgets.
6. **Layouts**: Verify admin, user, blank, CMS, and page layouts render correctly with `baselayouts.Options`.
7. **Social Sharing**: Verify blog post social share buttons still render via `github.com/dracory/social`.
8. **Integration Tests**: Run `go test -tags=integration ./...`.

---

## 📝 Additional Notes

- **`base` bumped v0.39.0 → v0.42.3**: This is a large bump. The `base` module now contains `session`, `blogblocks`, `payment`, `tz`, `layouts`, and `testutils` sub-packages that the blueprint previously duplicated internally.
- **`neat` bumped v0.39.0 → v0.41.0**: Introduces the `neatcontracts.Driver` type used in `database_config.go`.
- **`stripe-go` v73 → v86 (indirect)**: The direct `stripe-go/v73` dependency is removed. `stripe-go/v86` is now an indirect dependency. If you have custom Stripe integration, update your imports to `v86`.
- **`govalidator`, `smithy-go`, `gjson` moved to indirect**: These are no longer direct dependencies. If you import them directly, add them to your `go.mod` as direct requires.
- **`useradmin` callback evolution**: The v0.39.0 `useradmin` used `Registry`, `VaultTokenizer`, `SessionStore`, `BlindIndexStore`, `GeoStore`, and `SecureCookie`. The v0.40.0 `useradmin` uses `UserPiiSeal`/`UserPiiUnseal`/`UsersPiiUnseal`, `OnUserImpersonate`, `OnUserSearch`, `GeoResolver`, and derives the Secure flag from TLS. This is a deliberate decoupling: the external package no longer knows about host-project types like `app.AppInterface`, `vaultstore`, `sessionstore`, `blindindexstore`, or `geostore`.
- **Context key aliasing**: The `config.*ContextKey` → `basesession.*ContextKey` aliasing fixes a latent bug. If you previously had code that stored a user under `config.AuthenticatedUserContextKey{}` and retrieved it under `basesession.AuthenticatedUserContextKey{}` (or vice versa), it would have returned `nil`. It now works correctly.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: layouts.Options`
**Cause**: `layouts.Options` was moved to `baselayouts.Options`.
**Solution**: Add `baselayouts "github.com/dracory/base/layouts"` and replace `layouts.Options` with `baselayouts.Options`.

### Issue: `undefined: helpers.GetAuthUser`
**Cause**: `helpers.GetAuthUser` was moved to `basesession.GetAuthUser`.
**Solution**: Add `basesession "github.com/dracory/base/session"` and replace `helpers.GetAuthUser` with `basesession.GetAuthUser`.

### Issue: `undefined: testutils.NewRequest`
**Cause**: `testutils.NewRequest` was moved to `test.NewRequest`.
**Solution**: Add `"github.com/dracory/test"` and replace `testutils.NewRequest` with `test.NewRequest`.

### Issue: `useradmin.AdminOptions` has no field `Registry` / `SecureCookie` / `AuthUserID`
**Cause**: These fields were removed in `useradmin` v0.2.0.
**Solution**: Remove the fields and use the new callback adapters (see Breaking Change #2).

### Issue: `admin.Handle(w, r)` returns too many values / unused variable `html`
**Cause**: `Handle()` no longer returns a string.
**Solution**: Change `html := admin.Handle(w, r)` to `admin.Handle(w, r)` and remove the write block.

### Issue: `cannot use driver (variable of type string) as neatcontracts.Driver`
**Cause**: `db.ConnectionConfig.Driver` type changed to `neatcontracts.Driver`.
**Solution**: Cast with `neatcontracts.Driver(driver)`.

### Issue: `undefined: helpers.BlogPostBlocksToString` / `helpers.GenerateCartCacheKey` / `helpers.TimezoneFromRequest`
**Cause**: These helpers were moved to `github.com/dracory/base/*`.
**Solution**: Use `blogblocks.BlocksToString`, `basepayment.GenerateCartCacheKey`, and `basetz.FromRequest` respectively.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions about the external `useradmin` package: https://github.com/dracory/useradmin
- For questions about the external `social` package: https://github.com/dracory/social
- For questions about the `base` module: https://github.com/dracory/base
