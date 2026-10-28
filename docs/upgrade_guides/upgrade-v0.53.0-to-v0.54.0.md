# Upgrade Guide: v0.53.0 to v0.54.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.53.0 to v0.54.0.

## Overview

This is a maintenance release. The visitor analytics admin section now uses the standalone `github.com/dracory/statsadmin` package instead of the `admin` subpackage bundled inside `statsstore`. This aligns stats with the blog/shop pattern: the admin UI lives in its own module (`statsadmin`), while `statsstore` remains the pure data store. The internal controller drops its hand-rolled `adminLayout` adapter in favor of the shared `adapters.NewLayoutFunc`, and several dracory dependencies were bumped.

Two follow-up changes landed on main after the v0.54.0 release and are reflected in this guide:

- **Per-request country name caching** — the stats controller no longer resolves country names with one `geostore.CountryFindByIso2` call per visitor IP. It now fetches the full country list once per request (`geostore.CountryList`) and serves all lookups from an in-memory iso2→name map. See [Step 3.5](#step-35-adopt-per-request-country-name-caching) below.
- **Dependency bumps** — `geostore` v1.12.0→v1.13.0, `statsadmin` v0.3.0→v0.8.0, `statsstore` v1.32.1→v1.33.3, `entitystore` v1.22.1→v1.23.0, plus assorted third-party updates.

**Key Changes:**
- `internal/controllers/admin/stats/stats_controller.go` — rewritten to wrap `statsadmin.New(statsadmin.AdminOptions{...})`; the local `adminLayout` type and the `statsstore/admin` + `statsstore/admin/shared` imports were removed.
- `internal/controllers/admin/stats/stats_controller_test.go` — `adminLayout` tests removed; controller tests retained.
- Lint hardening — the codebase is now clean under `nilaway` and `errcheck` (`task nilaway`, `task errcheck`). See [Lint Fixes](#-lint-fixes-nilaway--errcheck) below.
- Dependency updates:
  - `statsadmin` added (v0.3.0 at release; current v0.8.0)
  - `entitystore` v1.22.0→v1.23.0
  - `geostore` v1.12.0→v1.13.0
  - `llm` v1.9.0→v1.10.0
  - `rtr` v1.9.0→v1.10.1
  - `settingstore` v1.13.0→v1.14.0
  - `shopadmin` v0.6.0→v0.8.0
  - `statsstore` v1.32.0→v1.33.3
  - Other bumps include `blogstore` v1.38.0→v1.39.0, `neat` v0.57.0→v0.58.0, `logadmin` v0.1.0→v0.2.0
  - Third-party bumps include `aws-sdk-go-v2` modules, `google.golang.org/api`, `google.golang.org/genai`, `google.golang.org/genproto`, `dromara/carbon/v2`, `gabriel-vasile/mimetype`, `modernc.org/sqlite`, and `golang.org/x/net`.

---

## ⚠️ Breaking Changes

### 1. `statsstore/admin` is replaced by `github.com/dracory/statsadmin`

**Change**: The stats admin dashboard moved from `github.com/dracory/statsstore/admin` to the dedicated `github.com/dracory/statsadmin` module. Its entry point follows the same conventions as `blogadmin`/`shopadmin`: `statsadmin.New(statsadmin.AdminOptions{...})` returns an `AdminInterface` (an `http.Handler` with a `Handle` method), configured with `FuncLayout`, `AuthUserID`, `FlashError`, `HomeURL`, `BaseURL`, and `CountryNameByIso2` callbacks.

**Old Usage**:
```go
import (
    statsAdmin "github.com/dracory/statsstore/admin"
    statsAdminShared "github.com/dracory/statsstore/admin/shared"
)

visitorAnalyticsAdmin, err := statsAdmin.New(statsAdmin.Options{
    Logger:            app.GetLogger(),
    Store:             app.GetStatsStore(),
    Layout:            myLayout, // implements statsAdminShared.LayoutInterface
    HomeURL:           links.Admin().Home(),
    WebsiteUrl:        links.Website().Home(),
    CountryNameByIso2: countryNameFunc,
})
```

**New Usage**:
```go
import (
    statsadmin "github.com/dracory/statsadmin"
)

admin, err := statsadmin.New(statsadmin.AdminOptions{
    Store:             app.GetStatsStore(),
    Logger:            app.GetLogger(),
    FuncLayout:        adapters.NewLayoutFunc(app),
    HomeURL:           links.Admin().Home(),
    WebsiteUrl:        links.Website().Home(),
    BaseURL:           links.Admin().Stats(),
    CountryNameByIso2: countryNameFunc,
    AuthUserID:        authUserIDFunc,
    FlashError:        flashErrorFunc,
})
```

**Action Required**:
- If your project copied the old stats admin controller, replace the `statsstore/admin` imports and `statsAdmin.Options` with `statsadmin` and `statsadmin.AdminOptions`.
- Delete any local `LayoutInterface` implementation (setters + `Render`); `FuncLayout` takes `(w, r, title, body, options)` and returns the rendered page — Blueprint's `internal/controllers/admin/adapters.NewLayoutFunc` is reusable across all dracory admin packages.
- `BaseURL` defaults to `/admin/stats`; set it explicitly if your stats admin is mounted elsewhere.
- No schema or data migration is required; only the admin UI wiring changed.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.54.0"
```

### Step 2: Update dependencies

```bash
go get -u ./...
go mod tidy
```

Or pin the dracory changes explicitly:

```bash
go get github.com/dracory/statsadmin@v0.8.0 \
       github.com/dracory/entitystore@v1.23.0 \
       github.com/dracory/geostore@v1.13.0 \
       github.com/dracory/llm@v1.10.0 \
       github.com/dracory/rtr@v1.10.1 \
       github.com/dracory/settingstore@v1.14.0 \
       github.com/dracory/shopadmin@v0.8.0 \
       github.com/dracory/statsstore@v1.33.3
go mod tidy
```

### Step 3: Migrate copied stats admin code

If your project copied `internal/controllers/admin/stats/stats_controller.go`, mirror the new wiring:

```go
admin, err := statsadmin.New(statsadmin.AdminOptions{
    Store:             app.GetStatsStore(),
    Logger:            app.GetLogger(),
    FuncLayout:        adapters.NewLayoutFunc(app),
    HomeURL:           links.Admin().Home(),
    WebsiteUrl:        links.Website().Home(),
    BaseURL:           links.Admin().Stats(),
    CountryNameByIso2: countryNameByIso2,
    AuthUserID: func(r *http.Request) string {
        user := basesession.GetAuthUser(r)
        if user == nil {
            return ""
        }
        return user.GetID()
    },
    FlashError: func(w http.ResponseWriter, r *http.Request, message, redirectURL string, delaySeconds int) string {
        return helpers.ToFlashError(app.GetCacheStore(), w, r, message, redirectURL, delaySeconds)
    },
})
if err != nil {
    // log + flash error, same as before
    return
}
admin.Handle(w, r)
```

Then remove:
- The `github.com/dracory/statsstore/admin` and `.../admin/shared` imports
- Any local `adminLayout` struct implementing `statsAdminShared.LayoutInterface`

Routes stay unchanged — the controller is still mounted at `links.ADMIN_STATS` with a catch-all.

### Step 3.5: Adopt per-request country name caching

The `CountryNameByIso2` callback signature is unchanged
(`func(r *http.Request) func(iso2Code string) (string, error)`), but the
implementation changed. The statsadmin sessions page resolves a country name
once per unique visitor IP, so the old per-code
`geostore.CountryFindByIso2` lookup meant thousands of round trips per page
load. Blueprint now fetches the whole country list once per request and
serves every lookup from an in-memory iso2→name map:

```go
import (
    geostorepkg "github.com/dracory/geostore"
)

func (c *statsController) countryNameByIso2(r *http.Request) func(iso2Code string) (string, error) {
    var names map[string]string
    var loadErr error
    loaded := false

    load := func() error {
        if loaded {
            return loadErr
        }
        loaded = true
        names = map[string]string{}

        geostore := c.app.GetGeoStore()
        if geostore == nil {
            loadErr = errors.New("geostore is nil")
            return loadErr
        }
        countries, err := geostore.CountryList(r.Context(), geostorepkg.CountryQueryOptions{})
        if err != nil {
            loadErr = err
            return loadErr
        }
        for _, country := range countries {
            names[country.IsoCode2()] = country.Name()
        }
        return nil
    }

    return func(iso2Code string) (string, error) {
        if err := load(); err != nil {
            return "N/A", err
        }
        name, ok := names[iso2Code]
        if !ok {
            return "", errors.New("country not found for iso2: " + iso2Code)
        }
        return name, nil
    }
}
```

Key properties:

- **Lazy + memoized**: the `CountryList` query runs at most once per request, on the first lookup; subsequent calls reuse the map. Errors are cached too, so a failed load isn't retried within the same request.
- **Request-scoped**: the map lives in the closure returned per request, so no state leaks across requests and each request sees fresh data.
- **Requires `geostore` v1.13.0+** for `geostore.CountryQueryOptions`; if geostore is disabled, lookups return `("N/A", error)` exactly as before.

If your project copied the old per-code `CountryFindByIso2` implementation, replace it with the pattern above.

### Step 4: Verify

```bash
go mod tidy
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Unit tests**: `go test ./internal/controllers/admin/stats/...`
2. **Admin UI**: log in as an admin, visit `/admin/stats`, and confirm the visitor analytics dashboard renders inside the admin layout (navbar, menus).
3. **Country names**: confirm visitor rows show resolved country names rather than raw ISO2 codes when geostore is enabled. Blueprint's regression tests (`TestCountryNameByIso2_LoadsCountriesOnce`, `TestCountryNameByIso2_CacheScopedToRequest`) verify that 100 lookups trigger exactly one `CountryList` call and that the cache does not leak across requests.

---

## 📝 Additional Notes

- **Why a separate package**: keeping the admin UI out of `statsstore` matches the `blogadmin`/`shopadmin`/`useradmin`/`fileadmin`/`logadmin` convention — stores stay pure data layers, admins are self-contained UI modules.
- **`statsstore/admin` may still exist upstream** for backward compatibility, but Blueprint no longer imports it; new work should target `statsadmin`.
- **`FuncLayout` signature is shared** across dracory admin packages (`blogadmin`, `shopadmin`, `statsadmin`), so one layout adapter serves all of them.

---

## 🧹 Lint Fixes (nilaway & errcheck)

The v0.54.0 tree passes `nilaway` and `errcheck` cleanly. If you copied any of
the touched files, mirror the changes or expect the same diagnostics.

### nilaway

- `internal/app/app_implementation.go` — `GetSessionStore()` no longer has a
  `nil` receiver branch; it returns `r.sessionStore` directly like the other
  getters. The `nil` literal was flagged as a panic source at every call site
  (`remember_me_middleware.go`, `logout_controller.go`,
  `password_reset_controller.go`).
- `internal/controllers/auth/register/register_controller.go` — `PageHandler`
  and `handleSave` now nil-check `basesession.GetAuthUser(r)` immediately
  instead of relying on earlier guards (`GetAuthUser` may return `nil`).
- Tests — `last` recorder variables in throttle tests
  (`login_magiclink`, `login_otp`) and `updated` users from `UserFindByID`
  (`register_controller_test.go`) are nil-checked before dereference.

### errcheck

- `json.NewEncoder(w).Encode(data)` → `_ = json.NewEncoder(w).Encode(data)`
  in all `sendJSONResponse` helpers (`forgot_password`, `login_magiclink`,
  `login_otp`, `login_password`, `password_reset`, `register`,
  `register_password`).
- `fmt.Sscanf(...)` → `_, _ = fmt.Sscanf(...)` in the send/attempt counter
  helpers (`forgot_password`, `login_magiclink`, `login_otp`).
- `defer a.Close()` → `defer func() { _ = a.Close() }()` in
  `database/seeders/seeders_test.go`; same pattern for `rows.Close()` in
  `login_magiclink/login_controller_test.go`.

### Verification

```bash
task nilaway
task errcheck
```

---

## 🆘 Common Issues and Solutions

### Issue: Compile error `undefined: statsAdmin` or `cannot find package "github.com/dracory/statsstore/admin"`

**Cause**: Project-local code still imports `statsstore/admin` after the module reorganized.

**Solution**: Switch the import to `github.com/dracory/statsadmin` and update `statsAdmin.Options` → `statsadmin.AdminOptions`, replacing the `Layout` option with `FuncLayout`.

### Issue: Stats admin page renders without site chrome (bare Bootstrap page)

**Cause**: `AdminOptions.FuncLayout` was left nil.

**Solution**: Pass `FuncLayout: adapters.NewLayoutFunc(app)` (or your own layout function matching the signature).

### Issue: Stats admin 404s or links point to the wrong base path

**Cause**: `AdminOptions.BaseURL` does not match where the route is mounted.

**Solution**: Set `BaseURL` to the mounted path (e.g. `links.Admin().Stats()`); it defaults to `/admin/stats`.

### Issue: Stats sessions page is slow (or hammers geostore)

**Cause**: `CountryNameByIso2` still resolves names with a per-code `CountryFindByIso2` lookup, and the sessions page calls it once per unique visitor IP.

**Solution**: Adopt the per-request cached `CountryList` pattern from [Step 3.5](#step-35-adopt-per-request-country-name-caching). One query serves the whole request; verify with a counting wrapper that N lookups produce exactly one geostore call.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
