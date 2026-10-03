# Proposal: "Remember Me" via Long-Lived Session + Auto-Login Middleware

## Status

Implemented

## Context

All login methods funnel through `shared.SessionLoginUser`
(`internal/controllers/auth/shared/session_login.go`), which creates a
`sessionstore` session expiring in 2h (4h in dev) and sets the auth cookie
via `auth.AuthCookieSet` (default `MaxAge` 2h). Once the browser session or
the 2h window ends, users must log in again — there is no "remember me"
option.

`middlewares.AuthMiddleware` (wrapping `rtrMiddleware.AuthMiddleware` from
`github.com/dracory/rtr` v1.9.0) reads the session key from the auth cookie,
looks it up via `SessionFindByKey`, checks `IsExpired()`, and injects
user+session into the request context. On any failure it proceeds
anonymously — it has no fallback hook, so a persistent mechanism must run as
a separate middleware *before* it.

Notable facts surfaced while designing this:

- `logout` (`internal/controllers/auth/logout/logout_controller.go`) only
  removes the auth cookie — the server-side session stays valid until
  expiry (pre-existing bug; should be fixed as part of this work).
- `sessionstore` (`github.com/dracory/sessionstore` v1.19.1) is not just a
  session table — it also exposes a generic K/V API
  (`Set/Get/Has/Delete/Extend` with TTL, `SessionExpiryGoroutine` for
  cleanup) and `SessionList` supports querying by `UserID`. A long-lived
  session can therefore serve as the remember token itself, with
  `SessionFindByKey` providing built-in expiry validation.
- `internal/middlewares/ai_browser_auto_login.go` is an existing precedent
  for "no session → create one → downstream middleware sees it" inside a
  global before-middleware.

## Goal

Opt-in "Remember me" checkbox on login. When checked, the server creates a
second session with a 30-day `ExpiresAt` and places its session key in a
separate `remember_token` cookie (HttpOnly, Secure, SameSite=Lax). A new
global middleware detects requests with no valid auth session but a valid
remember session, rotates it, creates a fresh 2h session, and injects the
authenticated context — the user stays logged in across browser restarts
without extending normal session lifetimes.

## Key Design Decisions

