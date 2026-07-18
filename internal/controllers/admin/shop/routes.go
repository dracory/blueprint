package admin

import (
	"errors"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/rtr"
)

func ShopRoutes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	shop := rtr.NewRoute().
		SetName("Admin > Shop").
		SetPath(links.ADMIN_SHOP).
		SetHandler(NewShopAdminController(app).Handler)

	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(links.ADMIN_SHOP + links.CATCHALL).
		SetHandler(NewShopAdminController(app).Handler)

	return []rtr.RouteInterface{
		shop,
		shopCatchAll,
	}, nil
}
