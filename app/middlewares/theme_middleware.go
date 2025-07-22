package middlewares

import (
	"github.com/dracory/rtr"
	"github.com/gouniverse/dashboard"
)

func ThemeMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Theme Middleware").
		SetHandler(dashboard.ThemeMiddleware)
}
