# Upgrade Guide: v0.40.0 to v0.41.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.40.0 to v0.41.0.

## Overview

This release introduces **`IsEnabledXxxStore()` / `IsDisabledXxxStore()`** methods on `app.AppInterface` for every store, and replaces all direct nil-comparison guards (`Get*Store() == nil` / `!= nil`) across the codebase with these named predicates.

The motivation is to eliminate operator-inversion bugs. The old pattern `app.GetFeedStore() == nil` vs `!= nil` is a single-character difference that inverts logic, is easy to miss in code review, and has already caused bugs (e.g. a controller that forgot the guard entirely). The new `IsDisabledFeedStore()` / `IsEnabledFeedStore()` methods make intent explicit and remove the `==`/`!=` and `!` inversion hazards.

**Key Changes:**
- `AppInterface` gains two new methods per store: `IsEnabledXxxStore() bool` and `IsDisabledXxxStore() bool` (48 new methods across 24 stores)
- `appImplementation` implements all 48 methods; `IsEnabled*` checks `r != nil && r.xxxStore != nil`, `IsDisabled*` returns `!IsEnabled*`
- All ~160 call sites in the blueprint are migrated from `Get*Store() == nil` / `!= nil` to the new predicates
- `GetXxxStore()` / `SetXxxStore()` methods are unchanged — they stay as-is
- `stores_config.go` is unchanged — the new methods check runtime state (nil pointer), not config flags
- No dependency changes in `go.mod`

**Affected stores:** AuditStore, BlogStore, ChatStore, BlindIndexStoreEmail, BlindIndexStoreFirstName, BlindIndexStoreLastName, CacheStore, CmsStore, CustomStore, EntityStore, FeedStore, GeoStore, LogStore, MetaStore, SessionStore, SettingStore, ShopStore, SqlFileStorage, StatsStore, SubscriptionStore, TaskStore, UserStore, VaultStore.

---

## ⚠️ Breaking Changes

---

### 1. New `IsEnabled*Store()` / `IsDisabled*Store()` Methods on `AppInterface`

**Change**: `app.AppInterface` (in `internal/app/app_interface.go`) now declares two additional methods for every store accessor. Any type implementing `AppInterface` must implement these new methods or it will fail to compile.

For each existing `GetXxxStore()` / `SetXxxStore()` pair, the interface now also requires:

```go
// Feed store
GetFeedStore() feedstore.StoreInterface
SetFeedStore(s feedstore.StoreInterface)
IsEnabledFeedStore() bool
IsDisabledFeedStore() bool
```

The full list of added method pairs:

| Store | IsEnabled method | IsDisabled method |
|---|---|---|
| Audit | `IsEnabledAuditStore()` | `IsDisabledAuditStore()` |
| Blog | `IsEnabledBlogStore()` | `IsDisabledBlogStore()` |
| Chat | `IsEnabledChatStore()` | `IsDisabledChatStore()` |
| BlindIndex (Email) | `IsEnabledBlindIndexStoreEmail()` | `IsDisabledBlindIndexStoreEmail()` |
| BlindIndex (FirstName) | `IsEnabledBlindIndexStoreFirstName()` | `IsDisabledBlindIndexStoreFirstName()` |
| BlindIndex (LastName) | `IsEnabledBlindIndexStoreLastName()` | `IsDisabledBlindIndexStoreLastName()` |
| Cache | `IsEnabledCacheStore()` | `IsDisabledCacheStore()` |
| Cms | `IsEnabledCmsStore()` | `IsDisabledCmsStore()` |
| Custom | `IsEnabledCustomStore()` | `IsDisabledCustomStore()` |
| Entity | `IsEnabledEntityStore()` | `IsDisabledEntityStore()` |
| Feed | `IsEnabledFeedStore()` | `IsDisabledFeedStore()` |
| Geo | `IsEnabledGeoStore()` | `IsDisabledGeoStore()` |
| Log | `IsEnabledLogStore()` | `IsDisabledLogStore()` |
| Meta | `IsEnabledMetaStore()` | `IsDisabledMetaStore()` |
| Session | `IsEnabledSessionStore()` | `IsDisabledSessionStore()` |
| Setting | `IsEnabledSettingStore()` | `IsDisabledSettingStore()` |
| Shop | `IsEnabledShopStore()` | `IsDisabledShopStore()` |
| SqlFileStorage | `IsEnabledSqlFileStorage()` | `IsDisabledSqlFileStorage()` |
| Stats | `IsEnabledStatsStore()` | `IsDisabledStatsStore()` |
| Subscription | `IsEnabledSubscriptionStore()` | `IsDisabledSubscriptionStore()` |
| Task | `IsEnabledTaskStore()` | `IsDisabledTaskStore()` |
| User | `IsEnabledUserStore()` | `IsDisabledUserStore()` |
| Vault | `IsEnabledVaultStore()` | `IsDisabledVaultStore()` |

