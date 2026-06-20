# Proposal: AI Browser — Zero-Config Server with Auto-Login

## Status

Implemented

## Context

Blueprint projects require manual authentication to test any authenticated flow locally. Developers must register a user, verify the email, log in, and only then can they test admin pages, user dashboards, or protected routes. This friction slows development and makes automated testing (especially AI agent workflows) impractical.

CourseThread solved this with a dedicated `cmd/devserver/` entry point and `DevAutoLoginMiddleware` that automatically creates a session for a pre-provisioned dev user. More recently, the tinyfunnel project implemented `cmd/ai-browser/` — a more complete solution that includes hardcoded environment setup, vault/blind index support, subscription seeding, and session cookie output for AI agents.

This proposal adopts the `ai-browser` naming and incorporates lessons from both implementations.

## Goal

Add a `cmd/ai-browser/` entry point to Blueprint that starts a fully functional application with zero configuration and an auto-login middleware, enabling frictionless local development and AI agent testing.

## Key Design Decisions

1. **Separate entry point** — `cmd/ai-browser/` runs alongside `cmd/server/`, not replacing it. The AI browser is opt-in.
2. **Dedicated port** — AI browser runs on a different port (default `34756`) to avoid conflicts with the main server.
3. **Hardcoded environment** — All required environment variables are set programmatically so no `.env` file is needed. SQLite is used as the default database for zero-setup.
4. **Pre-provisioned user** — A dev user with administrator role is created on startup if it doesn't exist. The user email and admin role are configurable via CLI flags.
5. **Middleware-based auto-login** — `AiBrowserAutoLoginMiddleware` checks for an existing auth cookie; if none, it creates a session and sets the cookie. The standard `AuthMiddleware` then picks up the session naturally.
6. **AiBrowserRouter** — A separate `AiBrowserRouter()` function wraps the standard `Router()` with the auto-login middleware prepended, avoiding route duplication.
7. **Security guard** — The AI browser refuses to start in production environment.
8. **Non-fatal background processes** — Background process failures are logged and the server continues, since AI browser mode prioritizes HTTP availability over background processing.
9. **Session cookie output** — The session key and cookie string are printed to stdout so AI agents can use them in subsequent HTTP requests.

## Proposed Implementation

### New Files

#### `cmd/ai-browser/main.go`

A slim entry point that:
- Calls `setHardcodedEnv()` to set all required environment variables
- Loads config via `config.NewFromEnv()`
- Overrides port to `34756`
- Initializes app via `app.New(cfg)`
- Runs migrations
- Calls `initializeAiBrowserUser(app, email, isAdmin)` to create/update the dev user (if `-with-user` flag is provided)
- Registers tasks
- Starts background processes (non-fatal on failure)
- Starts web server with `routes.AiBrowserRouter(app)`
- Prints session credentials to stdout
- Handles graceful shutdown

**CLI flags:**

```
-with-user <email>   Auto-seed a test user with the given email
-admin               Assign administrator role to the auto-seeded user
```

#### `cmd/ai-browser/seed_user.go`

User initialization and session creation logic:

- `seedUserAndSession(app, email, isAdmin)` — orchestrates user creation, subscription seeding, and session creation
- `findUserByEmail(ctx, app, email)` — handles both vault-enabled and non-vault user stores
- `createUser(ctx, app, email, isAdmin)` — creates user with password hash, vault tokens, and blind index entries when vault is enabled
- `seedSubscription(ctx, app, userID)` — creates an active free trial subscription if none exists
- `findOrCreateSession(ctx, app, user)` — finds an existing non-expired session or creates a new 24-hour session

#### `cmd/ai-browser/background.go`

Minimal background process stub (same pattern as `cmd/server/`).

#### `cmd/ai-browser/background_processes.go`

Starts the same background processes as the main server (task workers, schedulers). Failures are logged but do not stop the server.

#### `cmd/ai-browser/cli_mode.go`

Same CLI mode detection as `cmd/server/`.

#### `internal/middlewares/ai_browser_auto_login.go`

