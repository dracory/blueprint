# Upgrade Guide: v0.22.0 to v0.23.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.22.0 to v0.23.0.

## Overview

This release introduces a version tracking system, significant blogadmin UI refactoring (Go components → JavaScript), security middleware improvements, and dependency upgrades.

## ⚠️ Breaking Changes

### 1. Version Tracking System Added
**Change**: Added version constant and tracking system to `internal/config/version.go`

**Old Usage**:
```go
// No version tracking existed
```

**New Usage**:
```go
import "project/internal/config"

version := config.GetVersion() // Returns "0.23.0"
```

**Action Required**:
- No action required for existing applications
- The version constant is automatically available
- This enables automated upgrade detection and migration tools

---

### 2. Security Headers Middleware Signature Changed
**Change**: `getScriptSources()` function no longer takes `isDevelopment` parameter; environment detection moved to `NewSecurityHeadersMiddleware()`

**Old Usage**:
```go
func getScriptSources(isDevelopment bool) []string {
    sources := []string{
        "'self'",
        "https://cdn.jsdelivr.net",
        // ...
    }
    return sources
}
```

**New Usage**:
```go
func getScriptSources() []string {
    sources := []string{
        "'self'",
        "https://cdn.jsdelivr.net",
        "https://cdn.tiny.cloud", // NEW: TinyCloud CDN added
        // ...
    }
    return sources
}

// Environment detection moved to NewSecurityHeadersMiddleware
func NewSecurityHeadersMiddleware(registry registry.RegistryInterface) rtr.MiddlewareInterface {
    isDevelopment := false
    if registry != nil {
        isDevelopment = registry.GetConfig().IsEnvDevelopment() || registry.GetConfig().IsEnvLocal()
    }
    // ... rest of implementation
}
```

**Action Required**:
- If you have custom implementations of `getScriptSources()`, remove the `isDevelopment` parameter
- Update any calls to `getScriptSources()` to not pass the parameter
- If you override security headers middleware, ensure environment detection uses registry config methods

---

### 3. Blogadmin UI Architecture Refactored (Go → JavaScript)
**Change**: Post update UI components migrated from Go to JavaScript with HTMX

**Files Removed**:
- `pkg/blogadmin/post_update/post_content_component.go`
- `pkg/blogadmin/post_update/post_content_component_test.go`
- `pkg/blogadmin/post_update/post_details_component.go`
- `pkg/blogadmin/post_update/post_details_component_test.go`
- `pkg/blogadmin/post_update/post_seo_component.go`
- `pkg/blogadmin/post_update/post_seo_component_test.go`
- `pkg/blogadmin/post_update/post_versioning_component.go`
- `pkg/blogadmin/post_update/post_versioning_test.go`

**Files Added**:
- `pkg/blogadmin/post_update/post_content.html`
- `pkg/blogadmin/post_update/post_content.js`
- `pkg/blogadmin/post_update/post_details.html`
- `pkg/blogadmin/post_update/post_details.js`
- `pkg/blogadmin/post_update/post_seo.html`
- `pkg/blogadmin/post_update/post_seo.js`
- `pkg/blogadmin/post_update/post_versioning.html`
- `pkg/blogadmin/post_update/post_versioning.js`

**Old Usage**:
```go
// Component was rendered server-side with Go
func (c *PostContentComponent) Render() string {
    return hb.Render(c.template, c.data)
}
```

**New Usage**:
```go
// Controller now handles API actions for JavaScript components
func (controller *postUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
    action := req.GetStringTrimmed(r, "action")
    
    switch action {
    case "load-content":
        return controller.handleLoadContent(w, r)
    case "save-content":
        return controller.handleSaveContent(w, r)
    // ... other actions
    }
    
    // Render view based on query parameter
    view := req.GetStringTrimmedOr(r, "view", "content")
    return controller.renderView(view, postID)
}
```

**Action Required**:
- If you have custom blogadmin components, migrate them to JavaScript/HTMX pattern
- Update any tests that relied on Go component rendering
- Review `pkg/blogadmin/post_update/post_update_controller.go` for new action handlers
- Ensure your templates use the new HTML/JS component structure
- **Note**: This change only affects blogadmin package; if you don't customize blogadmin, no action needed

---

### 4. Blogadmin Post Manager UI Reorganized
**Change**: Post manager UI reorganized with separate edit and view buttons

**Old Usage**:
```html
<!-- Single button for edit, magic button for AI -->
<a :href="getPostUpdateUrl(post.id)" target="_blank" class="btn btn-primary btn-sm">
  <i class="bi bi-pencil-square"></i>
</a>
<a @click="generateContent(post)" class="btn btn-info btn-sm">
  <i class="bi bi-magic"></i>
</a>
```

**New Usage**:
```html
<!-- Separate edit and view buttons, magic button moved -->
<a @click="generateContent(post)" class="btn btn-info btn-sm">
  <i class="bi bi-magic"></i>
</a>
<a :href="getPostUpdateUrl(post.id)" target="_blank" class="btn btn-primary btn-sm">
  <i class="bi bi-pencil-square"></i>
</a>
<a :href="getWebsitePostUrl(post.id, post.slug)" target="_blank" class="btn btn-info btn-sm">
  <i class="bi bi-eye"></i>
</a>
```

