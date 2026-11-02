# Upgrade Guide: v0.57.0 to v0.58.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.57.0 to v0.58.0.

## Overview

This release is a **performance release** for the public website surface, plus routine dependency bumps. Repeated blog/sitemap requests no longer hit the database on every call — results are served from the instance memory cache (`GetMemoryCache()`, `ttlcache`) with a 1h TTL, and card images are batch-resolved in one `MediaList` query instead of an N+1 per-post media lookup.

**Key Changes:**
- **Blog index caching** — `prepareData` in the blog home controller is cached per page (`blog:home:page:<n>`, 1h TTL)
- **N+1 fix on blog cards** — new `rules.PostImageURLs()` resolves card images for the whole page in a single `MediaList` query; `cardPost`/`postImage` now take a pre-resolved `imageURL`; new `shared.SizedThumbnailURLForImage()` variant
- **Post recommendations caching** — the "keep reading" pool (same 18 published posts on every post page) is cached in memory (1h TTL)
- **sitemap.xml caching** — the built XML is cached in memory (1h TTL) and `Cache-Control: public, max-age=3600` is sent so crawlers mostly never reach the origin
- `github.com/dracory/auth` v0.47.0 → v0.49.0 (renames `Passwordless` → `Otp` throughout — same rename documented in the v0.56→v0.57 guide; the v0.57.0 tag itself still carried the old names)
- Dependency bumps: `blogadmin` v0.3.2 → v0.5.0, `cmsblocks` v0.1.0 → v0.1.1, `email` v0.2.0 → v0.3.0, `entitystore` v1.23.2 → v1.24.0, `logadmin` v0.2.1 → v0.4.0, `statsadmin` v0.9.2 → v0.10.0, `taskstore` v1.33.2 → v1.35.0, `useradmin` v0.5.1 → v0.8.0, `modernc.org/libc` v1.79.1 → v1.80.0
- AI agent skills consolidated under `.agents/skills/` (`.windsurf` / `.devin` skill copies removed)

---

## ⚠️ Breaking Changes

### 1. `dracory/auth` renames `Passwordless` → `Otp` (v0.49.0)

**Change**: The auth package's OTP (passwordless) API surface was renamed in v0.49.0. This was already documented in the v0.56→v0.57 guide, but the `v0.57.0` tag itself still shipped auth v0.47.0 — if you did not already apply that section, apply it now.

| Old name | New name |
|----------|----------|
| `types.ConfigPasswordless` | `types.ConfigOtp` |
| `types.AuthPasswordlessInterface` | `types.AuthOtpInterface` |
| `auth.NewPasswordlessAuth` | `auth.NewOtpAuth` |
| `IsPasswordless()` | `IsOtp()` |
| `authadapters.BuildPasswordlessAuth` | `authadapters.BuildOtpAuth` |

**Action Required**:

```bash
grep -rln "Passwordless" internal/ | xargs sed -i \
  -e 's/AuthPasswordlessInterface/AuthOtpInterface/g' \
  -e 's/ConfigPasswordless/ConfigOtp/g' \
  -e 's/NewPasswordlessAuth/NewOtpAuth/g' \
  -e 's/BuildPasswordlessAuth/BuildOtpAuth/g' \
  -e 's/IsPasswordless()/IsOtp()/g'
```

### 2. `cardPost` / `postImage` signatures take a resolved `imageURL`

**Change**: The blog index card functions no longer resolve the image internally — they receive the batch-resolved URL so the page performs one `MediaList` query total instead of one per card.

**Old Usage** (`internal/controllers/website/blog/home`):

```go
func cardPost(ctx context.Context, app app.AppInterface, post blogstore.PostInterface) hb.TagInterface {
	// ...
	Child(postImage(ctx, app, post))
}

func postImage(ctx context.Context, app app.AppInterface, post blogstore.PostInterface) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(ctx, app, post, "300", "200", "80")
```

**New Usage**:

