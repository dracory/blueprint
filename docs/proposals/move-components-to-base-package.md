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

### 4. Shared UI & Domain Packages (High Priority)

#### Directories to Move:
- `pkg/blogai`
- `pkg/blogblocks`
- `pkg/blogtheme`
- `pkg/testimonials`
- `pkg/webtheme`

#### Reasoning:
These packages in `pkg/` are already structured as libraries but are currently local to Blueprint. They should be moved to their own repositories if they are intended to be reused.

#### Proposed Location:
`github.com/dracory/*` (Independent repositories)

---

### 5. Standardized Configuration & Context Keys (High Priority)

#### Logic to Move:
- `internal/config/constants.go` - Universal environment variable keys (APP_*, DB_*, MAIL_*, LLM_*, STRIPE_*)
- `internal/config/constants.go` - `AuthenticatedUserContextKey`, `AuthenticatedSessionContextKey`, etc.

#### Reasoning:
These constants and context keys are identical across all Dracory projects. Centralizing them ensures that middlewares and helpers from different packages can interoperate seamlessly.

#### Proposed Location:
`github.com/dracory/base/config` and `github.com/dracory/base/context` (or `github.com/dracory/auth`)

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

### 7. Generic Session & User Helpers (Medium Priority)

#### Files to Move:
- `internal/helpers/extend_session.go`
- `internal/helpers/get_auth_sesson.go`
- `internal/helpers/user_settings.go`
- `internal/ext/user.go` - `DisplayNameFull`, `IsClient`, etc.

#### Reasoning:
Once the context keys are standardized, these helpers become framework-agnostic utilities that work with any `sessionstore` or `userstore` implementation.

#### Proposed Location:
`github.com/dracory/base/session` and `github.com/dracory/base/user`

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