**Action Required**:
- No action required unless you customized the post manager UI
- The post title now links to edit page instead of website
- New "View" button added to view post on website
- Button order changed: Magic → Edit → View → Delete

---

### 5. Global Middlewares Rate Limits Adjusted
**Change**: Rate limit constants added and environment-based selection improved

**Old Usage**:
```go
// Hardcoded rate limits in globalMiddlewares()
rtrMiddleware.RateLimitByIPMiddleware(20, 1),      // per second
rtrMiddleware.RateLimitByIPMiddleware(120, 1*60),   // per minute
rtrMiddleware.RateLimitByIPMiddleware(12000, 60*60), // per hour
```

**New Usage**:
```go
// Constants defined at package level
const (
    defaultVisitsPerSec  = 20
    defaultVisitsPerMin  = 120
    defaultVisitsPerHour = 12000
    
    devVisitsPerSec  = 1000
    devVisitsPerMin  = 10000
    devVisitsPerHour = 100000
)

// Environment-based selection
perSec, perMin, perHour := getRateLimits(registry)
rtrMiddleware.RateLimitByIPMiddleware(perSec, 1),
rtrMiddleware.RateLimitByIPMiddleware(perMin, 1*60),
rtrMiddleware.RateLimitByIPMiddleware(perHour, 60*60),
```

**Action Required**:
- No action required for standard Blueprint applications
- If you have custom rate limit middleware, consider adopting the constant pattern
- Development environments now have much higher rate limits (1000x) for easier testing

---

### 6. HTTPS Redirect Middleware Added
**Change**: HTTPS redirect middleware added to global middlewares with webhook skip

**Old Usage**:
```go
// No automatic HTTPS redirect in global middlewares
```

**New Usage**:
```go
// Automatically added in production (not dev/local/testing)
if isNotTesting && isNotDevelopment && isNotLocal {
    globalMiddlewares = append(globalMiddlewares,
        httpsredirect.NewHTTPSRedirectMiddlewareWithConfig(httpsredirect.Config{
            SkipFunc: func(r *http.Request) bool {
                return r.URL.Path == "/api/internal/webhook"
            },
        }),
    )
}
```

**Action Required**:
- No action required for production deployments
- Webhook endpoint `/api/internal/webhook` is automatically excluded from HTTPS redirect
- If you have other webhook endpoints, add them to the SkipFunc
- If you don't want HTTPS redirect, remove this middleware from `internal/routes/global_middlewares.go`

---

### 7. Admin User Constants Removed
**Change**: Removed unused admin user route constants from `internal/links/constants.go`

**Old Usage**:
```go
import "project/internal/links"

// These constants existed but are now removed
links.ADMIN_USERS_USER_CREATE
links.ADMIN_USERS_USER_DELETE
links.ADMIN_USERS_USER_MANAGER
links.ADMIN_USERS_USER_UPDATE
```

**New Usage**:
```go
// Constants removed - use dynamic routing or registry methods
// ADMIN_USERS_USER_IMPERSONATE is still available
links.ADMIN_USERS_USER_IMPERSONATE
```

**Action Required**:
- Search your codebase for references to these removed constants:
  - `ADMIN_USERS_USER_CREATE`
  - `ADMIN_USERS_USER_DELETE`
  - `ADMIN_USERS_USER_MANAGER`
  - `ADMIN_USERS_USER_UPDATE`
- Replace with dynamic route generation or use `ADMIN_USERS` base constant
- Update any hardcoded URLs that used these constants

---

### 8. Session Store Function Signature Changed
**Change**: `newSessionStore()` function now requires registry parameter and uses environment-based session timeout

**Old Usage**:
```go
// Session store initialization
store, err := newSessionStore(db)
```

**New Usage**:
```go
// Session store initialization with registry parameter
store, err := newSessionStore(db, registry)

// Session timeout is now environment-based:
// - Production: 7200 seconds (2 hours)
// - Development: 14400 seconds (4 hours)
```

**Action Required**:
- If you call `newSessionStore()` directly, add the registry parameter
- If you have custom session timeout logic, review the new environment-based implementation
- Test session behavior in both development and production environments

---

### 9. Dependency Upgrades
**Change**: Multiple dependencies upgraded to latest versions

**Upgraded Dependencies**:
- `github.com/dracory/blogstore`: v1.12.0 → v1.18.0
- `github.com/flosch/pongo2/v6`: v6.0.0 → v6.1.0
- `github.com/go-sql-driver/mysql`: v1.9.3 → v1.10.0
- AWS SDK v2: v1.41.6 → v1.41.7
- Google Cloud SDK: Various minor version bumps
- `google.golang.org/genai`: v1.54.0 → v1.56.0
- `google.golang.org/grpc`: v1.80.0 → v1.81.0
- `github.com/mattn/go-isatty`: v0.0.21 → v0.0.22
- `modernc.org/libc`: v1.72.1 → v1.72.2

