# Upgrade Guide: v0.17.0 to v0.18.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.17.0 to v0.18.0.

**Key Statistics:**
- 127 files changed
- 8,673 insertions, 1,978 deletions
- Focus: Interface migrations, type safety improvements, security enhancements

---

## ⚠️ Breaking Changes

### 1. Blog Store Interface Migration

**Change**: The blog store `Post` struct has been converted to a `PostInterface` with getter methods. All direct property access must use getter methods.

**Old Usage**:
```go
import "github.com/dracory/store/blog"

func renderPost(post blogstore.Post) string {
    title := post.Title()
    id := post.ID()
    slug := post.Slug()
    content := post.Content()
    summary := post.Summary()
    imageURL := post.ImageUrlOrDefault()
    publishedAt := post.PublishedAtCarbon()
    editor := post.Editor()
    
    // Usage in templates/widgets
    return fmt.Sprintf("<h1>%s</h1>", post.Title())
}

// Type signatures
func processPosts(posts []blogstore.Post) {
    for _, post := range posts {
        // ...
    }
}
```

**New Usage**:
```go
import "github.com/dracory/store/blog"

func renderPost(post blogstore.PostInterface) string {
    title := post.GetTitle()
    id := post.GetID()
    slug := post.GetSlug()
    content := post.GetContent()
    summary := post.GetSummary()
    imageURL := post.GetImageUrlOrDefault()
    publishedAt := post.GetPublishedAtCarbon()
    editor := post.GetEditor()
    
    // Usage in templates/widgets
    return fmt.Sprintf("<h1>%s</h1>", post.GetTitle())
}

// Type signatures
func processPosts(posts []blogstore.PostInterface) {
    for _, post := range posts {
        // ...
    }
}
```

**Action Required**:
- Update all function signatures from `blogstore.Post` to `blogstore.PostInterface`
- Replace `post.Title()` with `post.GetTitle()`
- Replace `post.ID()` with `post.GetID()`
- Replace `post.Slug()` with `post.GetSlug()`
- Replace `post.Content()` with `post.GetContent()`
- Replace `post.Summary()` with `post.GetSummary()`
- Replace `post.ImageUrlOrDefault()` with `post.GetImageUrlOrDefault()`
- Replace `post.PublishedAt()` with `post.GetPublishedAt()`
- Replace `post.PublishedAtCarbon()` with `post.GetPublishedAtCarbon()`
- Replace `post.Editor()` with `post.GetEditor()`
- Replace `post.ContentType()` with `post.GetContentType()`
- Replace `post.CanonicalURL()` with `post.GetCanonicalURL()`
- Replace `post.Status()` with `post.GetStatus()`
- Replace `post.Featured()` with `post.GetFeatured()`
- Replace `post.ImageUrl()` with `post.GetImageUrl()`
- Replace `post.MetaDescription()` with `post.GetMetaDescription()`
- Replace `post.MetaKeywords()` with `post.GetMetaKeywords()`
- Replace `post.MetaRobots()` with `post.GetMetaRobots()`
- Replace `post.Memo()` with `post.GetMemo()`
- Replace `post.Website()` with `post.GetWebsite()`

**Automated Update Command**:
```bash
# Update type signatures
find . -name "*.go" -exec sed -i 's/blogstore\.Post/blogstore.PostInterface/g' {} \;

# Update getter methods (run in order)
find . -name "*.go" -exec sed -i 's/\.Title()/.GetTitle()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ID()/.GetID()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Slug()/.GetSlug()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Content()/.GetContent()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Summary()/.GetSummary()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ImageUrlOrDefault()/.GetImageUrlOrDefault()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.PublishedAt()/.GetPublishedAt()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.PublishedAtCarbon()/.GetPublishedAtCarbon()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Editor()/.GetEditor()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ContentType()/.GetContentType()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.CanonicalURL()/.GetCanonicalURL()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Status()/.GetStatus()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Featured()/.GetFeatured()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ImageUrl()/.GetImageUrl()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaDescription()/.GetMetaDescription()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaKeywords()/.GetMetaKeywords()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaRobots()/.GetMetaRobots()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Memo()/.GetMemo()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Website()/.GetWebsite()/g' {} \;
```

