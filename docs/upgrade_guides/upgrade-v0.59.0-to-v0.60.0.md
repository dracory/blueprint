# Upgrade Guide: v0.59.0 to v0.60.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.59.0 to v0.60.0.

## Overview

This release continues the standalone-package extraction: the shared asset-serving controllers (`thumb`, `file`, `media`, `resource`, `cdn`) and the generic CMS widgets moved out of the template into `github.com/dracory/base` (`fileserve`, `widgets` packages) and `github.com/dracory/uncdn`. The packages export **handler factories**, not controllers — the project decides naming and mounting, so all behaviour stays wired in `internal/controllers/shared/routes.go` and `internal/widgets/widgets_registry.go`.

**Key Changes:**

- `internal/controllers/shared/thumb/` (722 + 1,541 LOC) → `fileserve.ThumbHandler(fileserve.ThumbOptions{...})` in `dracory/base`
- `internal/controllers/shared/file/` and `media/` (identical twins) → `fileserve.FileHandler(storage)` — one implementation mounted at both `links.FILES` and `links.MEDIA`
- `internal/controllers/shared/resource/` → `fileserve.ResourceHandler(readResource, notFound)`
- `internal/controllers/shared/cdn/` → `uncdn.CdnHandler()` in `dracory/uncdn`
- `internal/widgets/` generic widgets (authenticated, unauthenticated, visible, print, blockeditor, terms-of-use, contact-form) → `dracory/base/widgets`; `Widget` is now an alias for `basewidgets.Widget` and `WidgetRegistry` delegates with injected `RegistryOptions`
- Dependencies: `dracory/base v0.44.0`, `dracory/uncdn v0.12.0`
- ~2,400 lines of controller code + ~2,100 lines of tests removed from the template
- New `e2e/` folder: sandboxed app instance (`e2e/browser`) plus two interchangeable browser-test suites — Playwright (`e2e/tests`) and go-rod (`e2e/rod`)

---

## ⚠️ Breaking Changes

Only relevant if your project **customized** the moved packages or called their constructors directly. The public HTTP behaviour (URLs, headers, output) is unchanged.

### 1. `internal/controllers/shared/thumb` removed

**Old usage**:
```go
thumb.NewThumbController(app).Handler
```

**New usage**:
```go
fileserve.ThumbHandler(fileserve.ThumbOptions{
    Logger:         app.GetLogger(),
    Cache:          app.GetFileCache(),
    Storage:        app.GetSqlFileStorage(),
    ResourceReader: resources.ToBytes,
})
```

Notes:

- `ThumbOptions.Cache` accepts any `Contains/Fetch/Save` implementation (`cachego.Cache` satisfies it); `Storage` accepts `filesystem.StorageInterface`; `ResourceReader` defaults to `os.ReadFile`; `Logger` defaults to `slog.Default()`.
- The handler no longer depends on `rtr` — it parses `/th/{extension}/{size}/{quality}/{path...}` positionally, so it works with any router.

### 2. `internal/controllers/shared/file` and `media` removed

**Old usage**:
```go
file.NewFileController(app.GetSqlFileStorage()).Handler   // links.FILES
media.NewMediaController(app.GetSqlFileStorage()).Handler // links.MEDIA
```

**New usage** (one handler, mounted twice):
```go
fileserve.FileHandler(app.GetSqlFileStorage())
```

### 3. `internal/controllers/shared/resource` removed

**Old usage**:
```go
resource.NewResourceController().Handler
```

**New usage**:
```go
fileserve.ResourceHandler(resources.Resource, page_not_found.PageNotFoundController(nil).Handler)
```

The reader and the not-found handler are injected; passing `nil` for `notFound` falls back to a plain "404 page not found" message.

### 4. `internal/controllers/shared/cdn` removed

**Old usage**:
```go
cdn.NewCdnController().Handler // void http.HandlerFunc
```

**New usage**:
```go
uncdn.CdnHandler() // func(w http.ResponseWriter, r *http.Request)
```

The implementation moved into the `dracory/uncdn` module alongside the embedded assets it serves.

### 5. Generic widgets moved to `dracory/base/widgets`

**Old usage**:
```go
widgets.NewAuthenticatedWidget(app)
widgets.NewUnauthenticatedWidget(app)
widgets.NewVisibleWidget(app)
widgets.NewPrintWidget(app)
widgets.NewBlockeditorWidget(app)
widgets.NewTermsOfUseWidget(app)
widgets.NewContactFormWidget(app)
```

**New usage**:
```go
basewidgets.NewAuthenticatedWidget()
basewidgets.NewUnauthenticatedWidget()
basewidgets.NewVisibleWidget(func(environment string) bool {
    return app.GetConfig().GetAppEnv() == environment
})
basewidgets.NewPrintWidget()
basewidgets.NewBlockeditorWidget(links.Website().Widget)
basewidgets.NewTermsOfUseWidget()
basewidgets.NewContactFormWidget()
```

