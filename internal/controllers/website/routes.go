package website

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"

	"project/internal/controllers/shared/page_not_found"
	"project/internal/controllers/website/blog"
	"project/internal/controllers/website/cms"
	"project/internal/controllers/website/contact"
	"project/internal/controllers/website/home"
	"project/internal/controllers/website/seo"
	"project/internal/controllers/website/swagger"
)

func Routes(app types.RegistryInterface) []rtr.RouteInterface {
	if app == nil || app.GetConfig() == nil {
		return []rtr.RouteInterface{}
	}

	homeRoute := rtr.NewRoute().
		SetName("Website > Home Controller").
		SetPath(links.HOME).
		SetHTMLHandler(home.NewHomeController(app).Handler)

	pageNotFoundRoute := rtr.NewRoute().
		SetName("Shared > Page Not Found Controller").
		SetPath(links.CATCHALL).
		SetHTMLHandler(page_not_found.PageNotFoundController().Handler)

	faviconRoute := rtr.NewRoute().
		SetName("Website Favicon").
		SetPath("/favicon.svg").
		SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
			w.Header().Add("Content-Type", "image/svg+xml .svg .svgz")
			return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 32 32"><circle cx="20" cy="8" r="1" fill="currentColor"></circle><circle cx="23" cy="8" r="1" fill="currentColor"></circle><circle cx="26" cy="8" r="1" fill="currentColor"></circle><path d="M28 4H4a2.002 2.002 0 0 0-2 2v20a2.002 2.002 0 0 0 2 2h24a2.002 2.002 0 0 0 2-2V6a2.002 2.002 0 0 0-2-2zm0 2v4H4V6zM4 12h6v14H4zm8 14V12h16v14z" fill="currentColor"></path></svg>`
		})

	// These are custom routes for the website, that cannot be served by the CMS
	websiteRoutes := []rtr.RouteInterface{
		faviconRoute,
	}

	// Comment if you do not use the blog routes
	websiteRoutes = append(websiteRoutes, blog.Routes(app)...)
	websiteRoutes = append(websiteRoutes, contact.Routes(app)...)

	// Comment if you do not use the payment routes
	// websiteRoutes = append(websiteRoutes, paymentRoutes...)
	websiteRoutes = append(websiteRoutes, seo.Routes(app)...)
	websiteRoutes = append(websiteRoutes, swagger.Routes()...)

	isCmsUsed := app.GetConfig().GetCmsStoreUsed() && app.GetCmsStore() != nil

	if isCmsUsed {
		websiteRoutes = append(websiteRoutes, cms.Routes(app)...)
	} else {
		websiteRoutes = append(websiteRoutes, homeRoute)
		websiteRoutes = append(websiteRoutes, pageNotFoundRoute)
	}

	return websiteRoutes
}