---

### 2. Entity Store Interface Changes

**Change**: The entity store has changed from concrete `entitystore.Entity` type to `entitystore.EntityInterface`, and method signatures have changed.

**Old Usage**:
```go
import "github.com/dracory/store/entity"

func processEntity(entity entitystore.Entity) error {
    entityType := entity.Type()
    entityID := entity.ID
    
    attributes, err := entity.GetAttributes()
    for _, attr := range attributes {
        key := attr.AttributeKey()
        value := attr.AttributeValue()
        // ...
    }
}
```

**New Usage**:
```go
import "github.com/dracory/store/entity"

func processEntity(store entitystore.StoreInterface, entity entitystore.EntityInterface) error {
    entityType := entity.GetType()
    entityID := entity.GetID()
    createdAt := entity.GetCreatedAt()
    updatedAt := entity.GetUpdatedAt()
    
    // Fetch attributes from store
    attributes, err := store.EntityAttributeList(context.Background(), entity.GetID())
    for _, attr := range attributes {
        key := attr.GetKey()
        value := attr.GetValue()
        // ...
    }
}
```

**Action Required**:
- Update `entity.Type()` to `entity.GetType()`
- Update `entity.ID` to `entity.GetID()`
- Update `entity.GetAttributes()` to `store.EntityAttributeList(ctx, entity.GetID())`
- Update `attr.AttributeKey()` to `attr.GetKey()`
- Update `attr.AttributeValue()` to `attr.GetValue()`
- Add `store entitystore.StoreInterface` parameter where needed
- Use `entitystore.COLUMN_CREATED_AT` and `entitystore.COLUMN_UPDATED_AT` constants

---

### 3. Testimonials Package API Changes

**Change**: The testimonials package has been updated to work with the new entity store interface.

**Old Usage**:
```go
import "project/pkg/testimonials"

// Creating testimonial from entity
testimonial, err := testimonials.NewTestimonialFromEntity(entity)

// Listing testimonials
testimonialList, err := testimonials.TestimonialList(store)
for _, entry := range result {
    testimonial, err := testimonials.NewTestimonialFromEntity(entry)
    // ...
}
```

**New Usage**:
```go
import "project/pkg/testimonials"

// Creating testimonial from entity - now requires store parameter
testimonial, err := testimonials.NewTestimonialFromEntity(store, entity)

// Listing testimonials
testimonialList, err := testimonials.TestimonialList(store)
for _, entry := range result {
    testimonial, err := testimonials.NewTestimonialFromEntity(store, entry)
    // ...
}
```

**Action Required**:
- Pass `store entitystore.StoreInterface` as first parameter to `NewTestimonialFromEntity()`

---

### 4. Thumb Controller File Storage Changes

**Change**: The thumb controller now uses SQL file storage for local files instead of generating URLs.

**Old Usage**:
```go
// Files with "files/" prefix were converted to URLs
if strings.HasPrefix(data.path, "files/") {
    data.path = links.URL(data.path, nil)
    data.isURL = true
}
```

**New Usage**:
```go
// Files with "files/" prefix now use SQL file storage
if strings.HasPrefix(data.path, "files/") {
    data.isFiles = true
    data.path = strings.TrimPrefix(data.path, "files/")
}

// Later in generateThumb:
if data.isFiles {
    storage := controller.registry.GetSqlFileStorage()
    if storage == nil {
        return "", "file storage not initialized"
    }
    
    exists, err := storage.Exists(data.path)
    if err != nil || !exists {
        return "", "file not found"
    }
    
    imgBytes, err = storage.ReadFile(data.path)
    // ...
}
```

**Action Required**:
- No changes required for existing usage, but ensure `GetSqlFileStorage()` is properly initialized in your registry

