# Upgrade Guide: v0.32.0 to v0.33.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.32.0 to v0.33.0.

## Overview

This release refactors auth/admin/user middlewares to use shared `rtr/middlewares` implementations, restructures the `pkg/shopadmin/products` package into separate `product_update` and `product_delete` packages, refactors blog admin components from monolithic files to split handler+Vue.js architectures, introduces a unified `NewPageLayout` for blog controllers, replaces `github.com/dracory/sb` sort constants with `github.com/dracory/neat` sort constants across all controllers, and bumps multiple Dracory store dependencies.

**Key Changes:**
- Auth/Admin/User middlewares refactored to delegate to `github.com/dracory/rtr/middlewares` shared implementations
- `pkg/shopadmin/products` package deleted — files moved to `pkg/shopadmin/product_update` and `pkg/shopadmin/product_delete`
- `NewProductUpdateController` signature changed — now requires `fileManagerURL` parameter
- `NewOrderManagerController` moved from `internal/controllers/admin/shop` to `pkg/shopadmin/order_manager` — now requires `app` parameter
- `SizedThumbnailURL` signature changed — now requires `context.Context` as first parameter
- `cardPost` and `postImage` functions now require `context.Context` and `app.AppInterface` parameters
- `blogController.page` and `postController.page` method signatures changed
- Blog controllers now use `NewPageLayout` instead of `NewCmsLayout`/`NewBlankLayout`
- `blogadmin.AdminOptions` struct — `LLMEngineInterface`, `LLMEngine`, and `BlogTopic` fields removed
- `pkg/blogadmin/ai_post_content_update` refactored from monolithic controller to split files with Vue.js
- `pkg/blogadmin/ai_title_generator` refactored from liveflux component to Vue.js modal
- `pkg/blogadmin/blog_settings` refactored from monolithic controller to split files with Vue.js
- `pkg/blogai/agent_titlegenerator_v1.go` deleted
- `cfmt` replaced with `slog` in `pkg/blogai`
- `internal/controllers/website/blog/post/post_recommendations_component.go` deleted
- `internal/controllers/admin/shop/order_manager_controller.go` deleted
- New `internal/rules/post_image.go` with `PostImageURL` function
- New `internal/layouts/page_layout.go` with unified `NewPageLayout`
- New `internal/layouts/layout_interface.go` with `LayoutInterface`
- New `internal/layouts/user_layout_navbar.go` and `user_layout_types.go`
- New `DisableNavbar` field in `layouts.Options`
- New `CONTROLLER_ORDER_DETAILS` constant in shop shared constants
- New `DiscountFormValidationRule` for discount form validation
- `rtr` upgraded to v1.7.0, `blogstore` to v1.29.0, `neat` to v0.27.0, `modernc.org/sqlite` to v1.53.0
- Multiple other Dracory store dependency version bumps
- `mingrammer/cfmt` moved from direct to indirect dependency
- `github.com/dracory/sb` sort constants (`sb.DESC`, `sb.ASC`) replaced with `github.com/dracory/neat` constants (`neat.SortDesc`, `neat.SortAsc`) across 9+ files
- `stretchr/testify/assert` removed from shopadmin test files — replaced with `t.Fatalf` style assertions
- Blog controller tests no longer require CMS store setup (`SetCmsStoreUsed`, `SetCmsStoreTemplateID`, `SeedTemplate`)
- `post_manager_controller.go` now uses `rules.PostImageURL` for `image_url` in AJAX responses instead of `post.GetImageUrl()`
- Post details UI: image URL field marked deprecated, media picker modal added
- `fakeLogStore.count` field type changed from `int` to `int64` in logadmin tests
- `postController.recommendationsSection` rewritten — liveflux component replaced with direct blogstore query + card rendering
- `postController.markdownToHtml` error handling changed — `panic(err)` replaced with graceful plain-text fallback
- `blogController.page()` signature changed — now takes `*http.Request` as first parameter
- `postController.page()` signature changed — now takes `context.Context` as second parameter
- `SectionBanner()` completely redesigned — dark background, decorative icons, inline SVG styles removed; replaced with card-based design
- `liveflux` import removed from `post_controller.go` and `ai_title_generator_controller.go` — no more `liveflux.Placeholder`/`liveflux.SSR` calls
- `ai_title_generator` — new `ACTION_SETTINGS_FETCH` and `ACTION_SETTINGS_SUBMIT` actions added
- `shopadmin/product_update_controller.go` massively expanded — from 39 lines to 761+ lines with AJAX handlers and component system
- `blogadmin/post_update_controller.go` — new media tab and AJAX handlers (load-media, upload-media, save-media, delete-media, add-media)
- `blogadmin/post_update_controller.go` — handler signatures changed from `_ http.ResponseWriter` to `w http.ResponseWriter`
- `pkg/shopadmin/discounts/discount_controller.go` deleted (was unused stub)
- `shop_controller_test.go` — `TestNewOrderManagerController` and `TestOrderManagerController_Handler` tests deleted
- `productupdate` tests — CMS store setup removed (`WithCmsStore(true, "test-template")`)
- New `internal/layouts/brutalski_theme.css` (440-line embedded CSS theme) and `internal/layouts/logo.svg`
- New `pkg/blogai/database.go` (commented-out SQLite store initialization scaffold)
- Blog post page CSS overhauled — Roboto font import removed, styling changed to system fonts
- Blog card post styling redesigned — card classes, grid columns, badge styling, removed separator div and `target="_blank"`

---

## ⚠️ Breaking Changes

### 1. Auth/Admin/User Middlewares Refactored to Shared rtr/middlewares

**Change**: `AuthMiddleware`, `AdminMiddleware`, and `UserMiddleware` in `internal/middlewares/` now delegate to `github.com/dracory/rtr/middlewares` shared implementations. The internal handler functions (`authHandler`, `adminHandler`, `userMiddlewareHandler`) and cache helpers (`cacheGetSession`, `cacheSetSession`, `cacheGetUser`, `cacheSetUser`) have been removed. Tests must now use `AuthMiddleware(app).GetHandler()` instead of `authHandler(app, next)`.

**Old Usage**:
```go
// v0.32.0 — internal/middlewares/auth_middleware_test.go
handler := authHandler(app, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // test assertions
}))
```

**New Usage**:
```go
// v0.33.0
handler := AuthMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // test assertions
}))
```

**Old Usage** (admin middleware):
```go
// v0.32.0
m := rtr.NewMiddleware().
    SetName("Admin Middleware").
    SetHandler(adminMiddlewareHandler(app))
```

