# Version Workflow

This document describes the version workflow for the Blueprint rapid application development (RAD) starter template.

## Overview

Blueprint uses a hardcoded version constant in `internal/config/version.go` to track the template version. This enables AI tools and upgrade scripts to easily identify the current version and determine the appropriate migration path.

## Version Constant

The version is defined in `internal/config/version.go`:

```go
const Version = "0.23.0"
```

This constant is:
- **Hardcoded** for easy AI parsing during upgrades
- **Accessible** via `config.GetVersion()` function
- **Logged** on server startup for visibility
- **Used** in upgrade guide generation

## Release Workflow

### 1. Start New Version

When beginning work on a new version:

```bash
# Create release branch
git checkout -b release/v0.23.0

# Update version constant in internal/config/version.go
# Change: const Version = "0.23.0"

# Commit version bump
git add internal/config/version.go
git commit -m "feat: bump version to 0.23.0"
```

### 2. Develop Features

- Work on features in the release branch
- Do not modify the version constant during development
- Follow existing coding standards and patterns

### 3. Prepare Release

When the version is ready for release:

```bash
# Ensure all tests pass
go test ./...

# Build the application
go build -o ./tmp/main ./cmd/server

# Generate upgrade guide using AI skill
# (See .windsurf/workflows/upgrade-guide-generator.md)
```

### 4. Release

```bash
# Switch to main branch
git checkout main

# Merge release branch
git merge release/v0.23.0

# Tag the release
git tag v0.23.0

# Push to remote
git push
git push --tags
```

### 5. Start Next Version

Immediately after releasing:

```bash
# Create next release branch
git checkout -b release/v0.24.0

# Update version constant in internal/config/version.go
# Change: const Version = "0.24.0"

# Commit version bump
git add internal/config/version.go
git commit -m "feat: bump version to 0.24.0"

# Push branch
git push -u origin release/v0.24.0
```

## Version Visibility

The version is visible in multiple places:

1. **Server Startup**: Displayed in logs when the application starts
   ```
   Starting Blueprint v0.23.0
   ```

2. **Code**: Accessible via `config.GetVersion()`

3. **Upgrade Guides**: Used to generate upgrade guides between versions

4. **Git Tags**: Each release is tagged with the version number

## Benefits

- **AI-Friendly**: Hardcoded version allows AI to easily identify project version during upgrades
- **Clear Workflow**: Branch-based versioning prevents working directly on main
- **Upgrade Automation**: Version constant enables automated upgrade guide generation
- **Visibility**: Version in logs provides easy identification of running version
- **Traceability**: Each release has clear branch and tag for tracking changes

## Upgrade Guide Generation

When generating upgrade guides, the AI can:

1. Read the current version from `internal/config/version.go`
2. Compare with previous version tags
3. Generate appropriate upgrade guide in `docs/upgrade_guides/upgrade-v{FROM}-to-v{TO}.md`
4. Follow the template in `docs/upgrade_guides/upgrade-guide-prompt.md`

## Example

Current state:
- Latest released version: v0.22.0
- Current development version: v0.23.0
- Branch: release/v0.23.0
- Version constant: `const Version = "0.23.0"`

When v0.23.0 is ready:
1. Generate upgrade guide: `upgrade-v0.22.0-to-v0.23.0.md`
2. Merge to main
3. Tag: `git tag v0.23.0`
4. Create branch: `release/v0.24.0`
5. Update version constant: `const Version = "0.24.0"`
