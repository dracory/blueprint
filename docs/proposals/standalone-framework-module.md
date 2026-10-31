# Proposal: Externalizing Blueprint Features into Standalone Dracory Packages

## Status

Revised proposal (supersedes the August 2026 draft). Grounded in the codebase as of v0.55.0.

## Executive Summary

Blueprint is a starter template: every project that clones it inherits **~537 files / ~57k lines of Go** (~32.5k non-test, ~24.3k test). Most of that code is not project-specific — it is framework functionality that happens to live inside the application repository.

The project has already proven the extraction model. Two waves shipped successfully:

1. **Store layer** — all persistence lives in external `dracory/*store` packages (`userstore`, `sessionstore`, `shopstore`, `cmsstore`, …, ~20 packages).
2. **Feature layer (recent)** — `dracory/auth` replaced ~5,000 lines of hand-built auth controllers (v0.55.0), and the admin panels were replaced by `dracory/useradmin`, `dracory/blogadmin`, `dracory/fileadmin`, `dracory/logadmin`, `dracory/shopadmin`, `dracory/statsadmin`, each wrapped by a thin local controller plus the `internal/controllers/admin/adapters` bridge.

**This proposal continues that exact pattern** — vertical feature extraction into `dracory/*` packages with strongly-typed option structs and callback adapters — rather than the earlier monolithic `blueprint-core` idea. The goal: shrink a cloned project to its genuinely app-specific code (config choices, links, layouts, routes glue, migrations, seeders) so the working codebase is small and every feature upgrades via `go get`.

### The extraction contract (already established)

Every extraction follows the same shape, proven by auth and the admin packages:

- The package owns the feature end-to-end: HTML pages, JSON APIs, business logic, tests.
- The package takes a **strongly-typed options struct** — stores, callbacks (`FuncLayout`, `FlashRedirect`, `OnUserUpdated`, …), URL constants — never `interface{}` maps or stringly-typed lookups.
- The app keeps a **thin adapter/controller** that injects its services and mounts the package mux (see `internal/authadapters/`, `internal/controllers/admin/adapters/`).
- App-owned routes that need project rules mount *before* the package catch-all (how `register` and `logout` coexist with `dracory/auth` today).

---

## Part 1 — Vertical feature extractions (proven pattern, biggest wins)

Ordered by size of win and lowest risk. Each section lists the exact files that move.

### 1.1 Registration → `dracory/auth` (finish the auth extraction)

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/controllers/auth/register/` | 496 / 480 |
| `internal/controllers/auth/register_password/` | 348 / — |
| `internal/controllers/auth/shared/session_login.go` | 395 / 109 |
| `internal/controllers/auth/logout/` | 66 / 130 |
| `internal/controllers/auth/routes.go` | 207 / 196 |
| `internal/rules/auth/` (register form validation) | 83 / 133 |
| `internal/authadapters/` (the glue — stays in app) | 910 / 159 |

**Why**: `dracory/auth` already ships registration, password-reset, and impersonation pages; Blueprint still owns ~2,400 lines of register/logout/session-login code and a second validation rule set. Moving registration into the package (as opt-in methods alongside the existing login methods) completes the extraction: the entire `/auth` surface becomes one package mux with project hooks.

**How**: Extend `types.Method`/`Config*` with `RegisterMethod` options (`FuncRegisterFields`, `PostRegisterHooks` for profile completion, `FuncValidationRules` for project rules like `rules/auth`). The app's `routes.go` shrinks to composing methods; `register`/`register_password`/`logout` logic moves into the package behind the same hooks that kept them app-owned (remember-cookie cleanup, project rules). `internal/authadapters` stays — it is the injection layer, not feature code.

**Expected saving**: ~2,400 lines (10 files), leaving only `authadapters` + a short `routes.go`.

### 1.2 Media manager admin → `dracory/mediaadmin`

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/controllers/admin/media/media_manager_controller.go` | 1,026 / 601 |
| `internal/controllers/admin/media/routes.go` | 23 / — |

