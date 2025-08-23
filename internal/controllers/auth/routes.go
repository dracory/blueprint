package auth

import (
    "project/internal/links"
    "project/internal/types"

    "github.com/dracory/rtr"
)

func Routes(application types.AppInterface) []rtr.RouteInterface {
    return []rtr.RouteInterface{
        rtr.NewRoute().
            SetName("Auth > Auth Controller").
            SetPath(links.AUTH_AUTH).
            SetHTMLHandler(NewAuthenticationController(application).Handler),
        rtr.NewRoute().
            SetName("Auth > Login Controller").
            SetPath(links.AUTH_LOGIN).
            SetHTMLHandler(NewLoginController(application).Handler),
        rtr.NewRoute().
            SetName("Auth > Logout Controller").
            SetPath(links.AUTH_LOGOUT).
            SetHTMLHandler(NewLogoutController(application).AnyIndex),
        rtr.NewRoute().
            SetName("Auth > Register Controller").
            SetPath(links.AUTH_REGISTER).
            SetHTMLHandler(NewRegisterController(application).Handler),
    }
}
