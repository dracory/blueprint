# Upgrade Guide: v0.9.0 to v0.10.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.9.0 to v0.10.0.

## ⚠️ Major Breaking Changes

### 1. Task Package Reorganization
**Change**: Task files moved from flat structure to subdirectories by feature

**Old Structure**:
```bash
internal/tasks/
├── blind_index_rebuild_task.go
├── clean_up_task.go
├── email_admin_task.go
├── hello_world_task.go
└── stats/
```

**New Structure**:
```bash
internal/tasks/
├── blind_index_rebuild/
│   ├── blind_index_rebuild_task.go
│   └── doc.go
├── clean_up/
│   └── clean_up_task.go
├── email_admin/
│   └── email_admin_task.go
├── hello_world/
│   └── hello_world_task.go
├── stats/
│   └── stats_visitor_enhance_task.go
├── register_tasks.go
└── doc.go
```

**Action Required**:
- Update import paths for task files
- Update any direct file references to tasks
- Task registration via `register_tasks.go` remains unchanged

### 2. Database SSL Mode Configuration
**Change**: Added configurable SSL mode for database connections

**Old Configuration**:
```go
// SSL mode was hardcoded or not configurable
```

**New Configuration**:
```go
// SSL mode now configurable via config interface
cfg.GetDatabaseSSLMode() // Returns "require" as default
```

**Environment Variable**:
```bash
# Add to .env.example
DB_SSL_MODE="require"  # or "disable", "verify-full", etc.
```

**Action Required**:
- Add `DB_SSL_MODE` to environment configuration
- Update database connection code to use `GetDatabaseSSLMode()`
- SQLite databases automatically skip SSL configuration

### 3. Environment Variable Standardization
**Change**: Environment variable naming standardized across stores

**Old Variables**:
```bash
# Inconsistent naming patterns
```

**New Variables**:
```bash
# Standardized store configuration section
# All store configs follow pattern: [STORE]_STORE_USED="yes"
```

**Action Required**:
- Review `.env.example` for new standardized variable names
- Update any environment variable references
- Add new `FEED_STORE` configuration section if needed

### 4. Cache Helper Functions Reorganization
**Change**: Cache helper functions reorganized and renamed for consistency

**Old Usage**:
```go
// Various cache helper functions with inconsistent naming
```

**New Usage**:
```go
// Standardized cache helper function names
// Improved consistency across cache operations
```

**Action Required**:
- Update calls to renamed cache helper functions
- Review cache-related code for new function names
- Test cache functionality after migration

### 5. Component Interface Method Standardization
**Change**: Component interface method names standardized

**Old Interface**:
```go
// Inconsistent method naming across components
```

**New Interface**:
```go
// Standardized method naming conventions
// Improved consistency across all components
```

**Action Required**:
- Update component implementations to use standardized method names
- Review any custom components for interface compliance
- Test component functionality after updates

### 6. Schedule Files Reorganization
**Change**: Schedule files reorganized and cleaned up

**Old Structure**:
```bash
internal/schedules/
├── [various schedule files with inconsistent organization]
```

**New Structure**:
```bash
internal/schedules/
├── [cleaned up and reorganized schedule files]
```

**Action Required**:
- Update import paths for schedule files
- Review schedule registration code
- Test schedule execution after migration

### 7. CMS Route Standardization
**Change**: CMS route naming and paths standardized

**Old Routes**:
```go
// Inconsistent CMS route naming and paths
```

**New Routes**:
```go
// Standardized CMS route names and paths
// Improved consistency across CMS functionality
```

**Action Required**:
- Update any direct CMS route references
- Review CMS controller implementations
- Test CMS functionality after migration

### 8. Liveflux Controller Context Handling
**Change**: Improved context handling in Liveflux controller

**Old Implementation**:
```go
// Component registration in controller initialization
// String literal context keys
```

**New Implementation**:
```go
// Component registration moved to component constructors
// Typed context key constants for type safety
```

**Action Required**:
- Update Liveflux controller implementations
- Replace string literal context keys with typed constants
- Move component registration to appropriate constructors

