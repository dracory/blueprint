---
name: project-upgrader
description: Upgrade a Blueprint-based Go project to the latest Blueprint version by evaluating and applying version-specific upgrade guides. Use when the user wants to upgrade a project that follows the Blueprint architecture, migrate through multiple Blueprint versions, or apply breaking changes from a Blueprint release.
---

# Project Upgrader Skill

## Overview

This skill upgrades a Blueprint-based Go project to the latest Blueprint version. It reads the project's current version, determines the correct upgrade path (including intermediate versions), and applies each upgrade guide step-by-step while avoiding duplicate changes.

## When to Use

- A user wants to upgrade a Blueprint-based project to the latest Blueprint version.
- A user needs to migrate through multiple intermediate Blueprint versions.
- A user asks to apply changes from a new Blueprint release.
- A user mentions upgrading, migrating, or updating a Blueprint project.

## Parameters

The user must provide:
- `<project_path>` — The absolute path to the Blueprint-based project to upgrade.
- `<blueprint_path>` *(optional)* — The absolute path to the Blueprint reference repository. Defaults to `D:\PROJECTs\dracory.com\blueprint` if not provided.

If `<project_path>` is not known, ask the user to provide it.

## Prerequisites

1. The project to upgrade must be a Blueprint-based application.
2. The project must have a version constant in `<project_path>/internal/config/version.go`.
3. You must have access to the Blueprint reference repository (default: `<blueprint_path>`).

## Execution Steps

### 1. Identify Versions

1. **Read the project's current version**:
   - File: `<project_path>/internal/config/version.go`
   - Read the `Version` constant.
   - Example: `const Version = "0.22.0"`
   - If the version file or constant is not found, treat the project as version `0.16.0`.

2. **Read the latest Blueprint version**:
   - File: `<blueprint_path>/internal/config/version.go`
   - Read the `Version` constant.
   - Example: `const Version = "0.23.0"`

3. **Structural validation check** (catches version constant mismatches with actual code structure across the whole project):
   - For each key directory that commonly changes between Blueprint versions, list `.go` files in both the project and the blueprint reference. Key directories to check include:
     - `internal/config`
     - `internal/app`
     - `cmd/server`
     - `database/migrations`
     - `internal/controllers` (and sub-packages)
   - Compare the file lists. Flag any of these discrepancies:
     - **Files present in blueprint but missing in project** — suggests the project never applied that guide's file changes.
     - **Files present in project but missing in blueprint** — suggests old files that should have been removed.
     - **Different file counts** — a strong indicator of a structural mismatch.
     - **Different package names or directory renames** — e.g., `internal/registry` vs `internal/app`.
   - **Naming consistency check** — When a package, type, or concept is renamed in an upgrade (e.g., `registry` → `app`), verify that ALL occurrences were updated:
     - File names: `registry_*.go` should become `app_*.go` (especially test files).
     - Function parameter names: `registry AppInterface` should become `app AppInterface`.
     - Variable names and comments referencing the old name.
     - Error message strings referencing the old name.
   - **Completeness check** — When a feature is added to blueprint (e.g., a new middleware, a test file for a migration), check if the project has it. Missing test files are a common sign of an incomplete upgrade.
   - If structural discrepancies are found, **do not trust the version constant alone**. Read the relevant upgrade guide(s) to determine what file changes were expected, and treat the project as if it needs those guides applied.
   - Examples of structural mismatches:
     - The project claims `0.29.0` but still has `internal/config/load.go`, `functions.go`, and `defaults.go` while blueprint has `app_config.go`, `database_config.go`, etc. — the project likely missed the v0.18.0→v0.19.0 config restructure.
     - The project still has `internal/registry/registry_interface.go` while blueprint has `internal/app/app_interface.go` — the project likely missed the v0.28.0→v0.29.0 registry→app rename.
     - The project renamed `registry_implementation.go` → `app_implementation.go` but left `registry_close_test.go`, `registry_datastores_test.go`, and `registry_logger_test.go` unchanged — the rename was incomplete.
     - The project renamed `registry AppInterface` parameters to `app AppInterface` in `app_implementation.go` but left them as `registry AppInterface` in `stores_*.go` files — the rename was inconsistent across the package.

4. **Determine the upgrade path**:
   - If current == latest AND structural validation passes, report that no upgrade is needed.
   - If structural discrepancies exist despite matching version constants, plan to apply the relevant guides anyway.
   - Identify all relevant upgrade guides between the current and latest versions.
   - Guides are typically applied from oldest to newest, but the order is flexible:
     - Some guides may be skipped if the project's packages already have newer methods that make the older guide's changes irrelevant or conflicting.
     - When in doubt, compare the guide's proposed changes against the current state of the project files and the `<blueprint_path>` reference.
   - After completing the applied guides, update the project's version constant to the latest version.

### 2. Baseline Test Check

Before planning or applying any upgrades, verify the project is in a healthy state:

1. Run the project's tests:
   ```bash
   cd <project_path>
   go test ./...
   ```
2. **If any tests fail**, stop immediately and reject the upgrade. Report the failing tests to the user and do not proceed until the project is in a passing state.
3. Only continue to planning if `go test ./...` passes completely.

### 3. Create Upgrade Plan

Before making any changes to the project, create a detailed upgrade plan:

1. **List all upgrade guides** between the current and latest versions.
2. **Evaluate each guide's applicability**:
   - Check if the guide's proposed changes are still relevant given the current state of the project's packages.
   - If a guide's changes are already present or would conflict with newer package APIs, mark it to be skipped.
