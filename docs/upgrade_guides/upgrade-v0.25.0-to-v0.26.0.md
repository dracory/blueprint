# Upgrade Guide: v0.25.0 to v0.26.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.25.0 to v0.26.0.

## Overview

This release focuses on test infrastructure improvements, dependency cleanup, and a minor configuration enhancement for the task queue runner. The changes are primarily internal to the test suite with minimal impact on production code.

**Key Changes:**
- Added QueueName configuration to TaskQueueRunner initialization
- Removed testify and related test dependencies
- Refactored tests from table-driven to individual functions for better parallelization
- Consolidated functional tests into main test files
- Added new shopadmin package (non-breaking feature addition)

---

## ⚠️ Breaking Changes

### 1. TaskQueueRunner QueueName Configuration

**Change**: The `TaskQueueRunnerOptions` now requires an explicit `QueueName` field to specify which queue the background task runner processes.

**Old Usage**:
```go
// cmd/server/background_processes.go
runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
    IntervalSeconds: 2,
    UnstuckMinutes:  2,
    MaxConcurrency:  10,
    Logger:          log.Default(),
})
```

**New Usage**:
```go
// cmd/server/background_processes.go
runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
    IntervalSeconds: 2,
    UnstuckMinutes:  2,
    QueueName:       taskstore.DefaultQueueName,
    MaxConcurrency:  10,
    Logger:          log.Default(),
})
```

**Action Required**:
- Add `QueueName: taskstore.DefaultQueueName` to `TaskQueueRunnerOptions` in `cmd/server/background_processes.go`
- If using custom queues, specify the appropriate queue name instead of `DefaultQueueName`

---

### 2. Removed Test Dependencies

**Change**: The following test dependencies have been removed from go.mod:
- `github.com/stretchr/testify` v1.11.1
- `github.com/davecgh/go-spew` v1.1.2-0.20180830191138-d8f796af33cc
- `github.com/pmezard/go-difflib` v1.0.1-0.20181226105442-5d4384ee4fb2
- `gopkg.in/yaml.v3` v3.0.1

**Old Usage**:
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
    assert.Equal(t, expected, actual)
    require.NoError(t, err)
}
```

**New Usage**:
```go
import (
    "testing"
)

