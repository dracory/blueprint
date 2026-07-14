package config

// Version is the current version of the Blueprint rapid application development (RAD) starter template.
// This constant is used by AI tools and upgrade scripts to identify the template version.
// When upgrading, AI can compare this constant with upgrade guides to determine the appropriate migration path.
//
// Version Workflow:
// 1. Development happens on release branches (e.g., release/v0.23.0)
// 2. Version constant is updated when creating a new release branch
// 3. When ready, generate upgrade guide, merge to main, and tag the release
// 4. Immediately create the next release branch and increment the version
const Version = "0.34.0"

// GetVersion returns the current version of the Blueprint rapid application development (RAD) starter template.
func GetVersion() string {
	return Version
}
