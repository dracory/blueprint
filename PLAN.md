# Blueprint Improvement Plan

## Overview
This document captures the current set of application and configuration refinements queued for the `internal/app/` and `internal/config/` packages. It focuses on tightening store lifecycle management, strengthening configuration validation, building test coverage, and clarifying defaults so future contributors can ramp quickly.

## Workstreams
- **Configuration hardening**
  - Factor repeated validation patterns in `internal/config/load.go` into small helpers with typed errors.
  - Surface missing-variable context in error messages and ensure SQLite branches stay lean.
- **Documentation & defaults**
  - Move shared defaults (translation list, daily analysis symbols) from `config.Load()` into a dedicated `defaults.go` for discoverability.
  - Expand comments around the `KEY_*_STORE_USED` values in `internal/config/constants.go` to explain operational impact.
- **Testing enhancements**
  - Add table-driven tests for `deriveEnvEncKey()` in `internal/config/envenc.go` to cover happy path and validation errors.
  - Introduce regression tests for the new config helpers once extracted.

## Milestones
1. Complete migration helper renames and add targeted unit tests for store initialization flows.
2. Refactor configuration validation into helpers, update error handling, and introduce associated tests.
3. Extract defaults, update documentation, and ensure CI covers the new test suites.

## Open Questions
- Are there timelines or release windows that dictate sequencing across environments?
- Should configuration errors surface as structured logs or remain plain errors for the calling layer to wrap?
