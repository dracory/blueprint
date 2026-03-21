# Proposal: Move Components to Base Package

## Overview

This proposal identifies components from the Blueprint project that should be moved to the `github.com/dracory/base` package or other shared Dracory modules to improve code reusability and reduce duplication across projects.

## Components to Move from Blueprint

### 1. Vault & Tokenization Helpers (Medium Priority)

#### Files to Move:
- `internal/ext/vault.go` - `VaultTokenUpsert` function
- `internal/helpers/untokenize.go` - Batch untokenization utility

#### Reasoning:
These functions provide a higher-level API over the `vaultstore` package for common operations like "create or update token" and "untokenize a map of values".

#### Proposed Location:
`github.com/dracory/base/vault`

---

### 4. Config Encryption Loader (Medium Priority)

#### Logic to Move:
- `internal/config/load.go` - `initializeEnvEncVariables` logic for hydrating environment variables from `.vault` files using `envenc`.

#### Reasoning:
The logic for locating, reading, and decrypting environment secrets from vault files is generic and could be part of the `envenc` or `base/config` package.

#### Proposed Location:
`github.com/dracory/envenc`

---

### 5. Block Editor Renderer (Low Priority)

#### Files to Move:
- `internal/helpers/blog_post_blocks_to_string.go` - Logic for converting JSON blocks to HTML.

#### Reasoning:
If the same block editor format is used across multiple Dracory projects (e.g., in CMS, Blog, etc.), the renderer should be centralized.

#### Proposed Location:
`github.com/dracory/blockeditor`

---

### 6. Shared UI & Domain Packages (High Priority)

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