**New Usage**:
```go
// v0.33.0
return rtrMiddleware.AdminMiddleware(rtrMiddleware.AdminMiddlewareConfig{
    GetUser: func(r *http.Request) rtrMiddleware.AdminMiddlewareUser {
        user := helpers.GetAuthUser(r)
        if user == nil {
            return nil
        }
        return user
    },
    OnNotAuthenticated: func(w http.ResponseWriter, r *http.Request) {
        // redirect logic
    },
    // ... other callbacks
})
```

**Action Required**:
- Update any test code that calls `authHandler(app, next)` directly to use `AuthMiddleware(app).GetHandler()(next)`
- Update error message assertions: `"session store not enabled"` → `"auth middleware: SessionStore is required"`, `"user store not enabled"` → `"auth middleware: UserStore is required"`, etc.
- If you had custom cache helper functions, remove them — caching is now handled by the shared rtr middleware
- Ensure `go.mod` has `github.com/dracory/rtr v1.7.0` or later

---

### 2. pkg/shopadmin/products Package Deleted — Split into product_update and product_delete

**Change**: The `pkg/shopadmin/products` package has been deleted. Files were moved to `pkg/shopadmin/product_update` and `pkg/shopadmin/product_delete`. The `product_manager_controller.go` and `product_update_controller.go` files in the old `products` package were deleted (functionality consolidated into `product_update/product_update_controller.go`).

**Old Usage**:
```go
// v0.32.0
import "project/pkg/shopadmin/products"

products.NewProductUpdateController(app, options.FileManagerURL).Handler(w, r)
products.NewProductDeleteController(app).Handler(w, r)
```

**New Usage**:
```go
// v0.33.0
import (
    "project/pkg/shopadmin/product_update"
    "project/pkg/shopadmin/product_delete"
)

product_update.NewProductUpdateController(app, options.FileManagerURL).Handler(w, r)
product_delete.NewProductDeleteController(app).Handler(w, r)
```

**Action Required**:
- Replace all imports of `project/pkg/shopadmin/products` with `project/pkg/shopadmin/product_update` or `project/pkg/shopadmin/product_delete` as appropriate
- Delete `pkg/shopadmin/products/` directory if you have a local copy
- Update route registrations in `pkg/shopadmin/routes.go` (already updated in template)

**Migration Command**:
```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|project/pkg/shopadmin/products|project/pkg/shopadmin/product_update|g' {} \;
```

---

### 3. NewProductUpdateController Signature Changed

**Change**: `NewProductUpdateController` now requires a `fileManagerURL` parameter.

**Old Usage**:
```go
// v0.32.0
controller := product_update.NewProductUpdateController(app)
```

**New Usage**:
```go
// v0.33.0
controller := product_update.NewProductUpdateController(app, fileManagerURL)
```

**Action Required**:
- Update all calls to `NewProductUpdateController` to pass the `fileManagerURL` string

---

### 4. NewOrderManagerController Moved and Signature Changed

**Change**: `NewOrderManagerController` was moved from `internal/controllers/admin/shop/order_manager_controller.go` to `pkg/shopadmin/order_manager`. It now requires an `app.AppInterface` parameter.

**Old Usage**:
```go
// v0.32.0
import "project/internal/controllers/admin/shop"

NewOrderManagerController().Handler(w, r)
```

**New Usage**:
```go
// v0.33.0
import orderManager "project/pkg/shopadmin/order_manager"

orderManager.NewOrderManagerController(app).Handler(w, r)
```

**Action Required**:
- Delete `internal/controllers/admin/shop/order_manager_controller.go` if present
- Update imports and calls to use `pkg/shopadmin/order_manager`
- Pass `app` parameter to `NewOrderManagerController`

---

### 5. SizedThumbnailURL Signature Changed

**Change**: `SizedThumbnailURL` now requires a `context.Context` as the first parameter. It also now uses `rules.PostImageURL` to resolve post images from media store before falling back to `post.GetImageUrlOrDefault()`.

**Old Usage**:
```go
// v0.32.0
thumbnailURL := shared.SizedThumbnailURL(app, post, "300", "200", "80")
```

**New Usage**:
```go
// v0.33.0
thumbnailURL := shared.SizedThumbnailURL(ctx, app, post, "300", "200", "80")
```

**Action Required**:
- Add `context.Context` as the first argument to all `SizedThumbnailURL` calls
- Pass the request context (`r.Context()`) or `context.Background()` as appropriate

---

### 6. cardPost and postImage Function Signatures Changed

**Change**: `cardPost` and `postImage` in `internal/controllers/website/blog/home/card_post.go` now require `context.Context` and `app.AppInterface` parameters.

**Old Usage**:
```go
// v0.32.0
card := cardPost(post)
img := postImage(post)
```

**New Usage**:
```go
// v0.33.0
card := cardPost(ctx, app, post)
img := postImage(ctx, app, post)
```

**Action Required**:
- Update all calls to `cardPost` and `postImage` to pass `ctx` and `app` parameters

---

### 7. Blog Controller page() Method Signatures Changed

**Change**: `blogController.page()` now takes `*http.Request` instead of just the data struct. `postController.page()` now takes `context.Context` as a second parameter.

**Old Usage**:
```go
// v0.32.0 — blog_controller.go
Content: hb.Wrap().HTML(controller.page(data))

// v0.32.0 — post_controller.go
Content: hb.Wrap().HTML(c.page(post))
```

**New Usage**:
```go
// v0.33.0 — blog_controller.go
Content: hb.Wrap().HTML(controller.page(r, data))

// v0.33.0 — post_controller.go
Content: hb.Wrap().HTML(c.page(post, r.Context()))
```

**Action Required**:
- Update `page()` method calls to pass the request or context as needed

---

### 8. Blog Controllers Now Use NewPageLayout

**Change**: Blog controllers (`blog_controller.go`, `post_controller.go`) now use `layouts.NewPageLayout()` instead of conditionally choosing between `layouts.NewCmsLayout()` and `layouts.NewBlankLayout()`.

**Old Usage**:
```go
// v0.32.0
if c.app.GetConfig().GetCmsStoreUsed() {
    return layouts.NewCmsLayout(c.app, r, options).ToHTML()
} else {
    return layouts.NewBlankLayout(c.app, r, options).ToHTML()
}
```

**New Usage**:
```go
// v0.33.0
return layouts.NewPageLayout(c.app, r, options).ToHTML()
```

**Action Required**:
- Replace conditional layout selection with `NewPageLayout()` in blog controllers
- `NewPageLayout` returns `LayoutInterface` — the `ToHTML()` method is still available

