package auth

import (
	"net/http"
	"project/app/links"
	"project/config"
	"project/internal/helpers"
	"strings"

	"github.com/dracory/base/req"
)

type loginController struct{}

func NewLoginController() *loginController {
	return &loginController{}
}

func (controller *loginController) Handler(w http.ResponseWriter, r *http.Request) string {
	homeURL := links.Website().Home()
	userURL := links.User().Home()

	if !config.UserStoreUsed || config.UserStore == nil {
		return helpers.ToFlashError(w, r, `user store is required`, homeURL, 5)
	}

	if config.VaultStoreUsed && config.VaultStore == nil {
		return helpers.ToFlashError(w, r, `vault store is required`, homeURL, 5)
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
