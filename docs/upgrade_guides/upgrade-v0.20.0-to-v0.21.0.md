# Upgrade Guide: v0.20.0 to v0.21.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.20.0 to v0.21.0.

## Summary

Version 0.21.0 introduces several major architectural improvements including package migrations for blog and log admin interfaces, new CMS block types, a load testing tool, and various code quality improvements. The most significant changes involve moving admin controllers from `internal` to `pkg` packages for better reusability.

**Key Changes:**
- Blog admin controllers migrated from `internal/controllers/admin/blog` to `pkg/blogadmin`
- Log admin controllers migrated from `internal/controllers/admin/logs` to `pkg/logadmin`
- New CMS block types (blog post, blog post list, search)
- New command-line load testing tool (`cmd/loadtest`)
- Main server refactored with background processes moved to separate files
- Multiple dependency updates (blogstore, cmsstore, cdn, envenc, vaultstore, versionstore)
- Blog taxonomy enabled by default
- MCP debug logging removed from website routes
- Widget registry typo fix (`NewBlockeditotWidget` -> `NewBlockeditorWidget`)
- Links initialization now thread-safe with `sync.Once`

---

## Breaking Changes

### 1. Blog Admin Package Migration

**Change**: All blog admin controllers have been moved from `internal/controllers/admin/blog` to `pkg/blogadmin` to make them reusable as a standalone package.

**Old Import Paths** (v0.20.0):
```go
import (
    "project/internal/controllers/admin/blog/ai_post_content_update"
    "project/internal/controllers/admin/blog/ai_post_editor"
    "project/internal/controllers/admin/blog/ai_post_generator"
    "project/internal/controllers/admin/blog/ai_test"
    "project/internal/controllers/admin/blog/ai_title_generator"
    "project/internal/controllers/admin/blog/ai_tools"
    "project/internal/controllers/admin/blog/blog_settings"
    "project/internal/controllers/admin/blog/category_manager"
    "project/internal/controllers/admin/blog/dashboard"
    "project/internal/controllers/admin/blog/post_create"
    "project/internal/controllers/admin/blog/post_delete"
    "project/internal/controllers/admin/blog/post_manager"
    "project/internal/controllers/admin/blog/post_update"
    "project/internal/controllers/admin/blog/shared"
    "project/internal/controllers/admin/blog/tag_manager"
)
```

**New Import Paths** (v0.21.0):
```go
import (
    aIPostContentUpdate "project/pkg/blogadmin/ai_post_content_update"
    aIPostEditor "project/pkg/blogadmin/ai_post_editor"
    aIPostGenerator "project/pkg/blogadmin/ai_post_generator"
    aITest "project/pkg/blogadmin/ai_test"
    aITitleGenerator "project/pkg/blogadmin/ai_title_generator"
    aITools "project/pkg/blogadmin/ai_tools"
    "project/pkg/blogadmin/blog_settings"
    "project/pkg/blogadmin/category_manager"
    "project/pkg/blogadmin/dashboard"
    "project/pkg/blogadmin/post_create"
    "project/pkg/blogadmin/post_delete"
    "project/pkg/blogadmin/post_manager"
    "project/pkg/blogadmin/post_update"
    "project/pkg/blogadmin/shared"
    "project/pkg/blogadmin/tag_manager"
)
```

**Action Required**:
- Update all imports referencing the old blog admin controller paths
- The `pkg/blogadmin` package now provides a clean API with `blogadmin.New()` and `blogadmin.Routes()`
- The old `internal/controllers/admin/blog` directory is removed

**Migration Command**:
```bash
# Update all Go file imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/blog|pkg/blogadmin|g' {} \;
```

---

### 2. Log Admin Package Migration

**Change**: The log admin interface has been moved from `internal/controllers/admin/logs` to `pkg/logadmin` following the same pattern as blog admin.

**Old Import Paths** (v0.20.0):
```go
import (
    "project/internal/controllers/admin/logs/log_manager"
    "project/internal/controllers/admin/logs/shared"
)
```

**New Import Paths** (v0.21.0):
```go
import (
    "project/pkg/logadmin/log_manager"
    "project/pkg/logadmin/shared"
)
```

