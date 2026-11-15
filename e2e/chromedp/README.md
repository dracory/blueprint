# e2e/chromedp

Go-native alternative to the Playwright suite in `e2e/tests/`, alongside
`e2e/rod/`. Projects using Blueprint that prefer an all-Go toolchain can use
either Go suite instead of adding a Node dependency.

## How it works

- [chromedp](https://github.com/chromedp/chromedp) drives a headless Chrome
  over CDP. Unlike go-rod it does not download a browser — it uses Chrome on
  PATH, `CHROME_BIN`, or Playwright's `chromium_headless_shell` when present
  (Windows, WSL, and Linux paths are probed).
- `TestMain` builds `../browser` to `tmp/e2e-browser-bin`, starts it with
  `-reset` (fresh SQLite DB), waits for `http://127.0.0.1:34757`, then kills
  it after the run.
- One shared headless browser serves the whole run; each test gets a fresh
  tab via `newTab(t)` with a 30s deadline.
- Same seeded credentials as Playwright: `e2e-admin@blueprint.local` (admin)
  and `e2e-user@blueprint.local` (regular user), both `password123`.

## Run

```bash
go test -v ./e2e/chromedp     # or: task e2e:chromedp
```

## Gotchas learned

- Build the sandbox to a binary and run that — killing `go run` only kills
  the wrapper and leaves the child server holding the DB.
- The dracory/auth login button is `button.ButtonLogin` (JS onclick, not a
  submit input).
- The login page waits ~2s on success before redirecting — poll the URL
  (`waitForPath`) instead of asserting immediately.
- An admin-seeded user lands on `/admin`, not `/user`.
