# Upgrade Guide: v0.44.0 to v0.45.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.44.0 to v0.45.0.

## Overview

This release migrates user roles from the single `role` column on the user record to the new **userstore roles tables** (`role` + `user_role`). `userstore` now stores roles as separate entities identified by a `handle`, and users are linked to roles via a `user_role` join table.

**Key Changes:**
- `NewUserStore` in `internal/config/store_builders.go` now sets `RolesEnabled: true` with `RoleTableName` (`snv_users_role`) and `UserRoleTableName` (`snv_users_user_role`)
- New migration `2026_09_23_0001_store_user_roles_migrate` creates the role tables, seeds four default roles — `administrator`, `manager`, `superuser`, `user` — and backfills `user_role` assignments from each user's legacy `role` column
- New `internal/helpers/user_roles.go` provides `UserActiveRoleHandles`, `UserHasActiveRole`, `UserHasAnyActiveRole`, and `UserActiveRoleAssign`
- All `user.IsAdministrator()` / `user.IsSuperuser()` / `user.IsManager()` / `user.GetRole()` / `user.SetRole()` call sites were replaced with store-based role checks
- New users created via `shared.SessionLogin` are automatically assigned the `user` role
- Dependency: `userstore` v1.19.1 pseudo-version (roles/groups tables support)

---

## ⚠️ Breaking Changes

---

### 1. Roles moved to a separate table

**Change**: Role membership is no longer read from `user.GetRole()`. `UserInterface.IsAdministrator()`, `IsSuperuser()`, and `IsManager()` still exist in `userstore` but only read the legacy `role` column — they no longer reflect real role membership once you stop writing that column.

**Old Usage**:
```go
if authUser.IsAdministrator() || authUser.IsSuperuser() {
	// admin-only logic
}

user := userstore.NewUser().
	SetRole(userstore.USER_ROLE_ADMINISTRATOR)
```

**New Usage**:
```go
if helpers.UserHasAnyActiveRole(ctx, app, authUser,
	userstore.USER_ROLE_ADMINISTRATOR,
	userstore.USER_ROLE_SUPERUSER) {
	// admin-only logic
}

// Assign a role (idempotent):
err := helpers.UserActiveRoleAssign(ctx, app, userID, userstore.USER_ROLE_ADMINISTRATOR)
```

**Action Required**:
- Replace every `IsAdministrator()` / `IsSuperuser()` / `IsManager()` / `GetRole()` / `SetRole()` call with the helpers. All role checks now need a `context.Context` and the app instance (they query the `user_role` table).
- Find call sites:
  ```bash
  grep -rn "IsAdministrator\|IsSuperuser\|IsManager\|\.GetRole()\|\.SetRole(" --include="*.go" .
  ```
- Functions that previously only received a `userstore.UserInterface` (e.g. layout menu builders) now also need `app.AppInterface` and `context.Context` parameters.

---

### 2. New migration creates role tables and seeds default roles

**Change**: `database/migrations/2026_09_23_0001_store_user_roles_migrate.go` runs `store.MigrateUp` (which now creates `snv_users_role` and `snv_users_user_role`), seeds the four default roles by handle, and backfills `user_role` rows from the legacy `role` column.

The seeded roles:

| Handle | Name |
|---|---|
| `administrator` | Administrator |
| `manager` | Manager |
| `superuser` | Superuser |
| `user` | User |

**Action Required**:
- Copy the migration file and register it in `getStoreMigrations` in `database/migrations/migrate.go`:
  ```go
  if cfg.GetUserStoreUsed() {
  	migrations = append(migrations, &StoreUserMigrate{app: reg})
  	migrations = append(migrations, &StoreUserRolesMigrate{app: reg})
  }
  ```
- The backfill maps each user's existing `role` column value to the role with the same handle. Users with an empty or unknown role value get no assignment.
- `Down()` drops only the two new role tables; the user table is owned by `StoreUserMigrate`.

---

### 3. `NewUserStore` requires role table names

**Change**: With `RolesEnabled: true`, `userstore.NewStore` returns an error unless `RoleTableName` and `UserRoleTableName` are set.

