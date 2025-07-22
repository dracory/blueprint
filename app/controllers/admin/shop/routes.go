package admin

import (
	"net/http"
	"project/app/links"

	"github.com/dracory/base/req"
	"github.com/dracory/rtr"

	shopDiscounts "project/app/controllers/admin/shop/discounts"
	shopProducts "project/app/controllers/admin/shop/products"
	"project/app/controllers/admin/shop/shared"
)

func ShopRoutes() []rtr.RouteInterface {
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.Value(r, "controller")

		if controller == shared.CONTROLLER_DISCOUNTS {
			return shopDiscounts.NewDiscountController().AnyIndex(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_CREATE {
			return shopProducts.NewProductCreateController().Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_DELETE {
			return shopProducts.NewProductDeleteController().Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCTS {
			return shopProducts.NewProductManagerController().Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_UPDATE {
			return shopProducts.NewProductUpdateController().Handler(w, r)
		}

		if controller == shared.CONTROLLER_ORDERS {
			return NewOrderManagerController().Handler(w, r)
		}

		return NewHomeController().Handler(w, r)
	}

	shopOrders := rtr.NewRoute().
		SetName("Admin > Shop > Orders").
		SetPath(links.ADMIN_SHOP).
		SetHTMLHandler(handler)

	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(links.ADMIN_USERS + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		shopOrders,
		shopCatchAll,
	}
}
