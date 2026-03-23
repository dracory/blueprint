# Proposal: Move Components to Base Package

## Overview

This proposal identifies components from the Blueprint project that should be moved to the `github.com/dracory/base` package or other shared Dracory modules to improve code reusability and reduce duplication across projects.

## Components to Move from Blueprint

### 3. Block Editor Renderer (Low Priority)

#### Files to Move:
- `internal/helpers/blog_post_blocks_to_string.go` - Logic for converting JSON blocks to HTML.

#### Reasoning:
If the same block editor format is used across multiple Dracory projects (e.g., in CMS, Blog, etc.), the renderer should be centralized.

#### Proposed Location:
`github.com/dracory/blockeditor`

---

### 6. Registry Base & Store Manager (Medium Priority)

#### Logic to Move:
- `internal/registry/registry_implementation.go` - `cacheDirectory`, `New`, and `Close` boilerplate.
- `internal/registry/stores_*.go` - The `Initialize`, `Migrate`, and `NewStore` pattern for all database stores.

#### Reasoning:
All projects use the same orchestration logic for caches, loggers, and database connections. A generic Store Manager could handle conditional initialization and auto-migration of all standard stores (Audit, Blog, etc.) using a simple configuration map.

#### Proposed Location:
`github.com/dracory/base/registry` or `github.com/dracory/registry`

---

### 8. Email Styling Constants (Low Priority)

#### Files to Move:
- `internal/emails/consts.go` - Standardized inline CSS for email templates.

#### Reasoning:
Generic styles for headings, paragraphs, and buttons in emails should be shared to maintain a consistent look across all transactional emails.

#### Proposed Location:
`github.com/dracory/base/email`

---

## Components to Keep in Blueprint

### Application-Specific Components:
- **Registry Implementation**: The specific wiring of all stores is unique to the application's feature set.
- **Config Definitions**: The list of environment variables and configuration keys specific to Blueprint's features.
- **Routes & Controllers**: Business logic and UI endpoints.
- **Middlewares**: `AuthMiddleware`, `AdminMiddleware`, `SubscriptionMiddleware` (as they rely on app-specific User/Session models).

## Benefits

1. **Faster Project Bootstrapping**: Common security (jailbots), utility (htmx, vault), and infra (database optimizations) are available out-of-the-box.
2. **Unified Security**: Improvements to the bot-jailing blacklist benefit all projects instantly.
3. **Consistency**: All Dracory projects will handle SQLite and environment encryption in the same optimized way.
