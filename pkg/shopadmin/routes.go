package shopadmin

import (
	"errors"
	"net/http"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/pkg/shopadmin/category_create"
	"project/pkg/shopadmin/category_manager"
	"project/pkg/shopadmin/category_update"
	"project/pkg/shopadmin/discount_manager"
	"project/pkg/shopadmin/home"
	"project/pkg/shopadmin/order_details"
	"project/pkg/shopadmin/order_manager"
	"project/pkg/shopadmin/product_delete"
	"project/pkg/shopadmin/product_manager"
	"project/pkg/shopadmin/product_update"
	"project/pkg/shopadmin/shared"
)

func Routes(app app.AppInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error) {
	var options AdminOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return home.NewHomeController(app, options.FileManagerURL).Handler(w, r)
		case shared.CONTROLLER_PRODUCTS:
			return product_manager.NewProductManagerController(app).Handler(w, r)
		case shared.CONTROLLER_PRODUCT_UPDATE:
			return product_update.NewProductUpdateController(app, options.FileManagerURL).Handler(w, r)
		case shared.CONTROLLER_PRODUCT_DELETE:
			return product_delete.NewProductDeleteController(app).Handler(w, r)
		case shared.CONTROLLER_CATEGORIES:
			return category_manager.NewCategoryManagerController(app).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_CREATE:
			return category_create.NewCategoryCreateController(app).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_UPDATE:
			return category_update.NewCategoryUpdateController(app).Handler(w, r)
		case shared.CONTROLLER_DISCOUNTS:
			return discount_manager.NewDiscountManagerController(app).Handler(w, r)
		case shared.CONTROLLER_ORDERS:
			return order_manager.NewOrderManagerController(app).Handler(w, r)
		case shared.CONTROLLER_ORDER_DETAILS:
			return order_details.NewOrderDetailsController(app).Handler(w, r)
		}

		// Default to home
		return home.NewHomeController(app, options.FileManagerURL).Handler(w, r)
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
