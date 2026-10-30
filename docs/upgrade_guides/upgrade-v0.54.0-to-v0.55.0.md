# Upgrade Guide: v0.54.0 to v0.55.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.54.0 to v0.55.0.

## Overview

This release migrates authentication to the `github.com/dracory/auth` package (multi-method login pages, JSON APIs, and middleware served from the package instead of hand-built per-method controllers) and fixes a bug where the user-layout navbar dropdown ("To Admin Dashboard") used the legacy `role` column check (`user.IsAdministrator()` / `user.IsSuperuser()`) while the admin middleware and the page layout's menu items check the v0.45 role tables via `helpers.UserHasAnyActiveRole`.

**Key Changes:**
- `internal/authadapters/` — new package bridging the app's session/cache/user/email stores and layout to the auth package's callback contract.
- `internal/controllers/auth/routes.go` — rewritten around `auth.NewMultiMethodAuth("/auth", methods)`; the package mux serves all login pages/APIs.
- `internal/layouts/auth_layout.go` — new blank layout + theme switcher for the package's built-in pages.
- Legacy per-method auth controllers deleted (`authentication_authknight`, `forgot_password`, `login_authknight`, `login_magiclink`, `login_otp`, `login_password`, `password_reset`, `shared/login_alternatives.go`) — ~5k lines replaced by the package.
- `internal/tasks/email_delivery/` — new generic auth-email delivery task consuming package delivery keys; the three per-flow email tasks (`email_otp`, `email_magic_link`, `email_password_reset`) are deleted.
- `internal/middlewares/remember_me_middleware.go` — now a thin wrapper over the package's `middlewares.RememberMeMiddleware`.
- `internal/layouts/user_layout_user_menu_items.go` — `userLayoutUserMenuItems` now takes `(application app.AppInterface, ctx context.Context, authUser)` and gates the "To Admin Dashboard" item on `helpers.UserHasAnyActiveRole(..., USER_ROLE_ADMINISTRATOR, USER_ROLE_SUPERUSER)` instead of the legacy column checks.
- `internal/layouts/user_layout_navbar.go` — call site updated to pass `app` and `r.Context()`.

---

## ⚠️ Breaking Changes

### 1. Authentication is served by `dracory/auth`; legacy auth routes are gone

**Change**: All login pages and auth APIs are served by the package's multi-method mux mounted at `/auth`. The hand-built per-method controllers and their assets were deleted.

**Old Behavior**: Each login method had its own controller and URL:

| Legacy route | New route |
|---|---|
| `/auth/login` (password) | `/auth/login` (unchanged — primary method) |
| `/auth/otp-login` | `/auth/login-otp` |
| `/auth/magiclink-login` | `/auth/login-magiclink` |
| `/auth/password-login` | `/auth/login-password` |
| `/auth/authknight-login` | `/auth/login-authknight` |
| `/auth/forgot-password` | `/auth/password-restore` |
| `/auth/password-reset` | `/auth/password-reset` (served by package) |

**New Behavior**: `internal/links` constants already point at the new URLs. The package renders its own login/register/restore/reset/logout pages inside the app layout (`authadapters.NewAuthLayoutFunc`).

**Action Required**:
- Update `github.com/dracory/auth` to the version used by this release (v0.41.0 or later — see the note about the `replace` directive in Migration Steps).
- Set `AUTH_LOGIN_METHODS` (comma-separated: `password`, `otp`, `magiclink`, `authknight`) in `.env`. `AUTH_LOGIN_METHOD` (singular) still works but is deprecated and logs a warning.
- If any code or emails referenced the old URLs directly (`/auth/otp-login`, `/auth/forgot-password`, …), update them — `routes_test.go` includes `TestRoutesNoLegacyPaths` which fails if the old paths resurface.
- If your project forked or extended the deleted controllers, port that logic into adapters/hooks (`types.PostAuthHooks`, `LoginAlternatives`, `FuncLayout`) or keep it as an app-owned route mounted before the package mux (this is how `register` and `logout` remain app-owned).

### 2. Auth email tasks replaced by a single generic delivery task

**Change**: `AuthEmailDeliveryTask` (alias `AuthEmailDeliveryTask`) resolves opaque delivery keys produced by the auth package's `FuncEmailEnqueue`. The per-flow tasks (`EmailOTPTask`, `EmailMagicLinkTask`, `EmailPasswordResetTask`) were removed along with `emails/user_email_magic_link.go`.

**Action Required**:
- If your project registers tasks by copying `internal/tasks/register_tasks.go`, swap the three registrations for `email_delivery.AuthEmailDeliveryTask`.
- When the task store is disabled, emails send synchronously via `FuncEmailSend` — no queue needed.
- Queued task payloads contain only an opaque `delivery_key`; credentials never enter the queue. Queued tasks enqueued by a previous version's per-flow tasks have no consumer anymore — drain or delete them.

### 3. Remember-me middleware delegates to the auth package

**Change**: `internal/middlewares/remember_me_middleware.go` is now a thin wrapper around `middlewares.RememberMeMiddleware` (public in the auth package). Blueprint's `SessionStoreAdapter` gained `SessionUpdate` to support grace-window rotation.

**Action Required**:
- If you copied `shared/session_login.go`: `SessionLogin`, `IssueRememberSession`, `RememberMeCheckbox`, and the `rememberMe` parameter were removed — session + remember issuance now happens in the package's `PostAuth` pipeline via `PostAuthHooks.IssueRememberSession`/`IssueAuthSession` (wired in `authadapters/builder.go`).
- `AUTH_REMEMBER_ME_ENABLED` / `AUTH_REMEMBER_ME_DAYS` are unchanged.