---

### 5. Blog Store Configuration Changes

**Change**: The blog store configuration now includes taxonomy support fields.

**Old Usage**:
```go
st, err := blogstore.NewStore(blogstore.NewStoreOptions{
    DB:                  db,
    PostTableName:       "snv_blogs_post",
    VersioningEnabled:   true,
    VersioningTableName: "snv_blogs_version",
    AutomigrateEnabled:  true,
})
```

**New Usage**:
```go
st, err := blogstore.NewStore(blogstore.NewStoreOptions{
    DB:                  db,
    PostTableName:       "snv_blogs_post",
    TaxonomyEnabled:     false,        // NEW
    TaxonomyTableName:   "snv_blogs_taxonomy", // NEW
    TermTableName:       "snv_blogs_term",     // NEW
    VersioningEnabled:   true,
    VersioningTableName: "snv_blogs_version",
    AutomigrateEnabled:  true,
})
```

**Action Required**:
- Add `TaxonomyEnabled`, `TaxonomyTableName`, and `TermTableName` fields to your blog store initialization
- Set `TaxonomyEnabled: false` if you don't need taxonomy features

---

### 6. Rate Limiting Environment Awareness

**Change**: Rate limits are now environment-aware with different limits for development vs production.

**Old Usage**:
```go
// Fixed rate limits regardless of environment
rtrMiddleware.RateLimitByIPMiddleware(20, 1),        // 20 req per second
rtrMiddleware.RateLimitByIPMiddleware(180, 1*60),    // 180 req per minute
rtrMiddleware.RateLimitByIPMiddleware(12000, 60*60), // 12000 req per hour
```

**New Usage**:
```go
// Environment-aware rate limits
perSec, perMin, perHour := getRateLimits(registry)

rtrMiddleware.RateLimitByIPMiddleware(perSec, 1),      // per second
rtrMiddleware.RateLimitByIPMiddleware(perMin, 1*60), // per minute
rtrMiddleware.RateLimitByIPMiddleware(perHour, 60*60), // per hour

// Helper function
func getRateLimits(registry registry.RegistryInterface) (perSec, perMin, perHour int) {
    if registry.GetConfig() != nil {
        isDevelopment := registry.GetConfig().IsEnvDevelopment()
        isLocal := registry.GetConfig().IsEnvLocal()
        
        if isDevelopment || isLocal {
            return 1000, 10000, 100000 // Relaxed limits for dev
        }
    }
    return 20, 10, 12000 // Production limits
}
```

**Action Required**:
- Update middleware configuration to use environment-aware rate limits
- Add the `getRateLimits()` helper function to your middleware setup

---

### 7. Jail Bots Middleware Configuration Changes

**Change**: Additional exclude paths have been added to the jail bots middleware.

**Old Usage**:
```go
rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
    Exclude:      []string{"/new"},
    ExcludePaths: []string{"/blog*", "/th*", "/liveflux*"},
}),
```

**New Usage**:
```go
rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
    Exclude: []string{"/new"},
    ExcludePaths: []string{
        "/blog*",
        "/th*",
        "/liveflux*",
        "/admin/cms*",      // NEW
        "/admin/*cms*",     // NEW
        "/assets*",         // NEW
        "*/assets/*",       // NEW
        "/files/*",         // NEW
    },
}),
```

**Action Required**:
- Add new exclude paths if you want to allow CMS and asset routes without bot protection

---

### 8. Security Headers CSP Updates

**Change**: Content Security Policy headers have been updated with additional sources.

**Old CSP Configuration**:
```go
csp := middlewares.CSPDirectives{
    DefaultSrc: []string{"'self'"},
    ScriptSrc:  getScriptSources(isDevelopment),
    StyleSrc:   getStyleSources(isDevelopment),
    FontSrc: []string{
        "'self'",
        "https://cdn.jsdelivr.net",
        "https://fonts.googleapis.com",
        "https://fonts.gstatic.com",
        "https://cdnjs.cloudflare.com",
        "https://maxcdn.bootstrapcdn.com",
    },
    ImgSrc: []string{
        "'self'",
        "data:",
        "https://*",
        "http://*",
    },
}
```

