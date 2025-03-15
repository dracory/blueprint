package router

// import (
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/gouniverse/router"
// )

// // Router represents the application router
// type Router struct {
// 	routes      []router.RouteInterface
// 	middlewares []router.Middleware
// 	chiRouter   *chi.Mux
// }

// // New creates a new Router instance
// func New() *Router {
// 	return &Router{
// 		routes:      []router.RouteInterface{},
// 		middlewares: []router.Middleware{},
// 	}
// }

// // Get registers a new GET route
// func (r *Router) Get(path string, handler http.HandlerFunc) {
// 	route := &router.Route{
// 		Path:    path,
// 		Methods: []string{http.MethodGet},
// 		Handler: handler,
// 	}
// 	r.routes = append(r.routes, route)
// }

// // Post registers a new POST route
// func (r *Router) Post(path string, handler http.HandlerFunc) {
// 	route := &router.Route{
// 		Path:    path,
// 		Methods: []string{http.MethodPost},
// 		Handler: handler,
// 	}
// 	r.routes = append(r.routes, route)
// }

// // Put registers a new PUT route
// func (r *Router) Put(path string, handler http.HandlerFunc) {
// 	route := &router.Route{
// 		Path:    path,
// 		Methods: []string{http.MethodPut},
// 		Handler: handler,
// 	}
// 	r.routes = append(r.routes, route)
// }

// // Delete registers a new DELETE route
// func (r *Router) Delete(path string, handler http.HandlerFunc) {
// 	route := &router.Route{
// 		Path:    path,
// 		Methods: []string{http.MethodDelete},
// 		Handler: handler,
// 	}
// 	r.routes = append(r.routes, route)
// }

// // Group creates a new route group with the given prefix
// func (r *Router) Group(prefix string) *RouteGroup {
// 	return &RouteGroup{
// 		router: r,
// 		prefix: prefix,
// 	}
// }

// // Use adds middleware to the router
// func (r *Router) Use(middleware func(http.Handler) http.Handler) {
// 	r.middlewares = append(r.middlewares, router.Middleware{
// 		Name:    "Custom Middleware",
// 		Handler: middleware,
// 	})
// }

// // ServeHTTP implements the http.Handler interface
// func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	// Lazy initialization of the Chi router if it hasn't been created yet
// 	if r.chiRouter == nil {
// 		r.chiRouter = router.NewChiRouter(r.middlewares, r.routes)
// 	}
// 	r.chiRouter.ServeHTTP(w, req)
// }

// // Static serves static files from the given directory
// func (r *Router) Static(path, dir string) {
// 	// Add a static file server route
// 	fileServer := http.FileServer(http.Dir(dir))
// 	r.Get(path+"/*", http.StripPrefix(path, fileServer).ServeHTTP)
// }

// // RouteGroup represents a group of routes with a common prefix
// type RouteGroup struct {
// 	router *Router
// 	prefix string
// }

// // Get registers a new GET route in the group
// func (g *RouteGroup) Get(path string, handler http.HandlerFunc) {
// 	g.router.Get(g.prefix+path, handler)
// }

// // Post registers a new POST route in the group
// func (g *RouteGroup) Post(path string, handler http.HandlerFunc) {
// 	g.router.Post(g.prefix+path, handler)
// }

// // Put registers a new PUT route in the group
// func (g *RouteGroup) Put(path string, handler http.HandlerFunc) {
// 	g.router.Put(g.prefix+path, handler)
// }

// // Delete registers a new DELETE route in the group
// func (g *RouteGroup) Delete(path string, handler http.HandlerFunc) {
// 	g.router.Delete(g.prefix+path, handler)
// }

// // Use adds middleware to the route group
// func (g *RouteGroup) Use(middleware func(http.Handler) http.Handler) {
// 	// This is a simplified implementation that applies middleware to all routes
// 	// In a real implementation, we would need to track which middlewares apply to which route groups
// 	g.router.Use(middleware)
// }