**Action Required**:
- Update any custom code that imported log admin controllers
- The `pkg/logadmin` package provides `logadmin.New()` and `logadmin.Routes()`

**Migration Command**:
```bash
# Update all Go file imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/logs|pkg/logadmin|g' {} \;
```

---

### 3. CMS Block Types Initialization

**Change**: CMS block types are now initialized automatically in `startBackgroundProcesses()`. The commented-out code is now enabled.

**Old Usage** (v0.20.0):
```go
// In cmd/server/main.go or background_processes.go
// cmsblocks.CmsAddBlockTypes(registry)    // Add CMS block types - was commented out
```

**New Usage** (v0.21.0):
```go
// In cmd/server/background_processes.go
cmsblocks.CmsAddBlockTypes(registry)    // Add CMS block types - now enabled
```

**Action Required**:
- If you had custom initialization for CMS block types, verify it doesn't conflict with the automatic initialization
- The new block types registered are: BlogPostList, BlogPost, and Search

---

### 4. Main Server Refactoring

**Change**: The `cmd/server/main.go` file has been refactored. Background process management and the `isCliMode()` function have been moved to separate files.

**Old Structure** (v0.20.0):
```
cmd/server/main.go (contained everything)
```

**New Structure** (v0.21.0):
```
cmd/server/main.go           - Simplified main entry point
cmd/server/background.go     - backgroundGroup struct and methods
cmd/server/background_processes.go - startBackgroundProcesses function
cmd/server/cli_mode.go       - CLI mode detection (new file)
```

**Action Required**:
- If you had custom modifications to `main.go` for background processes, move them to `background_processes.go`
- The `isCliMode()` function was removed; CLI mode is now handled by `cmd/server/cli_mode.go`

---

### 5. Dependency Updates

**Change**: Multiple dependencies have been updated to newer versions.

**Updated Dependencies**:
```
github.com/aws/smithy-go         v1.24.3 -> v1.25.1
github.com/dracory/blogstore     v1.10.0 -> v1.12.0
github.com/dracory/cdn           v1.10.0 -> v1.11.0
github.com/dracory/cmsstore      v1.23.0 -> v1.29.0
github.com/dracory/envenc        v1.2.0  -> v1.4.1
github.com/dracory/logstore      v1.13.0 -> v1.14.0
github.com/dracory/vaultstore    v0.37.0 -> v0.38.0
github.com/dracory/versionstore  v0.6.0  -> v0.9.0
modernc.org/sqlite               v1.48.1 -> v1.49.1
```

**New Dependency**:
```
github.com/stretchr/testify v1.11.1
```

**Action Required**:
- Run `go mod tidy` and `go mod download` after upgrading
- The testify dependency is now a direct dependency (previously may have been indirect)

---

### 6. Widget Registry Typo Fix

**Change**: Fixed typo in widget registry function name.

**Old** (v0.20.0):
```go
NewBlockeditotWidget(registry) // Typo in function name
```

**New** (v0.21.0):
```go
NewBlockeditorWidget(registry) // Corrected spelling
```

**Action Required**:
- If you have custom code referencing `NewBlockeditotWidget`, update to `NewBlockeditorWidget`

---

### 7. Blog Taxonomy Enabled by Default

**Change**: The blog store now has taxonomy enabled by default.

**Old** (v0.20.0):
```go
st, err := blogstore.NewStore(blogstore.NewStoreOptions{
    DB:              db,
    PostTableName:   "snv_blogs_post",
    TaxonomyEnabled: false,  // Taxonomy disabled
    // ...
})
```

**New** (v0.21.0):
```go
st, err := blogstore.NewStore(blogstore.NewStoreOptions{
    DB:              db,
    PostTableName:   "snv_blogs_post",
    TaxonomyEnabled: true,   // Taxonomy now enabled
    // ...
})
```

**Action Required**:
- If you don't use taxonomy features, no action needed
- If you have custom blog store initialization, verify taxonomy tables exist in your database
- The taxonomy tables (`snv_blogs_taxonomy`, `snv_blogs_term`) must exist in your database schema

---