```go
const aiBrowserUserID = "ai-browser-user"

func AiBrowserAutoLoginMiddleware(app app.AppInterface) rtr.MiddlewareInterface
```

The middleware:
1. Checks for an existing auth cookie via `r.Cookie(auth.CookieName)` — same pattern as `authHandlerSessionKey` in `auth_middleware.go`
2. If a cookie exists, passes through to the next handler (letting `AuthMiddleware` validate it normally)
3. If no cookie exists, looks up `ai-browser-user` in the user store
4. Creates a 24-hour session for the user
5. Sets the auth cookie on the response
6. Proceeds to the next handler — `AuthMiddleware` will pick up the new cookie on this same request via the cookie header

**Important**: The middleware must NOT use `helpers.GetAuthUser(r)` to check authentication state, because `GetAuthUser` reads from the request context which is populated by `AuthMiddleware` — and `AuthMiddleware` runs *after* this middleware in the chain. Checking the cookie directly is the correct approach.

#### `internal/middlewares/ai_browser_auto_login_test.go`

Tests for the auto-login middleware:
- Passes through when a valid auth cookie is present
- Creates a session and sets cookie when no auth cookie is present
- Creates a session when the auth cookie is present but invalid
- Does not create duplicate sessions on subsequent requests
- Handles missing user store / session store gracefully

### Modified Files

#### `internal/routes/router.go`

Add `AiBrowserRouter()` function:

```go
func AiBrowserRouter(app app.AppInterface) rtr.RouterInterface {
    r := Router(app) // reuse existing router — no route duplication
    r.AddBeforeMiddlewares(middlewares.AiBrowserAutoLoginMiddleware(app))
    return r
}
```

This wraps the existing `Router()` rather than duplicating the route setup, ensuring `AiBrowserRouter` always stays in sync with `Router`.

### Hardcoded Environment Setup

The `setHardcodedEnv()` function in `cmd/ai-browser/main.go` sets all required environment variables:

```go
func setHardcodedEnv() {
    // App configuration
    os.Setenv("APP_NAME", "Blueprint")
    os.Setenv("APP_URL", "http://127.0.0.1:34756")
    os.Setenv("APP_HOST", "127.0.0.1")
    os.Setenv("APP_PORT", "34756")
    os.Setenv("APP_ENV", "development")
    os.Setenv("APP_DEBUG", "true")

    // Database (SQLite for zero-setup)
    os.Setenv("DB_DRIVER", "sqlite")
    os.Setenv("DB_DATABASE", "tmp/ai-browser.db")

    // Authentication
    os.Setenv("AUTH_REGISTRATION_ENABLED", "yes")

    // Session
    os.Setenv("SESSION_SECRET", "ai-browser-session-secret-change-me")

    // EnvEnc (disabled)
    os.Setenv("ENVENC_USED", "no")
    os.Setenv("ENVENC_KEY_PRIVATE", "")

    // Vault store
    os.Setenv("VAULT_STORE_KEY", "ai-browser-vault-key-32-chars-long!!")

    // LLM providers (all disabled)
    os.Setenv("ANTHROPIC_API_USED", "no")
    os.Setenv("GEMINI_API_USED", "no")
    os.Setenv("OPENAI_API_USED", "no")
    os.Setenv("OPENROUTER_API_USED", "no")
    os.Setenv("VERTEX_AI_API_USED", "no")

    // Payments (disabled)
    os.Setenv("STRIPE_KEY_PRIVATE", "")
    os.Setenv("STRIPE_KEY_PUBLIC", "")

    // i18n
    os.Setenv("TRANSLATION_LANGUAGE_DEFAULT", "en")
}
```

### User Initialization

The `initializeAiBrowserUser` function in `cmd/ai-browser/seed_user.go`:

- If `-with-user` flag is not provided, uses default email `ai-browser@blueprint.local`
- Checks if the user exists (via blind index search if vault is enabled, or direct email lookup otherwise)
- If not found, creates user with:
  - Email, first name "AI", last name "Browser"
  - Role `ADMINISTRATOR` (or `USER` if `-admin` flag is not set)
  - Status `ACTIVE`
  - Password hash set via `user.SetPasswordAndHash("password123")`
  - Vault tokens for email, first name, last name (when `UserStoreVaultEnabled` is true)
  - Blind index search values for email, first name, last name (when vault is enabled)