**New CSP Configuration**:
```go
csp := middlewares.CSPDirectives{
    DefaultSrc: []string{"'self'"},
    ScriptSrc:  getScriptSources(isDevelopment),
    StyleSrc:   getStyleSources(isDevelopment),
    ConnectSrc: []string{              // NEW SECTION
        "'self'",
        "https://cdnjs.cloudflare.com",
        "http://cdnjs.cloudflare.com",
    },
    FontSrc: []string{
        "'self'",
        "https://cdn.jsdelivr.net",
        "https://fonts.googleapis.com",
        "https://fonts.gstatic.com",
        "https://cdnjs.cloudflare.com",
        "http://cdnjs.cloudflare.com",     // NEW
        "https://maxcdn.bootstrapcdn.com",
    },
    ImgSrc: []string{
        "'self'",
        "data:",
        "https://*",
        "http://*",
    },
}

// Script sources also updated
func getScriptSources(isDevelopment bool) []string {
    sources := []string{
        // ... existing sources
        "https://cdnjs.cloudflare.com",
        "http://cdnjs.cloudflare.com",     // NEW
        // ...
    }
}

// Style sources also updated
func getStyleSources(isDevelopment bool) []string {
    sources := []string{
        // ... existing sources
        "https://cdnjs.cloudflare.com",
        "http://cdnjs.cloudflare.com",     // NEW
        "https://code.jquery.com",         // NEW
        // ...
    }
}
```

**Action Required**:
- Add `ConnectSrc` directive to your CSP configuration
- Add `http://cdnjs.cloudflare.com` to FontSrc, ScriptSrc, and StyleSrc
- Add `https://code.jquery.com` to StyleSrc

---

### 9. Test Dependency Changes

**Change**: Tests have been updated to remove dependencies on `testify/assert` and `testify/require` in favor of standard library testing.

**Old Usage (testify/require)**:
```go
import "github.com/stretchr/testify/require"

func TestSomething(t *testing.T) {
    result, err := someOperation()
    require.NoError(t, err)
    require.NotNil(t, result)
    require.Equal(t, expected, result)
}
```

**Old Usage (testify/assert)**:
```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result, err := someOperation()
    assert.False(t, nextCalled, "next handler should not be called")
    assert.Equal(t, expected, actual)
    assert.True(t, condition)
}
```

**New Usage**:
```go
func TestSomething(t *testing.T) {
    result, err := someOperation()
    if err != nil {
        t.Fatalf("expected no error, got: %v", err)
    }
    if result == nil {
        t.Fatal("result should not be nil")
    }
    if result != expected {
        t.Fatalf("expected %v, got %v", expected, result)
    }
    if nextCalled {
        t.Error("next handler should not be called")
    }
}
```

**Action Required**:
- Replace `require.NoError(t, err)` with `if err != nil { t.Fatalf(...) }`
- Replace `require.NotNil(t, x)` with `if x == nil { t.Fatal(...) }`
- Replace `require.Equal(t, expected, actual)` with `if actual != expected { t.Fatalf(...) }`
- Replace `assert.Equal(t, expected, actual)` with `if actual != expected { t.Errorf(...) }`
- Replace `assert.False(t, condition, msg)` with `if condition { t.Error(msg) }`
- Replace `assert.True(t, condition, msg)` with `if !condition { t.Error(msg) }`
- Remove `github.com/stretchr/testify/require` imports where no longer needed
- Remove `github.com/stretchr/testify/assert` imports where no longer needed

---

### 10. Shop Admin Product Update Controller Moved

**Change**: The product update controller in the shop admin has been moved from `products` package to a new `productupdate` subdirectory package.

