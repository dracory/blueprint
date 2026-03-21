# Proposal: Move Components to Base Package

## Overview

This proposal identifies components from the Blueprint project that should be moved to the `github.com/dracory/base` package to improve code reusability and reduce duplication across projects.

## Current Base Package Structure

The `github.com/dracory/base` package currently provides:
- Array utilities (`github.com/dracory/arr`)
- BBCode to HTML conversion
- Environment variables (`github.com/dracory/env`)
- Database functionality (`github.com/dracory/database`)
- Email functionality
- Markdown to HTML conversion
- Object property management
- Router functionality (`github.com/dracory/rtr`)
- Server functionality
- String manipulation
- Test utilities
- Timezone handling
- Workflow utilities

## Components to Move from Blueprint

### 1. HTTP Utilities (High Priority) ✅ COMPLETED

#### Files to Move:
- ~~`internal/utils/response_utils.go`~~ - Safe HTTP response body closing ✅ MOVED
- ~~`internal/helpers/redirect.go`~~ - HTTP redirect helper ✅ MOVED

#### Reasoning:
These are generic HTTP utilities that can be used by any Go web application. They don't depend on Blueprint-specific architecture.

#### Proposed Location:
`github.com/dracory/base/http` ✅ COMPLETED

#### Status:
- ~~`response_utils.go`~~ has been successfully moved to `github.com/dracory/base/http` ✅
- ~~`redirect.go`~~ has been successfully moved to `github.com/dracory/base/http` ✅
- Both functions include comprehensive tests and documentation
- Ready for use in other projects

---

### 2. File System Utilities (High Priority) ✅

#### Files to Move:
- ~~`internal/helpers/filesystem.go`~~ - Embedded file system helpers

#### Status:
- ~~Moved to `github.com/dracory/base/filesystem`~~ ✅
- Functions `EmbeddedFileToBytes` and `EmbeddedFileToString` successfully migrated
- Comprehensive tests added with embedded test data
- Ready for use in other projects

#### Reasoning:
Embedded file system operations are common across many Go applications and should be available in the base package.

#### Proposed Location:
~~`github.com/dracory/base/filesystem`~~ ✅

---

### 3. URL Building Utilities (Medium Priority) ✅

#### Files to Move:
- ~~`internal/links/links.go`~~ - Core URL building functions (moved to base)
- ~~`internal/links/url.go`~~ - URL construction helper (refactored to use base)

#### Status:
- ~~Moved to `github.com/dracory/base/url`~~ ✅
- **Dependency Injection**: URLBuilder struct with configurable root URL
- **Clean Architecture**: No environment variable dependencies in base package
- **Functions Migrated**: `RootURL()`, `BuildURL()`, `BuildQuery()`, `HttpBuildQuery()`
- **Comprehensive Tests**: Full test coverage with dependency injection patterns
- **Blueprint Integration**: Updated to use base utilities with backward compatibility

#### Reasoning:
URL building is a common need, but the application-specific link files should remain in Blueprint.

#### Proposed Location:
~~`github.com/dracory/base/url`~~ ✅

#### Keep in Blueprint:
- `admin_links.go`, `auth_links.go`, `user_links.go`, `website_links.go` (application-specific)
- `constants.go` (application-specific route constants)

---

### 4. Generic Types (Medium Priority) ✅

#### Files to Move:
- ~~`internal/types/flash_message.go`~~ - Flash message structure

#### Status:
- ~~Moved to `github.com/dracory/base/types`~~ ✅
- **FlashMessage Structure**: Standard flash message type for web applications
- **Comprehensive Tests**: Full test coverage including zero values and common types
- **Documentation**: Complete README with usage examples for authentication and forms
- **Blueprint Integration**: Updated testutils to use base FlashMessage type
- **Clean Separation**: Generic type now available for all Dracory projects

#### Reasoning:
Flash messaging is a common web pattern and the type definition is generic.

#### Proposed Location:
~~`github.com/dracory/base/types`~~ ✅

---

### 5. Security Middleware (High Priority) ✅

#### Files to Move:
- ~~`internal/middlewares/https_redirect_middleware.go`~~ - Moved to RTR with customization
- ~~`internal/middlewares/security_headers_middleware.go`~~ - Moved to RTR with customization

#### Status:
- ~~Moved to `github.com/dracory/rtr/middleware/security`~~ ✅
- **Highly Configurable**: Both middlewares support extensive customization per project
- **HTTPS Redirect**: Configurable localhost skipping, custom skip functions, trusted proxies
- **Security Headers**: Full CSP, HSTS, frame options, XSS protection, permissions policy
- **Blueprint Integration**: Updated to use RTR middlewares with project-specific configs
- **Comprehensive Tests**: Full test coverage for all configuration options
- **Documentation**: Complete README with examples for different use cases

#### Analysis:
**Location: RTR Package (not Base)**
- Both middlewares return `rtr.MiddlewareInterface` and use `rtr.NewMiddleware()`
- Designed specifically for RTR router framework integration
- HTTP middleware layer is a router concern, not base package concern
- Base package should remain framework-agnostic

#### Proposed Location:
~~`github.com/dracory/rtr/middleware/security`~~ ✅