```go
func cardPost(ctx context.Context, app app.AppInterface, post blogstore.PostInterface, imageURL string) hb.TagInterface {
	// ...
	Child(postImage(ctx, app, post, imageURL))
}

func postImage(ctx context.Context, app app.AppInterface, post blogstore.PostInterface, imageURL string) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURLForImage(ctx, app, imageURL, "300", "200", "80")
```

**Action Required**:
- Update both signatures and the internal `postImage` call
- These are unexported helpers — only your own copies/tests need changing
- `shared.SizedThumbnailURL` is unchanged and still resolves per-post for single-post callers (e.g. the post page hero image)

### 3. `blogControllerData` gained an `imageURLs` field

**Change**: `internal/controllers/website/blog/home.blogControllerData` now carries `imageURLs map[string]string` (postID → image URL), populated by `rules.PostImageURLs`. Tests that build the struct directly and assert on rendered thumbnails must populate it.

**Old Usage** (test):

```go
data := blogControllerData{
	postList:  []blogstore.PostInterface{post},
	postCount: 25,
	page:      0,
	perPage:   12,
}
```

**New Usage**:

```go
data := blogControllerData{
	postList:  []blogstore.PostInterface{post},
	postCount: 25,
	page:      0,
	perPage:   12,
	imageURLs: map[string]string{"1": "http://example.com/img.png"},
}
```

---

## 🔄 Migration Steps

### Step 1: Update dependencies

```bash
go get github.com/dracory/auth@v0.49.0
go get github.com/dracory/blogadmin@v0.5.0
go get github.com/dracory/cmsblocks@v0.1.1
go get github.com/dracory/email@v0.3.0
go get github.com/dracory/entitystore@v1.24.0
go get github.com/dracory/logadmin@v0.4.0
go get github.com/dracory/statsadmin@v0.10.0
go get github.com/dracory/taskstore@v1.35.0
go get github.com/dracory/useradmin@v0.8.0
go mod tidy
```

### Step 2: Apply the auth rename (if not already done under v0.57)

Run the sed script in Breaking Change 1. Blueprint's own usage was in `internal/authadapters/builder.go`, `internal/authadapters/email_adapter.go`, `internal/controllers/auth/routes.go`, and `internal/authadapters/adapters_test.go` — check the same places in your project.

### Step 3: Add `rules.PostImageURLs` and `SizedThumbnailURLForImage`

Add the batch resolver to `internal/rules/post_image.go` (see `rules.PostImageURLs` in this version — one `MediaList` with `EntityIDIn`, first media per entity wins, falls back to `post.GetImageUrl()`).

Add the pre-resolved variant to `internal/controllers/website/blog/shared/sized_thumbnail_url.go`:

```go
// SizedThumbnailURLForImage is SizedThumbnailURL with a pre-resolved image
// URL — used by list pages that batch-resolve post images in one query
// instead of hitting the media table per card.
func SizedThumbnailURLForImage(ctx context.Context, app app.AppInterface, postImageURL, width, height, quality string) string {
	extension := ImageExtension(postImageURL)

	if app == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, app.GetFileCache())
}
```

Keep `SizedThumbnailURL` as a thin wrapper that resolves the image then delegates.

### Step 4: Cache `prepareData` in the blog home controller

After the blog-store nil check, short-circuit on the memory cache; at the end, batch-resolve images and store the whole `blogControllerData`:

```go
// Blog content changes rarely, so each page's post list/count/images is
// served from the memory cache.
cacheKey := "blog:home:page:" + cast.ToString(page)
if mc := controller.app.GetMemoryCache(); mc != nil {
	if item := mc.Get(cacheKey); item != nil {
		if cached, ok := item.Value().(blogControllerData); ok {
			return cached, ""
		}
	}
}
// ... existing PostList/PostCount queries ...
data.imageURLs = rules.PostImageURLs(r.Context(), blogStore, postList)

if mc := controller.app.GetMemoryCache(); mc != nil {
	mc.Set(cacheKey, data, blogHomeCacheTTL) // time.Hour
}
```

If your blog index has extra query params (e.g. a `q` search filter), include them in the key so filtered listings cache independently.

