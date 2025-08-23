package seo

import (
	"net/http"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	adsRoute := rtr.NewRoute().
		SetName("Website > ads.txt").
		SetPath("/ads.txt").
		SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
			//return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
			return "google.com, pub-YOURNUMBER, DIRECT, YOURSTRING"
		})

	robotsRoute := rtr.NewRoute().
		SetName("Website > RobotsTxt").
		SetPath("/robots.txt").
		SetHTMLHandler(NewRobotsTxtController().Handler)

	securityRoute := rtr.NewRoute().
		SetName("Website > SecurityTxt").
		SetPath("/security.txt").
		SetHTMLHandler(NewSecurityTxtController().Handler)

	sitemapRoute := rtr.NewRoute().
		SetName("Website > Sitemap").
		SetPath("/sitemap.xml").
		SetHTMLHandler(NewSitemapXmlController(app.GetBlogStore()).Handler)

	return []rtr.RouteInterface{
		adsRoute,
		robotsRoute,
		securityRoute,
		sitemapRoute,
	}
}
