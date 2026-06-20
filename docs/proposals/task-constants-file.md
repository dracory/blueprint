# Proposal: Task Constants File

## Status

Implemented

## Context

Blueprint's task aliases are hardcoded as string literals throughout the codebase. For example, `"CourseEnrollmentTask"` appears in both the task handler's `Alias()` method and in any controller that enqueues that task. This leads to:

1. **Typo risk** — A misspelled alias string silently fails to match at runtime.
2. **No compile-time safety** — String literals are not checked by the compiler.
3. **Greppability issues** — Finding all references to a task requires searching for the exact string, not a symbol.
4. **Inconsistency** — Some tasks define their alias inline; others might use different casing or naming.

CourseThread solved this with a centralized `internal/tasks/constants.go` file that defines all task aliases as exported constants.

## Goal

Centralize all task alias strings into a single constants file (`internal/tasks/constants.go`) to provide compile-time safety and a single source of truth.

## Proposed Implementation

### New File

#### `internal/tasks/constants.go`

```go
package tasks

// Task aliases used throughout the application.
// These constants are referenced by task handlers and any code
// that enqueues tasks via TaskDefinitionEnqueueByAlias().
const (
    // BlindIndexRebuildTaskAlias is the alias for the blind index rebuild task.
    BlindIndexRebuildTaskAlias = "BlindIndexRebuildTask"

    // CleanUpTaskAlias is the alias for the cleanup task.
    CleanUpTaskAlias = "CleanUpTask"

    // EmailTestTaskAlias is the alias for the email test task.
    EmailTestTaskAlias = "EmailTestTask"

    // EmailToAdminTaskAlias is the alias for the admin notification email task.
    EmailToAdminTaskAlias = "EmailToAdminTask"

    // EmailToAdminOnNewContactFormSubmittedTaskAlias is the alias for the
    // contact form submission admin notification task.
    EmailToAdminOnNewContactFormSubmittedTaskAlias = "EmailToAdminOnNewContactFormSubmittedTask"

    // EmailToAdminOnNewUserRegisteredTaskAlias is the alias for the
    // new user registration admin notification task.
    EmailToAdminOnNewUserRegisteredTaskAlias = "EmailToAdminOnNewUserRegisteredTask"

    // HelloWorldTaskAlias is the alias for the hello world task.
    HelloWorldTaskAlias = "HelloWorldTask"

    // StatsVisitorEnhanceTaskAlias is the alias for the stats visitor
    // enhancement task.
    StatsVisitorEnhanceTaskAlias = "StatsVisitorEnhanceTask"
)
```

### Modified Files

Each task handler's `Alias()` method is updated to reference the constant instead of a string literal:

#### `internal/tasks/blind_index_rebuild/blind_index_rebuild_task.go`

```go
// Before:
func (h *blindIndexRebuildTask) Alias() string {
    return "BlindIndexRebuildTask"
}

// After:
func (h *blindIndexRebuildTask) Alias() string {
    return tasks.BlindIndexRebuildTaskAlias
}
```

**Note:** This creates an import cycle (`tasks` → `blind_index_rebuild` → `tasks`). To resolve this, the constants file should live in a sub-package or separate package. Two options:

**Option A: `internal/taskconstants/` package (recommended)**

Create `internal/taskconstants/taskconstants.go` with the constants. Both `internal/tasks/` sub-packages and external callers import it. No cycle.

**Option B: Keep constants in each task sub-package**

Each task sub-package exports its own `Alias` constant. Less centralized but no cycle.

### Recommended: Option A — `internal/taskconstants/`

#### `internal/taskconstants/taskconstants.go`

```go
// Package taskconstants defines task alias constants used throughout
// the application for task registration and enqueueing.
package taskconstants

const (
    BlindIndexRebuildTaskAlias                          = "BlindIndexRebuildTask"
    CleanUpTaskAlias                                    = "CleanUpTask"
    EmailTestTaskAlias                                  = "EmailTestTask"
    EmailToAdminTaskAlias                               = "EmailToAdminTask"
    EmailToAdminOnNewContactFormSubmittedTaskAlias      = "EmailToAdminOnNewContactFormSubmittedTask"
    EmailToAdminOnNewUserRegisteredTaskAlias            = "EmailToAdminOnNewUserRegisteredTask"
    HelloWorldTaskAlias                                 = "HelloWorldTask"
    StatsVisitorEnhanceTaskAlias                        = "StatsVisitorEnhanceTask"
)
```

Each task handler imports `taskconstants` and returns the constant:

```go
import "project/internal/taskconstants"

func (h *blindIndexRebuildTask) Alias() string {
    return taskconstants.BlindIndexRebuildTaskAlias
}
```

Any controller that enqueues a task also uses the constant:

```go
// Before:
app.GetTaskStore().TaskDefinitionEnqueueByAlias(ctx, "default", "HelloWorldTask", params)

// After:
app.GetTaskStore().TaskDefinitionEnqueueByAlias(ctx, "default", taskconstants.HelloWorldTaskAlias, params)
```

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/taskconstants/taskconstants.go` | Create | Centralized task alias constants |
| `internal/tasks/blind_index_rebuild/blind_index_rebuild_task.go` | Modify | Use constant for alias |
| `internal/tasks/clean_up/clean_up_task.go` | Modify | Use constant for alias |
| `internal/tasks/email_test/email_test_task.go` | Modify | Use constant for alias |
| `internal/tasks/email_admin/email_admin_task.go` | Modify | Use constant for alias |
| `internal/tasks/email_admin_new_contact/email_admin_new_contact_task.go` | Modify | Use constant for alias |
| `internal/tasks/email_admin_new_user_registered/email_admin_new_user_registered_task.go` | Modify | Use constant for alias |
| `internal/tasks/hello_world/hello_world_task.go` | Modify | Use constant for alias |
| `internal/tasks/stats/stats_visitor_enhance_task.go` | Modify | Use constant for alias |
| `internal/tasks/cms_transfer_task.go` | Modify | Use constant for alias (if applicable) |
| Any controller enqueuing tasks | Modify | Replace string literals with constants |

## Testing

- `go test ./internal/taskconstants/...` — verify constants are non-empty
- `go test ./internal/tasks/...` — existing tests pass with constant references
- `go build ./...` — compile-time verification that all references resolve

## Verification

- `go build ./...` succeeds
- `go test ./...` passes
- `grep -r "TaskDefinitionEnqueueByAlias" --include="*.go"` — no hardcoded task alias strings remain
- Each task handler's `Alias()` returns the constant, not a literal

## Benefits

- **Compile-time safety** — Misspelled constant names fail to compile
- **Single source of truth** — All task aliases in one file
- **Greppable** — Find all usages of a task by searching for the constant name
- **Self-documenting** — Each constant has a doc comment explaining its purpose
- **Consistent** — Naming convention enforced by the constants file

## Risks

- **Import cycle** — Mitigated by using a separate `taskconstants` package
- **Mechanical change** — Many files modified but each change is a simple string → constant replacement
