package shopadmin

import (
	"errors"
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/pkg/shopadmin/category_manager"
	"project/pkg/shopadmin/discount_manager"
	"project/pkg/shopadmin/home"
	"project/pkg/shopadmin/order_manager"
	"project/pkg/shopadmin/product_manager"
	"project/pkg/shopadmin/products"
	"project/pkg/shopadmin/shared"
)

func Routes(registry registry.RegistryInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error) {
	var options AdminOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	if registry == nil {
		return nil, errors.New("registry cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return home.NewHomeController(registry, options.FileManagerURL).Handler(w, r)
		case shared.CONTROLLER_PRODUCTS:
			return product_manager.NewProductManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_PRODUCT_CREATE:
			return products.NewProductCreateController(registry).Handler(w, r)
		case shared.CONTROLLER_PRODUCT_UPDATE:
			return products.NewProductUpdateController(registry, options.FileManagerURL).Handler(w, r)
		case shared.CONTROLLER_PRODUCT_DELETE:
			return products.NewProductDeleteController(registry).Handler(w, r)
		case shared.CONTROLLER_CATEGORIES:
			return category_manager.NewCategoryManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_CREATE:
			return category_manager.NewCategoryManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_UPDATE:
			return category_manager.NewCategoryManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_DISCOUNTS:
			return discount_manager.NewDiscountManagerController(registry).Handler(w, r)
		case shared.CONTROLLER_ORDERS:
			return order_manager.NewOrderManagerController(registry).Handler(w, r)
		}

		// Default to home
		return home.NewHomeController(registry, options.FileManagerURL).Handler(w, r)
	}

	shop := rtr.NewRoute().
		SetName("Admin > Shop").
		SetPath(links.ADMIN_SHOP).
		SetHTMLHandler(handler)

	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(links.ADMIN_SHOP + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		shop,
		shopCatchAll,
	}, nil
}
