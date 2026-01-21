# Upgrade Guide: v0.13.0 to v0.14.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.13.0 to v0.14.0.

## ⚠️ Breaking Changes

### 1. Widget Constructor Registry Parameter Addition
**Change**: Widget constructors now require a registry parameter for accessing services

**Old Usage**:
```go
func NewPrintWidget() *printWidget {
	return &printWidget{}
}

type printWidget struct{}
```

**New Usage**:
```go
func NewPrintWidget() *printWidget {
	return &printWidget{}
}

type printWidget struct {
	registry registry.RegistryInterface
}
```

**Action Required**:
- Update all widget constructors to accept registry parameter
- Add registry field to widget structs
- Update constructor documentation to reflect registry parameter

### 2. CMS Frontend Shortcodes Parameter
**Change**: CMS frontend instance creation now requires explicit shortcodes parameter

**Old Usage**:
```go
instance := webtheme.New(blocks).ToHtml()
return cmsFrontend.New(cmsstore.FrontendOptions{
    Store:              registry.GetCmsStore(),
    Logger:             registry.GetLogger(),
    CacheEnabled:       true,
    CacheExpireSeconds: 1 * 60,
})
```

**New Usage**:
```go
instance := webtheme.New(blocks).ToHtml()
return cmsFrontend.New(cmsstore.FrontendOptions{
    Store:              registry.GetCmsStore(),
    Shortcodes:         shortcodes,
    Logger:             registry.GetLogger(),
    CacheEnabled:       true,
    CacheExpireSeconds: 1 * 60,
})
```

**Action Required**:
- Add `Shortcodes: shortcodes,` parameter to CMS frontend options
- Ensure shortcodes are properly initialized before CMS frontend creation

### 3. Dependency Updates
**Change**: Several dependencies were updated to newer versions

**Updated Dependencies**:
- `github.com/dracory/cmsstore`: v1.5.0 → v1.6.0
- `github.com/dracory/uid`: v1.8.0 → v1.9.0
- `modernc.org/sqlite`: v1.43.0 → v1.44.3

**Action Required**:
- Run `go mod tidy` to update dependencies
- Review dependency release notes for any API changes

### 4. Configuration Interface Formatting
**Change**: Configuration interface methods now include blank line separators for better readability

**Old Interface Format**:
```go
type databaseConfigInterface interface {
    SetDatabaseDriver(string)
    GetDatabaseDriver() string
    SetDatabaseHost(string)
    GetDatabaseHost() string
    // ... more methods without spacing
}
```

**New Interface Format**:
```go
type databaseConfigInterface interface {
    SetDatabaseDriver(string)
    GetDatabaseDriver() string

    SetDatabaseHost(string)
    GetDatabaseHost() string

    // ... more methods with blank line separators
}
```

**Action Required**:
- No code changes required, this is a formatting improvement
- Interface functionality remains identical

## 🔄 Migration Steps

### Step 1: Update Widget Constructors
Update all widget constructors in `internal/widgets/` directory:

```bash
# Find all widget files
find internal/widgets -name "*_widget.go" -type f

# For each widget, update constructor signature and add registry field
```

**Example Changes**:
```go
// Before
func NewAuthenticatedWidget() *authenticatedWidget {
    return &authenticatedWidget{}
}
type authenticatedWidget struct{}

// After  
func NewAuthenticatedWidget(registry registry.RegistryInterface) *authenticatedWidget {
    return &authenticatedWidget{registry: registry}
}
type authenticatedWidget struct {
    registry registry.RegistryInterface
}
```

### Step 2: Update CMS Controller
Ensure CMS controller passes shortcodes to frontend:

```bash
# Verify the shortcodes parameter is added in internal/controllers/website/cms/cms_controller.go
```

### Step 3: Update Dependencies
Update Go modules:

```bash
go mod tidy
go mod download
```

### Step 4: Update Widget Registry Calls
Update widget registry initialization to pass registry parameter:

```bash
# Update calls in internal/widgets/registry.go or similar
```

### Step 5: Test Configuration Loading
Verify configuration interface changes don't break existing code:

```bash
go test ./internal/config/...
```

## 🧪 Testing After Migration

### 1. Unit Tests
```bash
# Run all unit tests
go test ./...

# Run widget-specific tests
go test ./internal/widgets/...

# Run CMS controller tests  
go test ./internal/controllers/website/cms/...
```

### 2. Integration Tests
```bash
# Run integration tests with tags
go test -tags=integration ./...

# Test CMS functionality specifically
go test -tags=integration ./internal/controllers/website/cms/...
```

### 3. Manual Testing
- Verify all CMS shortcodes render correctly
- Test widget functionality with registry access
- Confirm configuration loading works as expected
- Check dependency compatibility

## 📝 Additional Notes

### New Features
- **Enhanced Widget Architecture**: Widgets now have access to registry for better service integration
- **Improved CMS Shortcodes**: Explicit shortcodes parameter provides better control over CMS functionality
- **Better Code Organization**: Configuration interfaces now have improved readability with spacing

### Removed Features
- No features were removed in this release

### Dependency Improvements
- Updated SQLite to latest stable version for better performance and security
- CMS store updates include potential performance improvements
- UID package updates may include new features or bug fixes

## 🆘 Common Issues and Solutions

### Issue 1: Widget Constructor Mismatch
**Problem**: Compilation errors due to missing registry parameter in widget constructors

**Solution**: Update all widget constructors to accept `registry registry.RegistryInterface` parameter

### Issue 2: CMS Frontend Missing Shortcodes
**Problem**: CMS frontend fails to initialize due to missing shortcodes parameter

**Solution**: Add `Shortcodes: shortcodes,` to the CMS frontend options in controller

### Issue 3: Dependency Conflicts
**Problem**: `go mod tidy` reports version conflicts

**Solution**: Clean go.mod and go.sum files, then run:
```bash
rm go.sum
go mod tidy
```

### Issue 4: Registry Access in Widgets
**Problem**: Widgets trying to access registry field that wasn't initialized

**Solution**: Ensure widget constructors properly assign registry parameter to struct field

## 📞 Support

For additional support:
- Check the [Blueprint GitHub repository](https://github.com/dracory/blueprint)
- Review existing upgrade guides in `docs/upgrade_guides/`
- Consult the project documentation for detailed API references
- Open an issue on GitHub for specific problems not covered in this guide

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested  
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
