---
description: Upgrade a Blueprint-based project to the latest version
---

# Project Upgrader

This workflow upgrades a Blueprint-based project to the latest version by following the appropriate upgrade guides.

## When to Use

Use this workflow when:
- You need to upgrade a Blueprint-based project to the latest version
- A new Blueprint version has been released and you want to apply the changes
- You need to migrate through multiple intermediate versions

## Prerequisites

1. The project to upgrade must be a Blueprint-based application
2. You must have access to the Blueprint repository (D:\PROJECTs\dracory.com\blueprint)
3. The project must have a version constant in `<project_path>/internal/config/version.go`

## Steps to Upgrade Project

### 1. Identify Project Path

Specify the project path to upgrade. If the project path is not known, ask the user to provide it:

```
<project_path>
```

The project is a web application that follows the pattern set by Blueprint.

### 2. Version Comparison

**Read the project's current version**:
- Check the file: `<project_path>/internal/config/version.go`
- Read the `Version` constant to determine the current version
- Example: `const Version = "0.22.0"`

**Read the latest Blueprint version**:
- Check the file: `D:\PROJECTs\dracory.com\blueprint\internal/config/version.go`
- Read the `Version` constant to determine the latest version
- Example: `const Version = "0.23.0"`

**Determine the upgrade path**:
- Compare the two versions
- If versions are consecutive (e.g., v0.22.0 to v0.23.0), use the single upgrade guide: `upgrade-v{CURRENT}-to-v{LATEST}.md`
- If there are multiple version gaps (e.g., v0.20.0 to v0.23.0), you must upgrade through each intermediate version sequentially:
  - First: `upgrade-v0.20.0-to-v0.21.0.md`
  - Then: `upgrade-v0.21.0-to-v0.22.0.md`
  - Finally: `upgrade-v0.22.0-to-v0.23.0.md`
- After each upgrade step, update the project's version constant before proceeding to the next step

### 3. Locate Upgrade Guides

The upgrade guides can be found at:

```
D:\PROJECTs\dracory.com\blueprint\docs\upgrade_guides\upgrade-[FROM_VERSION]-to-[TO_VERSION].md
```

### 4. Apply Upgrade Guide

For each upgrade guide in the sequence:

**Follow the upgrade guide step-by-step to**:
- Update dependencies in go.mod
- Apply breaking changes
- Update import paths
- Modify configuration
- Run tests to verify the upgrade

**Important**: Some features or changes mentioned in the upgrade guide may already be implemented in the project. Before applying each change, check if the feature already exists in the project files. If it does, skip that step to avoid duplication or conflicts.

### 5. Make Necessary Changes

Make the necessary changes to the project files to upgrade it to the latest version. You may compare the project files with the upgrade guides for Blueprint to understand the changes that need to be made.

**Reference Implementation**: The Blueprint project files should be used as a reference for the latest version of Blueprint. If in doubt, refer to the Blueprint project files at:

```
D:\PROJECTs\dracory.com\blueprint
```

### 6. Update Version Constant

After completing each upgrade step, update the project's version constant in `<project_path>/internal/config/version.go` to reflect the new version before proceeding to the next upgrade guide (if upgrading through multiple versions).

### 7. Test the Project

After making the necessary changes, test the project to ensure that it is working correctly:

```bash
# Navigate to project directory
cd <project_path>

# Run tests
go test ./...

# Build the application
go build -o ./tmp/main ./cmd/server

# Run the application (manual testing)
go run ./cmd/server
```

### 8. Verify Upgrade

Verify that the upgrade was successful by:
- Checking that the version constant matches the latest Blueprint version
- Running all tests and ensuring they pass
- Building the application successfully
- Testing key functionality manually

## Quality Checklist

- [ ] Project version constant updated to latest Blueprint version
- [ ] All upgrade guides applied in correct order (for multi-version upgrades)
- [ ] Dependencies updated in go.mod
- [ ] Breaking changes applied correctly
- [ ] Import paths updated
- [ ] Configuration modified as needed
- [ ] All tests pass
- [ ] Application builds successfully
- [ ] Key functionality tested manually
- [ ] No duplicate features added (checked before applying each change)

## Example

To upgrade a project from v0.22.0 to v0.23.0:

1. Identify project path: `D:\PROJECTs\myproject`
2. Read project version: `D:\PROJECTs\myproject\internal/config/version.go` → `0.22.0`
3. Read latest Blueprint version: `D:\PROJECTs\dracory.com\blueprint\internal/config/version.go` → `0.23.0`
4. Determine upgrade path: v0.22.0 → v0.23.0 (consecutive, single guide)
5. Locate upgrade guide: `D:\PROJECTs\dracory.com\blueprint\docs\upgrade_guides\upgrade-v0.22.0-to-v0.23.0.md`
6. Follow the upgrade guide step-by-step
7. Check if features already exist before applying changes
8. Update project version constant to `0.23.0`
9. Run tests: `cd D:\PROJECTs\myproject && go test ./...`
10. Build and test manually

**Example for Multi-Version Upgrade** (v0.20.0 to v0.23.0):

1. Identify project path: `D:\PROJECTs\myproject`
2. Read project version: `0.20.0`
3. Read latest Blueprint version: `0.23.0`
4. Determine upgrade path: v0.20.0 → v0.21.0 → v0.22.0 → v0.23.0 (3 intermediate versions)
5. Apply upgrade-v0.20.0-to-v0.21.0.md
6. Update version constant to `0.21.0`
7. Apply upgrade-v0.21.0-to-v0.22.0.md
8. Update version constant to `0.22.0`
9. Apply upgrade-v0.22.0-to-v0.23.0.md
10. Update version constant to `0.23.0`
11. Run tests and verify

## Common Issues

### Issue: Build fails after dependency update
**Solution**: Run `go mod tidy` to resolve dependency conflicts

### Issue: Import path errors after upgrade
**Solution**: Ensure all import paths have been updated according to the upgrade guide

### Issue: Tests fail after upgrade
**Solution**: Review the breaking changes in the upgrade guide and ensure all API changes have been applied correctly

### Issue: Feature already exists in project
**Solution**: Skip that step in the upgrade guide to avoid duplication

## Support

For additional help:
- Review the upgrade guides in `D:\PROJECTs\dracory.com\blueprint\docs\upgrade_guides/`
- Compare with the Blueprint reference implementation at `D:\PROJECTs\dracory.com\blueprint`
- Check the version workflow: `D:\PROJECTs\dracory.com\blueprint\docs\version_workflow.md`
