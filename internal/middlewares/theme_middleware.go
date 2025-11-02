package middlewares

import (
	"github.com/dracory/dashboard"
	"github.com/dracory/rtr"
)

func ThemeMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Theme Middleware").
		SetHandler(dashboard.ThemeMiddleware)
}