### 9. Environment Encryption Validation
**Change**: Enhanced validation for environment encryption

**Old Validation**:
```go
// Basic validation with debug logging
```

**New Validation**:
```go
// Explicit validation for missing private key when encryption enabled
// Removed debug logging statements
// Improved variable naming consistency
```

**Action Required**:
- Update environment encryption initialization code
- Remove any dependencies on debug logging
- Test encryption functionality with new validation

### 10. Database Connection Optimization
**Change**: Added database connection pool settings for production

**Old Configuration**:
```go
// Basic database connection without pool optimization
```

**New Configuration**:
```go
// Production-ready database connection pool settings
// Configurable max open/idle connections and lifetime
// Template caching with thread-safe sync.RWMutex
```

**Action Required**:
- Update database connection code to use new pool settings
- Configure connection pool parameters for production
- Test database performance with new settings

## 🔄 Migration Steps

### Step 1: Update Task Import Paths
```bash
# Update task file imports from flat structure to subdirectories
# Manual review required for each task import
```

### Step 2: Add Database SSL Configuration
```bash
# Add to .env file
DB_SSL_MODE="require"
```

### Step 3: Update Environment Variables
```bash
# Review and update environment variable names
# Follow new standardized naming conventions
# Add new FEED_STORE configuration if needed
```

### Step 4: Update Cache Helper Functions
```bash
# Update calls to renamed cache helper functions
# Manual review required for each cache usage
```

### Step 5: Update Component Interfaces
```bash
# Update component implementations to use standardized method names
# Review custom components for interface compliance
```

### Step 6: Update Schedule Imports
```bash
# Update import paths for reorganized schedule files
# Review schedule registration code
```

### Step 7: Update CMS Routes
```bash
# Update CMS route references
# Review CMS controller implementations
```

### Step 8: Update Liveflux Controller
```bash
# Update context handling
# Replace string literals with typed constants
# Move component registration
```

### Step 9: Update Environment Encryption
```bash
# Update encryption validation code
# Remove debug logging dependencies
```

### Step 10: Configure Database Pool
```bash
# Update database connection code
# Configure pool settings for production
# Test performance improvements
```

## 🧪 Testing After Migration

1. **Task Execution**: Test all tasks with new import paths
2. **Database Connections**: Verify SSL mode and pool settings
3. **Cache Operations**: Test renamed cache helper functions
4. **Component Functionality**: Test standardized component interfaces
5. **Schedule Execution**: Test reorganized schedule files
6. **CMS Routes**: Test standardized CMS functionality
7. **Liveflux Controller**: Test improved context handling
8. **Environment Encryption**: Test enhanced validation
9. **Database Performance**: Test connection pool optimization

## 📝 Additional Notes

### New Features Added
- Configurable database SSL mode
- Database connection pool optimization
- Template caching with thread-safe implementation
- Enhanced environment encryption validation
- Improved context handling in Liveflux controller

### Improved Organization
- Task files organized by feature in subdirectories
- Cache helper functions standardized
- Component interface methods consistent
- Schedule files cleaned up
- CMS routes standardized

### Code Quality Improvements
- Removed debug logging statements
- Improved variable naming consistency
- Better error handling and validation
- Type-safe context keys

## 🆘 Common Issues

### Issue: Task Import Path Not Found
**Solution**: Update import paths from flat structure to subdirectories

### Issue: Database SSL Configuration Error
**Solution**: Add `DB_SSL_MODE` environment variable and update connection code

### Issue: Cache Helper Function Not Found
**Solution**: Update function calls to use new standardized names

### Issue: Component Interface Method Not Found
**Solution**: Update component implementations to use standardized method names

### Issue: Liveflux Context Handling Error
**Solution**: Replace string literal context keys with typed constants

## 📞 Support

For issues during migration:
1. Check this guide first
2. Review the updated documentation in `docs/overview.md`
3. Check the architectural review in `docs/review.md`
4. Run tests to identify specific issues

---

**Version**: v0.10.0  
**Previous Version**: v0.9.0  
**Release Date**: November 2025