### 4. Admin-dashboard menu item visibility now follows role tables

**Change**: "To Admin Dashboard" is hidden unless the user has an **active** `administrator` or `superuser` role assignment in `snv_users_user_role`.

**Old Behavior**: The item appeared whenever the legacy `users.role` column read `administrator` or `superuser`, even if the user had no matching row in `snv_users_user_role` (and would be denied by the admin middleware).

**New Behavior**: Menu visibility and the admin middleware agree — both consult the role tables.

**Action Required**:
- No schema or configuration changes are required.
- If users report a missing admin link after upgrading, verify they have an active row in `snv_users_user_role` pointing at the `administrator` (or `superuser`) role — the menu is now telling the truth, the account likely lacks the assignment.
- If your project forked `user_layout_user_menu_items.go`, apply the same edit rather than copying the file.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.55.0"
```

### Step 2: Update copied layout code

If your project copied Blueprint's layout code, mirror the role-check change:

```go
// Before
if authUser.IsAdministrator() || authUser.IsSuperuser() {
    items = append(items, adminDashboardMenuItem)
}

// After — signature gains app + ctx
if helpers.UserHasAnyActiveRole(ctx, application, authUser,
    userstore.USER_ROLE_ADMINISTRATOR, userstore.USER_ROLE_SUPERUSER) {
    items = append(items, adminDashboardMenuItem)
}
```

Apply to any project-local equivalents of:

- `internal/layouts/user_layout_user_menu_items.go` — `userLayoutUserMenuItems` (signature becomes `(application app.AppInterface, ctx context.Context, authUser userstore.UserInterface)`)
- `internal/layouts/user_layout_navbar.go` — call site: `userLayoutUserMenuItems(app, r.Context(), authUser)`

Also audit for the same pattern elsewhere: search for `IsAdministrator()`, `IsSuperuser()`, `IsManager()`, and `GetRole()` outside `internal/helpers/user_roles.go` and migrate each to `UserHasActiveRole` / `UserHasAnyActiveRole`.

### Step 3: Update the auth dependency

Update `github.com/dracory/auth` to v0.41.0 or later:

```bash
go get github.com/dracory/auth@latest
go mod tidy
```

> **Note**: during development of this release Blueprint pinned a local checkout
> (`replace github.com/dracory/auth => ../../_modules_dracory/auth`) because the
> multi-method fixes it depends on were not yet tagged. Once a release containing
> them is published (v0.42.0+), remove the `replace` line and `go get` the tag.

### Step 4: Port the auth integration (if your project forked auth code)

Projects that stayed close to Blueprint can copy the new files wholesale:

- `internal/authadapters/` (builder, session/tempkey/user/email adapters, layout func)
- `internal/layouts/auth_layout.go`
- `internal/controllers/auth/routes.go`
- `internal/tasks/email_delivery/` + the updated `register_tasks.go` and `constants.go`

Then delete your copies of the legacy controllers (`login_otp`, `login_magiclink`, `login_password`, `login_authknight`, `authentication_authknight`, `forgot_password`, `password_reset`, `shared/login_alternatives.go`) and the old email tasks.

Keep (app-owned): `register`, `register_password`, `logout` controllers, your user/vault lookup in `user_adapter.go`, branding in `applyBranding`.

### Step 5: Configure login methods

```env
AUTH_LOGIN_METHODS=password,otp,magiclink,authknight
```

`AUTH_LOGIN_METHOD` (singular) is deprecated — it still works as a fallback but logs a warning; remove it.

### Step 6: Verify

```bash
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

- `/auth/login` shows the primary method's page with "OR — use another method" buttons for each configured secondary method; `/auth/login-otp`, `/auth/login-magiclink`, `/auth/login-password`, `/auth/login-authknight` serve their methods.
- Log in with a password account: the package sets the auth cookie and the post-auth hook redirects per role.
- Log in with OTP: the code email arrives (via the queue when the task store is enabled, synchronously otherwise); the 8-box code input accepts the letters in the code; verify redirects to the app.
- Check "remember me", restart the browser, and confirm `AUTH_REMEMBER_ME_*` keeps you signed in.
- `/auth/logout` renders the package logout page and clears the session.
- `AUTH_FORGOT_PASSWORD` (`/auth/password-restore`) sends a reset email and `/auth/password-reset` completes it.
- Log in as a user whose **legacy** `role` column says `administrator` but who has **no** `snv_users_user_role` row: the "To Admin Dashboard" item must not appear, matching the admin middleware's denial.
- Log in as a user with an active `administrator` role assignment: the item must appear and `/admin` must be reachable.
- Log out: only "Login" (and public items) show.

---

## 📋 Checklist

- [ ] `internal/config/version.go` updated to `0.55.0`
- [ ] `github.com/dracory/auth` updated to v0.41.0+ (no stray `replace` pointing at a local checkout)
- [ ] `AUTH_LOGIN_METHODS` set in `.env`; `AUTH_LOGIN_METHOD` removed
- [ ] `internal/authadapters/` and `internal/layouts/auth_layout.go` present
- [ ] Legacy auth controllers and the three per-flow email tasks deleted
- [ ] `email_delivery.AuthEmailDeliveryTask` registered
- [ ] No remaining references to `/auth/otp-login`, `/auth/magiclink-login`, `/auth/password-login`, `/auth/authknight-login`, `/auth/forgot-password`
- [ ] `user_layout_user_menu_items.go` uses `UserHasAnyActiveRole` (signature takes `app` + `ctx`)
- [ ] `user_layout_navbar.go` call site updated
- [ ] No remaining `IsAdministrator()` / `IsSuperuser()` / `IsManager()` calls outside `internal/helpers/user_roles.go`
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
