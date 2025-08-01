package routes

import (
	"project/config"
	"project/internal/middlewares"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// globalMiddlewares returns a list of middlewares to be applied to all routes
func globalMiddlewares() []rtr.MiddlewareInterface {
	globalMiddlewares := []rtr.MiddlewareInterface{
		middlewares.JailBotsMiddleware(middlewares.JailBotsConfig{Exclude: []string{}}),
		rtrMiddleware.CompressMiddleware(5, "text/html", "text/css"),
		rtrMiddleware.GetHead(),
		rtrMiddleware.CleanPathMiddleware(),
		rtrMiddleware.RedirectSlashesMiddleware(),
		//router.NewNakedDomainToWwwMiddleware([]string{"localhost", "127.0.0.1", "http://sinevia.local"}),
		rtrMiddleware.TimeoutMiddleware(30),                 // 30s timeout
		rtrMiddleware.RateLimitByIPMiddleware(20, 1),        // 20 req per second
		rtrMiddleware.RateLimitByIPMiddleware(180, 1*60),    // 180 req per minute
		rtrMiddleware.RateLimitByIPMiddleware(12000, 60*60), // 12000 req hour
	}

	if config.AppEnvironment != config.APP_ENVIRONMENT_TESTING {
		globalMiddlewares = append(globalMiddlewares, rtrMiddleware.LoggerMiddleware())
		globalMiddlewares = append(globalMiddlewares, rtrMiddleware.RecoveryMiddleware())
	}

	globalMiddlewares = append(globalMiddlewares, middlewares.LogRequestMiddleware())
	globalMiddlewares = append(globalMiddlewares, middlewares.ThemeMiddleware())
	globalMiddlewares = append(globalMiddlewares, middlewares.AuthMiddleware())

	return globalMiddlewares
}
