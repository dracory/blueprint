package admin

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	shopadmin "project/pkg/shopadmin"
)

// shopAdminController wraps the pkg/shopadmin package for integration
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
		Registry:       controller.app,
		AdminHomeURL:   links.Admin().Home(),
		ShopAdminURL:   links.Admin().Shop(),
		FileManagerURL: links.Admin().FileManager(),
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
