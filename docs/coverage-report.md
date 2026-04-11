# Test Coverage Report

**Generated:** April 11, 2026  
**Project:** Blueprint  
**Status:** ⚠️ Build Failures Preventing Full Coverage Analysis

## Executive Summary

The blueprint project has **build failures** preventing complete test coverage analysis. The failures are related to missing methods in the `blogstore` interface that are being called in the blog management controllers. A total of **5 methods are missing** from the blogstore interface.

## Build Failures

### Critical Issues

The following build errors must be resolved before full coverage analysis can be completed:

#### 1. **Missing blogstore Interface Methods**

**Affected Files:**
- `internal/controllers/admin/blog/category_manager/category_manager_controller.go:295`
- `internal/controllers/admin/blog/tag_manager/tag_manager_controller.go:175`
- `internal/controllers/admin/blog/post_update/post_update_controller.go:316, 366, 401, 475, 524, 559`

**Missing Methods:**
1. `SetSequence()` - Used in category_manager_controller.go:295
2. `PostListByTermID()` - Used in tag_manager_controller.go:175
3. `TermListByPostID()` - Used in post_update_controller.go:316, 475
4. `PostAddTerm()` - Used in post_update_controller.go:366, 524
5. `PostRemoveTerm()` - Used in post_update_controller.go:559

**Error Details:**
```
term.SetSequence undefined (type blogstore.TermInterface has no field or method SetSequence)
blogStore.PostListByTermID undefined (type blogstore.StoreInterface has no field or method PostListByTermID)
blogStore.TermListByPostID undefined (type blogstore.StoreInterface has no field or method TermListByPostID)
blogStore.PostAddTerm undefined (type blogstore.StoreInterface has no field or method PostAddTerm)
blogStore.PostRemoveTerm undefined (type blogstore.StoreInterface has no field or method PostRemoveTerm)
```

#### 2. **CMS Controller Return Value Mismatch**

**File:** `internal/controllers/website/cms/cms_controller.go:65`

**Error:**
```
not enough return values
have (string)
want (bool, string)
```

#### 3. **Undefined CDN Lazy Loading**

**File:** `internal/controllers/website/blog/home/blog_controller.go:51`

**Error:**
```
undefined: cdn.Slazy_0_5_0
```

## Partial Coverage Results

Based on the test output before build failures, the following packages have coverage data:

