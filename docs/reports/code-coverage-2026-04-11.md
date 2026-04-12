# Test Coverage Report

**Generated:** April 12, 2026, 12:15 UTC+01:00
**Last Updated:** April 12, 2026, 15:08 UTC+01:00
**Project:** Blueprint
**Status:** ⚠️ Some Test Failures - Stage 0 Complete, Stage 1 Complete, Stage 2 Complete, Stage 3 In Progress

## Desired Coverage

The desired coverage at stage 0 for this project is **0%**. [✅ COMPLETE]
The desired coverage at stage 1 for this project is **>5% to ≤10%**. [✅ COMPLETE]
The desired coverage at stage 2 for this project is **>10% to ≤20%**. [✅ COMPLETE]
The desired coverage at stage 3 for this project is **>20% to ≤30%**. [IN PROGRESS]
The desired coverage at stage 4 for this project is **>30% to ≤40%**. [PENDING]
The desired coverage at stage 5 for this project is **>40% to ≤50%**. [PENDING]
The desired coverage at stage 6 for this project is **>50% to ≤60%**. [PENDING]
The desired coverage at stage 7 for this project is **>60% to ≤70%**. [PENDING]
The desired coverage at stage 8 for this project is **>70%**. [PENDING]

## Current Coverage Summary

| Stage | Target | Packages in Range |
|-------|--------|-------------------|
| Stage 0 | 0% | 2 packages ✅ |
| Stage 1 | >0% to ≤10% | 0 packages ✅ |
| Stage 2 | >10% to ≤20% | 0 packages ✅ |
| Stage 3 | >20% to ≤30% | 8 packages 🔄 |
| Stage 4 | >30% to ≤40% | 4 packages ⏳ |
| Stage 5 | >40% to ≤50% | 5 packages ⏳ |
| Stage 6 | >50% to ≤60% | 11 packages ⏳ |
| Stage 7 | >60% to ≤70% | 12 packages ⏳ |
| Stage 8 | >70% | 33 packages ✅ |
| Failed | - | 10 packages ⚠️ |
| No Tests | - | 0 packages ✅ |

## Stage 0 Packages (0% Coverage)

| Package | Reason for 0% Coverage |
|---------|----------------------|
| `cmd/envenc` | Main file is a thin CLI wrapper (119 bytes) calling external `github.com/dracory/envenc` library. Tests exist but test library functions, not local code. |
| `cmd/loadtest` | Main file contains CLI entry point. Tests exist but use mock HTTP servers to test generic HTTP client behavior, not the actual load testing logic in main.go. |

**Note:** Both packages have test files but coverage is reported as 0% because the tests don't exercise the local package code - they test external libraries or generic patterns.

**Last Updated:** April 12, 2026, 15:35 UTC+01:00
