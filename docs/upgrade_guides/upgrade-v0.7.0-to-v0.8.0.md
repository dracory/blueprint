# Upgrade Guide: v0.7.0 to v0.8.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.7.0 to v0.8.0.

## ⚠️ Breaking Changes

### 1. Link Constructor Deprecations and Removals
**Change**: Deprecated constructors removed and replaced with factory methods

**Old Usage**:
```go
// User links
links.NewUserLinks().Home()
links.NewUserLinks().Profile()

// Website links  
links.NewWebsiteLinks().Home()
links.NewWebsiteLinks().Blog(params)

// Admin links
links.NewAdminLinks().Home()
links.NewAdminLinks().Blog()

// Auth links
links.NewAuthLinks().Login(backUrl)
links.NewAuthLinks().Register()
```

**New Usage**:
```go
// User links
links.User().Home()
links.User().Profile()

// Website links
links.Website().Home()
links.Website().Blog(params)

// Admin links
links.Admin().Home()
links.Admin().Blog()

// Auth links
links.Auth().Login(backUrl)
links.Auth().Register()
```

**Action Required**:
- Replace all `NewUserLinks()` calls with `User()`
- Replace all `NewWebsiteLinks()` calls with `Website()`
- Replace all `NewAdminLinks()` calls with `Admin()`
- Replace all `NewAuthLinks()` calls with `Auth()`
- Update any custom code that extends these link types

### 2. Layout Options Request Field Removal
**Change**: Request field removed from layouts.Options struct

**Old Usage**:
```go
// In layout creation
options := layouts.Options{
    Request:        r,  // *http.Request
    AppName:        "My App",
    WebsiteSection: "dashboard",
    Title:          "Dashboard",
    Content:        content,
}
```

**New Usage**:
```go
// In layout creation - Request field removed
options := layouts.Options{
    // Request:        r,  // Removed
    AppName:        "My App",
    WebsiteSection: "dashboard",
    Title:          "Dashboard",
    Content:        content,
}
```

**Action Required**:
- Remove `Request: r` from all layouts.Options struct initializations
- Update layout functions to not rely on Request field
- Pass request as separate parameter if needed

### 3. Website Links Blog Method Signature Change
**Change**: Blog method parameter handling simplified

**Old Usage**:
```go
// In website_links.go v0.7.0
func (l *websiteLinks) Blog(params map[string]string) string {
    return URL(BLOG, params)
}

// Usage
links.Website().Blog(map[string]string{"page": "1"})
```

**New Usage**:
```go
// In website_links.go v0.8.0
func (l *websiteLinks) Blog(params ...map[string]string) string {
    p := lo.FirstOrEmpty(params)
    return URL(BLOG, p)
}

// Usage
links.Website().Blog(map[string]string{"page": "1"})
// or
links.Website().Blog()
```

**Action Required**:
- Update calls to Website().Blog() to use variadic parameters
- Ensure parameter handling is updated in custom link extensions

### 4. New Email Allowlist Middleware
**Change**: New middleware added for email-based access control

**Old Usage**:
```go
// No email allowlist middleware in v0.7.0
// Middleware stack was:
r.Use(middlewares.NewAuthMiddleware(app))
r.Use(middlewares.NewSessionMiddleware(app))
```

**New Usage**:
```go
// Email allowlist middleware available in v0.8.0
// Can be added to middleware stack:
r.Use(middlewares.NewEmailAllowlistMiddleware(app))
r.Use(middlewares.NewAuthMiddleware(app))
r.Use(middlewares.NewSessionMiddleware(app))
```

**Action Required**:
- No breaking change, but new middleware available
- Consider adding email allowlist middleware if needed for access control
- Update allowed emails map in the middleware for your requirements

### 5. SQL File Store Configuration Toggle
**Change**: New environment variable for SQL file store control

**Old Usage**:
```go
// In v0.7.0 - SQL file store was always enabled if configured
// No environment toggle available
```

**New Usage**:
```bash
# In .env.example - new toggle available
# ===============================================
# == SQL FILE STORE
# ===============================================
# SQL_FILE_STORE_USED="yes"
```

**Action Required**:
- Add `SQL_FILE_STORE_USED="yes"` to .env file if using SQL file store
- Update configuration loading to handle this new environment variable
- No breaking change if SQL file store is not used

### 6. Store Initialization Options in Test Framework
**Change**: Test framework enhanced with store initialization options

**Old Usage**:
```go
// In v0.7.0 - basic test setup
app := setupTestApp()
```

**New Usage**:
```go
// In v0.8.0 - enhanced test setup with store options
app := setupTestApp(setup.WithStoreInitialization(true))
```

**Action Required**:
- Update test setup code to use new store initialization options
- No breaking change for existing tests, but new options available

## 🔄 Migration Steps