**Old Usage**:
```go
import "project/internal/controllers/admin/shop/products"

// In routes or handlers
controller := products.NewProductUpdateController(registry)
```

**New Usage**:
```go
import productupdate "project/internal/controllers/admin/shop/products/productupdate"

// In routes or handlers
controller := productupdate.NewProductUpdateController(registry)
```

**Files Changed**:
- `internal/controllers/admin/shop/products/product_update_controller.go` moved to:
- `internal/controllers/admin/shop/products/productupdate/product_update_controller.go`

**New Components Added**:
- `productupdate/detailscomponent/` - Product details management
- `productupdate/mediacomponent/` - Product media management  
- `productupdate/metadatacomponent/` - Product metadata management
- `productupdate/tagscomponent/` - Product tags management

**Action Required**:
- Update import paths from `products.NewProductUpdateController` to `productupdate.NewProductUpdateController`
- Update route registrations in `internal/controllers/admin/shop/routes.go`
- The new controller uses component-based architecture with separate handlers for details, media, metadata, and tags

---

### 11. Logo HTML Changes

**Change**: The default logo HTML has been changed from styled text to an image.

**Old Usage**:
```go
// Generated styled text logo with "Blue" and "Print"
left := hb.Span().Style("padding: 5px;").Style("color: white; font-size: 20px;").Text("Blue")
right := hb.Span().Style("padding: 5px;").Style("background: white; color: orange;").Text("Print")
frame := hb.Div().Style("display: inline-block; ...").Child(left).Child(right)
return frame.ToHTML()
```

**New Usage**:
```go
// Simple image logo
img := hb.Image("https://dracory.com/assets/images/logo.png").ToHTML()
return img
```

**Action Required**:
- Update `internal/layouts/logo_html.go` to customize your logo
- Either restore the old styled text implementation or use your own image URL

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

```bash
# Update go.mod dependencies
go get -u ./...

# Clean up unused dependencies
go mod tidy

# Verify build
go build ./...
```

### Step 2: Update Blog Store Configuration

Locate your blog store initialization (usually in `internal/registry/stores_blog.go`):

```go
// Add these fields to your blog store options
st, err := blogstore.NewStore(blogstore.NewStoreOptions{
    DB:                  db,
    PostTableName:       "snv_blogs_post",
    TaxonomyEnabled:     false,        // ADD THIS
    TaxonomyTableName:   "snv_blogs_taxonomy", // ADD THIS
    TermTableName:       "snv_blogs_term",     // ADD THIS
    VersioningEnabled:   true,
    VersioningTableName: "snv_blogs_version",
    AutomigrateEnabled:  true,
})
```

### Step 3: Run Automated Blog Store Updates

```bash
# Run all blog store updates in sequence
cd /d/PROJECTs/dracory.com/blueprint

# Update type signatures
find . -name "*.go" -exec sed -i 's/blogstore\.Post/blogstore.PostInterface/g' {} \;

# Update getter methods
find . -name "*.go" -exec sed -i 's/\.Title()/.GetTitle()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ID()/.GetID()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Slug()/.GetSlug()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Content()/.GetContent()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Summary()/.GetSummary()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ImageUrlOrDefault()/.GetImageUrlOrDefault()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.PublishedAt()/.GetPublishedAt()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.PublishedAtCarbon()/.GetPublishedAtCarbon()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Editor()/.GetEditor()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ContentType()/.GetContentType()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.CanonicalURL()/.GetCanonicalURL()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.IsPublished()/.IsPublished()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Status()/.GetStatus()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Featured()/.GetFeatured()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.ImageUrl()/.GetImageUrl()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaDescription()/.GetMetaDescription()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaKeywords()/.GetMetaKeywords()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.MetaRobots()/.GetMetaRobots()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Memo()/.GetMemo()/g' {} \;
find . -name "*.go" -exec sed -i 's/\.Website()/.GetWebsite()/g' {} \;
```

### Step 4: Update Testimonials Package Usage