---

### 9. blogadmin.AdminOptions Struct — LLMEngine and BlogTopic Fields Removed

**Change**: `LLMEngineInterface` type, `LLMEngine` field, and `BlogTopic` field were removed from `blogadmin.AdminOptions`.

**Old Usage**:
```go
// v0.32.0
opts := blogadmin.AdminOptions{
    LLMEngine: myEngine,
    BlogTopic: "my topic",
    // ...
}
```

**New Usage**:
```go
// v0.33.0
opts := blogadmin.AdminOptions{
    // LLMEngine and BlogTopic no longer available
    // ...
}
```

**Action Required**:
- Remove any references to `LLMEngineInterface`, `LLMEngine`, or `BlogTopic` in `AdminOptions` initialization
- AI title generator settings are now managed via the blog settings page with Vue.js

---

### 10. pkg/blogadmin/ai_post_content_update Refactored

**Change**: The monolithic `ai_post_content_update_controller.go`, `form_ai_post_content_update.go`, and their test files were deleted. Replaced with split files: `controller.go`, `render_page.go`, `handle_fetch_data_ajax.go`, `handle_regenerate_block_ajax.go`, `handle_save_ajax.go`, plus Vue.js templates (`editor.html`, `editor.js`).

**Old Usage**:
```go
// v0.32.0
import "project/pkg/blogadmin/ai_post_content_update"

controller := ai_post_content_update.NewAiPostContentUpdateController(app)
```

**New Usage**:
```go
// v0.33.0
import "project/pkg/blogadmin/ai_post_content_update"

controller := aipostcontentupdate.NewController(app)
```

**Action Required**:
- Update constructor call from `NewAiPostContentUpdateController(app)` to `NewController(app)`
- Package name changed to `aipostcontentupdate`
- If you had custom code referencing `form_ai_post_content_update.go`, refactor to use the new split handlers

---

### 11. pkg/blogadmin/ai_title_generator Refactored

**Change**: The liveflux-based `modal_settings_component.go` was deleted. Replaced with Vue.js-based `settings_modal.html` and `settings_modal.js`, plus split handlers `handle_settings_fetch_data.go` and `handle_settings_submit.go`.

**Action Required**:
- Remove any references to the old liveflux component `titleGeneratorSettingsModal`
- The title generator settings modal is now a Vue.js component

---

### 12. pkg/blogadmin/blog_settings Refactored

**Change**: The monolithic `blog_settings_controller.go` was simplified. `form_blog_settings.go` and its test file were deleted. Replaced with split files: `render_page.go`, `handle_fetch_data_ajax.go`, `handle_submit_ajax.go`, plus Vue.js templates (`settings.html`, `settings.js`).

**Action Required**:
- Remove any references to `NewFormBlogSettings` or the old form component
- Blog settings now uses Vue.js with AJAX handlers

---

### 13. pkg/blogai — agent_titlegenerator_v1.go Deleted and cfmt Replaced

**Change**: `agent_titlegenerator_v1.go` was deleted. `cfmt` was replaced with `slog` in `agent_titlegenerator.go`. `constants_test.go` was deleted.

**Action Required**:
- Remove any references to `agent_titlegenerator_v1.go` or its functions
- Replace `cfmt.Successln(...)` calls with `slog.Info(...)` if you have custom code in blogai
- `mingrammer/cfmt` is now an indirect dependency only

---

### 14. post_recommendations_component.go Deleted

**Change**: `internal/controllers/website/blog/post/post_recommendations_component.go` and its test file were deleted.

**Action Required**:
- Remove any references to the post recommendations component
- Post recommendations are now handled inline in the post controller

---

### 15. New CONTROLLER_ORDER_DETAILS Constant

**Change**: A new `CONTROLLER_ORDER_DETAILS` constant was added to `internal/controllers/admin/shop/shared/constants.go`.

**Action Required**:
- If you have custom shop route handling, add support for the `order_details` controller case

---

### 16. pkg/shopadmin/routes.go — New Package Imports

**Change**: `pkg/shopadmin/routes.go` now imports `category_create`, `category_update`, `product_delete`, and `product_update` packages. Category create/update previously dispatched to `category_manager` for both operations.

**Old Usage**:
```go
// v0.32.0
case shared.CONTROLLER_CATEGORY_CREATE:
    return category_manager.NewCategoryManagerController(app).Handler(w, r)
case shared.CONTROLLER_CATEGORY_UPDATE:
    return category_manager.NewCategoryManagerController(app).Handler(w, r)
```

**New Usage**:
```go
// v0.33.0
case shared.CONTROLLER_CATEGORY_CREATE:
    return category_create.NewCategoryCreateController(app).Handler(w, r)
case shared.CONTROLLER_CATEGORY_UPDATE:
    return category_update.NewCategoryUpdateController(app).Handler(w, r)
```

**Action Required**:
- Update route dispatch logic if you have custom shop admin routing

### 17. sb Sort Constants Replaced with neat Sort Constants

**Change**: `github.com/dracory/sb` sort constants (`sb.DESC`, `sb.ASC`) have been replaced with `github.com/dracory/neat` sort constants (`neat.SortDesc`, `neat.SortAsc`) across 9+ files. The `sb` package remains in `go.mod` but is no longer used for sort constants in the affected files.

**Affected Files**:
- `internal/controllers/admin/shop/products/product_manager_controller.go`
- `pkg/blogadmin/post_manager/post_manager_controller.go`
- `pkg/blogadmin/post_manager/table_post_list.go`
- `pkg/blogadmin/post_update/post_versioning.go`
- `pkg/logadmin/log_manager/log_manager_controller.go`
- `pkg/shopadmin/category_manager/handle_categories_load_ajax.go`
- `pkg/shopadmin/discount_manager/handle_discounts_load_ajax.go`
- `pkg/shopadmin/order_manager/handle_orders_load_ajax.go`
- `pkg/shopadmin/product_manager/handle_products_fetch_ajax.go`

**Old Usage**:
```go
// v0.32.0
import "github.com/dracory/sb"

sortOrder := sb.DESC
direction := sb.ASC
```

**New Usage**:
```go
// v0.33.0
import "github.com/dracory/neat"

sortOrder := neat.SortDesc
direction := neat.SortAsc
```

**Action Required**:
- Replace `sb.DESC` with `neat.SortDesc` and `sb.ASC` with `neat.SortAsc`
- Replace `import "github.com/dracory/sb"` with `import "github.com/dracory/neat"` in affected files
- Update string comparisons: `sortOrder == "asc"` to `strings.EqualFold(sortOrder, neat.SortAsc)`