### High Coverage (>70%)

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/controllers/auth/login` | 100.0% | ✅ |
| `internal/controllers/auth/logout` | 100.0% | ✅ |
| `internal/controllers/user` | 100.0% | ✅ |
| `internal/controllers/website/home` | 95.0% | ✅ |
| `internal/controllers/website/pages/indexnow` | 95.8% | ✅ |
| `internal/controllers/shared/resource` | 93.3% | ✅ |
| `internal/controllers/shared/flash` | 90.3% | ✅ |
| `internal/controllers/admin/shop/products/productupdate/tagscomponent` | 78.7% | ✅ |
| `internal/controllers/admin/shop/products/productupdate/descriptioncomponent` | 79.6% | ✅ |
| `internal/ext` | 73.8% | ✅ |
| `internal/registry` | 71.4% | ✅ |
| `internal/website/seo` | 76.3% | ✅ |

### Medium Coverage (50-70%)

| Package | Coverage | Status |
|----------|----------|--------|
| `internal/controllers/user/account` | 64.2% | ⚠️ |
| `internal/tasks/email_admin` | 66.7% | ⚠️ |
| `internal/tasks/email_admin_new_contact` | 66.7% | ⚠️ |
| `internal/middlewares` | 59.9% | ⚠️ |
| `internal/tasks/clean_up` | 60.0% | ⚠️ |
| `internal/tasks/email_admin_new_user_registered` | 60.0% | ⚠️ |
| `internal/tasks/email_test` | 60.5% | ⚠️ |
| `internal/controllers/auth/authentication` | 48.6% | ⚠️ |
| `internal/controllers/auth/register` | 53.7% | ⚠️ |
| `pkg/social` | 55.1% | ⚠️ |

### Low Coverage (<50%)

| Package | Coverage | Status |
|----------|----------|--------|
| `internal/controllers/user/home` | 36.4% | ❌ |
| `internal/controllers/shared/thumb` | 30.1% | ❌ |
| `internal/controllers/shared/media` | 23.1% | ❌ |
| `internal/controllers/shared/file` | 23.1% | ❌ |
| `internal/tasks/blind_index_rebuild` | 23.5% | ❌ |
| `internal/controllers/admin/users/user_impersonate` | 28.6% | ❌ |
| `internal/helpers` | 28.7% | ❌ |
| `internal/tasks/hello_world` | 70.6% | ⚠️ |

### No Coverage (0%)

The following packages have no test coverage:

- `internal/controllers/admin/blog/*` (all blog admin controllers)
- `internal/controllers/admin/shop/shared`
- `internal/controllers/admin/stats`
- `internal/controllers/admin/tasks`
- `internal/controllers/admin/users` (except user_impersonate and user_update)
- `internal/controllers/admin/users/user_create`
- `internal/controllers/admin/users/user_delete`
- `internal/controllers/admin/users/user_manager`
- `internal/controllers/liveflux`
- `internal/controllers/shared/cdn`
- `internal/controllers/user/partials`
- `internal/controllers/website/*` (except home, seo, pages/indexnow)
- `internal/emails`
- `internal/layouts`
- `internal/links`
- `internal/resources`
- `internal/schedules`
- `internal/tasks/stats`
- `internal/widgets`
- `pkg/blogai`
- `pkg/testimonials`

## Recommendations

### Priority 1: Fix Build Failures

1. **Investigate blogstore package version/API**
   - Check if the blogstore package version is compatible with the code
   - Verify that the required methods exist in the installed version
   - Update the blogstore dependency if needed or update the controller code to use available methods

2. **Fix CMS Controller**
   - Review `internal/controllers/website/cms/cms_controller.go:65`
   - Ensure the function returns both (bool, string) as expected

3. **Fix CDN Lazy Loading**
   - Investigate the undefined `cdn.Slazy_0_5_0` reference
   - Check if this is a generated file or if there's a missing import

### Priority 2: Increase Coverage

After build failures are resolved, focus on:

1. **Blog Admin Controllers** (0% coverage)
   - Add comprehensive tests for blog management functionality
   - Test category and tag management
   - Test post creation and updates

2. **Admin User Management** (0% coverage)
   - Add tests for user creation, deletion, and management
   - Test user impersonation functionality

3. **Website Controllers** (0% coverage)
   - Add tests for blog home, CMS, and contact pages
   - Test SEO functionality

4. **Low Coverage Packages** (<50%)
   - Increase coverage for authentication and registration flows
   - Add more tests for media and file handling
   - Improve coverage for helper functions

### Priority 3: Target Coverage Goals

- **Overall Target:** 80%+ coverage
- **Critical Paths:** 90%+ coverage (auth, user management, core business logic)
- **UI Controllers:** 70%+ coverage (acceptable for view-heavy code)

## Test Execution Summary

**Total Packages Analyzed:** 50+  
**Packages with Coverage:** ~35  
**Packages with 0% Coverage:** ~15  
**Packages with Build Failures:** 5+

## Next Steps

1. Resolve the 3 critical build failures
2. Re-run test suite with `go test -coverprofile=coverage.out ./...`
3. Generate HTML coverage report with `go tool cover -html=coverage.out`
4. Create targeted test plans for low-coverage packages
5. Establish CI/CD pipeline to track coverage over time

---

**Report Status:** Incomplete - Build failures prevent full analysis  
**Action Required:** Fix build failures before proceeding with coverage improvements