Find all usages of `NewTestimonialFromEntity`:

```bash
grep -r "NewTestimonialFromEntity" --include="*.go" .
```

Update each call to include the store parameter:

```go
// Old
testimonial, err := testimonials.NewTestimonialFromEntity(entity)

// New
testimonial, err := testimonials.NewTestimonialFromEntity(store, entity)
```

### Step 5: Update Entity Store Usage

Find and update entity store usages:

```bash
# Find entity type usages
grep -r "entity\.Type()" --include="*.go" .
grep -r "entity\.ID" --include="*.go" .
grep -r "\.GetAttributes()" --include="*.go" .
grep -r "AttributeKey()" --include="*.go" .
grep -r "AttributeValue()" --include="*.go" .
```

### Step 6: Update Global Middlewares

Update `internal/routes/global_middlewares.go`:

```go
// Add rate limit helper
func getRateLimits(registry registry.RegistryInterface) (perSec, perMin, perHour int) {
    if registry.GetConfig() != nil {
        isDevelopment := registry.GetConfig().IsEnvDevelopment()
        isLocal := registry.GetConfig().IsEnvLocal()
        
        if isDevelopment || isLocal {
            return 1000, 10000, 100000
        }
    }
    return 20, 10, 12000
}

// Update middleware setup
func globalMiddlewares(registry registry.RegistryInterface) []rtr.MiddlewareInterface {
    perSec, perMin, perHour := getRateLimits(registry)
    
    globalMiddlewares := []rtr.MiddlewareInterface{
        rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
            Exclude: []string{"/new"},
            ExcludePaths: []string{
                "/blog*",
                "/th*",
                "/liveflux*",
                "/admin/cms*",
                "/admin/*cms*",
                "/assets*",
                "*/assets/*",
                "/files/*",
            },
        }),
        // ... other middlewares ...
        rtrMiddleware.RateLimitByIPMiddleware(perSec, 1),
        rtrMiddleware.RateLimitByIPMiddleware(perMin, 1*60),
        rtrMiddleware.RateLimitByIPMiddleware(perHour, 60*60),
    }
    // ...
}
```

### Step 7: Update Security Headers Middleware

Update `internal/middlewares/security_headers_middleware.go`:

```go
func NewSecurityHeadersMiddleware() rtr.MiddlewareInterface {
    return middlewares.NewSecurityHeadersMiddleware(middlewares.SecurityHeadersConfig{
        CSP: middlewares.CSPDirectives{
            DefaultSrc: []string{"'self'"},
            ScriptSrc:  getScriptSources(isDevelopment),
            StyleSrc:   getStyleSources(isDevelopment),
            ConnectSrc: []string{              // ADD THIS
                "'self'",
                "https://cdnjs.cloudflare.com",
                "http://cdnjs.cloudflare.com",
            },
            FontSrc: []string{
                "'self'",
                "https://cdn.jsdelivr.net",
                "https://fonts.googleapis.com",
                "https://fonts.gstatic.com",
                "https://cdnjs.cloudflare.com",
                "http://cdnjs.cloudflare.com",     // ADD THIS
                "https://maxcdn.bootstrapcdn.com",
            },
            // ... rest of config
        },
    })
}
```

### Step 8: Update Logo HTML (Optional)

Update `internal/layouts/logo_html.go` if you want to keep the old styled logo:

```go
func LogoHTML() string {
    primaryColor := "orange"
    secondaryColor := "white"
    
    left := hb.Span().
        Style("padding: 5px;").
        Style("color: " + secondaryColor + "; font-size: 20px;").
        Text("Blue")
    
    right := hb.Span().
        Style("padding: 5px;").
        Style("background: " + secondaryColor + "; color: " + primaryColor + ";").
        Text("Print")
    
    frame := hb.Div().
        Style("display: inline-block; justify-content: space-between; align-items: center; width: fit-content;").
        Style("padding: 0px;").
        Style("border: 3px solid " + primaryColor + "; background: " + primaryColor + "; color: " + secondaryColor + ";").
        Style("font-family: sans-serif; font-size: 20px; letter-spacing: 2px;").
        Child(left).
        Child(right)
    
    return frame.ToHTML()
}
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

### 2. Integration Tests

```bash
# Build the application
go build ./cmd/server

