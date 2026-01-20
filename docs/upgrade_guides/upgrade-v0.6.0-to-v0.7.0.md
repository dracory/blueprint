# Upgrade Guide: v0.6.0 to v0.7.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.6.0 to v0.7.0.

## ⚠️ Breaking Changes

### 1. Blog Store Configuration Addition
**Change**: New blog store configuration constants and environment variables added

**Old Usage**:
```go
// In v0.6.0 - no blog store configuration
// Store configurations started from CMS store

// Environment variables in .env.example:
# CMS_STORE_USED="yes"
# CUSTOM_STORE_USED="yes"
# ENTITY_STORE_USED="yes"
# ... (no BLOG_STORE_USED)
```

**New Usage**:
```go
// In v0.7.0 - blog store configuration added
const KEY_BLOG_STORE_USED = "BLOG_STORE_USED"

// Environment variables in .env.example:
# ===============================================
# == BLOG STORE
# ===============================================
# BLOG_STORE_USED="yes"

# ===============================================
# == CMS STORE
# ===============================================
# CMS_STORE_USED="yes"
# ... (other stores follow)
```

**Action Required**:
- Add `BLOG_STORE_USED="yes"` to .env file if using blog functionality
- Update configuration loading to handle blog store initialization
- No breaking change if blog store is not used

### 2. Taskfile.yml LIVEURL Variable Addition
**Change**: New LIVEURL variable added to taskfile.yml for deployment URL handling

**Old Usage**:
```yaml
# In taskfile.yml v0.6.0
vars:
  APPNAME: The Dracory Blueprint Project
  DATETIME: '{{now | date "20060102_150405"}}'
```

**New Usage**:
```yaml
# In taskfile.yml v0.7.0
vars:
  APPNAME: The Dracory Blueprint Project
  DATETIME: '{{now | date "20060102_150405"}}'
  LIVEURL: https://example.com
```

**Action Required**:
- Update taskfile.yml to include LIVEURL variable
- Customize LIVEURL value for your deployment domain
- No breaking change for existing functionality

### 3. Blog Store Support in Test Setup Utilities
**Change**: Test framework enhanced with blog store support

**Old Usage**:
```go
// In v0.6.0 - basic test setup without blog store
app := setupTestApp()
// Blog store not available in test utilities
```

**New Usage**:
```go
// In v0.7.0 - enhanced test setup with blog store
app := setupTestApp(setup.WithBlogStore(true))
// Blog store now available in test utilities
```

**Action Required**:
- Update test setup code to use new blog store options if needed
- No breaking change for existing tests
- New functionality available for blog-related testing

## 🔄 Migration Steps

### Step 1: Update Environment Configuration
```bash
# Add blog store configuration to .env file
echo '# ===============================================
# == BLOG STORE
# ===============================================
BLOG_STORE_USED="yes"' >> .env
```

### Step 2: Update Taskfile Configuration
```bash
# Add LIVEURL variable to taskfile.yml
# Find the vars section and add:
#   LIVEURL: https://your-domain.com

# Or use sed to add it after DATETIME line
sed -i '/DATETIME:/a\  LIVEURL: https://example.com' taskfile.yml
```

### Step 3: Update Test Setup (Optional)
```bash
# Find test files that use setupTestApp
find . -name "*_test.go" -type f -exec grep -l "setupTestApp" {} \;

# Update test setup to include blog store if needed
# Example: app := setupTestApp(setup.WithBlogStore(true))
```

### Step 4: Verify Configuration Loading
```bash
# Test that configuration loads correctly with new blog store setting
go test -run TestConfigLoad ./internal/config/

# Verify blog store initialization
go test -run TestBlogStore ./internal/...
```

### Step 5: Update Deployment Scripts
```bash
# If using deployment scripts, update them to use LIVEURL variable
# Update any hardcoded URLs to use {{.LIVEURL}} from taskfile
```

## 🧪 Testing After Migration

### 1. Configuration Tests
```bash
# Test configuration loading with new blog store setting
go test ./internal/config/...

# Test specific blog store configuration
go test -run TestBlogStoreConfig ./internal/config/
```

### 2. Blog Store Tests
```bash
# Test blog store functionality
go test ./internal/controllers/blog/...

# Test blog store initialization
go test -run TestBlogStoreInit ./internal/...
```

### 3. Test Framework Tests
```bash
# Test enhanced test framework
go test ./internal/testutils/...

# Test blog store in test environment
go test -run TestWithBlogStore ./internal/testutils/
```

### 4. Integration Tests
```bash
# Test full application with blog store enabled
go test -tags=integration ./...

# Test blog functionality end-to-end
go test -run TestBlogIntegration ./...
```

### 5. Taskfile Tests
```bash
# Test taskfile commands with new LIVEURL variable
task deploy:staging
task deploy:production

# Verify LIVEURL is used correctly in deployment
```

### 6. Manual Testing
```bash
# Start application and test:
# 1. Blog functionality works if BLOG_STORE_USED="yes"
# 2. Configuration loads without errors
# 3. Deployment scripts use correct URLs
go run main.go
```

## 📝 Additional Notes

### New Features Added
- **Blog Store Configuration**: Environment variable to control blog store initialization
- **Deployment URL Handling**: LIVEURL variable in taskfile for better deployment management
- **Enhanced Test Framework**: Blog store support in test utilities
- **Better Configuration Management**: More granular control over store initialization

### Configuration Enhancements
- **Granular Store Control**: Blog store can be enabled/disabled independently
- **Deployment Flexibility**: LIVEURL variable makes deployment scripts more portable
- **Test Environment**: Better testing support for blog functionality

### Non-Breaking Changes
- All changes are additive and backward compatible
- Existing applications will continue to work without modification
- New features are opt-in through configuration

## 🆘 Common Issues and Solutions

### Issue: Blog Store Not Initializing
**Problem**: Blog store not loading despite BLOG_STORE_USED="yes"
**Solution**: 
```bash
# Verify environment variable is set
grep BLOG_STORE_USED .env

# Check configuration loading
go test -run TestBlogStoreConfig ./internal/config/
```

### Issue: Deployment URL Incorrect
**Problem**: Deployment using wrong URL
**Solution**: 
```bash
# Update LIVEURL in taskfile.yml
sed -i 's|LIVEURL: https://example.com|LIVEURL: https://your-domain.com|' taskfile.yml
```

### Issue: Test Failures with Blog Store
**Problem**: Tests failing after blog store addition
**Solution**: 
```bash
# Update test setup to include blog store options
# Or disable blog store in test environment
echo 'BLOG_STORE_USED="no"' > .env.test
```

### Issue: Configuration Loading Errors
**Problem**: Configuration fails to load with new settings
**Solution**: 
```bash
# Verify all required environment variables are set
go test -run TestConfigLoad ./internal/config/

# Check for missing constants in contstants.go
```

## 📞 Support

For additional support:
- Check the repository issues: https://github.com/dracory/blueprint/issues
- Review the documentation: https://github.com/dracory/blueprint/docs
- Compare with the example application in the repository
- Check blog store examples in the codebase

## Migration Checklist

- [ ] Add BLOG_STORE_USED to .env file if using blog functionality
- [ ] Update taskfile.yml with LIVEURL variable
- [ ] Update test setup if using blog store in tests
- [ ] Verify configuration loading works correctly
- [ ] Test blog store functionality
- [ ] Update deployment scripts to use LIVEURL
- [ ] Run full test suite to ensure compatibility
- [ ] Test deployment with new URL handling