### Step 1: Update Link Constructor Calls
```bash
# Find and replace deprecated link constructors
find . -name "*.go" -type f -exec sed -i 's/links\.NewUserLinks()/links.User()/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/links\.NewWebsiteLinks()/links.Website()/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/links\.NewAdminLinks()/links.Admin()/g' {} \;
find . -name "*.go" -type f -exec sed -i 's/links\.NewAuthLinks()/links.Auth()/g' {} \;
```

### Step 2: Update Layout Options
```bash
# Find and remove Request field from layouts.Options
find . -name "*.go" -type f -exec grep -l "layouts\.Options" {} \; | xargs sed -i '/Request:/d'
```

### Step 3: Update Website Links Blog Calls
```bash
# Find Blog method calls and ensure they work with new signature
find . -name "*.go" -type f -exec grep -n "\.Blog(" {} \;
# Manually review and update if needed
```

### Step 4: Update Environment Configuration
```bash
# Add SQL file store toggle to .env file
echo '# ===============================================
# == SQL FILE STORE
# ===============================================
SQL_FILE_STORE_USED="yes"' >> .env
```

### Step 5: Update Test Setup
```bash
# Find test files and update setup if needed
find . -name "*_test.go" -type f -exec grep -l "setupTestApp\|setupTestApp(" {} \;
# Review and update with new store initialization options
```

### Step 6: Add Email Allowlist Middleware (Optional)
```go
// In your routes setup, consider adding:
r.Use(middlewares.NewEmailAllowlistMiddleware(app))

// Customize allowed emails in the middleware file
var allowedEmails = map[string]struct{}{
    "admin@yourdomain.com": {},
    "user@yourdomain.com":  {},
}
```

## 🧪 Testing After Migration

### 1. Link Generation Tests
```bash
# Test all link generation methods
go test ./internal/links/...

# Test specific link methods
go test -run TestUserLinks ./internal/links/
go test -run TestWebsiteLinks ./internal/links/
go test -run TestAdminLinks ./internal/links/
go test -run TestAuthLinks ./internal/links/
```

### 2. Layout Rendering Tests
```bash
# Test layout rendering without Request field
go test ./internal/layouts/...

# Test specific layout functions
go test -run TestLayoutOptions ./internal/layouts/
```

### 3. Middleware Tests
```bash
# Test middleware stack
go test ./internal/middlewares/...

# Test new email allowlist middleware
go test -run TestEmailAllowlistMiddleware ./internal/middlewares/
```

### 4. Integration Tests
```bash
# Test full application with new link methods
go test -tags=integration ./...

# Test SQL file store configuration
go test -run TestSqlFileStore ./internal/...
```

### 5. Manual Testing
```bash
# Start application and test:
# 1. All navigation links work correctly
# 2. Layouts render without Request field
# 3. Email allowlist middleware (if enabled)
# 4. SQL file store functionality
go run main.go
```

## 📝 Additional Notes

### New Features Added
- **Email Allowlist Middleware**: Restrict access based on email domains
- **SQL File Store Toggle**: Environment variable to control SQL file store
- **Enhanced Test Framework**: Store initialization options for better testing
- **Simplified Link Constructors**: Cleaner factory method patterns

### Deprecated Features Removed
- **NewUserLinks()**: Replaced with User()
- **NewWebsiteLinks()**: Replaced with Website()  
- **NewAdminLinks()**: Replaced with Admin()
- **NewAuthLinks()**: Replaced with Auth()

### API Improvements
- **Consistent Link Patterns**: All link types now use factory methods
- **Better Parameter Handling**: Website links use variadic parameters
- **Cleaner Layout Options**: Removed unnecessary Request field

## 🆘 Common Issues and Solutions

### Issue: Link Constructor Not Found
**Problem**: Build fails with "undefined: links.NewUserLinks"
**Solution**: 
```bash
# Replace with new factory method
sed -i 's/links\.NewUserLinks()/links.User()/g' **/*.go
```

### Issue: Layout Options Compilation Error
**Problem**: "unknown field 'Request' in struct of type layouts.Options"
**Solution**: Remove Request field from all layouts.Options initializations

### Issue: Blog Method Parameter Error
**Problem**: "too many arguments to links.Website().Blog"
**Solution**: Ensure Blog method calls use variadic parameter format

### Issue: SQL File Store Not Working
**Problem**: SQL file store not initializing
**Solution**: Add `SQL_FILE_STORE_USED="yes"` to environment configuration

### Issue: Test Failures
**Problem**: Tests failing after migration
**Solution**: Update test setup to use new store initialization options

## 📞 Support

For additional support:
- Check the repository issues: https://github.com/dracory/blueprint/issues
- Review the documentation: https://github.com/dracory/blueprint/docs
- Compare with the example application in the repository
- Check middleware and link examples in the codebase