# Run the application
./server

# Test blog functionality
curl http://localhost:8080/blog
curl http://localhost:8080/blog/post/1/test-post

# Test admin panel
curl http://localhost:8080/admin

# Test file thumbnails
curl http://localhost:8080/thumb/files/test.jpg?w=300&h=200
```

### 3. Manual Verification Checklist

- [ ] Blog posts list displays correctly
- [ ] Individual blog posts load properly
- [ ] Blog post images load correctly
- [ ] Admin panel accessible
- [ ] File thumbnails generate properly
- [ ] Rate limiting works (check logs)
- [ ] Security headers present in responses
- [ ] Jail bots middleware not blocking legitimate traffic
- [ ] Testimonials display correctly (if applicable)
- [ ] Logo displays properly

---

## 📝 Additional Notes

### New Features in v0.18.0

1. **Blog Taxonomy Support**: The blog store now supports taxonomies (categories/tags). Set `TaxonomyEnabled: true` to enable.

2. **Environment-Aware Rate Limiting**: Development environments now have relaxed rate limits automatically.

3. **Enhanced Security Headers**: Additional CSP directives for better security.

4. **SQL File Storage Integration**: Thumb controller now properly integrates with SQL file storage for local files.

5. **Admin Shop Category Management**: New category manager with create, update, and list functionality. Access via `/admin?controller=shop_categories`.

6. **Admin File Manager**: New comprehensive file manager with:
   - Bulk delete, move operations
   - Directory create/delete
   - File upload (single and multiple)
   - File rename and view functionality
   - Modal-based UI components

7. **Shop Product Management Refactor**: Product update controller now uses component-based architecture with separate handlers for:
   - Product details
   - Product media
   - Product metadata
   - Product tags

### Dependency Updates

Key dependency changes in this release:
- AWS SDK v2 updated to latest versions
- SQL File Store updated to v1.3.0
- Google API clients updated
- Various indirect dependencies cleaned up

### Removed Features

None in this release.

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: blogstore.Post"

**Error**: `undefined: blogstore.Post`

**Solution**: Update to use `blogstore.PostInterface` instead:
```go
// Change this
func process(post blogstore.Post)

// To this
func process(post blogstore.PostInterface)
```

### Issue 2: "post.Title undefined (type blogstore.PostInterface has no field or method Title)"

**Error**: Method not found on PostInterface

**Solution**: Use getter methods with `Get` prefix:
```go
// Change this
post.Title()

// To this
post.GetTitle()
```

### Issue 3: "too many arguments to testimonials.NewTestimonialFromEntity"

**Error**: Function signature mismatch

**Solution**: Pass store as first argument:
```go
// Change this
testimonials.NewTestimonialFromEntity(entity)

// To this
testimonials.NewTestimonialFromEntity(store, entity)
```

### Issue 4: "entity.Type undefined"

**Error**: Entity interface method not found

**Solution**: Use `GetType()` instead:
```go
// Change this
entity.Type()

// To this
entity.GetType()
```

### Issue 5: Test failures after migration

**Error**: Tests using testify/require fail

**Solution**: Replace testify assertions with standard library:
```go
// Change this
require.NoError(t, err)

// To this
if err != nil {
    t.Fatalf("expected no error, got: %v", err)
}
```

---

## 📞 Support

For issues or questions regarding this upgrade:

1. Check the [Blueprint GitHub repository](https://github.com/dracory/blueprint) for issues
2. Review the [CHANGELOG.md](../CHANGELOG.md) for detailed change descriptions
3. Contact the Dracory team via GitHub discussions

---

**End of Upgrade Guide**

*Last Updated: April 2026*