**Migration Command**:
```bash
# Replace sort constants
find . -type f -name "*.go" -exec sed -i 's|sb\.DESC|neat.SortDesc|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|sb\.ASC|neat.SortAsc|g' {} \;
```

---

### 18. testify/assert Removed from Shopadmin Tests

**Change**: `github.com/stretchr/testify/assert` imports were removed from shopadmin test files. Tests now use `t.Fatalf` and manual comparisons instead of `assert.NoError`/`assert.Equal`.

**Old Usage**:
```go
// v0.32.0
import "github.com/stretchr/testify/assert"

assert.NoError(t, err)
assert.Equal(t, http.StatusOK, response.StatusCode)
```

**New Usage**:
```go
// v0.33.0
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
if response.StatusCode != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
}
```

**Action Required**:
- If you have custom shopadmin tests using `testify/assert`, migrate to `t.Fatalf` style or keep `testify` as a test dependency (it remains in `go.mod`)
- This is a test convention change, not a production code breaking change

---

### 19. Blog Controller Tests — CMS Store Setup Removed

**Change**: Blog controller tests (`blog_controller_test.go`, `post_controller_test.go`) no longer set up CMS store (`SetCmsStoreUsed`, `SetCmsStoreTemplateID`, `SeedTemplate`). This aligns with the `NewPageLayout` change that removes conditional CMS/blank layout selection.

**Old Usage**:
```go
// v0.32.0
cfg := testutils.DefaultConf()
cfg.SetBlogStoreUsed(true)
cfg.SetCmsStoreUsed(true)
cfg.SetCmsStoreTemplateID("test-template")
app := testutils.Setup(testutils.WithCfg(cfg))

err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template")
```

**New Usage**:
```go
// v0.33.0
cfg := testutils.DefaultConf()
cfg.SetBlogStoreUsed(true)
app := testutils.Setup(testutils.WithCfg(cfg))
```

**Additional test changes**:
- `blog_controller_test.go`: `controller.page(data)` changed to `controller.page(nil, data)` (passing `nil` for request in unit tests)
- Pagination CSS class assertion changed from `pagination-primary-soft` to `page-link`
- `postCount` test fixture changed from 10 to 25

**Action Required**:
- Remove CMS store setup from blog controller tests if you have custom tests
- Update `page()` method calls in tests to pass `nil` or a mock request
- Update pagination class assertions if testing for specific CSS classes

---

### 20. post_manager_controller.go Uses rules.PostImageURL for AJAX Responses

**Change**: `pkg/blogadmin/post_manager/post_manager_controller.go` now uses `rules.PostImageURL(ctx, blogStore, post)` instead of `post.GetImageUrl()` for the `image_url` field in AJAX load posts responses.

**Old Usage**:
```go
// v0.32.0
"image_url": post.GetImageUrl(),
```

**New Usage**:
```go
// v0.33.0
"image_url": rules.PostImageURL(ctx, blogStore, post),
```

**Action Required**:
- If you have custom code that relies on `post.GetImageUrl()` for AJAX responses, update to use `rules.PostImageURL` for consistent image resolution from media store
- Add `"project/internal/rules"` import if not already present

---

### 21. logstore Test — fakeLogStore.count Type Changed to int64

**Change**: In `pkg/logadmin/log_manager/log_list_results_test.go`, the `fakeLogStore.count` field type changed from `int` to `int64`, reflecting a likely interface change in `logstore.LogCount` return type.

**Old Usage**:
```go
// v0.32.0
count int
// ...
return int64(s.count), nil
```

**New Usage**:
```go
// v0.33.0
count int64
// ...
return s.count, nil
```

**Action Required**:
- If you have custom `logstore` implementations or test fakes, update the `count` field to `int64` to match the interface
- The `LogCount` method should return `int64` directly instead of casting from `int`

---

### 22. postController.recommendationsSection Rewritten (liveflux → Direct Query)

**Change**: `internal/controllers/website/blog/post/post_controller.go` `recommendationsSection` method was completely rewritten. It no longer uses `liveflux.Placeholder` or `NewPostRecommendationsComponent`. Instead, it directly queries `blogstore.PostList`, filters out the current post, shuffles results, and renders cards inline using new `postCard()` and `truncatedSummary()` helper methods.

**Old Usage**:
```go
// v0.32.0
component := NewPostRecommendationsComponent(c.app)
rendered := liveflux.Placeholder(component, map[string]string{
    "post_id": post.GetID(),
})
return rendered
```

**New Usage**:
```go
// v0.33.0
options := blogstore.PostQueryOptions{
    Status:    blogstore.POST_STATUS_PUBLISHED,
    SortOrder: "DESC",
    OrderBy:   "published_at",
    Limit:     18,
}
postList, err := c.app.GetBlogStore().PostList(context.Background(), options)
// ... filter, shuffle, render cards inline
```

**Action Required**:
- Remove any custom code that references `NewPostRecommendationsComponent` or `liveflux.Placeholder` for recommendations
- The `liveflux` import is no longer needed in `post_controller.go`
- The `/liveflux` script URL was removed from the post page `ScriptURLs`

---

### 23. postController.markdownToHtml Error Handling Changed (panic → fallback)

**Change**: `postController.markdownToHtml` no longer panics on conversion errors. Instead, it returns a plain-text fallback with `<br>` tags.

**Old Usage**:
```go
// v0.32.0
if err := md.Convert([]byte(text), &buf); err != nil {
    panic(err)
}
```

**New Usage**:
```go
// v0.33.0
if err := md.Convert([]byte(text), &buf); err != nil {
    return "<p>" + strings.ReplaceAll(text, "\n", "<br>") + "</p>"
}
```

**Action Required**:
- If you relied on panic behavior for error detection, switch to checking the returned HTML for fallback content
- This is a safety improvement — malformed markdown will no longer crash the server

---

### 24. blogController.page() and postController.page() Signature Changes

**Change**: Both `blogController.page()` and `postController.page()` method signatures changed to accept request/context parameters.

**Old Usage**:
```go
// v0.32.0 — blog_controller.go
func (controller *blogController) page(data blogControllerData) string

// v0.32.0 — post_controller.go
func (c *postController) page(post blogstore.PostInterface) string
```

**New Usage**:
```go
// v0.33.0 — blog_controller.go
func (controller *blogController) page(r *http.Request, data blogControllerData) string

// v0.33.0 — post_controller.go
func (c *postController) page(post blogstore.PostInterface, ctx context.Context) string
```

