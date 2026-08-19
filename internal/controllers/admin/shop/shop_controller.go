package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/controllers/admin/adapters"
	"project/internal/helpers"
	"project/internal/links"

	shopadmin "github.com/dracory/shopadmin"
)

// shopAdminController wraps the external github.com/dracory/shopadmin package
// for integration with the blueprint admin interface.
type shopAdminController struct {
	app app.AppInterface
}

// NewShopAdminController creates a new shop admin controller
func NewShopAdminController(app app.AppInterface) *shopAdminController {
	return &shopAdminController{app: app}
}

// Handler processes shop admin requests
func (controller *shopAdminController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := shopadmin.New(shopadmin.AdminOptions{
		Store:            controller.app.GetShopStore(),
		Logger:           controller.app.GetLogger(),
		CustomerResolver: adapters.NewUserStoreCustomerResolver(controller.app.GetUserStore()),
		FuncLayout:       adapters.NewLayoutFunc(controller.app),
		AdminHomeURL:     links.Admin().Home(),
		ShopAdminURL:     links.Admin().Shop(),
		FileManagerURL:   links.Admin().FileManager(),
		AuthUserID: func(r *http.Request) string {
			user := helpers.GetAuthUser(r)
			if user == nil {
				return ""
			}
			return user.GetID()
		},
	})

	if err != nil {
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At admin > shopAdminController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.app.GetLogger(); logger != nil {
				logger.Error("At admin > shopAdminController > Handler", "write_error", writeErr.Error())
			}
		}
		return
	}

	admin.Handle(w, r)
}
