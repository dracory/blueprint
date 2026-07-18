package cart

import (
	"project/internal/app"

	"github.com/dracory/rtr"
)

// Routes returns the cart routes
func Routes(app app.AppInterface) []rtr.RouteInterface {
	if app == nil || app.GetConfig() == nil {
		return []rtr.RouteInterface{}
	}

	cartController := NewCartController(app)

	cartAPIRoute := rtr.NewRoute().
		SetName("Website > Shop Cart API").
		SetPath("/shop/cart/api").
		SetHTMLHandler(cartController.Handler)

	return []rtr.RouteInterface{
		cartAPIRoute,
	}
}