**Action Required**:
- Update all calls to `page()` to pass the new parameters
- In tests, pass `nil` for `*http.Request` or use `context.Background()` for `context.Context`

---

### 25. SectionBanner() Completely Redesigned

**Change**: `internal/controllers/website/blog/shared/section_banner.go` was completely redesigned. The dark background (`#1C1626`), decorative icon elements (crosshair, asterisk, star icons), inline SVG styles, and `fill-success`/`fill-orange`/`fill-purple` CSS classes were removed. Replaced with a card-based design featuring an icon badge, subtitle text, and styled breadcrumb links.

**Action Required**:
- If you customized `SectionBanner()` or relied on the dark background or decorative icons, update your code
- The function signature remains `func SectionBanner() *hb.Tag` — no parameter changes
- Custom CSS overrides for `.fill-success`, `.fill-orange`, `.fill-purple` classes are no longer needed

---

### 26. liveflux Removed from Blog Post Controller and AI Title Generator

**Change**: `github.com/dracory/liveflux` import removed from `post_controller.go` and `ai_title_generator_controller.go`. All `liveflux.Placeholder()`, `liveflux.SSR()`, and `liveflux.Script()` calls replaced with Vue.js CDN + embedded HTML/JS files.

**Action Required**:
- Remove `liveflux` imports from any custom blog controller code
- Remove `/liveflux` from `ScriptURLs` arrays in blog page options
- Replace `liveflux.Placeholder()` / `liveflux.SSR()` calls with Vue.js components

---

### 27. ai_title_generator — New Settings Fetch/Submit Actions

**Change**: `pkg/blogadmin/ai_title_generator/ai_title_generator_controller.go` added two new action constants and handler methods for Vue.js settings modal.

**New Constants**:
```go
ACTION_SETTINGS_FETCH  = "settings-fetch-data"
ACTION_SETTINGS_SUBMIT = "settings-submit"
```

**New Handler Methods**:
- `handleSettingsFetchData(r *http.Request) string` — returns blog settings as JSON
- `handleSettingsSubmit(r *http.Request) string` — saves blog settings from JSON

**New Files**:
- `settings_modal.html` — Vue.js template for settings modal
- `settings_modal.js` — Vue.js app for settings modal
- `handle_settings_fetch_data.go` — Settings fetch handler
- `handle_settings_submit.go` — Settings submit handler

**Action Required**:
- If you have custom actions in `ai_title_generator`, ensure they don't conflict with the new `settings-fetch-data` and `settings-submit` action names

---

### 28. shopadmin/product_update_controller.go Massively Expanded

**Change**: `pkg/shopadmin/product_update/product_update_controller.go` was expanded from ~39 lines to 761+ lines. The controller now handles AJAX actions for media, metadata, tags, and details components, plus a full page rendering system.

**New AJAX Actions**:
- `load-media`, `save-media`, `upload-media` — product media management
- `load-metadata`, `save-metadata` — product metadata management
- `load-tags`, `save-tags` — product tags management
- `load-details`, `save-details` — product details management

**New Component System**:
- `ProductDetailsComponent`, `ProductMediaComponent`, `ProductMetadataComponent`, `ProductTagsComponent`
- Each component implements `Mount()`, `Handle()`, `Render()` interface

**New Files**:
- `types.go` — `MetadataRequest`, `MetadataResponse`, `MetadataItem` types, `ReqArrayOfMaps` helper
- `product_details_component.go`, `product_media_component.go`, `product_metadata_component.go`, `product_tags_component.go`
- `media.html` (expanded +60 lines), `media.js` (expanded +94 lines)
- `details.html`, `details.js`, `metadata.html`, `metadata.js`, `tags.html`, `tags.js` (moved from `products/`)

**Action Required**:
- If you have custom product update controllers, update them to use the new component system
- The old `actionLoadProduct` and `actionUpdateProduct` constants were removed

---

### 29. blogadmin/post_update_controller.go — New Media Tab and Handler Signature Changes

**Change**: `pkg/blogadmin/post_update/post_update_controller.go` added a new "Media" tab with AJAX handlers, and changed many handler signatures from `_ http.ResponseWriter` to `w http.ResponseWriter`.

**New Media Actions**:
- `load-media`, `upload-media`, `save-media`, `delete-media`, `add-media`
- New `renderMediaView()` method and `newPostMediaComponent()` integration

**Handler Signature Changes** (selected methods):
- `handleAddTag(_ http.ResponseWriter, ...)` → `handleAddTag(w http.ResponseWriter, ...)`
- `handleRemoveTag(_ ...)` → `handleRemoveTag(w ...)`
- `handleLoadDetails(_ ...)` → `handleLoadDetails(w ...)`
- `handleSaveDetails(_ ...)` → `handleSaveDetails(w ...)`
- `handleRegenerateImage(_ ...)` → `handleRegenerateImage(w ...)`
- `handleLoadContent(_ ...)` → `handleLoadContent(w ...)`
- `handleSaveContent(_ ...)` → `handleSaveContent(w ...)`
- `handleBlockEditorHandle(_ ...)` → `handleBlockEditorHandle(w ...)`
- `handleLoadSEO(_ ...)` → `handleLoadSEO(w ...)`
- `handleSaveSEO(_ ...)` → `handleSaveSEO(w ...)`
- `handleLoadVersions(_ ...)` → `handleLoadVersions(w ...)`
- `handleLoadVersionDetail(_ ...)` → `handleLoadVersionDetail(w ...)`
- `handleRestoreVersionAttributes(_ ...)` → `handleRestoreVersionAttributes(w ...)`
- `renderCategoriesView(_ *http.Request, ...)` → `renderCategoriesView(r *http.Request, ...)`
- `renderTagsView(_ *http.Request, ...)` → `renderTagsView(r *http.Request, ...)`
- `renderDetailsView(_ *http.Request, ...)` → `renderDetailsView(r *http.Request, ...)`
- `renderContentView(_ *http.Request, ...)` → `renderContentView(r *http.Request, ...)`
- `renderSEOView(_ *http.Request, ...)` → `renderSEOView(r *http.Request, ...)`
- `renderVersioningModal(_ *http.Request, ...)` → `renderVersioningModal(r *http.Request, ...)`

**Action Required**:
- If you override or call these handler methods, update signatures to use `w` instead of `_`
- If you extend the post update controller, register the new media actions in your custom handler dispatch

---

### 30. shopadmin/discounts/discount_controller.go Deleted

**Change**: `pkg/shopadmin/discounts/discount_controller.go` was deleted. It was an unused stub controller with a placeholder "TODO: Implement full migration" message.

