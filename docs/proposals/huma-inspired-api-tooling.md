# Proposal: Huma-Inspired API Tooling (Generated OpenAPI, Tag Validation, Problem+JSON)

## Status

Proposed

## Context

Blueprint serves a **hand-written** `swagger.yaml` via
`internal/controllers/website/swagger/` (embedded Swagger UI + static spec).
Because the spec is maintained manually, it can silently drift from the routes
registered in `internal/routes/router.go` — docs can lie.

[Huma](https://github.com/danielgtaylor/huma) is a Go API framework built around
one idea worth stealing: **the code is the spec**. Operations are declared as
typed Go structs; OpenAPI 3.1, request validation, and RFC 9457 error responses
are all generated from those declarations.

We do not propose adopting huma itself — it is a full framework and Blueprint is
committed to `github.com/dracory/rtr`. We propose adopting huma's *ideas* as a
thin layer on top of `rtr`.

## Goal

1. Replace the static `swagger.yaml` with an OpenAPI 3.1 document **generated at
   startup from the registered route list** — docs that cannot lie.
2. Add optional, opt-in request validation driven by struct tags.
3. Standardise API error responses on `application/problem+json` (RFC 9457).
4. Strengthen tests so they exercise routes through the real `rtr.RouterInterface`
   (huma's `humatest` approach), catching registration mistakes.

## Key Design Decisions

1. **No new framework** — everything is a layer over `rtr.RouteInterface`. Routes
   that declare no schema keep working exactly as today; adoption is incremental.
2. **Single generation point** — all routes already funnel through
   `routes()`/`RoutesList()` in `internal/routes/router.go`. Walking that list to
   emit `openapi.json` is one loop, not per-controller work.
3. **Metadata lives on the route** — extend `rtr.NewRoute()` (or wrap it) with
   optional schema setters: `SetRequestSchema(v any)`, `SetResponseSchema(v any)`,
   `SetSummary(string)`, `SetTags(...string)`, `SetOperationID(string)`.
4. **Validation is opt-in per route** — a route that sets a `RequestSchema` gets
   automatic validation of path/query params and JSON body from struct tags
   (`path:`, `query:`, `json:`, `required:`, `minLength:`, etc.) before the
   handler runs. HTML/form routes unaffected.
5. **Errors are RFC 9457 `application/problem+json`** — a small `ErrorModel`
   type (`type`, `title`, `status`, `detail`, `errors[]`) returned for
   validation failures; also usable by the API layer proposal for consistency.
6. **Keep the embedded Swagger UI** — only the YAML source changes: serve the
   generated document instead of `swagger.yaml`.
7. **humatest-style tests** — new tests build `Router(app)` and issue real
   HTTP requests through it, so routing, middleware, validation, and
   serialisation are verified together.

## Proposed Implementation

### Phase 1 — OpenAPI generation (highest value, lowest risk)

- Add optional schema metadata to `rtr.Route` (upstream change in
  `github.com/dracory/rtr`, or a `routespec` wrapper package under `pkg/`).
- New `internal/openapi/` package:
  - `Generate(routes []rtr.RouteInterface) *OpenAPI` — walks routes, reflects on
    schema structs, emits OpenAPI 3.1.
  - Component schemas generated from Go structs via reflection (reuse an
    existing schema lib such as `github.com/danielgtaylor/huma/v2`'s schema
    package as reference, or `invopop/jsonschema` to avoid a dependency on the
    huma framework itself).
- Replace `SwaggerYAMLController` with a handler that serves the generated doc
  (JSON + YAML via `Accept` negotiation); keep `/swagger` UI unchanged.
- Annotate a few existing JSON routes as a pilot.

### Phase 2 — Request validation

- `internal/apivalidation` middleware/handler wrapper for `PostJSON` routes with
  a declared `RequestSchema`.
- Tag grammar modelled on huma: `path`, `query`, `header`, `json`, `required`,
  `minLength`, `maxLength`, `minimum`, `maximum`, `enum`, `format`.
- Validation failures return `422` with `problem+json` body listing field
  errors.

### Phase 3 — Error model + test harness

- `internal/apierrors` (or extend `dracory/api`): `ProblemModel`,
  `NewProblem(status, title, detail)`, field-error list.
- `internal/testutils`: `TestRouter(t, app)` helper returning an
  `httptest`-backed client bound to the real `Router(app)` — mirrors
  `humatest`.

### Non-goals

- Adopting huma as the HTTP framework.
- Content negotiation beyond JSON/YAML for the spec (CBOR etc. — out of scope).
- `autopatch`, `sse`, `humacli` equivalents — note for future reference only.

## Migration / Compatibility

- Fully additive: existing routes compile and behave identically without
  schema metadata.
- `swagger.yaml` is deleted once generation covers the documented endpoints;
  `/swagger` and `/docs/swagger.yaml` URLs are preserved (latter redirects or
  serves generated YAML).
- Downstream Blueprint projects get generation for free on upgrade; they opt in
  to validation per route.

## Open Questions

1. Schema metadata on `rtr.Route` upstream vs. a Blueprint-side wrapper —
   depends on how much we want to push into `dracory/rtr`.
2. Do we want `humatest`-style helpers to live in `internal/testutils` or a
   public `pkg/` for reuse by downstream projects?