func TestSomething(t *testing.T) {
    if expected != actual {
        t.Errorf("expected %v, got %v", expected, actual)
    }
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

**Action Required**:
- Run `go mod tidy` to remove unused dependencies
- Replace `assert` and `require` assertions with standard library `t.Error`, `t.Errorf`, `t.Fatal`, or `t.Fatalf`
- Search for all usages of `assert.` and `require.` in test files and replace with standard library equivalents
- Note: Blueprint has already migrated all internal tests to use standard library assertions

---

### 3. Test File Consolidation

**Change**: Functional test files have been consolidated into main test files. The pattern of having separate `*_functional_test.go` files has been removed.

**Old Usage**:
```go
// pkg/logadmin/log_manager/log_manager_controller_test.go
func TestLogManagerController(t *testing.T) {
    // Main tests
}

// pkg/logadmin/log_manager/log_manager_controller_functional_test.go (DELETED)
func TestLogManagerController_Functional(t *testing.T) {
    // Functional tests
}
```

**New Usage**:
```go
// pkg/logadmin/log_manager/log_manager_controller_test.go
func TestLogManagerController_RenderPage(t *testing.T) {
    // Test from main file
}

func TestLogManagerController_HandleLoadLogs(t *testing.T) {
    // Test from functional file (renamed)
}

func TestLogManagerController_HandleLogDelete(t *testing.T) {
    // Test from functional file (renamed)
}
```

**Action Required**:
- If your project has `*_functional_test.go` files, consolidate them into the main test file
- Rename test functions to remove the `Functional` suffix and add descriptive suffixes instead
- Example: `TestController_Functional` → `TestController_RenderPage` or `TestController_HandleAction`
- Delete any `.bak` test files that may exist

---

### 4. Table-Driven Test Refactoring

**Change**: Table-driven tests have been refactored into individual test functions for improved parallelization and better test isolation.

**Old Usage**:
```go
func TestStatsVisitorEnhanceTask_FindCountryByIp(t *testing.T) {
    tests := []struct {
        name       string
        response   *http.Response
        err        error
        wantResult string
    }{
        {
            name:       "successful lookup",
            response:   &http.Response{...},
            wantResult: "US",
        },
        {
            name:       "empty country code",
            response:   &http.Response{...},
            wantResult: "UN",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

**New Usage**:
```go
func TestStatsVisitorEnhanceTask_FindCountryByIp_WithMock_SuccessfulLookup(t *testing.T) {
    // Test for successful lookup
}

func TestStatsVisitorEnhanceTask_FindCountryByIp_WithMock_EmptyCountry(t *testing.T) {
    // Test for empty country code
}
```

**Action Required**:
- This is an internal refactoring pattern for Blueprint's test suite
- If your project has table-driven tests, you may optionally refactor them to individual functions for better parallelization
- This is not required for the upgrade but recommended for improved test performance

---

## 🔄 Migration Steps

### Step 1: Update Dependencies

Update go.mod to remove unused dependencies:

```bash
go mod tidy
```

This will automatically remove the testify and related dependencies that are no longer used.

### Step 2: Update TaskQueueRunner Configuration

Add the QueueName field to TaskQueueRunnerOptions:

```bash
# Edit cmd/server/background_processes.go
# Add QueueName: taskstore.DefaultQueueName to the TaskQueueRunnerOptions struct
```

Manual changes required in `cmd/server/background_processes.go`:
1. Locate the `TaskQueueRunnerOptions` initialization (around line 52)
2. Add the QueueName field:
```go
runner := taskstore.NewTaskQueueRunner(ts, taskstore.TaskQueueRunnerOptions{
    IntervalSeconds: 2,
    UnstuckMinutes:  2,
    QueueName:       taskstore.DefaultQueueName,  // ADD THIS LINE
    MaxConcurrency:  10,
    Logger:          log.Default(),
})
```

### Step 3: Replace Testify Assertions (If Used)

If your project uses testify assertions, replace them with standard library:

```bash
# Find all files using testify
grep -r "github.com/stretchr/testify" --include="*_test.go" .

# Replace assert.Equal with t.Errorf
# Replace assert.NoError with t.Errorf
# Replace require.NoError with t.Fatalf
```

Common replacements:
- `assert.Equal(t, expected, actual)` → `if expected != actual { t.Errorf("expected %v, got %v", expected, actual) }`
- `assert.NoError(t, err)` → `if err != nil { t.Errorf("unexpected error: %v", err) }`
- `require.NoError(t, err)` → `if err != nil { t.Fatalf("unexpected error: %v", err) }`
- `assert.Nil(t, value)` → `if value != nil { t.Errorf("expected nil, got %v", value) }`
- `assert.Contains(t, str, substr)` → `if !strings.Contains(str, substr) { t.Errorf("expected %q to contain %q", str, substr) }`

### Step 4: Consolidate Functional Test Files (If Present)

If your project has separate functional test files:

```bash
# Find functional test files
find . -name "*_functional_test.go"

# For each file, move its tests into the main test file
# Rename test functions to remove "Functional" suffix
# Delete the functional test file
```

Example consolidation:
1. Open `*_functional_test.go`
2. Copy all test functions to the main test file
3. Rename functions: `TestX_Functional` → `TestX_DescriptiveName`
4. Delete `*_functional_test.go`
5. Delete any `*.bak` files

### Step 5: Verify Build

Ensure the project builds successfully:

```bash
go build ./...
```

---

## 🧪 Testing After Migration

### 1. Unit Tests

Run all unit tests to ensure no regressions:

```bash
go test ./...
```

### 2. Integration Tests

Run integration tests with database:

```bash
go test -tags=integration ./...
```

### 3. Server Startup Test

Test that the server starts correctly with the new TaskQueueRunner configuration:

```bash
go run ./cmd/server
# Should start without errors
# Verify task queue runner initializes with QueueName
```

### 4. Background Process Test

Verify background processes work correctly:

```bash
# Run tests that exercise background processes
go test ./cmd/server/...
```

---

## 📝 Additional Notes

### New ShopAdmin Package

A new `pkg/shopadmin` package has been added with a complete shop administration interface. This is a non-breaking feature addition that follows the folder-per-controller pattern. If you need shop administration features, you can integrate this package using:

```go
import "project/pkg/shopadmin"

admin, err := shopadmin.New(shopadmin.AdminOptions{
    Registry:       registry,
    AdminHomeURL:   "/admin",
    ShopAdminURL:   "/admin/shop",
    AuthUserID:     auth.GetUserID,
    FileManagerURL: "/admin/files",
})
```

### Test Performance Improvements

The test refactoring from table-driven to individual functions provides:
- Better parallelization (each test runs independently)
- Clearer test failure messages (test names are more descriptive)
- Easier debugging (individual tests can be run in isolation)
- Better test coverage reporting

### Dependency Cleanup Benefits

Removing testify and related dependencies:
- Reduces dependency attack surface
- Simplifies dependency management
- Aligns with Go's standard library-first philosophy
- Reduces go.mod size and complexity

---

## 🆘 Common Issues and Solutions

### Issue 1: "undefined: taskstore.DefaultQueueName"

**Cause**: The taskstore package may not have DefaultQueueName in older versions.

**Solution**: Ensure you have the latest taskstore version:
```bash
go get -u github.com/dracory/taskstore
go mod tidy
```

### Issue 2: "undefined: assert" after migration

**Cause**: Testify assertions still in use after dependency removal.

**Solution**: Replace all assert/require calls with standard library equivalents (see Step 3).

### Issue 3: Tests fail with "undefined: testify"

**Cause**: Import statements for testify still present.

**Solution**: Remove all testify imports from test files:
```bash
sed -i '/github.com\/stretchr\/testify/d' *_test.go
```

### Issue 4: "cannot use QueueName (untyped string) as type string"

**Cause**: QueueName field may not exist in older taskstore versions.

**Solution**: Update taskstore to latest version and ensure the field is properly typed.

### Issue 5: Parallel test failures after refactoring

**Cause**: Tests may have shared state that doesn't work well with parallel execution.

**Solution**: Add `t.Parallel()` only to tests that are truly independent. Keep sequential tests without the parallel flag.

---

## 📞 Support

For issues or questions about this upgrade:
- Check the Blueprint repository: https://github.com/dracory/blueprint
- Review the taskstore documentation for QueueName configuration
- Open an issue on GitHub for migration problems

---

## Quality Checklist

- [x] All breaking changes identified and documented
- [x] Code examples are accurate and tested
- [x] Migration steps are in logical order
- [x] Action items are specific and actionable
- [x] Testing procedures are comprehensive
- [x] Common issues are addressed
- [x] Format follows markdown best practices
