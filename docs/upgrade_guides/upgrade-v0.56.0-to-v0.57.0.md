# Upgrade Guide: v0.56.0 to v0.57.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.56.0 to v0.57.0.

## Overview

This release extracts more Blueprint-owned code into standalone `dracory/*` packages:

**Key Changes:**
- **Auth register & logout moved into `dracory/auth` v0.47.0** — the package now serves `/auth/register` (both the password sign-up form and the post-authentication profile-completion form) and `/auth/logout` itself. Blueprint's local controllers `internal/controllers/auth/register`, `internal/controllers/auth/register_password`, and `internal/controllers/auth/logout` are deleted
- **`internal/cmsblocks` block types extracted to `github.com/dracory/cmsblocks` v0.1.0** — `blogpost`, `blogpostlist`, and `search` blocks now live in the external package; `CmsAddBlockTypes` only injects stores and a post-URL builder
- `github.com/dracory/auth` v0.44.0 → v0.47.0
- `github.com/dracory/cmsstore` v1.44.2 → v1.45.0
- `github.com/dracory/cmsblocks` v0.1.0 (new dependency)
- `github.com/dracory/social` v0.1.0 moved from direct to indirect (pulled in by `dracory/auth`)

---

## ⚠️ Breaking Changes

### 1. Local auth register/logout controllers removed

**Change**: `internal/controllers/auth/register`, `internal/controllers/auth/register_password`, and `internal/controllers/auth/logout` (including their `app.html`/`app.css`/`app.js` assets and tests) are deleted. The `dracory/auth` package now owns the full `/auth` surface — login, register (sign-up and profile completion), register-code verify, and logout — through a single catch-all mux.

**Old Usage** (`internal/controllers/auth/routes.go`):

```go
func Routes(application app.AppInterface) []rtr.RouteInterface {
	loginMethods := application.GetConfig().GetLoginMethods()

	routes := []rtr.RouteInterface{}
	if application.GetConfig().GetRegistrationEnabled() {
		routes = append(routes, registerRoutes(application, loginMethods)...)
	}
	routes = append(routes, rtr.NewRoute().
		SetName("Auth > Logout Controller").
		SetPath(links.AUTH_LOGOUT).
		SetHTMLHandler(logout.NewLogoutController(application).AnyIndex))

	return append(routes, packageAuthRoutes(application)...)
}
```

**New Usage**:

```go
// Routes builds the /auth surface. The dracory/auth package serves every
// login method plus its built-in pages (login, register, register-code
// verify, logout) through one mux, mounted as a catch-all on /auth/*.
func Routes(application app.AppInterface) []rtr.RouteInterface {
	return packageAuthRoutes(application)
}
```

**Action Required**:
- Delete your local `internal/controllers/auth/register`, `register_password`, and `logout` directories
- Simplify `Routes()` to return `packageAuthRoutes(application)` directly
- The logout controller previously cleared an admin "remember" session/cookie — that cleanup is now handled inside the package's logout flow
- If your project customized the register templates, port those customizations via the package's branding/template options instead of local controllers

### 2. Coarse `/auth` rate limit relaxed (5/min → 30/min)

**Change**: The blanket IP rate limit on the `/auth` mux is raised because a single OTP journey legitimately needs ~5 requests. Per-endpoint protection moved inside the package (sensitive API endpoints are rate-limited and email sends are throttled internally).

**Old Usage**:

```go
rtrMiddleware.RateLimitByIPMiddleware(5, 60)
```

**New Usage**:

```go
rtrMiddleware.RateLimitByIPMiddleware(30, 60)
```

**Action Required**:
- Update the `stricter` middleware slice in `packageAuthRoutes`
- Configure the package's own email-send throttle via `ConfigRateLimiting` (see below) — do not rely on the coarse IP cap to stop mail floods

### 3. `CmsAddBlockTypes` — block implementations moved to `dracory/cmsblocks`

**Change**: `internal/cmsblocks/blogpost`, `internal/cmsblocks/blogpostlist`, and `internal/cmsblocks/search` are deleted. The same block types are now provided by `github.com/dracory/cmsblocks` v0.1.0. Registration requires injecting a post-URL builder via `SetPostURLFunc` (the package doesn't know your app's link helpers).

**Old Usage**:

```go
import (
	"project/internal/cmsblocks/blogpost"
	"project/internal/cmsblocks/blogpostlist"
	"project/internal/cmsblocks/search"
)

blogPostListBlock := blogpostlist.NewBlogPostListBlockType(app.GetBlogStore())
cmsstore.RegisterCustomBlockType(blogPostListBlock)
blogPostBlock := blogpost.NewBlogPostBlockType(app.GetBlogStore())
cmsstore.RegisterCustomBlockType(blogPostBlock)
searchBlock := search.NewSearchBlockType(app.GetCmsStore(), app.GetBlogStore())
cmsstore.RegisterCustomBlockType(searchBlock)
```

**New Usage**:

```go
import (
	"project/internal/links"

	"github.com/dracory/cmsblocks/blogpost"
	"github.com/dracory/cmsblocks/blogpostlist"
	"github.com/dracory/cmsblocks/search"
	"github.com/dracory/cmsstore"
)

postURL := func(postID, postSlug string) string {
	return links.Website().BlogPost(postID, postSlug)
}

blogPostListBlock := blogpostlist.NewBlogPostListBlockType(app.GetBlogStore())
blogPostListBlock.SetPostURLFunc(postURL)
cmsstore.RegisterCustomBlockType(blogPostListBlock)

blogPostBlock := blogpost.NewBlogPostBlockType(app.GetBlogStore())
blogPostBlock.SetPostURLFunc(postURL)
cmsstore.RegisterCustomBlockType(blogPostBlock)

searchBlock := search.NewSearchBlockType(app.GetCmsStore(), app.GetBlogStore())
searchBlock.SetPostURLFunc(postURL)
cmsstore.RegisterCustomBlockType(searchBlock)
```