- If found, ensures the role matches the `-admin` flag
- Seeds an active free trial subscription if none exists
- Creates or reuses a non-expired session (24-hour expiry)
- Prints credentials to stdout:

```
==================================================
AI Browser Auto-Seed Credentials
==================================================
Email:      ai-browser@blueprint.local
Password:   password123
Role:       administrator
User ID:    <user-id>
Session Key: <session-key>
Cookie:     auth_session=<session-key>
Login URL:  http://127.0.0.1:34756/login
Dashboard:  http://127.0.0.1:34756/user/home
==================================================
```

### Session Expiry Handling

The auto-login middleware creates sessions with a 24-hour expiry. When a session expires:
- The `AuthMiddleware` will not find a valid session (it checks `session.IsExpired()`)
- The auto-login middleware will see no valid cookie context and create a new session
- This is acceptable for dev mode — no manual intervention needed

For longer-running dev sessions, the expiry can be increased to 7 days by configuration.

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `cmd/ai-browser/main.go` | Create | AI browser entry point with hardcoded env and auto-login |
| `cmd/ai-browser/seed_user.go` | Create | User/session/subscription seeding logic |
| `cmd/ai-browser/background.go` | Create | Background group stub |
| `cmd/ai-browser/background_processes.go` | Create | Background process startup (non-fatal) |
| `cmd/ai-browser/cli_mode.go` | Create | CLI mode detection |
| `internal/middlewares/ai_browser_auto_login.go` | Create | Auto-login middleware |
| `internal/middlewares/ai_browser_auto_login_test.go` | Create | Tests for auto-login middleware |
| `internal/routes/router.go` | Modify | Add `AiBrowserRouter()` function |

## Security Considerations

- The AI browser must only run in `local` or `development` environments. Add a startup check that exits if environment is `production` or `testing`.
- The default user ID (`ai-browser-user`) is hardcoded and should never exist in production databases.
- The AI browser port (`34756`) should not be exposed publicly.
- The auto-login middleware must never be registered in the main server's middleware chain.
- The hardcoded session secret and vault key must never be used in production.
- The `setHardcodedEnv()` function must never be called outside of `cmd/ai-browser/`.

## Testing

- Unit test `AiBrowserAutoLoginMiddleware` — verifies it creates a session when no cookie is present and passes through when a cookie exists
- Unit test `initializeAiBrowserUser` — verifies user creation, role enforcement, vault token creation, and blind index entries
- Unit test `seedSubscription` — verifies subscription and plan creation
- Integration test — AI browser starts and pages are accessible without manual login

## Verification

- `go test ./internal/middlewares/...` passes
- `go run ./cmd/ai-browser` starts successfully on port 34756
- `go run ./cmd/ai-browser -with-user test@example.com -admin` creates an admin user and prints credentials
- Visiting any page auto-logs in as the seeded user with the specified role
- Existing `cmd/server` functionality is unaffected
- AI browser refuses to start in production environment
- Session cookie is printed to stdout and usable in subsequent HTTP requests

## Benefits

- **Frictionless local development** — no manual registration/login needed, no `.env` file required
- **AI agent testing** — automated agents can test authenticated flows without credentials; session cookie is printed for immediate use
- **Zero config** — SQLite database, all env vars, and dev user are auto-provisioned on first run
- **Configurable** — `-with-user` and `-admin` flags allow testing different roles and users
- **Non-invasive** — completely separate from the main server; no risk to production
- **Proven pattern** — based on working implementations in CourseThread and tinyfunnel `cmd/ai-browser/`
- **Full feature support** — vault encryption, blind index, subscriptions, and background processes all work out of the box

## References

- **tinyfunnel `cmd/ai-browser/`** — Working implementation with hardcoded env, vault support, subscription seeding, and session cookie output
- **CourseThread `cmd/devserver/`** — Original auto-login middleware concept
- **Blueprint `cmd/server/`** — Main server entry point pattern that this follows