**Why**: a self-contained file/media manager (upload, rename, delete, directory ops against `cfg.GetMediaRoot()` via `dracory/filesystem` + `dracory/cdn`) — identical in shape to `fileadmin`, which is already external. Largest single file left in `internal/controllers`.

**How**: `mediaadmin.New(Options{MediaRoot, FuncLayout, FlashRedirect, AdminHomeURL, ...})` + `admin.Handle(w, r)`, copying the `users_controller.go` wrapper (66 lines) almost verbatim.

**Expected saving**: ~1,650 lines.

### 1.3 Shared asset serving → `dracory/filesrv` (new)

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/controllers/shared/thumb/` | 722 / 1,541 |
| `internal/controllers/shared/file/` | 145 / 469 |
| `internal/controllers/shared/media/` | 145 / 297 |
| `internal/controllers/shared/cdn/` | 191 / 144 |
| `internal/controllers/shared/flash/` | 134 / 351 |
| `internal/controllers/shared/resource/` | ~68 / ~68 |
| `internal/controllers/shared/theme/` | ~60 / ~30 |
| `internal/controllers/shared/page_not_found/` | 93 / — |
| `internal/controllers/shared/routes.go` | 85 / 93 |

**Why**: thumbnail generation/proxying (imaging + caching + safe-path handling), file and media download, CDN asset proxying, flash-message endpoints, embedded resource serving — pure HTTP infrastructure with zero project specificity. `thumb` alone is the biggest non-test file outside admin plus the biggest test file in the repo (1,541 lines).

**How**: one package, e.g. `filesrv.New(Options{MediaRoot, CacheStore, Logger, FuncLayout, ThemeDir ...})` exposing `ThumbHandler`, `FileHandler`, `MediaHandler`, `CDNHandler`, `FlashHandler`, `NotFoundHandler`. The app's `shared/routes.go` becomes a route table only.

**Expected saving**: ~4,600 lines (17 files) — the single largest extraction available.

### 1.4 User account area → `dracory/useraccount`

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/controllers/user/account/` (profile + `form_profile_update`) | 753 / 1,215 |
| `internal/controllers/user/home/` | 95 / 115 |
| `internal/controllers/user/partials/` | 193 / 24 |
| `internal/controllers/user/routes.go` | 53 / 385 |

**Why**: profile view/update with PII vault integration, user home dashboard, quick-links partials — generic authenticated-user functionality identical across Blueprint projects. The 525-line `form_profile_update` + 1,100-line test is the heaviest form left in the app.

**How**: `useraccount.New(Options{UserStore, VaultStore, FuncLayout, FuncFields, PostUpdateHooks ...})`; vault seal/unseal plugs in through the same `UserPiiSeal`/`UserPiiUnseal` callbacks `useradmin` already uses.

**Expected saving**: ~2,800 lines (11 files).

### 1.5 Website frontends → companion packages for existing stores

| Area | Files | LOC (non-test / test) | Target |
| :--- | :--- | :--- | :--- |
| Blog frontend | `website/blog/` (post 457/393, home 146/255, card_post 81, routes 168, shared sections, MCP) | ~1,300 | `dracory/blog` (frontend companion to `blogstore`/`blogadmin`) |
| Shop cart | `website/shop/cart/` (765) | ~800 | `dracory/shop` or `shopstore`-adjacent frontend |
| Contact form | `website/contact/` (396/216 + 67) | ~700 | `dracory/contact` (+ moves `email_admin_new_contact` task) |
| SEO | `website/seo/` (sitemap 103/77, robots 58/39, security.txt) + `website/pages/indexnow` (169/43) | ~500 | `dracory/seo` |
| CMS pages | `website/cms/` (72/223, widget_controller, routes 146) | ~500 | `dracory/cms` (frontend companion to `cmsstore`) |
| Swagger | `website/swagger/` | ~100 | `dracory/api` or keep |