1. **The remember token IS a session** — instead of a bespoke
   selector/validator record, issue a real `sessionstore` session with
   `ExpiresAt = now + 30d` and store its key in the `remember_token`
   cookie. This gives us for free:
   - expiry validation via `SessionFindByKey` / `IsExpired`
   - `SessionList` by `UserID` → "revoke all sessions" (password reset,
     account deactivation, future "active sessions" UI)
   - `user_id`, `ip`, `user_agent` already captured on the row
   - rotation = create new session + `SessionDelete` the old one
   No new tables, no hashing bookkeeping, no migrations.

   Trade-off vs. selector/validator: the token is stored in plaintext in
   the DB (it's the session key). That matches the existing threat model —
   every live auth session is already a plaintext key in the same table —
   so there is no security regression.

2. **Two sessions per remembered login** — the normal 2h auth session
   *and* the 30-day remember session. The remember session is never used
   directly for authz; it only seeds the creation of short-lived auth
   sessions. Its key goes in `remember_token`, not the auth cookie.

3. **Remember sessions are marked** — `IssueRememberSession` sets the
   session's `session_value` column to `"remember"`
   (`shared.RememberSessionValue`). Two consequences enforced by the
   middleware:
   - It only accepts marked sessions, so a regular auth session key planted
     in the remember cookie cannot mint persistent sessions.
   - A marked key found in the *auth* cookie is stripped from the request
     and expired — `AuthMiddleware` (rtr, unmodified) accepts any valid
     session key, so without this a copied remember token would act as a
     30-day auth credential that never rotates.

4. **New `RememberMeMiddleware` before `AuthMiddleware`** — registered in
   `globalMiddlewares` (`internal/routes/global_middlewares.go`, immediately
   before `middlewares.AuthMiddleware`). Per request:
   - Auth cookie resolves to a valid *unmarked* session → do nothing.
   - Auth cookie holds a marked (remember) session → strip + expire it,
     fall through to the remember flow.
   - Read `remember_token`, `SessionFindByKey`, require the marker.
   - Missing/expired/unmarked → expire the cookie, proceed anonymously.
   - **Valid** → `UserFindByID`, verify active, create a normal auth
     session (shared `CreateSession` helper), set the auth cookie, rotate
     the remember session, and rewrite the request Cookie header so
     AuthMiddleware sees the new session on this same request.

5. **Request-cookie rewrite, not context injection** — cookies set on `w`
   are not visible to AuthMiddleware in the same request, and `r.AddCookie`
   would leave a stale client cookie first in line (`r.Cookie` returns the
   first match). The middleware rebuilds the `Cookie` header with the new
   session key instead — same net effect as the
   `AiBrowserAutoLoginMiddleware` pattern, without duplicating entries.

6. **Grace-window rotation** — on each successful restore the used remember
   session is *shortened* to a 60-second grace window (not deleted) and a
   fresh 30-day session is issued. Hard-deleting would make concurrent
   requests carrying the same cookie (typical on browser restore) fail
   lookup and their expired-cookie responses could clobber the rotated
   cookie. Replays outside the grace window fail — the theft signal.

7. **Login pipeline change** — `SessionLogin` / `SessionLoginUser` gain a
   `rememberMe` parameter (or an options struct — decide at implementation;
   the signature ripples to all callers). When true, after the normal
   session/cookie they create the remember session and set `remember_token`
   with `types.WithMaxAge(30d)` and the existing `WithSecure(false)` dev
   convention.

8. **Checkbox placement per login method** — the "Remember me" checkbox is
   shown at the step where the user proves intent to log in on *this*
   device:
   - **Password** — on the email+password login form.
   - **OTP** — on the code-entry screen (step 2), submitted with
     `otp-verify-ajax`; not on the email step, since the user may abandon
     before proving identity.
   - **Magic link** — on the email-entry screen; the `remember` flag is
     stored server-side inside the `magicLinkCacheValue` JSON entry (next
     to `email`, `ip`, `return`) so it travels to the
     `magiclink-callback` verify step without being part of the emailed
     URL — a forwarded link or link scanner can't strip or set it.
   - **AuthKnight** — no checkbox in v1; the flag is negotiated on the
     external provider's page, so it stays without remember-me.

9. **Logout revokes everything** — the logout controller must delete the
   auth session (pre-existing bug: it currently only removes the cookie),
   read `remember_token` and `SessionDelete` it, and expire both cookies.

10. **Password reset / deactivation revokes all sessions** — since remember
   tokens are sessions, `SessionList` by `UserID` + `SessionSoftDelete` on
   password reset and account deactivation covers both auth and remember
   sessions uniformly — no reverse index needed.

11. **No localStorage** — the token is a 30-day credential; HttpOnly cookie
   storage is non-negotiable (XSS exfiltration). This also avoids a
   client-side bootstrap exchange endpoint, since the app is
   server-rendered and auth is cookie-based.

## Proposed Implementation

### Config

New settings following the `config_implementation.go` env-loading pattern:

- `AuthRememberMeEnabled` (env `AUTH_REMEMBER_ME_ENABLED`, default `false`)
- `AuthRememberMeDays` (env `AUTH_REMEMBER_ME_DAYS`, default `30`)

### Files

| File | Change |
|---|---|
| `internal/middlewares/remember_me_middleware.go` | New: read cookie → `SessionFindByKey` → validate user → create auth session → rotate remember session → inject context |
| `internal/middlewares/remember_me_middleware_test.go` | Tests: valid, expired, unknown/deleted session, inactive user, no cookie, rotation |
| `internal/routes/global_middlewares.go` | Register before `AuthMiddleware` |
| `internal/controllers/auth/shared/session_login.go` | `rememberMe` param; extract session-create + cookie-set into a shared helper reused by the middleware; issue remember session + cookie |
| `internal/controllers/auth/login_password/{login_controller.go,app.html,app.js}` | Checkbox on login form; `remember` field → `SessionLoginUser` |
| `internal/controllers/auth/login_otp/{login_controller.go,app.html,app.js}` | Checkbox on code-entry screen; `remember` field in `otp-verify-ajax` → `SessionLogin` |
| `internal/controllers/auth/login_magiclink/login_controller.go` (+ form assets) | Checkbox on email-entry screen; persist `remember` in `magicLinkCacheValue`; read it in `Handler` → `SessionLogin` |
| `internal/controllers/auth/logout/logout_controller.go` | Delete auth session (bug fix), delete remember session, expire both cookies |
| `internal/config/` | `AuthRememberMeEnabled`, `AuthRememberMeDays` |
| `internal/links/links.go` | Cookie name constant `remember_token` (or in `config/constants.go`) |

### Session creation reuse

Extract the session-create + cookie-set block from `SessionLoginUser`
(lines ~100–130) into a shared helper so both login and the middleware
create sessions identically (same expiry rules, cookie options, IP /
user-agent capture).

## Security Considerations

- Token = 30-day credential → HttpOnly + Secure + SameSite=Lax only.
- Rotation on every use; replayed rotated cookies hit a deleted session.
- Optional v2: mark auto-restored sessions (e.g. session value/meta flag)
  and require re-authentication for sensitive actions (password change,
  payments).
- Optional: bind remember sessions to user-agent substring or IP subnet —
  `sessionstore` already records both; compare loosely on use to detect
  cross-device theft.
- Malformed/missing cookies fail fast before any store lookup; existing
  per-IP rate limiters already apply.

## Testing Plan

- Unit tests for the middleware, mirroring `auth_middleware_test.go`
  patterns (`internal/testutils/seed_session.go` exists for seeding
  sessions — extend or add a remember-session seeder).
- Integration test: login with `remember=true` → delete auth session →
  request with remember cookie → expect authenticated context, new auth
  cookie, rotated remember cookie, old remember session deleted.
- `go test ./...` green, `gofmt`, exported symbols documented per AGENTS.md.

## Open Questions

1. 30-day default, or shorter (e.g. 14d)?
2. Should auto-restored sessions carry a "remembered" marker so sensitive
   routes can demand a fresh login?