### 8. MCP Debug Logging Removed

**Change**: Debug logging for MCP (Model Context Protocol) requests and responses has been removed from website blog and CMS routes.

**Old** (v0.20.0):
```go
// In website blog and CMS routes, there were debug log functions:
logRequest(r.Method, r.URL.Path, r.Header, bodyBytes)
logResponse(tw.status, tw.buf.Bytes())
```

**New** (v0.21.0):
```go
// Debug logging removed - cleaner production code
```

**Action Required**:
- No action required unless you relied on the debug output
- The functions `logRequest()` and `logResponse()` have been removed from `internal/controllers/website/blog/routes.go` and `internal/controllers/website/cms/routes.go`

---

### 9. Links Initialization Thread Safety

**Change**: The `links` package initialization is now thread-safe using `sync.Once`.

**Old** (v0.20.0):
```go
var initialized bool

func initializeURLBuilder() {
    if !initialized {
        appURL := os.Getenv("APP_URL")
        // ...
        baseurl.SetDefaultURL(appURL)
        initialized = true
    }
}
```

**New** (v0.21.0):
```go
var initOnce sync.Once

func initializeURLBuilder() {
    initOnce.Do(func() {
        appURL := os.Getenv("APP_URL")
        // ...
        baseurl.SetDefaultURL(appURL)
    })
}
```

**Action Required**:
- No action required - this is an internal improvement
- If you had custom initialization logic around the `initialized` variable, update accordingly

---

### 10. Thumb Controller Input Sanitization

**Change**: The thumbnail controller now sanitizes input parameters by trimming whitespace and normalizing extension case.

**Old** (v0.20.0):
```go
data.extension, _ = rtr.GetParam(r, "extension")
size, _ := rtr.GetParam(r, "size")
quality, _ := rtr.GetParam(r, "quality")
data.path, _ = rtr.GetParam(r, "path")
```

**New** (v0.21.0):
```go
data.extension, _ = rtr.GetParam(r, "extension")
data.extension = strings.TrimSpace(strings.ToLower(data.extension))
size, _ := rtr.GetParam(r, "size")
size = strings.TrimSpace(size)
quality, _ := rtr.GetParam(r, "quality")
quality = strings.TrimSpace(quality)
data.path, _ = rtr.GetParam(r, "path")
data.path = strings.TrimSpace(data.path)
```

**Action Required**:
- No action required - this is an internal improvement
- URLs with uppercase extensions or extra whitespace will now work correctly

---

## Migration Steps

### Step 1: Update Go Module Dependencies

```bash
# Update all dependencies to match go.mod
go mod tidy
go mod download
```

### Step 2: Update Import Paths

Update any custom code that imports from the old blog or log admin paths:

```bash
# Update blog admin imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/blog|pkg/blogadmin|g' {} \;

# Update log admin imports
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/logs|pkg/logadmin|g' {} \;

# Fix widget typo if referenced directly
find . -type f -name "*.go" -exec sed -i 's|NewBlockeditotWidget|NewBlockeditorWidget|g' {} \;
```

### Step 3: Verify CMS Block Types

If you have custom CMS block type initialization, ensure it doesn't conflict with the new automatic initialization in `cmd/server/background_processes.go`.

```go
// In your custom code, you can now rely on the built-in initialization:
cmsblocks.CmsAddBlockTypes(registry) // This is called automatically
```

### Step 4: Clean Up Old Files

Remove old files that have been replaced:

```bash
# Remove old blog admin directory (if it still exists after merge)
rm -rf internal/controllers/admin/blog/

# Remove old log admin directory
rm -rf internal/controllers/admin/logs/
```

### Step 5: Verify Build

```bash
go build -o ./tmp/main ./cmd/server
```

---

## Testing After Migration

### 1. Unit Tests

Run the full test suite:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 2. Integration Tests

Test the blog admin interface:

```bash
# Start the application and test blog routes
go run ./cmd/server
# Visit: http://localhost:8080/admin/blog
```

Test the log admin interface:

```bash
# Visit: http://localhost:8080/admin/logs
```

### 3. CMS Block Types Test

Verify CMS block types are registered:

```bash
# Check that block types are available in the CMS editor
# Blog Post List, Blog Post, and Search blocks should be available
```

### 4. Load Testing Tool

Test the new load testing command:

```bash
# Run the load testing tool
go run ./cmd/loadtest -url http://localhost:8080 -c 10 -d 30s

# With rate limiting
go run ./cmd/loadtest -url http://localhost:8080 -c 10 -d 30s -r 5
```

### 5. Database Migration Check

If using blog taxonomy:

```bash
# Ensure taxonomy tables exist
# snv_blogs_taxonomy and snv_blogs_term tables must be present
```

---

## Additional Notes

### New Features

1. **Load Testing Tool** (`cmd/loadtest`): A command-line load testing tool with configurable concurrency, duration, timeout, and rate limiting.

2. **CMS Block Types**: Three new built-in CMS block types:
   - Blog Post List Block - displays a list of blog posts
   - Blog Post Block - displays a single blog post
   - Search Block - provides search functionality

3. **Package Migrations**: Blog admin and log admin packages are now in `pkg/`, making them reusable as standalone packages.

4. **Thread-Safe Initialization**: The links package initialization is now thread-safe.

5. **Input Sanitization**: Thumbnail controller now properly sanitizes input parameters.

### Removed Features

- `internal/controllers/admin/blog/` - moved to `pkg/blogadmin/`
- `internal/controllers/admin/logs/log_manager/` - moved to `pkg/logadmin/log_manager/`
- `isCliMode()` function from `main.go` - moved to `cmd/server/cli_mode.go`
- `startBackgroundProcesses()` from `main.go` - moved to `cmd/server/background_processes.go`
- `backgroundGroup` struct from `main.go` - moved to `cmd/server/background.go`
- MCP debug logging functions from website routes

### Configuration Behavior Changes

- **Blog Taxonomy**: Now enabled by default. Ensure database schema supports taxonomy tables.
- **CMS Block Types**: Now automatically registered at startup if CMS store is enabled.
- **Connection Pool**: SQLite automatically uses pool=1, other databases use configured values.

---

## Common Issues and Solutions

### Issue 1: Import Path Errors

**Symptom**: Compilation errors for blog or log admin imports

**Solution**: Update all imports from `internal/controllers/admin/blog/*` to `pkg/blogadmin/*` and from `internal/controllers/admin/logs/*` to `pkg/logadmin/*`

```bash
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/blog|pkg/blogadmin|g' {} \;
find . -type f -name "*.go" -exec sed -i 's|internal/controllers/admin/logs|pkg/logadmin|g' {} \;
```

### Issue 2: Blog Taxonomy Tables Missing

**Symptom**: Database errors when using blog store

**Solution**: Create the required taxonomy tables or disable taxonomy if not needed:

```sql
-- Create taxonomy tables if they don't exist
CREATE TABLE IF NOT EXISTS snv_blogs_taxonomy (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS snv_blogs_term (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL
);
```

### Issue 3: Build Failures After Refactoring

**Symptom**: Cannot find `isCliMode()` or `startBackgroundProcesses()` functions

**Solution**: These functions are now in separate files in `cmd/server/`. The imports are automatically resolved since they're in the same package. No code changes needed unless you had custom implementations in `main.go`.

### Issue 4: Widget Registry Typo

**Symptom**: `undefined: NewBlockeditotWidget`

**Solution**: Replace with `NewBlockeditorWidget`:

```bash
find . -type f -name "*.go" -exec sed -i 's|NewBlockeditotWidget|NewBlockeditorWidget|g' {} \;
```

---

## Support

For issues related to this upgrade:

1. **Documentation**: Review the package documentation in `pkg/blogadmin/README.md` and `pkg/logadmin/README.md`

2. **Git History**: Review the migration commits:
   ```bash
   git log --oneline v0.20.0..v0.21.0
   ```

3. **Reference Implementation**: Compare with the reference implementation in the v0.21.0 tag

4. **Issues**: Report issues at https://github.com/dracory/blueprint/issues

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices

---

*This upgrade guide was generated for Blueprint v0.20.0 to v0.21.0 migration.*
