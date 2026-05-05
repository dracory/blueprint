Your task is to upgrade a project to the latest version of Blueprint.

The project is located at:

<project_path>

The project is a web application based on Blueprint that manages its components and features.

## Version Comparison

1. **Read the project's current version**:
   - Check the file: `<project_path>/internal/config/version.go`
   - Read the `Version` constant to determine the current version
   - Example: `const Version = "0.22.0"`

2. **Read the latest Blueprint version**:
   - Check the file: `D:\PROJECTs\dracory.com\blueprint\internal/config/version.go`
   - Read the `Version` constant to determine the latest version
   - Example: `const Version = "0.23.0"`

3. **Determine the upgrade path**:
   - Compare the two versions
   - Identify the upgrade guide needed: `upgrade-v{CURRENT}-to-v{LATEST}.md`
   - Example: If project is v0.22.0 and latest is v0.23.0, use `upgrade-v0.22.0-to-v0.23.0.md`

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

You should make the necessary changes to the project files to upgrade it to the latest version.

After making the necessary changes, you should test the project to ensure that it is working correctly.

The blueprint project files should be used as a reference for the latest version of Blueprint, if in doubt, you should refer to the blueprint project files.

The blueprint project files can be found at:

D:\PROJECTs\dracory.com\blueprint