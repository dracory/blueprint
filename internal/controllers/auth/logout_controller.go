package auth

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/auth"
)

type logoutController struct {
	registry registry.RegistryInterface
}

func NewLogoutController(registry registry.RegistryInterface) *logoutController {
	return &logoutController{registry: registry}
}

func (controller *logoutController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	auth.AuthCookieRemove(w, r)

	return helpers.ToFlashSuccess(controller.registry.GetCacheStore(), w, r, "You have been logged out successfully", links.Website().Home(), 5)
}