#### Reasoning:
These security-focused middlewares are framework-agnostic in concept but RTR-specific in implementation, providing essential security functionality for any web application using the RTR router.

#### Customization Features:
- **HTTPS Redirect**: SkipLocalhost, TrustedProxies, CustomSkipFunc
- **Security Headers**: Full CSP configuration, HSTS settings, frame options, custom headers
- **Project-Specific**: Each project can override defaults while maintaining secure baseline

---

### 6. Test Utilities (Medium Priority) ✅

#### Files to Move:
- ~~`internal/testutils/testutils.go`~~ - Mock SMTP server setup
- ~~`internal/testutils/constants.go`~~ - Test constants

#### Status:
- ~~Moved to `github.com/dracory/test`~~ ✅
- **Mock SMTP Server**: `SetupMailServer()` function for email testing
- **Test Constants**: Generic test identifiers (ADMIN_01, USER_01, ORDER_01, etc.)
- **Comprehensive Tests**: Full test coverage for SMTP mock functionality
- **Blueprint Integration**: Updated to use test package constants
- **Clean Separation**: Test utilities now available across Dracory ecosystem

#### Analysis:
**Location: Test Package (not Base)**
- Mock SMTP server is a testing utility, not runtime application logic
- Test constants are specifically for testing scenarios
- Existing `github.com/dracory/test` package already has comprehensive testing infrastructure
- Base package should contain runtime utilities, not test helpers

#### Proposed Location:
~~`github.com/dracory/test`~~ ✅

#### Reasoning:
Generic test utilities that can be used across projects. The setup.go file should remain as it's Blueprint-specific.

#### Keep in Blueprint:
- `setup.go` (Blueprint-specific test setup)
- `seed_*.go` files (Blueprint-specific data seeding)
- `login_as.go`, `new_request.go` (application-specific test helpers)

---

### 7. Layout Components (Low Priority)

#### Files to Move:
- `internal/layouts/options.go` - Layout options struct
- `internal/layouts/breadcrumb.go` - Breadcrumb type

#### Reasoning:
These are generic layout structures that could be reused, but layout implementations are often application-specific.

#### Proposed Location:
`github.com/dracory/base/layout`

#### Keep in Blueprint:
- All specific layout implementations (`admin_layout.go`, `user_layout.go`, etc.)

---

## Components to Keep in Blueprint

### Application-Specific Components:
- `internal/helpers/flash.go` - Depends on Blueprint's cache store and links
- `internal/helpers/timezone_from_request.go` - Depends on Blueprint's auth system
- `internal/helpers/get_auth_*.go` - Blueprint-specific auth helpers
- `internal/helpers/extend_session.go` - Blueprint-specific session handling
- `internal/helpers/user_settings.go` - Blueprint-specific user management
- `internal/helpers/blog_post_blocks_to_string.go` - Blueprint-specific blog functionality
- `internal/helpers/untokenize.go` - Depends on Blueprint's vault store
- `internal/middlewares/` (except security middlewares) - Most are Blueprint-specific
- `internal/layouts/` (except generic types) - Application-specific layouts
- `internal/links/` (except core URL building) - Application-specific routes
- `internal/controllers/` - Application-specific controllers
- `internal/config/` - Blueprint-specific configuration
- `internal/registry/` - Blueprint-specific dependency injection

## Migration Strategy

### Phase 1: High Priority Components
1. Move HTTP utilities (`response_utils.go`, `redirect.go`)
2. Move file system utilities (`filesystem.go`)
3. Move security middleware (`https_redirect_middleware.go`, `security_headers_middleware.go`)

### Phase 2: Medium Priority Components
1. Move URL building utilities
2. Move generic types (`flash_message.go`)
3. Move test utilities

### Phase 3: Low Priority Components
1. Move layout components
2. Review and move any additional utilities identified

### Implementation Steps:
1. Create new packages in `github.com/dracory/base`
2. Move files with appropriate package name changes
3. Update import statements in Blueprint
4. Add comprehensive tests to base package
5. Update documentation
6. Release new versions of base package

## Benefits

1. **Code Reusability**: Components can be used across multiple projects
2. **Reduced Duplication**: Common utilities centralized in one location
3. **Better Testing**: Centralized components can have more comprehensive test coverage
4. **Easier Maintenance**: Bug fixes and improvements benefit all projects
5. **Cleaner Blueprint**: Blueprint becomes more focused on application-specific code

## Risks and Mitigations

### Risks:
1. **Breaking Changes**: Import changes may break existing code
2. **Dependency Management**: Base package may acquire too many dependencies
3. **Version Compatibility**: Need to maintain backward compatibility

### Mitigations:
1. **Gradual Migration**: Move components in phases to minimize disruption
2. **Clear Documentation**: Document what's moved and how to migrate
3. **Version Tags**: Use semantic versioning for breaking changes
4. **Deprecation Notices**: Provide clear deprecation paths

## Conclusion

Moving these components to the base package will improve code organization and reusability while keeping Blueprint focused on application-specific functionality. The proposed migration strategy minimizes risk while providing immediate benefits through the high-priority components.
