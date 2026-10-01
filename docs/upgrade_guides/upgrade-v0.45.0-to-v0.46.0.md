# Upgrade Guide: v0.45.0 to v0.46.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.45.0 to v0.46.0.

## Overview

This release removes the `github.com/dracory/dashboard` dependency entirely. All layouts are now rendered locally with `github.com/dracory/hb`, following the CourseThread model: a single shared page-layout scaffold with a pluggable navbar. The old `NewUserLayout`/`NewAdminLayout`/`NewAdminCrudLayout` constructors are gone — replaced by a unified public layout (`NewPageLayout`), an admin layout package under the admin controllers, and a local theme system backed by Bootlegance CDN themes.

**Key Changes:**
- `github.com/dracory/dashboard` removed from `go.mod` — zero imports remain
- `internal/layouts` now exports `NewPageLayout`, `NewPageLayoutWithConfig`, `NewBlankLayout`, `NewCmsLayout`, `MenuItem`, `ThemeDropdown`, `ThemeName`, `ThemeStyleURL`, and `UserDisplayNames`
- Admin layouts moved to `internal/controllers/admin/layout` (`adminlayout.New`, `adminlayout.NewCrud`, `adminlayout.Page`)
- Theme switching reimplemented locally: `GET /theme?theme=<name>&redirect=<path>` sets a `theme` cookie and redirects back
- Homepage now renders through `NewPageLayout` — same navbar, theme picker, and footer as the rest of the site
- Logo is now an embedded SVG data URI (no external host, CSP-friendly)
- 404 page rebuilt on `NewBlankLayout` — inherits the active theme
- Navbar controls use mode-adaptive classes (`btn-outline-secondary`, `btn-secondary`) so they stay visible under `data-bs-theme="dark"`

---

## ⚠️ Breaking Changes

---

### 1. `NewUserLayout` removed — use `NewPageLayout`

**Change**: `NewUserLayout` and `NewAdminLayout` in `internal/layouts` were merged into a single auth-aware layout. It renders the public navbar for both guests and authenticated users and adapts menu items automatically.

**Old Usage**:
```go
return layouts.NewUserLayout(app, r, baselayouts.Options{
	Title:   "My Page",
	Content: content,
}).ToHTML()
```

**New Usage**:
```go
return layouts.NewPageLayout(app, r, baselayouts.Options{
	Title:   "My Page",
	Content: content,
}).ToHTML()
```

**Action Required**:
- Find call sites:
  ```bash
  grep -rn "NewUserLayout\|NewAdminLayout\|NewAdminCrudLayout" --include="*.go" .
  ```
- Replace `NewUserLayout` with `NewPageLayout` — the signature is unchanged. This applies to user controllers, `shared/flash`, `website/contact`, and the CMS layout middleware.
- Custom dashboard-era options (logo HTML, theme handler URL, navbar color mode) are gone. Options now come from `baselayouts.Options` only; inject a custom navbar via `NewPageLayoutWithConfig` + `PageLayoutConfig{NavbarFn: myNavbar}`.

---

### 2. Admin layouts moved to `internal/controllers/admin/layout`

**Change**: `NewAdminLayout` / `NewAdminCrudLayout` / `AdminPage` were removed from `internal/layouts` and now live next to the admin controllers.

**Old Usage**:
```go
return layouts.NewAdminLayout(app, r, options).ToHTML()
return layouts.NewAdminCrudLayout(app, r, "Title", content, styleURLs, style, jsURLs, js)
layouts.AdminPage(elements...)
```

**New Usage**:
```go
import adminlayout "project/internal/controllers/admin/layout"

return adminlayout.New(app, r, options).ToHTML()
return adminlayout.NewCrud(app, r, "Title", content, styleURLs, style, jsURLs, js)
adminlayout.Page(elements...)
```

**Action Required**:
- Update imports and call sites across `internal/controllers/admin/**`, including `adapters/adapters.go` (`adapters.NewLayoutFunc` now returns `adminlayout.New`).

---

### 3. `dashboardTypes.MenuItem` replaced by `layouts.MenuItem`

**Change**: Menu files no longer return `dashboardTypes.MenuItem`. The shared type is `internal/layouts.MenuItem` (`Icon`, `Title`, `URL`, `Target`).

**Action Required**:
- If you customized `*_menu_items.go` files, change the return type to `[]layouts.MenuItem` and drop the `dashboard/types` import.
- Admin menu builders moved to `internal/controllers/admin/layout` (unexported `mainMenuItems`/`userMenuItems`) — override navbars there if you forked them.

---

### 4. `dashboard.ThemeHandler` / `dashboard.ThemeMiddleware` removed

**Change**: Theme switching is a plain controller + cookie now. `ThemeMiddleware` no longer exists — the theme is resolved per-request from the `theme` cookie by `layouts.ThemeName(r)` when rendering `<head>`.

**Old Usage**:
```go
SetHandler(dashboard.ThemeHandler(themeMap))
```

**New Usage**:
```go
SetHTMLHandler(theme.NewThemeController(app).Handler)
```

**Action Required**:
- Remove `ThemeMiddleware` from your middleware chain if present; nothing replaces it.
- The `/theme` route accepts `theme` (validated against the known Bootlegance theme list) and `redirect` (path-only, open-redirect guarded).
- Generate switch links with `links.Website().Theme(map[string]string{"theme": name, "redirect": path})`.
- Embed `layouts.ThemeDropdown(app, r)` in custom navbars for the picker.
- Theme CSS: `layouts.ThemeStyleURL(theme)` → `cdn.jsdelivr.net/gh/lesichkovm/bootlegance@v0.5.0/themes/{name}/theme.css`. Default theme: `vanguard`.

