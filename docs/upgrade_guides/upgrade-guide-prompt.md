# Upgrade Guide Generation Prompt

## Task
Generate a comprehensive upgrade guide for upgrading a Go web application from version [FROM_VERSION] to [TO_VERSION].

## Context
- Project: Blueprint Go web application starter
- Repository: https://github.com/dracory/blueprint
- Target audience: LLMs and developers performing upgrades
- Goal: Provide clear, actionable steps for migrating between versions

## Required Analysis Steps

### 1. Version Comparison
- Check out the [FROM_VERSION] tag: `git checkout [FROM_VERSION]`
- Compare with main branch: `git log --oneline [FROM_VERSION]..HEAD`
- Identify breaking changes, new features, and removals

### 2. Key Areas to Examine
- **Entry point changes** (main.go location, cmd/ structure)
- **Package reorganization** (internal/ directory structure)
- **Interface/struct renames** (types, registry, config)
- **API signature changes** (method parameters, return types)
- **Dependency updates** (go.mod changes, removed packages)
- **Configuration changes** (environment variables, config structure)
- **Architecture refactors** (global singletons → registry pattern)
- **Store/task API changes** (enqueue methods, initialization)

### 3. Breaking Change Categories
- **File location changes** (moved/renamed files)
- **Import path changes** (package reorganization)
- **API signature changes** (method parameters, order)
- **Configuration changes** (env vars, config methods)
- **Dependency removals** (removed packages)
- **Architecture changes** (patterns, globals)

## Output Structure

### Header
```markdown
# Upgrade Guide: v[FROM_VERSION] to v[TO_VERSION]

This guide helps LLMs and developers upgrade Blueprint applications from v[FROM_VERSION] to v[TO_VERSION].
```

### Breaking Changes Section
For each breaking change:
```markdown
### [Number]. [Change Title]
**Change**: [Brief description of what changed]

**Old Usage**:
```go
[Code example of old way]
```

**New Usage**:
```go
[Code example of new way]
```

**Action Required**:
- [Specific action items]
- [Files/locations to update]
```

### Migration Steps
```markdown
## 🔄 Migration Steps

### Step [Number]: [Step Title]
[Detailed instructions with commands]
```

### Testing Section
```markdown
## 🧪 Testing After Migration

1. **Unit Tests**: [Testing instructions]
2. **Integration Tests**: [Testing instructions]
3. [Additional test categories]
```

### Additional Sections
- 📝 Additional Notes (new features, removed features)
- 🆘 Common Issues and Solutions
- 📞 Support information

## Content Guidelines

### Code Examples
- Use realistic, copy-pasteable code snippets
- Include import statements where relevant
- Show before/after comparisons clearly

### Action Items
- Be specific about what needs to be changed
- Include file paths when possible
- Provide search/replace commands when applicable

### Commands
- Include bash commands for automated updates
- Use `find` and `sed` for bulk changes
- Provide manual review notes when needed

### Emojis
- ⚠️ for breaking changes
- 🔄 for migration steps  
- 🧪 for testing
- 📝 for notes
- 🆘 for troubleshooting

## Quality Checklist
- [ ] All breaking changes identified and documented
- [ ] Code examples are accurate and tested
- [ ] Migration steps are in logical order
- [ ] Action items are specific and actionable
- [ ] Testing procedures are comprehensive
- [ ] Common issues are addressed
- [ ] Format follows markdown best practices

## Usage Instructions
1. Replace `[FROM_VERSION]` and `[TO_VERSION]` with actual version numbers
2. Execute the analysis steps to gather information
3. Follow the output structure to generate the guide
4. Review against quality checklist before finalizing