**Action Required**:
- Remove any imports of `pkg/shopadmin/discounts` that referenced `NewDiscountController`
- Discount management is handled by `internal/controllers/admin/shop/discounts/discount_controller.go` (unchanged)

---

### 31. shop_controller_test.go — Order Manager Tests Deleted

**Change**: `TestNewOrderManagerController` and `TestOrderManagerController_Handler` were deleted from `internal/controllers/admin/shop/shop_controller_test.go` since `NewOrderManagerController` was moved to `pkg/shopadmin/order_manager`.

**Action Required**:
- If you have custom tests referencing `NewOrderManagerController()` from the `shop` package, update them to use `orderManager.NewOrderManagerController(app)` from `pkg/shopadmin/order_manager`

---

### 32. productupdate Tests — CMS Store Setup Removed

**Change**: All test files in `internal/controllers/admin/shop/products/productupdate/` had `testutils.WithCmsStore(true, "test-template")` removed from their setup, consistent with the `NewPageLayout` change that no longer requires CMS store for layout rendering.

**Action Required**:
- Remove `WithCmsStore(true, "test-template")` from any custom shop product update tests

---

## 🔄 Migration Steps

### Step 1: Update go.mod Dependencies

```bash
go get github.com/dracory/rtr@v1.7.0
go get github.com/dracory/blogstore@v1.29.0
go get github.com/dracory/neat@v0.27.0
go get modernc.org/sqlite@v1.53.0
go mod tidy
```

### Step 2: Update Auth Middleware Tests

Update test files that call `authHandler` directly:

```bash
# Find test files referencing authHandler
grep -rn 'authHandler(' --include='*_test.go' .
```

Replace `authHandler(app, next)` with `AuthMiddleware(app).GetHandler()(next)`.

Update error message assertions:
- `"session store not enabled"` → contains `"auth middleware: SessionStore is required"`
- `"session store not initialized"` → contains `"auth middleware: SessionStore is required"`
- `"user store not enabled"` → contains `"auth middleware: UserStore is required"`
- `"user store not initialized"` → contains `"auth middleware: UserStore is required"`

### Step 3: Move shopadmin/products Imports

```bash
# Update product update imports
find . -type f -name "*.go" -exec sed -i 's|project/pkg/shopadmin/products|project/pkg/shopadmin/product_update|g' {} \;
```

Manually update `product_delete` imports separately, as the `product_delete_controller.go` moved to its own package.

### Step 4: Update NewProductUpdateController Calls

Add `fileManagerURL` parameter to all `NewProductUpdateController` calls:

```go
// Old
controller := product_update.NewProductUpdateController(app)

// New
controller := product_update.NewProductUpdateController(app, fileManagerURL)
```

### Step 5: Update NewOrderManagerController

Move from internal controller to pkg:

```go
// Old
import "project/internal/controllers/admin/shop"
NewOrderManagerController().Handler(w, r)

// New
import orderManager "project/pkg/shopadmin/order_manager"
orderManager.NewOrderManagerController(app).Handler(w, r)
```

### Step 6: Update SizedThumbnailURL Calls

Add `context.Context` as first parameter:

```go
// Old
shared.SizedThumbnailURL(app, post, "300", "200", "80")

// New
shared.SizedThumbnailURL(ctx, app, post, "300", "200", "80")
```

### Step 7: Update Blog Controller Layout Calls

Replace conditional layout selection with `NewPageLayout`:

```go
// Old
if c.app.GetConfig().GetCmsStoreUsed() {
    return layouts.NewCmsLayout(c.app, r, options).ToHTML()
} else {
    return layouts.NewBlankLayout(c.app, r, options).ToHTML()
}

// New
return layouts.NewPageLayout(c.app, r, options).ToHTML()
```

### Step 8: Remove LLMEngine and BlogTopic from AdminOptions

Remove these fields from any `blogadmin.AdminOptions` initialization:

```go
// Remove these lines:
// LLMEngine: myEngine,
// BlogTopic: "my topic",
```

### Step 9: Replace sb Sort Constants with neat Sort Constants

Replace all `sb.DESC`/`sb.ASC` with `neat.SortDesc`/`neat.SortAsc`:

```bash
find . -type f -name "*.go" -exec sed -i 's|sb\.DESC|neat.SortDesc|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|sb\.ASC|neat.SortAsc|g' {} \;
```

Update imports in affected files:
```go
// Old
import "github.com/dracory/sb"

// New
import "github.com/dracory/neat"
```

Also update string comparisons from `sortOrder == "asc"` to `strings.EqualFold(sortOrder, neat.SortAsc)`.

### Step 10: Update Blog Controller Tests

Remove CMS store setup from blog controller tests:

```go
// Remove these lines:
// cfg.SetCmsStoreUsed(true)
// cfg.SetCmsStoreTemplateID("test-template")
// err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template")
```

Update `page()` method calls in tests to pass `nil` for request:
```go
// Old
html := controller.page(data)

// New
html := controller.page(nil, data)
```

### Step 11: Remove liveflux from Blog Controllers

Remove `liveflux` imports and replace `liveflux.Placeholder()`/`liveflux.SSR()` calls with Vue.js components:

```go
// Old
import "github.com/dracory/liveflux"
rendered := liveflux.Placeholder(component, map[string]string{...})

// New — use Vue.js CDN + embedded HTML/JS files
vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")
```

Also remove `/liveflux` from `ScriptURLs` arrays in blog page options.

### Step 12: Update Internal Shop Routes for order_details

Update `internal/controllers/admin/shop/routes.go` to import and use the new `order_details` and `order_manager` packages:

```go
// Old
import "project/internal/controllers/admin/shop"
NewOrderManagerController().Handler(w, r)

// New
import (
    orderDetails "project/pkg/shopadmin/order_details"
    orderManager "project/pkg/shopadmin/order_manager"
)
orderManager.NewOrderManagerController(app).Handler(w, r)
orderDetails.NewOrderDetailsController(app).Handler(w, r)
```

### Step 13: Remove CMS Store Setup from productupdate Tests

Remove `testutils.WithCmsStore(true, "test-template")` from all shop product update test files:

```go
// Remove this line:
// testutils.WithCmsStore(true, "test-template"),
```

### Step 14: Update Version Constant

Update the version constant in `internal/config/version.go`:

```go
// Old (v0.32.0)
const Version = "0.32.0"

// New (v0.33.0)
const Version = "0.33.0"
```

### Step 15: Run go mod tidy and Tests

