package admin

import (
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/internal/controllers/admin/shop/categories"
	"project/internal/controllers/admin/shop/categories/categorymanager"
	categoryupdate "project/internal/controllers/admin/shop/categories/categoryupdate"
	shopDiscounts "project/internal/controllers/admin/shop/discounts"
	shopProducts "project/internal/controllers/admin/shop/products"
	"project/internal/controllers/admin/shop/products/productupdate"
	"project/internal/controllers/admin/shop/shared"
)

func ShopRoutes(registry registry.RegistryInterface) []rtr.RouteInterface {
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		if controller == shared.CONTROLLER_CATEGORIES {
			return categorymanager.NewCategoryManagerController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_CATEGORY_CREATE {
			return categories.NewCategoryCreateController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_CATEGORY_UPDATE {
			return categoryupdate.NewCategoryUpdateController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_DISCOUNTS {
			return shopDiscounts.NewDiscountController(registry).AnyIndex(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_CREATE {
			return shopProducts.NewProductCreateController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_DELETE {
			return shopProducts.NewProductDeleteController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCTS {
			return shopProducts.NewProductManagerController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_PRODUCT_UPDATE {
			return productupdate.NewProductUpdateController(registry).Handler(w, r)
		}

		if controller == shared.CONTROLLER_ORDERS {
			return NewOrderManagerController().Handler(w, r)
		}

		return NewHomeController(registry).Handler(w, r)
	}

	shopOrders := rtr.NewRoute().
		SetName("Admin > Shop > Orders").
		SetPath(links.ADMIN_SHOP).
		SetHTMLHandler(handler)

	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(links.ADMIN_SHOP + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		shopOrders,
		shopCatchAll,
	}
}
