package website

import (
	"net/http"
	"project/internal/config"
	"project/internal/links"

	"project/internal/controllers/shared"

	"project/internal/controllers/website/blog"
	"project/internal/controllers/website/cms"
	"project/internal/controllers/website/contact"
	"project/internal/controllers/website/home"
	"project/internal/controllers/website/seo"
	"project/internal/controllers/website/swagger"

	"github.com/dracory/rtr"
	// paypalControllers "project/controllers/website/paypal"
)

func Routes() []rtr.RouteInterface {
	homeRoute := rtr.NewRoute().
		SetName("Website > Home Controller").
		SetPath(links.HOME).
		SetHTMLHandler(home.NewHomeController().Handler)

	pageNotFoundRoute := rtr.NewRoute().
		SetName("Shared > Page Not Found Controller").
		SetPath(links.CATCHALL).
		SetHTMLHandler(shared.PageNotFoundController().Handler)

	faviconRoute := rtr.NewRoute().
		SetName("Website Favicon").
		SetPath("/favicon.svg").
		SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
			w.Header().Add("Content-Type", "image/svg+xml .svg .svgz")
			return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 32 32"><circle cx="20" cy="8" r="1" fill="currentColor"></circle><circle cx="23" cy="8" r="1" fill="currentColor"></circle><circle cx="26" cy="8" r="1" fill="currentColor"></circle><path d="M28 4H4a2.002 2.002 0 0 0-2 2v20a2.002 2.002 0 0 0 2 2h24a2.002 2.002 0 0 0 2-2V6a2.002 2.002 0 0 0-2-2zm0 2v4H4V6zM4 12h6v14H4zm8 14V12h16v14z" fill="currentColor"></path></svg>`
		})

	contactRoute := rtr.NewRoute().
		SetName("Website > Contact Controller").
		SetPath(links.CONTACT).
		SetMethod(http.MethodGet).
		SetHTMLHandler(contact.NewContactController().AnyIndex)

	contactSubmitRoute := rtr.NewRoute().
		SetName("Website > Contact Submit Controller").
		SetPath(links.CONTACT).
		SetMethod(http.MethodPost).
		SetHTMLHandler(contact.NewContactController().AnyIndex)

		// Serve Swagger UI using a controller
	swaggerUiRoute := rtr.NewRoute().
		SetName("Swagger UI").
		SetPath("/swagger").
		SetHandler(swagger.SwaggerUIController).
		SetMethod(http.MethodGet)

	// Serve embedded YAML
	swaggerYamlRoute := rtr.NewRoute().
		SetName("Swagger YAML").
		SetPath("/docs/swagger.yaml").
		SetHandler(swagger.SwaggerYAMLController).
		SetMethod(http.MethodGet)

	// paymentSuccess := &router.Route{
	// 	Name:        "Website > Payment Success Controller",
	// 	Path:        links.PAYMENT_SUCCESS,
	// 	HTMLHandler: payment.NewPaymentSuccessController().Handler,
	// }

	// paymentCancel := &router.Route{
	// 	Name:        "Website > Payment Cancel Controller",
	// 	Path:        links.PAYMENT_CANCELED,
	// 	HTMLHandler: payment.NewPaymentCanceledController().Handler,
	// }

	// These are custom routes for the website, that cannot be served by the CMS
	websiteRoutes := []rtr.RouteInterface{
		faviconRoute,
		contactRoute,
		contactSubmitRoute,
		// paymentSuccess,
		// paymentCancel,
	}

	websiteRoutes = append(websiteRoutes, swaggerUiRoute)
	websiteRoutes = append(websiteRoutes, swaggerYamlRoute)

	// Comment if you do not use the blog routes
	websiteRoutes = append(websiteRoutes, blog.Routes()...)

	// Comment if you do not use the payment routes
	// websiteRoutes = append(websiteRoutes, paymentRoutes...)
	websiteRoutes = append(websiteRoutes, seo.Routes()...)

	if config.CmsStoreUsed {
		websiteRoutes = append(websiteRoutes, cms.Routes()...)
	} else {
		websiteRoutes = append(websiteRoutes, homeRoute)
		websiteRoutes = append(websiteRoutes, pageNotFoundRoute)
	}

	return websiteRoutes
}

// func paymentRoutes() []router.RouteInterface {
// 	paymentRoutes := []router.RouteInterface{
// 		&router.Route{
// 			Name:        "Website > Payment Canceled Controller > Handle",
// 			Path:        links.PAYMENT_CANCELED,
// 			HTMLHandler: website.NewPaymentCanceledController().Handle,
// 		},
// 		&router.Route{
// 			Name:        "Website > Payment Success Controller > Handle",
// 			Path:        links.PAYMENT_SUCCESS,
// 			HTMLHandler: website.NewPaymentSuccessController().Handle,
// 		},
// 		&router.Route{
// 			Name:        "Guest > Paypal Success Controller > Index",
// 			Path:        links.PAYPAL_SUCCESS,
// 			HTMLHandler: paypalControllers.NewPaypalSuccessController().AnyIndex,
// 		},
// 		&router.Route{
// 			Name:        "Guest > Paypal Cancel Controller > Index",
// 			Path:        links.PAYPAL_CANCEL,
// 			HTMLHandler: paypalControllers.NewPaypalCancelController().AnyIndex,
// 		},
// 		&router.Route{
// 			Name:        "Guest > Paypal Notify Controller > Index",
// 			Path:        links.PAYPAL_NOTIFY,
// 			HTMLHandler: paypalControllers.NewPaypalNotifyController().AnyIndex,
// 		},
// 	}

// 	return paymentRoutes
// }