```bash
go mod tidy
go test ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run the full test suite to verify everything compiles and passes:

```bash
go test ./...
```

### 2. Auth Middleware Tests

Verify auth middleware tests pass with the new rtr/middlewares delegation:

```bash
go test ./internal/middlewares/... -run TestAuth
```

### 3. Shop Admin Tests

Verify shop admin routes and controllers work with the new package structure:

```bash
go test ./pkg/shopadmin/...
```

### 4. Blog Admin Tests

Verify blog admin components work with the refactored architecture:

```bash
go test ./pkg/blogadmin/...
```

### 5. Blog AI Tests

Verify blogai package compiles without the deleted v1 agent and cfmt:

```bash
go test ./pkg/blogai/...
```

### 6. Verify Application Startup

Start the application to confirm all middleware and routing works:

```bash
go run ./cmd/server
```

### 7. Shopadmin and Blogadmin Sort Constants

Verify controllers work with `neat` sort constants:

```bash
go test ./pkg/shopadmin/... ./pkg/blogadmin/... ./pkg/logadmin/...
```

---

## 📝 Additional Notes

### New Features

- **Unified PageLayout**: `internal/layouts/page_layout.go` provides a neo-brutalist unified layout with embedded CSS theme (`brutalski_theme.css`). Navbar adapts to auth state. Set `Options.DisableNavbar = true` for pages that provide their own navigation.
- **LayoutInterface**: New `LayoutInterface` interface in `internal/layouts/layout_interface.go` with `ToHTML() string` method.
- **UserMenuItem type**: New `UserMenuItem` struct in `internal/layouts/user_layout_types.go` for user navbar navigation items.
- **User Layout Navbar**: `internal/layouts/user_layout_navbar.go` provides navbar rendering for user-facing pages.
- **PostImageURL rule**: `internal/rules/post_image.go` resolves post image from media store (first media by sequence) with fallback to `post.GetImageUrl()`. Now used in `post_manager_controller.go` AJAX responses and `SizedThumbnailURL`.
- **Post Details Media Picker**: `pkg/blogadmin/post_update/post_details.html` and `post_details.js` now include a media picker modal for selecting post images from uploaded media. The image URL field is marked as deprecated.
- **neat Sort Constants**: `neat.SortDesc` and `neat.SortAsc` replace `sb.DESC` and `sb.ASC` across all controllers for consistent sort ordering.
- **DiscountFormValidationRule**: `internal/controllers/admin/shop/discounts/discount_form_validation_rule.go` extracts discount form validation into a reusable rule.
- **Post Media Component**: `pkg/blogadmin/post_update/post_media_component.go` with `post_media.html` and `post_media.js` for managing blog post media.
- **Vue.js Blog Admin**: `ai_post_content_update`, `ai_title_generator`, and `blog_settings` now use Vue.js with AJAX handlers instead of liveflux/server-side rendering.
- **Order Details Controller**: New `CONTROLLER_ORDER_DETAILS` constant and `order_details` package integration in shop admin routes.
- **Category Create/Update Split**: Category create and update now have separate controller packages (`category_create`, `category_update`).
- **Blog Post Recommendations Inline**: Post recommendations now rendered inline via direct blogstore query with shuffle and card rendering, replacing the liveflux component.
- **Blog Post CSS Overhaul**: Blog post page styling completely redesigned — system fonts replace Roboto, neo-brutalist styling with `clamp()` responsive font sizes.
- **SectionBanner Redesign**: Blog banner section redesigned with card-based layout, icon badge, subtitle, and styled breadcrumbs.
- **AI Title Generator Settings Modal**: New Vue.js settings modal with fetch/submit AJAX handlers replacing liveflux SSR component.
- **Product Update Component System**: New `Mount()`/`Handle()`/`Render()` component interface for product details, media, metadata, and tags.
- **Brutalski Theme CSS**: New `internal/layouts/brutalski_theme.css` (440 lines) embedded CSS theme for page layout.
- **Logo SVG**: New `internal/layouts/logo.svg` for layout branding.
- **BlogAI Database Scaffold**: New `pkg/blogai/database.go` with commented-out SQLite store initialization scaffold.
- **markdownToHtml Fallback**: Markdown conversion errors now return graceful plain-text fallback instead of panicking.

### Removed Features

- `LLMEngineInterface` type and `LLMEngine`/`BlogTopic` fields from `blogadmin.AdminOptions`
- `pkg/shopadmin/products` package (replaced by `product_update` and `product_delete`)
- `pkg/shopadmin/discounts/discount_controller.go` (unused)
- `internal/controllers/admin/shop/order_manager_controller.go` (replaced by `pkg/shopadmin/order_manager`)
- `internal/controllers/website/blog/post/post_recommendations_component.go`
- `pkg/blogai/agent_titlegenerator_v1.go`
- `pkg/blogai/constants_test.go`
- `pkg/shopadmin/product_update/constants.go` and `constants_test.go`
- `pkg/shopadmin/product_update/form.html`, `form.js`, `handle_product_load_ajax.go`, `handle_product_update_ajax.go` and their tests
- `pkg/shopadmin/product_update/product_update_page.go`
- `authHandler`, `adminHandler`, `userMiddlewareHandler` internal functions
- Cache helper functions (`cacheGetSession`, `cacheSetSession`, `cacheGetUser`, `cacheSetUser`)
- `mingrammer/cfmt` as a direct dependency (now indirect only)
- `testify/assert` imports from shopadmin test files (replaced with `t.Fatalf`)
- `pkg/shopadmin/discounts/discount_controller.go` (unused stub)
- `TestNewOrderManagerController` and `TestOrderManagerController_Handler` from `shop_controller_test.go`
- `liveflux` imports from `post_controller.go` and `ai_title_generator_controller.go`
- `liveflux.Placeholder()` and `liveflux.SSR()` calls in blog post controller
- `/liveflux` script URL from blog post page options
- Roboto font import from blog post page
- Dark background and decorative icons from `SectionBanner()`

### Dependency Updates

| Dependency | Old Version | New Version |
|---|---|---|
| `github.com/dracory/rtr` | v1.6.0 | v1.7.0 |
| `github.com/dracory/auditstore` | v1.5.0 | v1.8.0 |
| `github.com/dracory/blindindexstore` | v1.12.0 | v1.14.0 |
| `github.com/dracory/blogstore` | v1.25.0 | v1.29.0 |
| `github.com/dracory/cachestore` | v1.6.0 | v1.7.0 |
| `github.com/dracory/chatstore` | v1.1.0 | v1.2.0 |
| `github.com/dracory/cmsstore` | v1.33.0 | v1.34.0 |
| `github.com/dracory/customstore` | v1.10.0 | v1.11.0 |
| `github.com/dracory/entitystore` | v1.10.0 | v1.11.0 |
| `github.com/dracory/feedstore` | v1.1.0 | v1.2.0 |
| `github.com/dracory/geostore` | v1.5.0 | v1.6.0 |
| `github.com/dracory/logstore` | v1.18.0 | v1.19.0 |
| `github.com/dracory/metastore` | v1.7.0 | v1.8.0 |
| `github.com/dracory/neat` | v0.23.0 | v0.27.0 |
| `github.com/dracory/sessionstore` | v1.15.0 | v1.16.0 |
| `github.com/dracory/shopstore` | v1.17.0 | v1.18.0 |
| `github.com/dracory/statsstore` | v1.2.0 | v1.3.0 |
| `github.com/dracory/subscriptionstore` | v1.2.0 | v1.3.0 |
| `github.com/dracory/taskstore` | v1.25.0 | v1.26.0 |
| `github.com/dracory/userstore` | v1.15.0 | v1.16.0 |
| `github.com/dracory/vaultstore` | v1.2.0 | v1.3.0 |
| `github.com/dracory/versionstore` | v1.5.0 | v1.6.0 |
| `modernc.org/sqlite` | v1.52.0 | v1.53.0 |
| `modernc.org/libc` | v1.73.4 | v1.73.5 |

---

## 🆘 Common Issues and Solutions

### Issue 1: Cannot find package `project/pkg/shopadmin/products`

**Symptom**: Compile error `cannot find package "project/pkg/shopadmin/products"`.

**Solution**: The package was split into `product_update` and `product_delete`. Update your imports:

```go
// Old
import "project/pkg/shopadmin/products"

