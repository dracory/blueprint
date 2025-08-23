package auth

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"
	"strings"

	"github.com/dracory/base/req"
)

type loginController struct {
	app types.AppInterface
}

func NewLoginController(app types.AppInterface) *loginController {
	return &loginController{app: app}
}

func (controller *loginController) Handler(w http.ResponseWriter, r *http.Request) string {
	homeURL := links.Website().Home()
	userURL := links.User().Home()

	if controller.app.GetUserStore() == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, `user store is required`, homeURL, 5)
	}

	if controller.app.GetConfig().GetVaultStoreUsed() && controller.app.GetVaultStore() == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, `vault store is required`, homeURL, 5)
	}

	backUrl := req.ValueOr(r, "back_url", userURL)

	// Ensure back_url is part of our domain (contains our root URL)
	if !strings.HasPrefix(backUrl, homeURL) {
		backUrl = userURL
	}

	loginUrl := links.Auth().AuthKnightLogin(backUrl)

	http.Redirect(w, r, loginUrl, http.StatusSeeOther)
	return ""
}
