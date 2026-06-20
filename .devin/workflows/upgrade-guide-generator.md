---
description: Generate comprehensive upgrade guide between two Blueprint versions
---

# Upgrade Guide Generator

This skill generates comprehensive upgrade guides for Blueprint rapid application development (RAD) starter template versions.

## When to Use

Use this skill when:
- A new Blueprint version is ready for release (see `docs/version_workflow.md` step 3)
- You need to document breaking changes between versions
- You want to help users upgrade from one version to another

**Timing**: Generate the upgrade guide **before** merging the release branch to main and tagging the release. This is step 3 in the version workflow.

## Prerequisites

1. The version constant in `internal/config/version.go` should be updated to the new version
2. The release branch (e.g., `release/v0.23.0`) should be ready for release
3. Previous version should have a git tag (e.g., `v0.22.0`)
4. **Note**: Generate upgrade guide **before** merging to main and tagging (see `docs/version_workflow.md` step 3)

## Steps to Generate Upgrade Guide

### 1. Identify Versions

Read the current version from `internal/config/version.go` to determine the new version.

Example:
- Current version: `0.23.0`
- Previous version: `0.22.0` (latest tag)

**Version Gap Handling**:
- If versions are consecutive (e.g., v0.22.0 → v0.23.0), generate a single upgrade guide
- If there are version gaps (e.g., v0.20.0 → v0.23.0), generate multiple upgrade guides for each intermediate version:
  - `upgrade-v0.20.0-to-v0.21.0.md`
  - `upgrade-v0.21.0-to-v0.22.0.md`
  - `upgrade-v0.22.0-to-v0.23.0.md`

### 2. Verify Git Tags

Verify that the previous version has a git tag:

```bash
# List all version tags
git tag -l "v*"

# Verify specific tag exists
git tag -l "v0.22.0"
```

If the tag doesn't exist, you may need to tag the previous version first or use the commit hash instead.

### 3. Compare Versions

Use git to compare the versions:

```bash
# Get list of commits between versions
git log --oneline v0.22.0..v0.23.0

# Get diff summary
git diff --stat v0.22.0..v0.23.0

# View detailed changes
git diff v0.22.0..v0.23.0
```

### 4. Analyze Changes

**Reference Previous Guides**: Review existing upgrade guides in `docs/upgrade_guides/` to maintain consistency in style, structure, and terminology.

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

**Entry Point Changes**
- main.go location changes
- cmd/ structure changes
- Application entry point modifications

**Store/Task API Changes**
- Enqueue method changes
- Initialization changes
- Store interface modifications

### 5. Identify Breaking Changes

Categorize changes as breaking if they:
- Require code changes in applications using Blueprint
- Change method signatures or return types
- Remove or rename existing APIs
- Change default behaviors
- Require environment variable changes

### 6. Generate Upgrade Guide

**Section Styling**: Use emojis for section headers to match existing guides:
- ⚠️ for Breaking Changes
- 🔄 for Migration Steps  
- 🧪 for Testing After Migration
- 📝 for Additional Notes
- 🆘 for Common Issues and Solutions

Create a new file: `docs/upgrade_guides/upgrade-v{FROM}-to-v{TO}.md`

Follow the template in `docs/upgrade_guides/upgrade-guide-prompt.md`

**Structure:**

```markdown
# Upgrade Guide: v{FROM} to v{TO}

This guide helps LLMs and developers upgrade Blueprint applications from v{FROM} to v{TO}.

## Overview

[Brief overview of major changes]

**Key Changes:**
- [List major changes]

---

## ⚠️ Breaking Changes

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

### 7. Include Migration Commands

For automated changes, provide bash commands:

```bash
# Update imports
find . -type f -name "*.go" -exec sed -i 's|old/path|new/path|g' {} \;

# Update method names
find . -type f -name "*.go" -exec sed -i 's|OldMethod|NewMethod|g' {} \;
```

### 8. Test the Guide

After generating the guide:
1. Review for accuracy
2. Test migration commands on a sample project
3. Verify all breaking changes are documented
4. Check code examples are correct
5. Ensure migration steps are in logical order

### 9. Save the Guide

Save to: `docs/upgrade_guides/upgrade-v{FROM}-to-v{TO}.md`

### 10. Request Review

After saving the guide, request review from the maintainer or team to ensure:
- All breaking changes are accurately documented
- Migration steps are complete and tested
- Code examples are correct and copy-pasteable
- Style and structure match existing guides
- No critical information is missing

## Quality Checklist

- [ ] All breaking changes identified and documented
- [ ] Code examples are accurate and tested
- [ ] Migration steps are in logical order
- [ ] Action items are specific and actionable
- [ ] Testing procedures are comprehensive
- [ ] Common issues are addressed
- [ ] Format follows markdown best practices
- [ ] File naming follows pattern: `upgrade-vX.Y.Z-to-vX.Y.Z.md`
- [ ] Emoji styling used consistently (⚠️, 🔄, 🧪, 📝, 🆘)
- [ ] Version gaps handled correctly (multiple guides if needed)
- [ ] Git tag verified for previous version
- [ ] Previous guides reviewed for consistency
- [ ] Quality checklist included in generated guide

## Example

To generate an upgrade guide from v0.22.0 to v0.23.0:

1. Read `internal/config/version.go` → Version is "0.23.0"
2. Verify git tag exists: `git tag -l "v0.22.0"`
3. Run: `git log --oneline v0.22.0..v0.23.0`
4. Review previous upgrade guides for consistency
5. Analyze commits and diffs for breaking changes
6. Document each breaking change with old/new usage
7. Create migration steps with commands
8. Apply emoji styling to sections
9. Save to `docs/upgrade_guides/upgrade-v0.22.0-to-v0.23.0.md`
10. Request review from maintainer

**Example for Version Gap** (v0.20.0 to v0.23.0):

1. Identify gap: v0.20.0 → v0.23.0 has 2 intermediate versions
2. Generate three separate guides:
   - `upgrade-v0.20.0-to-v0.21.0.md`
   - `upgrade-v0.21.0-to-v0.22.0.md`
   - `upgrade-v0.22.0-to-v0.23.0.md`
3. Each guide follows the same process as above
4. Users must apply guides in order
