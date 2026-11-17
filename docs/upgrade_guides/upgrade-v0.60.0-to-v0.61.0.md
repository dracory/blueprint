# Upgrade Guide: v0.60.0 to v0.61.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.60.0 to v0.61.0.

## Overview

v0.61.0 is a **maintenance release**. There are no breaking changes and no required code edits. The release refreshes dependencies (including Go 1.27.3 in `go.mod`), demotes unused direct requires to indirect, and ships the self-contained browser e2e scaffolding under `e2e/` (chromedp, Playwright, go-rod suites plus the `e2e/browser` sandbox).

**Key Changes:**

- `go` directive bumped from `1.27.0` to `1.27.3`
- Dependency bumps: `dracory/auth` v0.53.0, `dracory/bs` v0.17.0, `dracory/cmsstore` v1.46.2, `dracory/fileadmin` v0.2.0, `dracory/logadmin` v0.5.0, `dracory/neat` v0.62.2, `dracory/sessionstore` v1.21.0, `dracory/shopadmin` v0.20.0, `dracory/shopstore` v1.30.0, `dracory/statsadmin` v0.12.0, `dracory/taskstore` v1.36.0, plus routine bumps in the indirect tree (`golang.org/x/*`, `google.golang.org/*`, `modernc.org/*`, etc.)
- New direct dependencies for e2e: `github.com/chromedp/chromedp` v0.14.2, `github.com/go-rod/rod` v0.116.2
- `disintegration/imaging`, `dracory/blockeditor`, and `dracory/form` demoted from direct to indirect requires (still present — no action needed)
- New `e2e/` folder: `e2e/browser` sandbox + chromedp, Playwright, and go-rod suites (see the v0.59→v0.60 guide §🆕 for full details — this folder was introduced in that cycle and refined here)
- New Taskfile tasks: `task e2e-browser`, `task e2e`, `task e2e:chromedp`, `task e2e:rod`, `task e2e:install`
- `.vscode/settings.json`: removed `python.terminal.useEnvFile` (cosmetic, optional)

---

## ⚠️ Breaking Changes

None. Upgrading is dependency-only unless you want the e2e scaffolding.

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go get github.com/dracory/auth@v0.53.0
go get github.com/dracory/bs@v0.17.0
go get github.com/dracory/cmsstore@v1.46.2
go get github.com/dracory/fileadmin@v0.2.0
go get github.com/dracory/logadmin@v0.5.0
go get github.com/dracory/neat@v0.62.2
go get github.com/dracory/sessionstore@v1.21.0
go get github.com/dracory/shopadmin@v0.20.0
go get github.com/dracory/shopstore@v1.30.0
go get github.com/dracory/statsadmin@v0.12.0
go get github.com/dracory/taskstore@v1.36.0
go mod tidy
```

Or simply copy this release's `go.mod`/`go.sum` and run `go mod tidy`.

### Step 2 (optional): Adopt the e2e scaffolding

Copy `e2e/` into your project, add the five task blocks under the "Development" section of `taskfile.yml`, and add the Go suite dependencies you keep:

```bash
go get github.com/chromedp/chromedp@v0.14.2   # chromedp suite
go get github.com/go-rod/rod@v0.116.2         # rod suite
cd e2e/playwright && npm install && npx playwright install chromium  # playwright suite
```

Pick **one** suite and delete the other two — chromedp is the recommended default (pure Go, no Node, no auto-downloaded browsers). If you keep only Playwright you can skip the Go deps entirely.

### Step 3 (optional): Bump the Go directive

```bash
go mod edit -go=1.27.3
```

---

## 🧪 Testing After Migration

```bash
go build ./...
go test ./...
```

E2e smoke check (only if you adopted `e2e/`):

```bash
task e2e:chromedp   # or: task e2e / task e2e:rod
```

---

## 📝 Additional Notes

- The e2e suites are fully opt-in; nothing in the main app imports chromedp or rod outside `e2e/`.
- `e2e/browser` serves the regular router (no auto-login) on `http://127.0.0.1:34757` with a disposable SQLite DB (`tmp/e2e-browser.db`) and seeds `e2e-admin@blueprint.local` / `e2e-user@blueprint.local` (both `password123`). It refuses to start when `APP_ENV=production`.
- `AGENTS.md` gained a "Browser e2e" line under Testing — worth copying so agents discover the tasks.
- Docker builds need no change for the Go bump: `Dockerfile` uses `golang:1.27` (major.minor), which always resolves to the latest patch (1.27.3 included).

---

## 🆘 Common Issues and Solutions

### Issue 1: `go.mod` merge conflicts when cherry-picking dependency bumps

**Symptom**: conflicts or stale `go.sum` after copying requires.
**Solution**: copy the whole `go.mod`/`go.sum` from the release and run `go mod tidy` — let the toolchain recompute requires rather than hand-editing.

### Issue 2: `task e2e` fails with "playwright: command not found"

**Symptom**: Playwright task errors immediately.
**Solution**: run `task e2e:install` first (npm install + `playwright install chromium`), or use `task e2e:chromedp` which needs no install.

### Issue 3: go-rod suite flagged by antivirus

**Symptom**: Windows Defender quarantines `leakless.exe` when running `task e2e:rod`.
**Solution**: expected — the suite already runs `Leakless(false)`; the helper is only invoked if you use rod's managed browser download. Prefer the chromedp suite.

---

## Support

- Blueprint repository: <https://github.com/dracory/blueprint>
- e2e suite docs: `e2e/README.md`
