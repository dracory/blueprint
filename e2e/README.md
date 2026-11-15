# e2e

End-to-end browser tests plus the sandboxed app instance they run against.
Everything for e2e lives in this folder. Three interchangeable suites are
provided — pick whichever fits your project's toolchain:

- `playwright/` — TypeScript Playwright specs (requires Node).
- `rod/` — pure-Go suite via go-rod (no Node; auto-downloads a browser).
- `chromedp/` — pure-Go suite via chromedp (no Node; uses Chrome on PATH,
  `CHROME_BIN`, or Playwright's headless shell).

All three run against the same sandbox and cover the same flows (login page,
admin/user redirect, /admin authorization, home page).

If choosing between the Go suites: prefer `chromedp`. It manages Chrome as a
direct child process — no helper binary, so nothing triggers Windows
Defender (go-rod's `leakless.exe` is flagged as a trojan, go-rod#325, so the
rod suite runs `Leakless(false)` and cleans up manually). Pick `rod` only if
you need its interception/screenshot API or its auto-downloaded Chromium —
chromedp will not fetch a browser for you.

## Layout

- `browser/` — `e2e-browser` Go command: sandboxed app on
  `http://127.0.0.1:34757` with a disposable SQLite DB (`tmp/e2e-browser.db`).
  Unlike `cmd/ai-browser` it uses the **regular router** (no auto-login), so
  tests exercise real auth flows. Seeded users: `e2e-admin@blueprint.local`
  (admin) and `e2e-user@blueprint.local` (regular user), both
  `password123`. Flags: `-reset` wipes the DB before start; `-with-user`,
  `-admin` mirror ai-browser. Never runs when `APP_ENV=production`.
- `playwright/` — Playwright suite; `playwright.config.ts` auto-starts
  `go run ./e2e/browser` via `webServer`.
- `rod/` — go-rod suite; builds and starts `browser/` itself.
- `chromedp/` — chromedp suite; builds and starts `browser/` itself.

## Run

Playwright:

```bash
cd e2e/playwright
npm install          # once
npx playwright install chromium   # only needed if Chrome isn't already installed
npm test             # starts the app automatically
```

Go suites (no install needed):

```bash
go test -v ./e2e/rod         # or: task e2e:rod
go test -v ./e2e/chromedp    # or: task e2e:chromedp
```

Manual server (for iterating on selectors):

```bash
go run ./e2e/browser -reset   # from repo root, or: task e2e-browser
cd e2e/playwright && npx playwright test --ui
```