**Why**: public-facing controllers that render a store's content — blog posts, CMS pages, cart ops with price validation against `shopstore`, contact submission with admin email, sitemap/robots/security.txt/IndexNow. Each is the "read side" of an already-external store/admin pair and varies little between projects.

**How**: same options-struct + `FuncLayout` contract. Routes become `append(routes, blogfrontend.Routes(app)...)`.

**Expected saving**: ~3,900 lines across ~30 files.

### 1.6 CMS block types → `dracory/cmsblocks` (or into `cmsstore`)

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/cmsblocks/blogpost/` (block type + renderer) | 792 / 319 |
| `internal/cmsblocks/blogpostlist/` | 574 / 215 |
| `internal/cmsblocks/search/` | 638 / 288 |
| `internal/cmsblocks/cms_add_block_types.go` | ~60 / — |

**Why**: generic `cmsstore.RegisterCustomBlockType` implementations (blog post, blog post list, search) — they depend only on `cmsstore`/`blogstore`, not on the app.

**How**: move the three block types into a shared package (or upstream into `cmsstore` itself, which already defines the block-type contract); the app's registration call shrinks to a conditional import.

**Expected saving**: ~3,200 lines (12 files).

### 1.7 Widgets & shortcodes → `dracory/widgets`

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/widgets/` (registry, blockeditor/blog/contact/testimonials/auth-visibility/print widgets, shortcodes) | 1,625 / 498 |

**Why**: a generic widget registry + shortcode resolver used by CMS rendering. Widget implementations (blog post list, contact form, testimonials) belong with the feature packages above or in a shared widgets package.

**How**: extract `widgets_registry` + `cms_add_shortcodes` into `dracory/widgets`; feature widgets move with their packages (contact widget → `dracory/contact`, blog widgets → `dracory/blog`).

**Expected saving**: ~2,100 lines (19 files).

### 1.8 Background tasks → the packages that own their domain

| Task | LOC (non-test / test) | Target |
| :--- | :--- | :--- |
| `tasks/blind_index_rebuild/` | 487 / 545 | `dracory/blindindexstore` (it re-encrypts blind indexes — store-internal concern) |
| `tasks/stats/` (visitor enhance) | 222 / 260 | `dracory/statsstore` or `statsadmin` |
| `tasks/cms_transfer_task.go` | 477 / — | `dracory/cmsstore` / `dracory/cms` |
| `tasks/clean_up/` | 98 / 97 | `dracory/taskstore` (queue maintenance) |
| `tasks/email_*` (admin notify, new user, contact, delivery, test) | ~730 / ~640 | `dracory/email`/`dracory/auth` (delivery already bridges auth) |
| `tasks/hello_world/` | ~50 / 102 | drop (example) or keep as template |

**Why**: each task implements a *store's* background concern; only `register_tasks.go` wiring is app-specific.

**Expected saving**: ~2,500–3,000 lines; `register_tasks.go` stays (wiring is a project choice).

### 1.9 Generic middlewares → `dracory/middlewares`

| Middleware | LOC (non-test / test) | Target |
| :--- | :--- | :--- |
| `security_headers_middleware.go` | 156 / 83 | `dracory/middlewares` (new) or `dracory/rtr/middlewares` |
| `httpsredirect/` | 143 / 111 | same |
| `maintenance_middleware.go` | 214 / 342 | same (needs maintenance-state callback) |
| `log_request_mIddleware.go` | 99 / 332 | same |
| `stats_middleware.go` | ~50 / — | `dracory/statsstore` |
| `subscription_middleware.go` | ~60 / test | `dracory/subscriptionstore` |
| `email_allowlist_middleware.go` | ~90 / 327 | `dracory/auth` |
| `api_auth_middleware.go` | 103 / 282 | `dracory/auth` already ships `ApiAuthOrErrorMiddleware` — replace |
| `auth_`, `user_`, `admin_middleware` | ~300 / ~1,100 | thin app wrappers over `dracory/auth` middlewares + role checks — mostly stay (role-table glue) |
| `ai_browser_auto_login` | 102 / 177 | move with `cmd/ai-browser` (see 2.3) |
| `cms_add_middlewares`, `cms_layout_middleware` | ~90 | `dracory/cms` |

