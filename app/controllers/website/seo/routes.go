package seo

import (
	"net/http"

	"github.com/gouniverse/responses"
	"github.com/gouniverse/router"
)

func Routes() []router.RouteInterface {
	adsRoute := &router.Route{
		Name: "Website > ads.txt",
		Path: "/ads.txt",
		HTMLHandler: responses.HTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			//return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
			return "google.com, pub-YOURNUMBER, DIRECT, YOURSTRING"
		}),
	}

	robotsRoute := &router.Route{
		Name:        "Website > RobotsTxt",
		Path:        "/robots.txt",
		HTMLHandler: NewRobotsTxtController().Handler,
	}

	securityRoute := &router.Route{
		Name:        "Website > SecurityTxt",
		Path:        "/security.txt",
		HTMLHandler: NewSecurityTxtController().Handler,
	}

	sitemapRoute := &router.Route{
		Name:        "Website > Sitemap",
		Path:        "/sitemap.xml",
		HTMLHandler: NewSitemapXmlController().Handler,
	}

	return []router.RouteInterface{
		adsRoute,
		robotsRoute,
		securityRoute,
		sitemapRoute,
	}
}
