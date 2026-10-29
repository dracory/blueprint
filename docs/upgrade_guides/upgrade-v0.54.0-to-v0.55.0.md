# Upgrade Guide: v0.54.0 to v0.55.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.54.0 to v0.55.0.

## Overview

This is a bugfix release. The user-layout navbar dropdown ("To Admin Dashboard") used the legacy `role` column check (`user.IsAdministrator()` / `user.IsSuperuser()`), while the admin middleware and the page layout's menu items check the v0.45 role tables via `helpers.UserHasAnyActiveRole`. When the two role sources diverge — e.g. the migration backfill missed a user, or a role was revoked through the new role tables — the navbar showed the admin link to users the middleware would then reject.

**Key Changes:**
- `internal/layouts/user_layout_user_menu_items.go` — `userLayoutUserMenuItems` now takes `(application app.AppInterface, ctx context.Context, authUser)` and gates the "To Admin Dashboard" item on `helpers.UserHasAnyActiveRole(..., USER_ROLE_ADMINISTRATOR, USER_ROLE_SUPERUSER)` instead of the legacy column checks.
- `internal/layouts/user_layout_navbar.go` — call site updated to pass `app` and `r.Context()`.

---

## ⚠️ Breaking Changes

### 1. Admin-dashboard menu item visibility now follows role tables

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

### Step 3: Verify

```bash
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

- Log in as a user whose **legacy** `role` column says `administrator` but who has **no** `snv_users_user_role` row: the "To Admin Dashboard" item must not appear, matching the admin middleware's denial.
- Log in as a user with an active `administrator` role assignment: the item must appear and `/admin` must be reachable.
- Log out: only "Login" (and public items) show.

---

## 📋 Checklist

- [ ] `internal/config/version.go` updated to `0.55.0`
- [ ] `user_layout_user_menu_items.go` uses `UserHasAnyActiveRole` (signature takes `app` + `ctx`)
- [ ] `user_layout_navbar.go` call site updated
- [ ] No remaining `IsAdministrator()` / `IsSuperuser()` / `IsManager()` calls outside `internal/helpers/user_roles.go`
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
