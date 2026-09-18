# ai-browser

Sandboxed instance of the app for AI agents to drive via browser automation
(Playwright / Chrome DevTools MCP). Never use in production — it refuses to
start when `APP_ENV=production`.

## Run

```bash
go run ./cmd/ai-browser                            # default seeded admin user
go run ./cmd/ai-browser -with-user test@example.com -admin
```

- URL: `http://127.0.0.1:34756`
- DB: dedicated SQLite at `tmp/ai-browser.db` (disposable — delete to reseed)
- Seeded credentials: `ai-browser@blueprint.local` / `password123`
- Session key + ready-made `authtoken=...` cookie are printed on startup

## How auto-login works

- `routes.AiBrowserRouter` prepends `middlewares.AiBrowserAutoLoginMiddleware`,
  which creates a 24h session + auth cookie for any request that lacks one.
- The middleware looks up a **fixed user ID `ai-browser-user`**
  (`internal/middlewares/ai_browser_auto_login.go`). The seeder must create the
  user with `SetID("ai-browser-user")` or auto-login silently does nothing.
- `middlewares.NewEmailAllowlistMiddleware` allows the seeded user because
  `main.go` sets `AUTH_EMAILS_ALLOWED_ACCESS` to the seeded email before
  config load. Without it, `authConfig()` falls back to a hardcoded
  allowlist (`info@sinevia.com`, `lesichkovm@gmail.com`) that does not
  include the seeded email.

## Gotchas

- `tmp/` is deleted by `task dev` cleanup — if the app fails with
  "unable to open database file", recreate it: `mkdir -p tmp`.
- The auto-login cookie is set with `Secure: true`, which browsers will not
  store over `http://`. It still works because the middleware mints a fresh
  session on every cookie-less request, but the cookie never persists.
- All external APIs (Anthropic, OpenAI, Gemini, OpenRouter, Vertex, Stripe,
  mail) are disabled via hardcoded env vars in `setHardcodedEnv()`.