**Action Required**:
- Run `go mod download` to fetch updated dependencies
- Run `go mod tidy` to clean up go.mod and go.sum
- Review blogstore v1.18.0 changelog for any API changes
- Test your application thoroughly after dependency updates

---

## 🔄 Migration Steps

### Step 1: Update Version Constant
```bash
# Update internal/config/version.go to v0.23.0
# This should already be done if you're on the release branch
```

### Step 2: Update Dependencies
```bash
# Download updated dependencies
go mod download

# Clean up go.mod and go.sum
go mod tidy

# Verify build
go build ./...
```

### Step 3: Review Security Middleware Changes
```bash
# Check if you have custom security headers middleware
grep -r "getScriptSources" internal/

# If found, remove isDevelopment parameter from function signature
# and ensure environment detection uses registry config
```

### Step 4: Review Blogadmin Customizations
```bash
# Check if you have custom blogadmin components
ls -la pkg/blogadmin/post_update/*_component.go

# If you have custom components, migrate them to JavaScript/HTMX pattern
# Reference the new .html and .js files for the pattern
```

### Step 5: Search for Removed Constants
```bash
# Search for references to removed admin user constants
grep -r "ADMIN_USERS_USER_CREATE\|ADMIN_USERS_USER_DELETE\|ADMIN_USERS_USER_MANAGER\|ADMIN_USERS_USER_UPDATE" .

# Update any found references to use dynamic routing or base constant
```

### Step 6: Check Session Store Initialization
```bash
# Search for direct calls to newSessionStore
grep -r "newSessionStore" .

# If found, update to pass registry parameter:
# Old: newSessionStore(db)
# New: newSessionStore(db, registry)
```

### Step 7: Update HTTPS Redirect (if needed)
```bash
# If you need to exclude additional webhook endpoints
# Edit internal/routes/global_middlewares.go
# Modify the SkipFunc in httpsredirect.Config
```

### Step 8: Test Application
```bash
# Run unit tests
go test ./...

# Run integration tests (if applicable)
go test -tags=integration ./...

# Start application and test manually
go run ./cmd/server
```

---

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 2. Integration Tests
```bash
# Run integration tests
go test -tags=integration ./...
```

### 3. Manual Testing
- **Security Headers**: Verify CSP headers include TinyCloud CDN
- **Blogadmin**: Test post creation, editing, and all post update tabs (content, details, SEO, versioning)
- **HTTPS Redirect**: Test that HTTP requests redirect to HTTPS in production
- **Rate Limits**: Verify rate limiting works in both development and production
- **Webhook**: Ensure webhook endpoint bypasses HTTPS redirect

### 4. Dependency Verification
- Verify blogstore API calls still work correctly
- Test database connections with updated mysql driver
- Verify template rendering with updated pongo2

---

## 📝 Additional Notes

### New Features
- **Version Tracking**: Added `config.Version` constant and `config.GetVersion()` function for automated upgrade detection
- **TinyCloud CDN**: Added to CSP script sources for TinyMCE editor
- **HTTPS Redirect**: Automatic HTTPS redirect in production with webhook exclusion
- **Improved Rate Limits**: Development environments have 1000x higher rate limits for easier testing
- **Version Logging**: Application now logs version on startup in cmd/server/main.go
- **Session Timeout**: Environment-based session timeouts (2 hours production, 4 hours development)

### Removed Features
- **Go Components**: Blogadmin post update components migrated to JavaScript/HTMX
- **Admin User Constants**: Removed unused route constants (ADMIN_USERS_USER_CREATE, ADMIN_USERS_USER_DELETE, ADMIN_USERS_USER_MANAGER, ADMIN_USERS_USER_UPDATE)

### Other Changes
- **JavaScript App Mounting**: Changed from DOMContentLoaded event to immediate mounting for blogadmin components
- **Config Constants**: Fixed section comment placement for Mail Configurations
- **Blogadmin UI**: Post manager UI reorganized with separate edit and view buttons

### Documentation
- Added version workflow documentation in `docs/version_workflow.md`
- Added upgrade guide generator workflow in `.windsurf/workflows/upgrade-guide-generator.md`

---

## 🆘 Common Issues and Solutions

### Issue: Build fails after dependency update
**Solution**: Run `go mod tidy` to resolve dependency conflicts

### Issue: Blogadmin post update tabs not loading
**Solution**: Ensure the new JavaScript files are properly embedded and served. Check browser console for JavaScript errors.

### Issue: CSP errors with TinyMCE editor
**Solution**: Verify `https://cdn.tiny.cloud` is in your CSP script sources. The upgrade should have added this automatically.

### Issue: Webhook endpoint not working with HTTPS redirect
**Solution**: Add your webhook path to the SkipFunc in `internal/routes/global_middlewares.go`

### Issue: Rate limits too strict in development
**Solution**: Development environments should automatically use higher limits (1000x). Verify `IsEnvDevelopment()` or `IsEnvLocal()` returns true.

---

## 📞 Support

For additional help:
- Review the version workflow: `docs/version_workflow.md`
- Check existing upgrade guides in `docs/upgrade_guides/`
- Review the upgrade guide generator workflow: `.windsurf/workflows/upgrade-guide-generator.md`
