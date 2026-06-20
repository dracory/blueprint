# Proposal: API Layer with CORS + Token Auth

## Status

Proposed

## Context

Blueprint currently only serves HTML controllers. There is no API layer for external clients, mobile apps, or SPA frontends. Any project needing JSON APIs must build the infrastructure from scratch — CORS handling, token-based authentication, JSON request/response patterns, and API-specific middleware.

CourseThread implemented a clean API layer with:
- CORS middleware via `go-chi/cors`
- `APIAuthMiddleware` — token-based auth using the `Authorization` header with memory cache for performance
- Separate API context keys to distinguish web session auth from API token auth
- JSON response format using `github.com/dracory/api`
- Route-level middleware application (CORS + API auth applied per-route, not globally)

## Goal

Add a reusable API layer scaffold to Blueprint that provides CORS, token-based authentication, and a JSON controller pattern — ready for projects to extend with domain-specific API endpoints.

## Key Design Decisions

1. **Separate controller tree** — API controllers live in `internal/controllers/api/`, completely separate from HTML controllers.
2. **Token-based auth** — API auth uses the `Authorization` header (session token), not cookies. This enables external clients and SPAs.
3. **Separate context keys** — `APIAuthenticatedSessionContextKey` and `APIAuthenticatedUserContextKey` distinguish API auth from web auth, preventing accidental cross-context access.
4. **Memory-cached auth** — Validated tokens are cached in the memory cache for 1 minute to avoid repeated session store lookups.
5. **Per-route middleware** — CORS and API auth are applied per-route, not globally, so they don't affect HTML routes.
6. **JSON response format** — Uses `github.com/dracory/api` for consistent `Success`/`Error` JSON responses.
7. **CORS is configurable** — Allowed origins default to the app URL; configurable via existing config interface.

## Proposed Implementation

### New Files

#### `internal/controllers/api/routes.go`

```go
package api

func Routes(app app.AppInterface) []rtr.RouteInterface {
    apiRoutes := []rtr.RouteInterface{}

    // Example: health check endpoint
    // healthRoute := rtr.NewRoute().
    //     SetName("API > Health").
    //     SetPath("/api/health").
    //     SetJSONHandler(health.NewHealthController(app).Handler)
    // apiRoutes = append(apiRoutes, healthRoute)

    // Apply CORS and API auth middleware to all API routes
    for _, route := range apiRoutes {
        route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
            rtrMiddlewares.CORSMiddleware(cors.Options{
                AllowedOrigins:   []string{app.GetConfig().GetAppUrl()},
                AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
                AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
                AllowCredentials: true,
                MaxAge:           300,
            }),
            middlewares.NewAPIAuthMiddleware(app),
        })
    }

    return apiRoutes
}
```

#### `internal/controllers/api/doc.go`

Package documentation explaining the API layer conventions:
- POST-only pattern with JSON request/response bodies
- All routes protected by API auth middleware
- Authorization header with session token required
- Uses `github.com/dracory/api` for response formatting

#### `internal/middlewares/api_auth_middleware.go`

```go
type apiAuthCacheItem struct {
    session sessionstore.SessionInterface
    user    userstore.UserInterface
}

func NewAPIAuthMiddleware(app app.AppInterface) rtr.MiddlewareInterface
```

The middleware:
1. Extracts token from `Authorization` header
2. Checks memory cache for validated session (1-minute TTL)
3. If not cached, validates session via `SessionFindByKey`
4. Loads user via `UserFindByID`
5. Sets `APIAuthenticatedSessionContextKey` and `APIAuthenticatedUserContextKey` in context
6. Returns JSON error responses (not HTML) on auth failures

#### `internal/middlewares/api_auth_middleware_test.go`

Tests covering:
- Missing authorization header → 401 JSON error
- Invalid/expired token → 401 JSON error
- Valid token → context populated, next handler called
- Cache hit → no session store lookup
- Cache expiry → re-validates against session store
- Session missing user ID → 401 JSON error
- User not found → 401 JSON error

### Modified Files

#### `internal/config/z_config_constants.go`

Add API auth context keys:

```go
type APIAuthenticatedSessionContextKey struct{}
type APIAuthenticatedUserContextKey struct{}
```

#### `internal/routes/router.go`

Add API routes to the route list:

```go
func routes(app app.AppInterface) []rtr.RouteInterface {
    routes := []rtr.RouteInterface{}
    routes = append(routes, admin.Routes(app)...)
    routes = append(routes, api.Routes(app)...)  // NEW
    routes = append(routes, auth.Routes(app)...)
    // ... existing routes ...
    return routes
}
```

#### `go.mod`

Add dependency:

```
github.com/go-chi/cors/v5
```

## API Response Format

All API endpoints use `github.com/dracory/api` for consistent responses:

```go
// Success response
api.Success("OK").ToString()
// {"status":"success","message":"OK"}

// Error response
api.Error("Invalid token").ToString()
// {"status":"error","message":"Invalid token"}

// Data response
api.Success(data).SetData(payload).ToString()
// {"status":"success","data":{...}}
```

## Example API Controller Pattern

```go
package health

type HealthController struct {
    app app.AppInterface
}

func NewHealthController(app app.AppInterface) *HealthController {
    return &HealthController{app: app}
}

func (c *HealthController) Handler(w http.ResponseWriter, r *http.Request) string {
    user := helpers.GetAPIAuthUser(r)
    return api.Success("Healthy").
        SetData(map[string]any{
            "user": user.GetEmail(),
        }).ToString()
}
```

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/controllers/api/routes.go` | Create | API route registration with CORS + auth middleware |
| `internal/controllers/api/doc.go` | Create | Package documentation |
| `internal/middlewares/api_auth_middleware.go` | Create | Token-based API auth middleware with caching |
| `internal/middlewares/api_auth_middleware_test.go` | Create | Comprehensive auth middleware tests |
| `internal/config/z_config_constants.go` | Modify | Add API auth context keys |
| `internal/routes/router.go` | Modify | Register API routes |
| `internal/helpers/auth.go` | Modify | Add `GetAPIAuthUser()` helper |
| `go.mod` | Modify | Add `go-chi/cors` dependency |

## Security Considerations

- **Token in Authorization header** — Not in URL params or cookies, preventing CSRF and URL leakage.
- **CORS restricted to app URL** — Not wildcard (`*`); configurable per environment.
- **Cache TTL is short** — 1-minute cache means revoked sessions are invalidated within 60 seconds.
- **Separate context keys** — API auth context is distinct from web auth context, preventing cross-context access.
- **Rate limiting** — API routes inherit global rate limiting middleware; additional per-route limits can be added.

## Testing

- Unit test `APIAuthMiddleware` — all auth failure paths return JSON, not HTML
- Unit test cache behavior — cache hit/miss/expiry
- Integration test — API route returns JSON with valid token
- Integration test — CORS headers present on API responses
- Integration test — OPTIONS preflight handled correctly

## Verification

- `go test ./internal/middlewares/...` passes
- `go test ./internal/controllers/api/...` passes
- API routes return JSON responses (not HTML)
- CORS headers present on API responses
- Auth failures return JSON error with 401 status
- Existing HTML routes unaffected

## Benefits

- **Enables API-first projects** — Blueprint becomes suitable for SPAs, mobile apps, and external integrations
- **Clean separation** — API controllers and auth are completely separate from HTML controllers
- **Performance** — Memory-cached auth avoids repeated session store lookups
- **Consistent responses** — `dracory/api` package ensures uniform JSON format
- **Proven pattern** — Already tested in CourseThread production
