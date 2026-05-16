Your task is to upgrade a Blueprint based project to the latest version.

The project is located at:

<project_path>

The project is a web application that follows the pattern set by Blueprint.

## Version Comparison

1. **Read the project's current version**:
   - Check the file: `<project_path>/internal/config/version.go`
   - Read the `Version` constant to determine the current version
   - If the version is not found, consider the project to be at version 0.16.0
   - Example: `const Version = "0.22.0"`

2. **Read the latest Blueprint version**:
   - Check the file: `D:\PROJECTs\dracory.com\blueprint\internal/config/version.go`
   - Read the `Version` constant to determine the latest version
   - Example: `const Version = "0.23.0"`

3. **Determine the upgrade path**:
   - Compare the two versions
   - If versions are consecutive (e.g., v0.22.0 to v0.23.0), use the single upgrade guide: `upgrade-v{CURRENT}-to-v{LATEST}.md`
   - If there are multiple version gaps (e.g., v0.20.0 to v0.23.0), you must upgrade through each intermediate version sequentially:
     - First: `upgrade-v0.20.0-to-v0.21.0.md`
     - Then: `upgrade-v0.21.0-to-v0.22.0.md`
     - Finally: `upgrade-v0.22.0-to-v0.23.0.md`
   - After each upgrade step, update the project's version constant before proceeding to the next step

## Upgrade Process

You may compare the project files with the upgrade guides for Blueprint to understand the changes that need to be made.

The upgrade guides can be found at pattern:

D:\PROJECTs\dracory.com\blueprint\docs\upgrade_guides\upgrade-[FROM_VERSION]-to-[TO_VERSION].md

Follow the upgrade guide step-by-step to:
- Update dependencies in go.mod
- Apply breaking changes
- Update import paths
- Modify configuration
- Run tests to verify the upgrade

**Important**: Some features or changes mentioned in the upgrade guide may already be implemented in the project. Before applying each change, check if the feature already exists in the project files. If it does, skip that step to avoid duplication or conflicts.

You should make the necessary changes to the project files to upgrade it to the latest version.

After making the necessary changes, you should test the project to ensure that it is working correctly.

The blueprint project files should be used as a reference for the latest version of Blueprint, if in doubt, you should refer to the blueprint project files.

The blueprint project files can be found at:

D:\PROJECTs\dracory.com\blueprint