**Why**: security headers, HTTPS redirect, maintenance gate, request logging, per-domain middlewares are framework plumbing; `dracory/auth` already exports auth middlewares so several local ones become thin aliases.

**Expected saving**: ~2,500–3,000 lines.

### 1.10 Remaining small packages

| Area | LOC | Target |
| :--- | :--- | :--- |
| `internal/ext/` (user vault tokenize/untokenize) | 332 / 338 | `dracory/vaultstore` or `userstore` — pure vault helpers |
| `internal/testutils/` (Setup + fixtures) | 675 / 449 | extend `dracory/test` (already a dependency) with the Setup/options pattern |
| `internal/emails/` (send/template machinery) | ~400 of 702 | `dracory/email` (templates themselves stay app-owned) |
| `internal/controllers/liveflux/` | 73 / 256 | **delete** — deprecated per AGENTS.md |
| `internal/cli/` + `internal/cmds/` | 390 / 531 | `dracory/cmd`-based dispatcher (job runner) — see 2.2 |
| `internal/schedules/` | 202 / 121 | gocron glue — moves with the tasks they run |
| `pkg/testimonials/` | 174 / 123 | `dracory/testimonials` or drop |

---

## Part 2 — Tooling and the framework core

### 2.1 `cmd/deploy` → standalone `dracory/blueprint-deploy` CLI

`cmd/deploy/` (9 files, ~1,530 lines incl. SSH + Cloud Run deploy logic) is an operator tool, not runtime code. Publish it as `go install github.com/dracory/blueprint-deploy@latest`; remove from the app repo entirely.

### 2.2 `cmd/envenc`, `cmd/loadtest`, `cmd/ai-browser`

- `cmd/envenc` (222 lines) just wraps `dracory/envenc` — replace with `task envenc` in `taskfile.yml` calling `go run github.com/dracory/envenc@latest`. Delete.
- `cmd/loadtest` (464 lines) — a dev tool; move to a standalone repo or `dracory/test`.
- `cmd/ai-browser` (621 lines) — sandboxed auto-login instance for agent testing; move to `dracory/test` alongside `testutils` so every Blueprint project gets it for free.

### 2.3 The framework core: `internal/app` + `internal/config` (+ `cli`/`cmds`/`routes`)

| Files | LOC (non-test / test) |
| :--- | :--- |
| `internal/app/` (AppInterface, implementation, datastores, db open) | 1,370 / 654 |
| `internal/config/` (env loading, 497-line interface aggregate, store builders) | 3,406 / 739 |
| `internal/cli/`, `internal/cmds/` | 390 / 531 |
| `internal/routes/` | 184 / 80 |

**Why this is last**: this is the genuinely hard part — `configImplementation` mixes framework fields (DB, host, auth toggles, store flags) with project fields (Stripe keys, LLM providers, media roots), and `datastores.go` wires ~20 optional stores. Every project edits both.

**How** (staged, avoiding the monolith trap):

1. **`dracory/blueprint-config`**: extract the *mechanism* — env loading, the getter/setter interface-aggregate pattern, `database_config` pool builders — while the project's config keeps only its domain struct embedding `BaseConfig` + the slices it needs (`AuthConfig`, `LLMConfig`, `PaymentConfig`, …). Same composable-interface pattern the old draft proposed, but only for the stable machinery.
2. **`dracory/blueprint` (app core)**: extract `AppInterface` getters/setters + `datastores.go` conditional builders + `database_open`/`dbdrivers` into a core package; the project's `app` becomes a thin struct embedding it and registering only the stores it uses. Strictly typed `RegisterUserStore(...)`-style methods — no store maps.
3. **`dracory/blueprint-cli`**: maintenance-mode dispatcher + `cmds.ExecuteJob` → `dracory/cmd`-based package; `cmd/server/main.go` shrinks to ~40 lines (build app → register migrations/tasks → `engine.Run()`).