```go
return userstore.NewStore(userstore.NewStoreOptions{
	DB:                db,
	RolesEnabled:      true,
	UserTableName:     "snv_users_user",
	RoleTableName:     "snv_users_role",
	UserRoleTableName: "snv_users_user_role",
})
```

**Action Required**:
- Apply the same options in your `NewUserStore` (adjust table names to your project's naming convention).

---

### 4. New users get the `user` role automatically

**Change**: `userCreate` in `internal/controllers/auth/shared/session_login.go` calls `helpers.UserActiveRoleAssign(ctx, app, userID, userstore.USER_ROLE_USER)` after `UserCreate`.

**Action Required**:
- If you have custom user-creation code (seeders, admin controllers, CLI tools), assign roles via `UserActiveRoleAssign` or `UserRoleFindByUserIDAndRoleIDOrCreate` instead of `SetRole`.
- Test seeders: `testutils.SeedUser` now assigns `user`/`administrator` role rows for `test.USER_01`/`test.ADMIN_01` via `RoleFindByHandleOrCreate`.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.45.0"
```

### Step 2: Enable roles on the user store

Update `NewUserStore` in `internal/config/store_builders.go` as shown in Breaking Change #3.

### Step 3: Add the roles migration

Copy `2026_09_23_0001_store_user_roles_migrate.go` into `database/migrations` and register it after `StoreUserMigrate` in `migrate.go`.

### Step 4: Add the role helpers

Copy `internal/helpers/user_roles.go` and replace all legacy role call sites (Breaking Change #1). Notable sites in stock Blueprint:

- `internal/controllers/auth/shared/session_login.go` — `calculateRedirectURL` and `userCreate`
- `internal/middlewares/admin_middleware.go` — `adminUserAdapter.HasRole` now queries the store
- `internal/middlewares/subscription_middleware.go`
- `internal/controllers/website/blog/post/post_controller.go`
- `internal/layouts/user_layout_user_menu_items.go`, `admin_layout_user_menu_items.go` (signatures gained `app` + `ctx` params)
- `internal/testutils/seed_user.go`, `cmd/ai-browser/seed_user.go`

### Step 5: Update dependencies

```bash
go mod tidy
go build ./...
```

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **Roles seeded**: after startup, verify the `snv_users_role` table contains the four default handles.
4. **Backfill**: verify existing admin users have a `user_role` row linking them to the `administrator` role.
5. **Admin access**: log in as an admin user and confirm `/admin` is reachable; confirm a regular user is redirected away.

---

## 📝 Additional Notes

- The legacy `role` column still exists on the user table (userstore keeps it in the schema) — it is simply no longer written or read by Blueprint. You may drop it in a future custom migration once you are confident nothing else reads it.
- `userstore` also ships optional `GroupsEnabled` group tables — Blueprint does not enable them; set `GroupsEnabled: true` + table names if you need groups.
- The admin user adapter in `admin_middleware.go` resolves roles per-request via `r.Context()` — keep that pattern rather than caching roles on the user object.

---

## 🆘 Common Issues and Solutions

### Issue: `user store: RoleTableName is required when RolesEnabled is true`

**Cause**: `RolesEnabled: true` without the role table names.

**Solution**: Add `RoleTableName` and `UserRoleTableName` to `NewStoreOptions` (Breaking Change #3).

### Issue: Admin users redirected away from `/admin` after upgrade

**Cause**: Roles were migrated but the migration didn't run, or the backfill didn't match (e.g. legacy `role` value had no matching handle).

**Solution**: Confirm `2026_09_23_0001_store_user_roles_migrate` ran (check `migration_tracker`), then verify a `user_role` row exists linking the user ID to the `administrator` role ID. Assign manually via `UserActiveRoleAssign` if needed.

### Issue: Tests fail because seeded users have no roles

**Cause**: Custom seeders still call `SetRole` instead of assigning a role row.

**Solution**: Use `RoleFindByHandleOrCreate` + `UserRoleFindByUserIDAndRoleIDOrCreate`, or `helpers.UserActiveRoleAssign`.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
