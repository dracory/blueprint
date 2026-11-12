# e2e/rod

Go-native alternative to the Playwright suite in `e2e/tests/`. Projects using
Blueprint that prefer an all-Go toolchain can use this instead of adding a
Node dependency.

## How it works

- [go-rod](https://go-rod.github.io/) drives a headless Chromium over CDP.
  It reuses Playwright's `chromium_headless_shell` when present, otherwise
  rod downloads a managed Chromium; override with `CHROME_BIN`.
- `TestMain` builds `../browser` to `tmp/e2e-browser-bin`, starts it with
  `-reset` (fresh SQLite DB), waits for `http://127.0.0.1:34757`, then kills
  it after the run.
- Same seeded credentials as Playwright: `e2e-browser@blueprint.local` /
  `password123`.

## Run

```bash
go test -v ./e2e/rod     # or: task e2e:rod
```

## Gotchas learned

- `page.Timeout(d)` returns a *new* page value — assign it once
  (`page = page.Timeout(...)`) or later calls retry forever.
- `MustX` methods panic — avoid `MustClose` on an already-dead page.
- `Leakless(false)` disables rod's `leakless.exe` helper, which Windows
  Defender flags as a trojan (false positive, go-rod#325); cleanup is done
  via `launcher.Kill`.
- The dracory/auth login button is `button.ButtonLogin` (JS onclick, not a
  submit input).
- An admin-seeded user lands on `/admin`, not `/user`.
