package routes

import (
	"project/internal/middlewares"
	"project/internal/types"
	"time"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// globalMiddlewares returns a list of middlewares to be applied to all routes
func globalMiddlewares(app types.AppInterface) []rtr.MiddlewareInterface {
	globalMiddlewares := []rtr.MiddlewareInterface{
		// Exclude generic patterns that could match legit routes like /user/news
		middlewares.JailBotsMiddleware(middlewares.JailBotsConfig{
			Exclude:      []string{"/new",},
			ExcludePaths: []string{"/blog*", "/th*", "/liveflux*"},
		}),
		rtrMiddleware.CompressMiddleware(5, "text/html", "text/css"),
		rtrMiddleware.GetHead(),
		rtrMiddleware.CleanPathMiddleware(),
		rtrMiddleware.RedirectSlashesMiddleware(),
		// router.NewNakedDomainToWwwMiddleware([]string{"localhost", "127.0.0.1", "http://sinevia.local"}),
		rtrMiddleware.TimeoutMiddleware(30 * time.Second),   // 30s timeout
		rtrMiddleware.RateLimitByIPMiddleware(20, 1),        // 20 req per second
		rtrMiddleware.RateLimitByIPMiddleware(180, 1*60),    // 180 req per minute
		rtrMiddleware.RateLimitByIPMiddleware(12000, 60*60), // 12000 req hour
	}

	// Conditionally add logger and recovery when not running tests
	if app.GetConfig() != nil {
		if app.GetConfig().IsEnvTesting() {
			globalMiddlewares = append(globalMiddlewares,
				rtrMiddleware.LoggerMiddleware(),
				rtrMiddleware.RecoveryMiddleware(),
			)
		}
	}

	globalMiddlewares = append(globalMiddlewares,
		middlewares.LogRequestMiddleware(app),
		middlewares.ThemeMiddleware(),
		middlewares.AuthMiddleware(app),
	)

	return globalMiddlewares
}