**Expected saving**: ~5,000 lines, but schedule it after Part 1 so each step stays independently shippable.

---

## What stays in the app (by design)

- `internal/links/` (511 lines) — project URL constants; they differ per site.
- `internal/layouts/` (1,399) — theme/branding; projects customize.
- `internal/emails/` templates — project copy.
- `internal/controllers/admin/adapters/` + `internal/authadapters/` (~1,400) — the injection glue that makes extraction possible.
- `internal/rules/` non-auth rules, `internal/helpers/`, `internal/routes/` — project wiring.
- `database/migrations|seeders|models` (~2,700) — schema is owned by the app (the boilerplate `store_*_migrate.go` wrappers could later be generated by `dracory/neat` tooling, not extracted).
- `internal/resources/`, `internal/cache/` — embed FS + tiny cache, project-scoped.

---

## Impact summary

| Extraction | Files | LOC removed (non-test + test) | Pattern precedent |
| :--- | :--- | :--- | :--- |
| 1.1 Registration → `dracory/auth` | ~10 | ~2,400 | auth v0.55.0 |
| 1.2 `dracory/mediaadmin` | 3 | ~1,650 | fileadmin/useradmin |
| 1.3 `dracory/filesrv` (shared serving) | ~17 | ~4,600 | — |
| 1.4 `dracory/useraccount` | ~11 | ~2,800 | useradmin |
| 1.5 Website frontends (blog/shop/contact/seo/cms/swagger) | ~30 | ~3,900 | store+admin pairs |
| 1.6 `dracory/cmsblocks` | 12 | ~3,200 | cmsstore block types |
| 1.7 `dracory/widgets` | 19 | ~2,100 | — |
| 1.8 Tasks → owning packages | ~20 | ~2,800 | email_delivery → auth |
| 1.9 Middlewares | ~15 | ~2,800 | dracory/auth middlewares |
| 1.10 Small (ext, testutils, emails mech, liveflux del, cli) | ~25 | ~2,500 | dracory/test exists |
| 2.1–2.2 Tool CLIs (deploy/envenc/loadtest/ai-browser) | ~16 | ~2,800 | — |
| 2.3 Framework core (app/config/cli) | ~30 | ~5,000 | — |
| **Total externalizable** | **~200 files** | **~34,000 LOC** | |

**Result**: a cloned Blueprint project drops from ~537 files / ~57k LOC to roughly **~200 files / ~23k LOC** — almost entirely config, links, layouts, routes wiring, migrations, and adapters — while each feature upgrades independently via `go get`.

## Suggested phasing

| Phase | Contents | Why first/last |
| :--- | :--- | :--- |
| **A** | 1.1 registration→auth, 1.2 mediaadmin, 1.10 liveflux delete | Same pattern as shipped extractions; registration completes the auth story |
| **B** | 1.3 filesrv, 1.4 useraccount, 1.5 frontends | Big, self-contained, few cross-deps |
| **C** | 1.6 cmsblocks, 1.7 widgets, 1.8 tasks, 1.9 middlewares | Touch many packages; easier once frontends exist |
| **D** | 2.x tool CLIs, then 2.3 framework core | The core is the highest-coupling extraction; do it last |

## Advantages over the previous draft

1. **Continues what already shipped** — the admin/auth adapter pattern, not a new `blueprint-core` monolith.
2. **Every item names real files and sizes** — measured against `main`, not estimates.
3. **Vertical slices keep strict typing** — stores and callbacks injected via options structs; no `any` registries.
4. **Each phase ships independently** — no flag day; each extraction is one release + one upgrade-guide entry, exactly like v0.54/v0.55.