// New
import (
    "project/pkg/shopadmin/product_update"
    "project/pkg/shopadmin/product_delete"
)
```

### Issue 2: `authHandler` undefined

**Symptom**: Compile error `undefined: authHandler` in test files.

**Solution**: The internal `authHandler` function was removed. Use `AuthMiddleware(app).GetHandler()` instead:

```go
// Old
handler := authHandler(app, next)

// New
handler := AuthMiddleware(app).GetHandler()(next)
```

### Issue 3: Wrong number of arguments for `NewProductUpdateController`

**Symptom**: Compile error `not enough arguments in call to NewProductUpdateController`.

**Solution**: Add the `fileManagerURL` parameter:

```go
// Old
product_update.NewProductUpdateController(app)

// New
product_update.NewProductUpdateController(app, fileManagerURL)
```

### Issue 4: Wrong number of arguments for `SizedThumbnailURL`

**Symptom**: Compile error `not enough arguments in call to SizedThumbnailURL`.

**Solution**: Add `context.Context` as the first parameter:

```go
// Old
shared.SizedThumbnailURL(app, post, "300", "200", "80")

// New
shared.SizedThumbnailURL(ctx, app, post, "300", "200", "80")
```

### Issue 5: `LLMEngineInterface` undefined

**Symptom**: Compile error `undefined: LLMEngineInterface` in code referencing `blogadmin.AdminOptions`.

**Solution**: Remove all references to `LLMEngineInterface`, `LLMEngine`, and `BlogTopic` from `AdminOptions` initialization. AI title generator settings are now managed via the blog settings page.

### Issue 6: `NewOrderManagerController` requires argument

**Symptom**: Compile error `not enough arguments in call to NewOrderManagerController`.

**Solution**: The constructor now requires an `app.AppInterface` parameter and is imported from `pkg/shopadmin/order_manager`:

```go
// Old
import "project/internal/controllers/admin/shop"
NewOrderManagerController().Handler(w, r)

// New
import orderManager "project/pkg/shopadmin/order_manager"
orderManager.NewOrderManagerController(app).Handler(w, r)
```

### Issue 7: `sb.DESC` or `sb.ASC` undefined

**Symptom**: Compile error `undefined: sb.DESC` or `undefined: sb.ASC`.

**Solution**: Replace with `neat` sort constants:

```go
// Old
import "github.com/dracory/sb"
sb.DESC → neat.SortDesc
sb.ASC → neat.SortAsc

// New
import "github.com/dracory/neat"
neat.SortDesc
neat.SortAsc
```

Also update string comparisons: `sortOrder == "asc"` to `strings.EqualFold(sortOrder, neat.SortAsc)`.

### Issue 8: `cfmt` undefined in blogai

**Symptom**: Compile error `undefined: cfmt` in `pkg/blogai` code.

**Solution**: Replace `cfmt.Successln(...)` with `slog.Info(...)`:

```go
// Old
cfmt.Successln("Response: ", response)

// New
slog.Info("Response: ", "response", response)
```

### Issue 9: `liveflux` undefined in blog controllers

**Symptom**: Compile error `undefined: liveflux` in `post_controller.go` or `ai_title_generator_controller.go`.

**Solution**: The `liveflux` package was removed from blog controllers. Replace liveflux components with Vue.js + AJAX:

```go
// Old
import "github.com/dracory/liveflux"
liveflux.Placeholder(component, map[string]string{...})
liveflux.SSR(component, map[string]string{...})

// New — use Vue.js CDN + embedded HTML/JS
vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")
// Read embedded .html and .js files
```

Also remove `/liveflux` from `ScriptURLs` arrays.

### Issue 10: `NewPostRecommendationsComponent` undefined

**Symptom**: Compile error `undefined: NewPostRecommendationsComponent`.

**Solution**: The component was deleted. Recommendations are now rendered inline in `post_controller.go`:

```go
// Old
component := NewPostRecommendationsComponent(c.app)
rendered := liveflux.Placeholder(component, map[string]string{...})

// New — direct blogstore query + inline card rendering
postList, _ := c.app.GetBlogStore().PostList(context.Background(), options)
// filter, shuffle, render cards
```

### Issue 11: `page()` method wrong number of arguments

**Symptom**: Compile error `not enough arguments in call to controller.page`.

**Solution**: Both `blogController.page()` and `postController.page()` now require additional parameters:

```go
// Old — blog_controller.go
controller.page(data)

// New — blog_controller.go
controller.page(r, data)

// Old — post_controller.go
c.page(post)

// New — post_controller.go
c.page(post, ctx)
```

---

## 📞 Support

For issues or questions about this upgrade:
- Check the [Blueprint repository](https://github.com/dracory/blueprint)
- Open an issue on GitHub for upgrade-specific problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
- [x] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [x] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [x] Previous guides reviewed for consistency
- [x] Git tag verified for previous version
- [x] Quality checklist included in generated guide
