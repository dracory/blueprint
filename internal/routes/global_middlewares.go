package routes

import (
	"project/internal/middlewares"
	"project/internal/registry"
	"time"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// Rate limit constants
const (
	defaultVisitsPerMin  = 10    // Default rate limit
	defaultVisitsPerSec  = 20    // Default rate limit per second
	defaultVisitsPerHour = 12000 // Default rate limit per hour

	devVisitsPerMin  = 10000  // Development rate limit per minute
	devVisitsPerSec  = 1000   // Development rate limit per second
	devVisitsPerHour = 100000 // Development rate limit per hour
)

// getRateLimits returns appropriate rate limits based on environment
func getRateLimits(registry registry.RegistryInterface) (perSec, perMin, perHour int) {
	if registry.GetConfig() != nil {
		isDevelopment := registry.GetConfig().IsEnvDevelopment()
		isLocal := registry.GetConfig().IsEnvLocal()

		if isDevelopment || isLocal {
			return devVisitsPerSec, devVisitsPerMin, devVisitsPerHour
		}
	}
	return defaultVisitsPerSec, defaultVisitsPerMin, defaultVisitsPerHour
}

// globalMiddlewares returns a list of middlewares to be applied to all routes
func globalMiddlewares(registry registry.RegistryInterface) []rtr.MiddlewareInterface {
	// Get rate limits based on environment
	perSec, perMin, perHour := getRateLimits(registry)

	globalMiddlewares := []rtr.MiddlewareInterface{
		// Exclude generic patterns that could match legit routes like /user/news
		rtrMiddleware.JailBotsMiddleware(rtrMiddleware.JailBotsConfig{
			Exclude: []string{"/new"},
			ExcludePaths: []string{
				"/blog*",
				"/th*",
				"/liveflux*",
				"/admin/cms*",
				"/admin/*cms*",
				"/assets*",
				"*/assets/*",
				"/files/*",
			},
		}),
		rtrMiddleware.CompressMiddleware(5, "text/html", "text/css"),
		rtrMiddleware.GetHead(),
		rtrMiddleware.CleanPathMiddleware(),
		rtrMiddleware.RedirectSlashesMiddleware(),
		// router.NewNakedDomainToWwwMiddleware([]string{"localhost", "127.0.0.1", "http://sinevia.local"}),
		rtrMiddleware.TimeoutMiddleware(30 * time.Second),     // 30s timeout
		rtrMiddleware.RateLimitByIPMiddleware(perSec, 1),      // per second
		rtrMiddleware.RateLimitByIPMiddleware(perMin, 1*60),   // per minute
		rtrMiddleware.RateLimitByIPMiddleware(perHour, 60*60), // per hour
	}

	// Conditionally add logger and recovery when not running tests
	if registry.GetConfig() != nil {
		isNotTesting := !registry.GetConfig().IsEnvTesting()
		if isNotTesting {
			globalMiddlewares = append(globalMiddlewares,
				rtrMiddleware.LoggerMiddleware(),
				rtrMiddleware.RecoveryMiddleware(),
			)
		}
	}

	// Add HTTPS redirect middleware only in production (not in development or testing)
	if registry.GetConfig() != nil &&
		!registry.GetConfig().IsEnvTesting() &&
		!registry.GetConfig().IsEnvDevelopment() {
		globalMiddlewares = append(globalMiddlewares,
			middlewares.NewHTTPSRedirectMiddleware(),
		)
	}

	globalMiddlewares = append(globalMiddlewares,
		middlewares.LogRequestMiddleware(registry),
		middlewares.NewSecurityHeadersMiddleware(),
		middlewares.ThemeMiddleware(),
		middlewares.AuthMiddleware(registry),
	)

	return globalMiddlewares
}
