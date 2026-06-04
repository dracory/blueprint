package seo

import (
	"net/http"
	"project/internal/controllers/website/pages/indexnow"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/rtr"
)

func Routes(app app.AppInterface) []rtr.RouteInterface {
	adsRoute := rtr.NewRoute().
		SetName("Website > ads.txt").
		SetPath("/ads.txt").
		SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
			//return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
			return "google.com, pub-YOURNUMBER, DIRECT, YOURSTRING"
		})

	robotsRoute := rtr.NewRoute().
		SetName("Website > RobotsTxt").
		SetPath(links.ROBOTS_TXT).
		SetHTMLHandler(NewRobotsTxtController().Handler)

	securityRoute := rtr.NewRoute().
		SetName("Website > SecurityTxt").
		SetPath(links.SECURITY_TXT).
		SetHTMLHandler(NewSecurityTxtController().Handler)

	sitemapRoute := rtr.NewRoute().
		SetName("Website > Sitemap").
		SetPath(links.SITEMAP_XML).
		SetHTMLHandler(NewSitemapXmlController(app).Handler)

	indexNowRoute := rtr.NewRoute().
		SetName("Website > IndexNow Controller").
		SetPath(links.INDEXNOW).
		SetHTMLHandler(indexnow.NewIndexNowController(app).Handler)

	indexNowKeyRoute := rtr.NewRoute().
		SetName("Website > IndexNow Key").
		SetPath("/" + app.GetConfig().GetIndexNowKey() + ".txt").
		SetHTMLHandler(NewIndexNowKeyController(app).Handler)

	return []rtr.RouteInterface{
		adsRoute,
		robotsRoute,
		securityRoute,
		sitemapRoute,
		indexNowRoute,
		indexNowKeyRoute,
	}
}
