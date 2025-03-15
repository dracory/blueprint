package routes

import (
	"net/http"

	"project/app/config"
	"project/app/controllers/auth"
	"project/app/controllers/shared/notfound"
	"project/app/controllers/website/about"
	"project/app/controllers/website/home"
	"project/app/middleware"
	"project/internal/platform/database"

	"github.com/dracory/base/router"
)

// SetupRoutes configures all routes and middleware for the application
func SetupRoutes(cfg *config.Config, db *database.Database) router.RouterInterface {
	// Create controllers
	homeController := home.NewController(cfg, db)
	aboutController := about.NewController(cfg, db)
	notFoundController := notfound.NewController(cfg, db)

	// Create router
	r := router.NewRouter()

	// Apply global middleware
	recoverMiddleware := func(next http.Handler) http.Handler {
		return middleware.Recover()(next)
	}
	loggerMiddleware := func(next http.Handler) http.Handler {
		return middleware.Logger()(next)
	}
	corsMiddleware := func(next http.Handler) http.Handler {
		return middleware.CORS(cfg)(next)
	}

	r.AddBeforeMiddlewares([]router.Middleware{recoverMiddleware, loggerMiddleware, corsMiddleware})

	// Register routes
	indexRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/").
		SetHandler(homeController.Index)

	aboutRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/about").
		SetHandler(aboutController.Index)

	r.AddRoute(indexRoute)
	r.AddRoute(aboutRoute)

	r.AddGroup(auth.Routes(cfg, db))

	// Serve static files
	staticFileHandler := http.FileServer(http.Dir("./web/static"))
	staticRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetHandler(func(w http.ResponseWriter, req *http.Request) {
			// Strip the prefix
			http.StripPrefix("/static/", staticFileHandler).ServeHTTP(w, req)
		})
	r.AddRoute(staticRoute)

	// Register 404 handler
	notFoundRoute := router.NewRoute().
		SetMethod("").
		SetPath("/*").
		SetHandler(notFoundController.Index)
	r.AddRoute(notFoundRoute)

	return r
}