3. **Summarize the key changes** from each applicable guide:
   - Dependency updates
   - Breaking changes
   - Import path changes
   - Configuration changes
   - New files to add / old files to remove
4. **Identify potential risks** or conflicts based on the project's current state.
5. **Present the plan to the user**, including which guides will be applied, which will be skipped, and why. Ask for explicit approval before proceeding.

**Do not proceed with any file modifications until the user approves the plan.**

### 4. Locate Upgrade Guides

Upgrade guides are stored at:

```
<blueprint_path>/docs/upgrade_guides/upgrade-v{FROM_VERSION}-to-v{TO_VERSION}.md
```

Read the required guide(s) in order.

### 5. Apply Each Upgrade Guide

For each applicable upgrade guide:

1. **Re-evaluate the guide's relevance** before applying it. If the project's packages have evolved in a way that makes the guide's instructions outdated or conflicting, skip the guide and document the reason.
2. **Read the guide carefully** to understand all changes required.
3. **Before applying each change**, check whether the feature or change already exists in the target project files. If it does, **skip that step** to avoid duplication or conflicts.
4. **Apply the changes** to the target project:
   - Update dependencies in `go.mod`.
   - Apply breaking changes.
   - Update import paths.
   - Modify configuration files.
   - Add, remove, or refactor code as instructed.
5. **Use Blueprint as reference** — if in doubt, compare with the corresponding files in `<blueprint_path>`.

### 6. Update Version Constant

After all applicable guides have been completed, update the project's version constant in:

```
<project_path>/internal/config/version.go
```

Set it to the latest Blueprint version.

### 7. Test and Verify

After all guides have been applied:

1. Run tests:
   ```bash
   cd <project_path>
   go test ./...
   ```
2. Build the application:
   ```bash
   go build -o ./tmp/main ./cmd/server
   ```
3. If build or test failures occur:
   - Run `go mod tidy` to resolve dependency issues.
   - Review the upgrade guide for missed breaking changes or import path updates.
   - Compare with the Blueprint reference implementation to identify discrepancies.

## Quality Checklist

Before finishing, confirm:

- [ ] Project version constant matches the latest Blueprint version.
- [ ] Structural validation passes: key directories match blueprint (no missing or stale files).
- [ ] Naming consistency passes: no old names in file names, parameter names, variable names, or error messages after a rename.
- [ ] Completeness passes: every production file that should have a test file has one; every new blueprint feature is present in the project.
- [ ] All applicable upgrade guides were applied; skipped guides were documented with justification.
- [ ] `go.mod` dependencies are updated and `go mod tidy` has been run.
- [ ] Breaking changes from each guide are applied.
- [ ] Import paths are updated.
- [ ] Configuration is modified as needed.
- [ ] No duplicate features were added (checked before each change).
- [ ] `go test ./...` passes.
- [ ] `go build -o ./tmp/main ./cmd/server` succeeds.

## Common Issues

### Issue: Version constant says latest but files are outdated
**Symptom**: The version constant matches the latest Blueprint version, but the project still has old files (e.g., `internal/config/load.go`, `internal/registry/registry_interface.go`) or is missing new files that were introduced in an earlier release.
**Solution**: Trust the structural validation over the version constant. Check the relevant directories against the Blueprint reference, read the upgrade guide for the version where those files were changed, and apply the file changes even if the version constant suggests it's already done.

### Issue: Rename was applied inconsistently across the codebase
**Symptom**: A type or package was renamed (e.g., `registry` → `app`) but some files still use the old name in file names, parameter names, variable names, or error messages.
**Solution**: When a rename happens, search the entire codebase for ALL occurrences — not just imports and type references. Rename files, parameters, variables, and error messages to maintain consistency.

### Issue: Missing test files after upgrade
**Symptom**: Blueprint added `*_test.go` files for migrations, cmd/server commands, or other components, but the project doesn't have them.
**Solution**: After applying file changes from an upgrade guide, always check if blueprint has corresponding test files. If the project is missing them, add them (or port the relevant ones) to maintain parity.

### Issue: Build fails after dependency update
**Solution**: Run `go mod tidy` to resolve dependency conflicts.

### Issue: Import path errors after upgrade
**Solution**: Ensure all import paths have been updated according to the upgrade guide.

### Issue: Tests fail after upgrade
**Solution**: Review the breaking changes in the upgrade guide and ensure all API changes have been applied correctly.

### Issue: Feature already exists in project
**Solution**: Skip that step in the upgrade guide to avoid duplication.

## Important Rules

- **Always run structural validation** after reading version constants. Compare the actual files in `internal/config` (and other key directories) against the Blueprint reference. Do not rely solely on the version constant.
- **Always check for existing features** before applying a change from an upgrade guide. Skip steps that are already present in the project.
- **Evaluate guide applicability before applying**. A guide may be skipped if its proposed changes are already present, conflict with newer package APIs, or are superseded by a later guide's approach. Document the reason for any skipped guides.
- **Update the version constant once at the end** to the latest Blueprint version after all applicable guides are completed.
- **Prefer Blueprint files as reference** when the upgrade guide is ambiguous.
- **Ask the user for clarification** when uncertain about how to apply a change, which guide to prioritize, or whether to skip a step. Do not guess or make assumptions that could break the project.
- **Avoid bulk changes**. Apply changes incrementally and verify each step compiles and passes tests before moving on. Large sweeping changes across many files at once increase the risk of creating an unrecoverable mess.
