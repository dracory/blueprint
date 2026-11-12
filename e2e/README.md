# e2e

End-to-end browser tests (Playwright) plus the sandboxed app instance they run
against. Everything for e2e lives in this folder.

## Layout

- `browser/` — `e2e-browser` Go command: sandboxed app on
  `http://127.0.0.1:34757` with a disposable SQLite DB (`tmp/e2e-browser.db`).
  Unlike `cmd/ai-browser` it uses the **regular router** (no auto-login), so
  tests exercise real auth flows. Seeded users: `e2e-admin@blueprint.local`
  (admin) and `e2e-user@blueprint.local` (regular user), both
  `password123`. Flags: `-reset` wipes the DB before start; `-with-user`,
  `-admin` mirror ai-browser. Never runs when `APP_ENV=production`.
- `tests/` — Playwright specs.
- `rod/` — alternative go-rod test suite (pure Go, `go test -v ./e2e/rod`,
  headless Chromium). Builds and starts `browser/` itself.
- `playwright.config.ts` — auto-starts `go run ./e2e/browser` via `webServer`.

## Run

```bash
cd e2e
npm install          # once
npx playwright install chrome   # only needed if Chrome isn't already installed
npm test             # starts the app automatically
```

Manual server (for iterating on selectors):

```bash
go run ./e2e/browser -reset   # from repo root
cd e2e && npx playwright test --ui
```