**Action Required**:
- Delete `internal/cmsblocks/blogpost`, `blogpostlist`, and `search` directories
- Update imports to `github.com/dracory/cmsblocks/{blogpost,blogpostlist,search}`
- Add `SetPostURLFunc(...)` on each block — without it, post links in rendered blocks will be wrong/empty
- If you edited the local block renderers, port your changes as a fork or upstream PR — the local copies are gone

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go get github.com/dracory/auth@v0.47.0
go get github.com/dracory/cmsblocks@v0.1.0
go get github.com/dracory/cmsstore@v1.45.0
go mod tidy
```

### Step 2: Remove local auth register/logout controllers

```bash
rm -rf internal/controllers/auth/register
rm -rf internal/controllers/auth/register_password
rm -rf internal/controllers/auth/logout
```

Update `internal/controllers/auth/routes.go`: drop the `registerRoutes`/`logout` mounting and the `registerRoutes` helper (see Breaking Change 1).

### Step 3: Wire the new package hooks in `internal/authadapters/builder.go`

For each of `BuildPasswordlessAuth`, `BuildPasswordAuth`, and `BuildAuthKnightAuth`:

```go
cfg := types.ConfigPasswordless{ // or ConfigUsernameAndPassword
	ConfigRateLimiting: types.ConfigRateLimiting{
		// 3 sends per email per 5-minute sliding window (shared across
		// login-otp, register, magic-link, password-restore sends).
		MaxEmailSends:   3,
		EmailSendWindow: 5 * time.Minute,
	},
	// ...
}
```

And on the impl setter blocks:

```go
impl.SetFuncUserProfileUpdate(funcUserProfileUpdate(application))
impl.SetFuncUserEmailGet(funcUserEmailGet(application))
```

`funcUserProfileUpdate` persists first/last name for an already-authenticated user (vault-tokenizing when the vault store is enabled) and evicts the `auth:user:<id>` cache entry so the middleware sees fresh data. `funcUserEmailGet` resolves the plaintext email for the completion form (vault-aware). See `internal/authadapters/builder.go` for the full implementations.

### Step 4: Replace the post-auth registration redirect with the package gate

Remove the manual `IsRegistrationCompleted` flash-redirect from `postAuthHooks.RedirectURL` and add the package-owned gate instead:

```go
types.PostAuthHooks{
	// ...
	RegistrationComplete: func(ctx context.Context, userID string) (bool, error) {
		user, err := application.GetUserStore().UserFindByID(ctx, userID)
		if err != nil || user == nil {
			return false, err
		}
		return user.IsRegistrationCompleted(), nil
	},
	// ...
}
```

The package now redirects incomplete registrations to `/auth/register` itself; the old flash-message workaround is gone.

### Step 5: Migrate cmsblocks

Update `internal/cmsblocks/cms_add_block_types.go` imports and registration (see Breaking Change 3), then delete the extracted directories:

```bash
rm -rf internal/cmsblocks/blogpost internal/cmsblocks/blogpostlist internal/cmsblocks/search
```

### Step 6: Tidy and build

```bash
go mod tidy
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Unit Tests**: `go test ./...` — all suites should pass unchanged
2. **Auth flows**: verify register (password sign-up), register-as-profile-completion after OTP/magic-link login, login, logout, and remember-me on `dracory/auth` v0.47.0
3. **Rate limiting**: confirm repeated OTP/magic-link sends get throttled after 3 sends per email per 5 minutes
4. **CMS blocks**: render a page containing Blog Post, Blog Post List, and Search blocks; confirm post links point at your site's `links.Website().BlogPost` URLs
5. **Registration completion**: log in with a user missing first/last name → should be redirected to `/auth/register`, submit names → redirected onward, names persisted (and vault-tokenized if `UserStoreVaultEnabled`)

---

## 📝 Additional Notes

- `github.com/dracory/social` is now an indirect dependency (used by `dracory/auth`); `go mod tidy` moves it automatically — no action needed
- `cmd/server/main.go` no longer carries the commented-out `internal/cmsblocks` import
- Obsolete design docs (`docs/proposals/api-layer-cors-token-auth.md`, `multiple-login-methods.md`, `remember-me.md`) were removed — the features they proposed are now implemented in `dracory/auth`

## 🆘 Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| `undefined: register.NewRegisterController` | Local register controllers were deleted — let the package serve `/auth/register` (Step 2) |
| Incomplete-registration users loop back to `/auth/register` after submitting | The auth middleware caches the user for ~5 min; ensure `funcUserProfileUpdate` deletes `auth:user:<id>` from the memory cache (Step 3) |
| Block post links render as `#` or empty | Missing `SetPostURLFunc` on the cmsblocks types (Step 5) |
| OTP emails stopped sending after 3 attempts | Expected — `ConfigRateLimiting` throttles sends; adjust `MaxEmailSends`/`EmailSendWindow` if needed |
| `cannot find package github.com/dracory/cmsblocks/...` | Run `go get github.com/dracory/cmsblocks@v0.1.0` |

## 📞 Support

- Blueprint repository: https://github.com/dracory/blueprint
- dracory/auth: https://github.com/dracory/auth
- dracory/cmsblocks: https://github.com/dracory/cmsblocks
