# e2e

End-to-end browser tests plus the sandboxed app instance they run against.
Everything for e2e lives in this folder. Three interchangeable suites are
provided, in order of our recommendation:

1. `chromedp/` — pure-Go suite via chromedp (no Node; uses Chrome on PATH,
   `CHROME_BIN`, or Playwright's headless shell). **Recommended default.**
2. `playwright/` — TypeScript Playwright specs (requires Node; strongest
   debugging/CI tooling).
3. `rod/` — pure-Go suite via go-rod (no Node; auto-downloads a browser).

All three run against the same sandbox and cover the same flows (login page,
admin/user redirect, /admin authorization, home page).

**This is a blueprint — it ships all three so projects can compare them. A
real project should pick one suite and delete the other two** (they are
independent; deleting a suite folder is all it takes). If you keep Playwright
you can also drop the Go e2e dependencies; if you keep a Go suite you can
drop `e2e/playwright/` and the npm install.

## Choosing a suite

|                    | chromedp         | playwright        | rod            |
|--------------------|------------------|-------------------|----------------|
| Language           | Go               | TypeScript        | Go             |
| Extra runtime      | none             | Node + npm        | none           |
| Browser setup      | Chrome on PATH / `CHROME_BIN` / Playwright's shell — never downloads | `playwright install` | auto-downloads, or reuses Playwright's |
| API style          | `chromedp.Run(ctx, actions...)` task list | locators + auto-waiting assertions | `Must*` chainable (panics on error) |
| Debugging          | basic            | best: `--ui`, traces, videos, inspector | basic |
| Process cleanup    | direct child process, killed on context cancel | handled | `leakless.exe` helper — flagged as trojan by Windows Defender (go-rod#325) |
| Speed              | fast (shared browser, tab per test) | fast | fast (shared browser, tab per test) |
| Extras             | minimal, stable API surface | codegen, multi-browser (Firefox/WebKit), devices | request interception, rich element API |

Rules of thumb:

- **chromedp** — the default. Pick for a pure-Go project; simplest and
  safest, no extra runtime, no helper binaries.
- **Playwright** — pick if your project is fine with a Node toolchain or you
  want the strongest debugging/CI tooling.
- **rod** — pick only if you need its interception/screenshot API or
  auto-downloaded browser.

## Layout

- `browser/` — `e2e-browser` Go command: sandboxed app on
  `http://127.0.0.1:34757` with a disposable SQLite DB (`tmp/e2e-browser.db`).
  Unlike `cmd/ai-browser` it uses the **regular router** (no auto-login), so
  tests exercise real auth flows. Seeded users: `e2e-admin@blueprint.local`
  (admin) and `e2e-user@blueprint.local` (regular user), both
  `password123`. Flags: `-reset` wipes the DB before start; `-with-user`,
  `-admin` mirror ai-browser. Never runs when `APP_ENV=production`.
- `chromedp/` — chromedp suite; builds and starts `browser/` itself.
- `playwright/` — Playwright suite; `playwright.config.ts` auto-starts
  `go run ./e2e/browser` via `webServer`.
- `rod/` — go-rod suite; builds and starts `browser/` itself.

## Run

Go suites (no install needed):

```bash
go test -v ./e2e/chromedp    # or: task e2e:chromedp
go test -v ./e2e/rod         # or: task e2e:rod
```

Playwright:

```bash
cd e2e/playwright
npm install          # once
npx playwright install chromium   # only needed if Chrome isn't already installed
npm test             # starts the app automatically
```

Manual server (for iterating on selectors):

```bash
go run ./e2e/browser -reset   # from repo root, or: task e2e-browser
cd e2e/playwright && npx playwright test --ui
```