### Step 5: Cache the recommendations pool in the post controller

```go
const recommendationsCacheKey = "blog:recommendations"
var postList []blogstore.PostInterface
if mc := c.app.GetMemoryCache(); mc != nil {
	if item := mc.Get(recommendationsCacheKey); item != nil {
		postList, _ = item.Value().([]blogstore.PostInterface)
	}
}
if postList == nil {
	var err error
	postList, err = c.app.GetBlogStore().PostList(context.Background(), options)
	if err != nil {
		return hb.Div()
	}
	if mc := c.app.GetMemoryCache(); mc != nil {
		mc.Set(recommendationsCacheKey, postList, time.Hour)
	}
}
```

### Step 6: Cache the sitemap body

Wrap `buildSitemapXML` in `sitemap_xml_controller.go` behind a memory-cache check (`website:sitemap.xml:mem`, 1h) and add `w.Header().Set("Cache-Control", "public, max-age=3600")` to the handler. Guard the cache read/write with `c.app != nil` — tests construct the controller without an app.

### Step 7: Move skills (optional)

If your project copied Blueprint's `.windsurf/` or `.devin/skills/` skill directories, they now live under `.agents/skills/`. Update any references (e.g. `docs/version_workflow.md` now points at `.agents/skills/upgrade-guide-generator/SKILL.md`).

### Step 8: Tidy, build, test

```bash
go mod tidy
gofmt -w ./internal/rules ./internal/controllers/website
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Unit Tests**: `go test ./...` — all suites should pass; update `blogControllerData` fixtures to include `imageURLs` where thumbnail HTML is asserted
2. **Blog index**: load `/blog` twice — the second request should not issue `PostList`/`PostCount`/`MediaList` queries (check store debug logging or query counts)
3. **Card images**: confirm thumbnails still render — every card image URL should resolve via the single `PostImageURLs` call, falling back to `GetImageUrl()` when no media exists
4. **Post page**: verify "keep reading" recommendations render and the second post load does not re-query `PostList`
5. **Sitemap**: hit `/sitemap.xml` twice — confirm the `Cache-Control` header and that the second response is byte-identical and DB-free
6. **Auth**: if the rename was applied, verify OTP/magic-link login still works on `dracory/auth` v0.49.0

---

## 📝 Additional Notes

- **Cache invalidation**: all new caches use a 1h TTL with no explicit invalidation. If your project writes posts via an admin UI or import task and needs immediate visibility, either shorten the TTL or delete the keys (`blog:home:page:*`, `blog:recommendations`, `website:sitemap.xml:mem`) on publish.
- **No DB-backed layer**: unlike a cache-store solution, memory-cache entries don't survive restart — first request after boot pays the queries. That's deliberate: content correctness on cold start beats warmth.
- **Admin packages jumped two+ minors** (`blogadmin`, `useradmin`, `logadmin`, `statsadmin`, `email`) but required zero code changes in Blueprint — they're version bumps only. If you call their APIs directly, skim each package's changelog for renamed options.
- `docs/version_workflow.md` now references `.agents/skills/` for the guide-generator skill.

## 🆘 Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| `undefined: BuildPasswordlessAuth` / `types.ConfigPasswordless` | Apply the Otp rename (Breaking Change 1) and `go get github.com/dracory/auth@v0.49.0` |
| Blog card test fails on missing thumbnail URL | Populate `imageURLs` in the `blogControllerData` fixture (Breaking Change 3) |
| Sitemap test panics with nil app | The cache path guards on `c.app != nil`; keep your controller constructible without an app or pass a test app via `testutils.Setup()` |
| Stale blog content after publishing a post | Expected — 1h TTL. Delete the cache keys on publish or shorten `blogHomeCacheTTL` |
| `unknown field EntityIDIn in MediaQueryOptions` | Your `blogstore` is too old — requires v1.39.x which supports `EntityIDIn` |

## 📞 Support

- Blueprint repository: https://github.com/dracory/blueprint
- dracory/auth: https://github.com/dracory/auth
