package admin

import (
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/cdn"
)

// == CONTROLLER ==============================================================

type userManagerController struct{ registry registry.RegistryInterface }

// == CONSTRUCTOR =============================================================

func NewUserManagerController(registry registry.RegistryInterface) *userManagerController {
	return &userManagerController{registry: registry}
}

func (controller *userManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, errorMessage, links.Admin().Home(), 10)
	}

	if data.action == ActionModalUserFilterShow {
		return controller.onModalUserFilterShow(data).ToHTML()
	}

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Users | User Manager",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}