---

### 5. `PageNotFoundController` now takes the app instance

**Change**: The 404 controller renders through `NewBlankLayout` and needs `app` for the page title/theme.

**Old Usage**:
```go
SetHTMLHandler(page_not_found.PageNotFoundController().Handler)
```

**New Usage**:
```go
SetHTMLHandler(page_not_found.PageNotFoundController(app).Handler)
```

**Action Required**:
- Pass `app` (may be `nil` where no app is in scope — the layout tolerates it). Call sites in stock Blueprint: `shared/routes.go`, `website/routes.go`, `shared/resource`.

---

### 6. Logo is an embedded data URI

**Change**: `layouts.LogoHTML()` returns an `<img>` with a `data:image/svg+xml;base64` source instead of hotlinking `https://dracory.com/assets/images/logo.png` (which violated `img-src` CSP).

**Action Required**:
- If you customized the logo, replace `internal/layouts/logo.svg` / `logo_html.go` — or override the navbar via `PageLayoutConfig.NavbarFn`.
- No `img-src` CSP additions needed; `data:` is already allowed.

---

### 7. Fixed dark-mode utilities replaced

**Change**: `btn-outline-dark` and `btn-dark` were removed from navbars, dropdowns, and blog buttons in favor of `btn-outline-secondary`/`btn-secondary`, which adapt under `data-bs-theme="dark"`.

**Action Required**:
- If you copied navbar/button markup, audit for `btn-outline-dark`/`btn-dark`:
  ```bash
  grep -rn "btn-outline-dark\|btn-dark" --include="*.go" .
  ```
- Replace with the secondary variants so controls remain visible in dark mode.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.46.0"
```

### Step 2: Migrate layout call sites

- `NewUserLayout` → `layouts.NewPageLayout` (all public/user-aware pages)
- `NewAdminLayout` / `NewAdminCrudLayout` / `AdminPage` → `adminlayout.New` / `adminlayout.NewCrud` / `adminlayout.Page` under `internal/controllers/admin/layout`
- Auth/standalone pages → `layouts.NewBlankLayout` (now loads Bootstrap JS + theme CSS)
- CMS-rendered pages → `layouts.NewCmsLayout` (unchanged)

### Step 3: Replace the theme handler

- Remove `dashboard.ThemeHandler`/`ThemeMiddleware` usage
- Add `internal/controllers/shared/theme/theme_controller.go` and register `links.Website().Theme` route (`/theme`) in `shared/routes.go`
- Add `layouts.ThemeDropdown(app, r)` to any custom navbar

### Step 4: Update menu item types

Change `dashboardTypes.MenuItem` → `layouts.MenuItem` everywhere and remove the `dashboard/types` import.

### Step 5: Clean up dependencies

```bash
go mod tidy
go build ./...
```

`github.com/dracory/dashboard` should disappear from `go.mod`.

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **Homepage**: `http://127.0.0.1:<port>/` — shares the site navbar with theme picker
4. **Theme switch**: use the palette dropdown → page reloads with the Bootlegance theme applied
5. **Dark mode**: moon/sun toggle flips `data-bs-theme`; all buttons/dropdowns must remain visible
6. **404**: visit a nonexistent path — page should inherit the active theme, not show a purple gradient

---

## 📝 Additional Notes

- `internal/layouts` is now the only shared scaffold; `controllers/*/layout` packages may import it, never the reverse (avoids import cycles). `layouts.NavbarFunc` and `PageLayoutConfig` are the extension hooks.
- The homepage lost its bespoke dark-gradient chrome intentionally — plain Bootstrap content (`lead`, `btn`, `row-cols` cards) restyles correctly under every Bootlegance theme.
- Bootlegance themes are loaded from jsdelivr — add `cdn.jsdelivr.net` to `style-src` if your CSP blocks it (stock Blueprint already allows it).
- `img-src` may still list legacy third-party origins from earlier projects — safe to prune to `'self' data:` plus hosts you actually use.

---

## 🆘 Common Issues and Solutions

### Issue: `undefined: layouts.NewUserLayout` / `NewAdminLayout` after upgrade

**Cause**: Old constructors removed.

**Solution**: `NewUserLayout` → `layouts.NewPageLayout`; admin constructors → `adminlayout` package (Breaking Changes #1, #2).

### Issue: Navbar controls invisible in dark mode

**Cause**: Fixed `btn-outline-dark`/`btn-dark` utilities don't adapt to `data-bs-theme="dark"`.

**Solution**: Use `btn-outline-secondary`/`btn-secondary` (Breaking Change #7).

### Issue: Logo blocked by CSP `img-src`

**Cause**: Old logo hotlinked `dracory.com`, which wasn't whitelisted.

**Solution**: v0.46.0 embeds the logo as a data URI — already allowed by `img-src 'self' data:`. If you host a custom logo externally, add its origin to `img-src`.

### Issue: Theme changes but the page doesn't restyle

**Cause**: Page uses bespoke CSS overrides (gradients, custom `.card`/`.btn` rules) that fight the theme.

**Solution**: Replace bespoke chrome with standard Bootstrap classes so the Bootlegance stylesheet controls the look — see the rewritten `home_controller.go` and 404 page for the pattern.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