**Action Required**:
- If your project has a custom type implementing `app.AppInterface` (a mock, fake, or alternative implementation), you must add all 48 methods. See Breaking Change #2 for the implementation pattern.
- The blueprint's own `appImplementation` already implements them; projects using `app.New()` directly are unaffected at the implementation level.
- If your project has test doubles implementing `AppInterface`, update them (see Breaking Change #2).

---

### 2. Implementing the New Methods on Custom `AppInterface` Types

**Change**: Any custom type satisfying `app.AppInterface` must implement the new methods. The reference implementation in `appImplementation` is:

```go
func (r *appImplementation) IsEnabledFeedStore() bool {
	return r != nil && r.feedStore != nil
}

func (r *appImplementation) IsDisabledFeedStore() bool {
	return !r.IsEnabledFeedStore()
}
```

**Old Usage** (a custom mock that previously satisfied `AppInterface`):
```go
type fakeApp struct{}

func (f *fakeApp) GetFeedStore() feedstore.StoreInterface { return nil }
func (f *fakeApp) SetFeedStore(s feedstore.StoreInterface) {}
// ... other Get/Set pairs ...
```

**New Usage**:
```go
type fakeApp struct {
	feedStore feedstore.StoreInterface
	// ... other store fields ...
}

func (f *fakeApp) GetFeedStore() feedstore.StoreInterface { return f.feedStore }
func (f *fakeApp) SetFeedStore(s feedstore.StoreInterface) { f.feedStore = s }
func (f *fakeApp) IsEnabledFeedStore() bool { return f != nil && f.feedStore != nil }
func (f *fakeApp) IsDisabledFeedStore() bool { return !f.IsEnabledFeedStore() }
// ... repeat for every store ...
```

**Action Required**:
- Search for custom `AppInterface` implementations:
  ```bash
  grep -rln "AppInterface" --include="*.go" .
  ```
- For each custom type, add the 48 methods following the pattern above. The `IsEnabled*` method returns `r != nil && r.<storeField> != nil`; the `IsDisabled*` method returns `!r.IsEnabled*Store()`.
- Note: the blueprint's test suite uses the real `app.New()` via `internal/testutils.Setup`, so it does not require mock updates. Your project may differ.

---

### 3. Call-Site Migration: Replace Nil-Comparison Guards

**Change**: All guard clauses that check store availability via `Get*Store() == nil` or `Get*Store() != nil` should be migrated to the new `IsDisabled*Store()` / `IsEnabled*Store()` methods. The blueprint has already migrated all ~160 of its own call sites; projects derived from blueprint should do the same.

There are four patterns to migrate:

#### Pattern A — Guard clause (early return on missing store)

**Old Usage**:
```go
if app.GetFeedStore() == nil {
	return errors.New("feed store not configured")
}
```

**New Usage**:
```go
if app.IsDisabledFeedStore() {
	return errors.New("feed store not configured")
}
```

#### Pattern B — Proceed clause (use store if available)

**Old Usage**:
```go
if app.GetFeedStore() != nil {
	app.GetFeedStore().LinkCount(ctx, query)
}
```

**New Usage**:
```go
if app.IsEnabledFeedStore() {
	app.GetFeedStore().LinkCount(ctx, query)
}
```

#### Pattern C — Init-then-check, non-nil (variable only used inside block)

**Old Usage**:
```go
if s := app.GetSessionStore(); s != nil {
	sessionStore = &authSessionStoreAdapter{store: s}
}
```

**New Usage**:
```go
if app.IsEnabledSessionStore() {
	s := app.GetSessionStore()
	sessionStore = &authSessionStoreAdapter{store: s}
}
```

#### Pattern D — Init-then-check, nil (variable used after the block)

**Old Usage**:
```go
if s := app.GetSessionStore(); s == nil {
	return errors.New("session store not configured")
}
// s is used here
```

**New Usage**:
```go
s := app.GetSessionStore()
if app.IsDisabledSessionStore() {
	return errors.New("session store not configured")
}
// s is used here
```

**Action Required**:
- Find all nil-comparison guards in your project:
  ```bash
  grep -rn "Get\w*Store()\s*==\s*nil" --include="*.go" .
  grep -rn "Get\w*Store()\s*!=\s*nil" --include="*.go" .
  grep -rn "GetSqlFileStorage()\s*==\s*nil" --include="*.go" .
  grep -rn "GetSqlFileStorage()\s*!=\s*nil" --include="*.go" .
  ```
- Apply the matching pattern (A, B, C, or D) to each occurrence.
- **Do NOT** introduce `!IsDisabled*Store()` or `!IsEnabled*Store()` calls — always use the positive form. `!IsDisabledFeedStore()` reintroduces the `!` inversion hazard this refactor is meant to eliminate; use `IsEnabledFeedStore()` instead.
- `GetXxxStore()` / `SetXxxStore()` calls that *use* the store (e.g. `app.GetFeedStore().Method()`) stay unchanged — only the nil-comparison guards change.

---

### 4. Compound Conditions

**Change**: Compound conditions combining multiple store checks (e.g. with `&&` / `||`) are migrated component-by-component.

**Old Usage**:
```go
if app.GetConfig().GetTaskStoreUsed() && app.GetTaskStore() == nil {
	return errors.New("task store is enabled but not initialized")
}

if a.GetUserStore() == nil || a.GetSessionStore() == nil {
	return errors.New("user and session stores are required")
}
```

**New Usage**:
```go
if app.GetConfig().GetTaskStoreUsed() && app.IsDisabledTaskStore() {
	return errors.New("task store is enabled but not initialized")
}

if a.IsDisabledUserStore() || a.IsDisabledSessionStore() {
	return errors.New("user and session stores are required")
}
```

**Action Required**:
- Search for compound conditions containing store nil checks and migrate each `== nil` to `IsDisabled*Store()` and each `!= nil` to `IsEnabled*Store()`.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.41.0"
```

### Step 2: Update custom `AppInterface` implementations

If your project has any custom type implementing `app.AppInterface`, add all 48 `IsEnabled*Store()` / `IsDisabled*Store()` methods following the pattern in Breaking Change #2.

```bash
grep -rln "AppInterface" --include="*.go" .
```

### Step 3: Migrate call sites

Replace all `Get*Store() == nil` / `!= nil` guards with the new predicates (see Breaking Change #3). Use these searches to find them:

```bash
grep -rn "Get\w*Store()\s*==\s*nil" --include="*.go" .
grep -rn "Get\w*Store()\s*!=\s*nil" --include="*.go" .
grep -rn "GetSqlFileStorage()\s*==\s*nil" --include="*.go" .
grep -rn "GetSqlFileStorage()\s*!=\s*nil" --include="*.go" .
```

### Step 4: Verify no inversions remain

Ensure no `!IsDisabled*Store()` or `!IsEnabled*Store()` calls were introduced:

```bash
grep -rn "!\s*Is\(Enabled\|Disabled\)\w*Store()" --include="*.go" .
```

This should return no matches.

### Step 5: Verify no nil-comparison guards remain

```bash
grep -rn "Get\w*Store()\s*\(==\|!=\)\s*nil" --include="*.go" .
grep -rn "GetSqlFileStorage()\s*\(==\|!=\)\s*nil" --include="*.go" .
```

These should return no matches (ignoring commented-out code).

---

## 🧪 Testing After Migration

1. **Build**: `go build ./...`
2. **Unit tests**: `go test ./...`
3. **Verify guard behavior**: Confirm that controllers/middlewares/tasks still skip work when a store is not configured, and still execute when it is. The new methods are semantically equivalent to the old nil checks.

---

## 📝 Additional Notes

- **No dependency changes**: This release does not add, remove, or bump any external dependencies. `go.mod` is unchanged.
- **No config changes**: `stores_config.go` and the config interface are unchanged. The new methods check runtime nil state, not config flags.
- **`GetXxxStore()` / `SetXxxStore()` unchanged**: The accessor and mutator methods remain as-is. Only the nil-comparison *guards* around them change.
- **Interface structure unchanged**: Stores remain flat on `AppInterface`; they are not grouped into sub-interfaces.
- **Semantic equivalence**: `IsEnabledXxxStore()` is exactly equivalent to `GetXxxStore() != nil` (with an additional nil-receiver guard), and `IsDisabledXxxStore()` is exactly equivalent to `GetXxxStore() == nil`. No behavioral change is intended — only readability and bug-prevention.

---

## 🆘 Common Issues and Solutions

### Issue: Compile error — missing `IsEnabled*Store()` / `IsDisabled*Store()` methods

**Cause**: A custom type implementing `app.AppInterface` does not implement the new methods.

**Solution**: Add all 48 methods to the custom type. See Breaking Change #2 for the implementation pattern. A quick way to find the offending type is to run `go build ./...` — the compiler will point at the type that fails to satisfy the interface.

### Issue: Accidentally wrote `!IsDisabledFeedStore()`

**Cause**: Habitual use of `!` when translating `!= nil`.

**Solution**: Replace `!IsDisabledFeedStore()` with `IsEnabledFeedStore()`. The whole point of this refactor is to avoid `!` inversions. Run the verification grep in Step 4 to catch these.

### Issue: Init-then-check variable used after the block

**Cause**: Pattern D — the variable declared in `if s := ...; s == nil` is used after the if block. Naively collapsing to `if app.IsDisabledSessionStore()` would drop the variable.

**Solution**: Move the `s := app.GetSessionStore()` declaration *before* the `if`, then use `if app.IsDisabledSessionStore()`. See Breaking Change #3, Pattern D.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
