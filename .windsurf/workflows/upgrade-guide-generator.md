---
description: Generate comprehensive upgrade guide between two Blueprint versions
---

# Upgrade Guide Generator

This skill generates comprehensive upgrade guides for Blueprint framework versions.

## When to Use

Use this skill when:
- A new Blueprint version is ready for release
- You need to document breaking changes between versions
- You want to help users upgrade from one version to another

## Prerequisites

1. Both versions should have git tags (e.g., `v0.22.0`, `v0.23.0`)
2. The version constant in `internal/config/version.go` should be updated
3. The release branch should be merged to main

## Steps to Generate Upgrade Guide

### 1. Identify Versions

Read the current version from `internal/config/version.go` to determine the new version.

Example:
- Current version: `0.23.0`
- Previous version: `0.22.0` (latest tag)

### 2. Compare Versions

Use git to compare the versions:

```bash
# Get list of commits between versions
git log --oneline v0.22.0..v0.23.0

# Get diff summary
git diff --stat v0.22.0..v0.23.0

# View detailed changes
git diff v0.22.0..v0.23.0
```

### 3. Analyze Changes

Examine the following areas:

**File Location Changes**
- Files moved or renamed
- Package reorganization
- New directories created

**Import Path Changes**
- Package imports that changed
- New packages added
- Old packages removed

**API Signature Changes**
- Method parameters changed
- Return types changed
- Method renames
- Interface changes

**Configuration Changes**
- New environment variables
- Config method changes
- Config structure changes

**Dependency Updates**
- Go module version changes in `go.mod`
- New dependencies added
- Dependencies removed

**Architecture Changes**
- Pattern changes (globals to registry, etc.)
- New architectural patterns introduced

### 4. Identify Breaking Changes

Categorize changes as breaking if they:
- Require code changes in applications using Blueprint
- Change method signatures or return types
- Remove or rename existing APIs
- Change default behaviors
- Require environment variable changes

### 5. Generate Upgrade Guide

Create a new file: `docs/upgrade_guides/upgrade-v{FROM}-to-v{TO}.md`

Follow the template in `docs/upgrade_guides/upgrade-guide-prompt.md`

**Structure:**

```markdown
# Upgrade Guide: v{FROM} to v{TO}

This guide helps LLMs and developers upgrade Blueprint applications from v{FROM} to v{TO}.

## Summary

[Brief summary of major changes]

**Key Changes:**
- [List major changes]

---

## Breaking Changes

### 1. [Change Title]

**Change**: [Description]

**Old Usage**:
```go
[Code example]
```

**New Usage**:
```go
[Code example]
```

**Action Required**:
- [Specific actions]

**Migration Command**:
```bash
[Commands if applicable]
```

---

## Migration Steps

### Step 1: [Step Title]
[Instructions with commands]

---

## Testing After Migration

### 1. Unit Tests
[Test instructions]

---

## Additional Notes

### New Features
[List new features]

### Removed Features
[List removed features]

---

## Common Issues and Solutions

### Issue 1: [Issue Title]
**Symptom**: [Description]
**Solution**: [Steps]

---

## Support

[Support information]
```

### 6. Include Migration Commands

For automated changes, provide bash commands:

```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|old/path|new/path|g' {} \;

# Update method names
find . -type f -name "*.go" -exec sed -i 's|OldMethod|NewMethod|g' {} \;
```

### 7. Test the Guide

After generating the guide:
1. Review for accuracy
2. Test migration commands on a sample project
3. Verify all breaking changes are documented
4. Check code examples are correct
5. Ensure migration steps are in logical order

### 8. Save the Guide

Save to: `docs/upgrade_guides/upgrade-v{FROM}-to-v{TO}.md`

## Quality Checklist

- [ ] All breaking changes identified and documented
- [ ] Code examples are accurate and tested
- [ ] Migration steps are in logical order
- [ ] Action items are specific and actionable
- [ ] Testing procedures are comprehensive
- [ ] Common issues are addressed
- [ ] Format follows markdown best practices
- [ ] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`

## Example

To generate an upgrade guide from v0.22.0 to v0.23.0:

1. Read `internal/config/version.go` → Version is "0.23.0"
2. Run: `git log --oneline v0.22.0..v0.23.0`
3. Analyze commits and diffs for breaking changes
4. Document each breaking change with old/new usage
5. Create migration steps with commands
6. Save to `docs/upgrade_guides/upgrade-v0.22.0-to-v0.23.0.md`
