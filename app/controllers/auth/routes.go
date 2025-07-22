package auth

import (
	"project/app/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Auth > Auth Controller").
			SetPath(links.AUTH_AUTH).
			SetHTMLHandler(NewAuthenticationController().Handler),
		rtr.NewRoute().
			SetName("Auth > Login Controller").
			SetPath(links.AUTH_LOGIN).
			SetHTMLHandler(NewLoginController().Handler),
		rtr.NewRoute().
			SetName("Auth > Logout Controller").
			SetPath(links.AUTH_LOGOUT).
			SetHTMLHandler(NewLogoutController().AnyIndex),
		rtr.NewRoute().
			SetName("Auth > Register Controller").
			SetPath(links.AUTH_REGISTER).
			SetHTMLHandler(NewRegisterController().Handler),
	}
}
