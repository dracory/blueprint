package auth

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"
	"strings"

	"github.com/dracory/req"
)

type loginController struct {
	registry registry.RegistryInterface
}

func NewLoginController(registry registry.RegistryInterface) *loginController {
	return &loginController{registry: registry}
}

func (controller *loginController) Handler(w http.ResponseWriter, r *http.Request) string {
	homeURL := links.Website().Home()
	userURL := links.User().Home()

	if controller.registry.GetUserStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, `user store is required`, homeURL, 5)
	}

	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, `vault store is required`, homeURL, 5)
	}

	backUrl := req.GetStringTrimmedOr(r, "back_url", userURL)

	// Ensure back_url is part of our domain (contains our root URL)
	if !strings.HasPrefix(backUrl, homeURL) {
		backUrl = userURL
	}

	loginUrl := links.Auth().AuthKnightLogin(backUrl)

	http.Redirect(w, r, loginUrl, http.StatusSeeOther)
	return ""
}