- `visible` widget: the `config`/`app` coupling is replaced by an injected `isEnvironment func(name string) bool`.
- `blockeditor` widget: `links.Website().Widget(alias, params)` is injected as `widgetURL`.
- `widgets.Widget` is now `type Widget = basewidgets.Widget`, structurally identical to `cmsstore.ShortcodeInterface` — existing `AddShortcodes` calls keep working.
- `widgets.WidgetRegistry(app)` keeps its signature and delegates to `basewidgets.WidgetRegistry(basewidgets.RegistryOptions{IsEnvironment, WidgetURL})`.
- Feature widgets (`blog_post*`, `testimonials`) remain in `internal/widgets` unchanged.

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go get github.com/dracory/base@v0.44.0
go get github.com/dracory/uncdn@v0.12.0
go mod tidy
```

### Step 2: Delete the moved packages

```bash
rm -rf internal/controllers/shared/thumb
rm -rf internal/controllers/shared/file
rm -rf internal/controllers/shared/media
rm -rf internal/controllers/shared/resource
rm -rf internal/controllers/shared/cdn
rm -f internal/widgets/authenticated_widget.go \
      internal/widgets/unauthenticated_widget.go \
      internal/widgets/visible_widget.go \
      internal/widgets/print_widget.go \
      internal/widgets/blockeditor_widget.go \
      internal/widgets/terms_of_use_widget.go \
      internal/widgets/contact_form_widget.go
```

### Step 3: Update routes

In `internal/controllers/shared/routes.go`, replace the deleted controllers with the handler factories shown above (this version's `routes.go` is the reference).

### Step 4: Update the widget registry

In `internal/widgets/widgets_registry.go`, delegate to `basewidgets.WidgetRegistry` with `RegistryOptions` (see §5). In `internal/widgets/widget.go`, replace the interface with the alias:

```go
type Widget = basewidgets.Widget
```

### Step 5: Fix any direct constructor calls

If custom code called the old `New*Controller`/`New*Widget` constructors, rewrite it per the "New usage" blocks above.

---

## 🧪 Testing After Migration

```bash
go build ./...
go test ./...
```

Manual smoke checks:

- `GET /cdn/swal.js` and `GET /cdn/bs523-swal.css` return gzipped assets
- `GET /th/jpg/100x0/70/<path>` returns a resized image
- `GET /files/<path>` and `GET /media/<path>` serve stored files
- `GET /resources/<path>` serves embedded assets; missing/private paths 404
- A CMS page containing `[x-authenticated]`, `[x-visible environment="..."]`, `[x-print]` renders as before

---

## 🆕 New: Browser e2e testing (`e2e/`)

Everything for end-to-end testing lives in `e2e/` and is **opt-in** — nothing
in the main app depends on it, so upgrading is just copying the folder.

```
e2e/
  browser/            # `e2e-browser` command — sandboxed app for tests
  rod/                # go-rod suite (pure Go)         → go test -v ./e2e/rod
  tests/              # Playwright suite (TypeScript)  → npx playwright test
  playwright.config.ts, package.json, README.md
```

- `go run ./e2e/browser` serves the app on `http://127.0.0.1:34757` with a
  disposable SQLite DB (`tmp/e2e-browser.db`), all external APIs disabled, a
  seeded admin (`e2e-browser@blueprint.local` / `password123`), and
  `AUTH_LOGIN_METHODS=password`. Flags: `-reset` (wipe DB), `-with-user`,
  `-admin`. Refuses to start when `APP_ENV=production`. Unlike
  `cmd/ai-browser` it uses the **regular router** — no auto-login — so tests
  exercise real auth flows.
- **Playwright suite**: `playwright.config.ts` auto-starts the server via
  `webServer` and uses bundled headless Chromium (`channel: 'chromium'`).
- **rod suite**: pure-Go alternative for projects that don't want a Node
  toolchain. `TestMain` builds and starts `e2e/browser` itself; uses
  Playwright's `chromium_headless_shell` if present, else rod's managed
  Chromium (override with `CHROME_BIN`).
- Taskfile tasks: `task e2e-browser`, `task e2e`, `task e2e:rod`,
  `task e2e:install`.

If you adopt it, copy `e2e/` and the four task blocks in `taskfile.yml`; no
other migration is needed.

---

## 📝 Additional Notes

- The extraction used **handler factories** (`XHandler(...)`) rather than exported controller types — the project owns naming and mounting.
- `base` gained no new dependencies except `robertkrimen/otto` (print widget); `uncdn` stayed dependency-free.
- `flash/`, `theme/`, and `page_not_found/` intentionally stay in the project — they render branded HTML pages through project `layouts`/`links`.

---

## 🆘 Common Issues and Solutions

### Issue 1: `undefined: thumb.NewThumbController` (or file/media/resource/cdn)

**Symptom**: compile error referencing the deleted packages.
**Solution**: apply §1–§4 replacements, then `go mod tidy`.

### Issue 2: `cannot use basewidgets.Widget as cmsstore.ShortcodeInterface`

**Symptom**: type error when registering shortcodes.
**Solution**: should not occur — `basewidgets.Widget` is structurally identical (`Alias`, `Description`, `Render`). If it does, ensure the local `Widget` alias is in place and each element is appended individually (slice types are not interchangeable).

### Issue 3: `x-visible` shortcode always hides environment-gated content

**Symptom**: content behind `environment="..."` never renders.
**Solution**: `IsEnvironment` wasn't wired — pass `func(env string) bool { return app.GetConfig().GetAppEnv() == env }` in `RegistryOptions`.

### Issue 4: Blockeditor widget renders but its endpoint 404s

**Symptom**: editor "handle" endpoint unreachable.
**Solution**: `WidgetURL` wasn't wired — pass `links.Website().Widget` in `RegistryOptions`.

---

## Support

- Blueprint repository: <https://github.com/dracory/blueprint>
- `dracory/base` (`fileserve`, `widgets`): <https://github.com/dracory/base>
- `dracory/uncdn` (`CdnHandler`): <https://github.com/dracory/uncdn>
