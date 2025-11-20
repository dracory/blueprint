package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONTROLLER ===============================================================

type homeController struct {
	app types.AppInterface
}

// == CONSTRUCTOR ==============================================================

func NewHomeController(app types.AppInterface) *homeController {
	return &homeController{app: app}
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Home",
		Content: controller.view(),
		ScriptURLs: []string{
			cdn.Jquery_3_7_1(),
			`https://cdnjs.cloudflare.com/ajax/libs/Chart.js/1.0.2/Chart.min.js`,
		},
		Styles: []string{},
	}).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) view() *hb.Tag {
	header := hb.Heading1().
		HTML("Admin Home").
		Style("margin-bottom:30px;margin-top:30px;")

	sectionTiles := hb.Section().
		Child(bs.Row().
			Class("g-4").
			Children(c.tiles()))

	sectionDailyVisitors := hb.Section().
		Style("margin-top:30px;margin-bottom:30px;").
		Child(bs.Row().
			Class("g-4").
			Child(bs.Column(12).
				Child(c.cardDailyVisitors())))

	return layouts.AdminPage(
		header,
		sectionTiles,
		sectionDailyVisitors,
	)
}

func (c *homeController) cardDailyVisitors() hb.TagInterface {
	if !c.app.GetConfig().GetStatsStoreUsed() {
		return nil
	}

	dates, visits, err := c.visitorsData()

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	labels := dates
	values := visits

	labelsJSON, err := json.MarshalIndent(labels, "", "  ")

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	valuesJSON, err := json.MarshalIndent(values, "", "  ")

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	script := hb.Script(`
			setTimeout(function () {
				generateVisitorsChart();
			}, 1000);
			function generateVisitorsChart() {
				var visitorData = {
					labels: ` + cast.ToString(labelsJSON) + `,
					datasets:
							[
								{
									fillColor: "rgba(172,194,132,0.4)",
									strokeColor: "#ACC26D",
									pointColor: "#fff",
									pointStrokeColor: "#9DB86D",
									data: ` + cast.ToString(valuesJSON) + `
								}
							]
				};

				var visitorContext = document.getElementById('VisitorsChart').getContext('2d');
				new Chart(visitorContext).Line(visitorData);
			}
		`)

	dailyReport := hb.Div().
		ID("DailyVisitorsReport").
		Class("card").
		Child(hb.Div().Class("card-header").
			Child(hb.H5().
				Text("Daily Visitors Report").
				Child(hb.Hyperlink().
					Class("card-link").
					Href(links.Admin().Stats(map[string]string{})).
					Text("View all").
					Style("float: right;")))).
		Child(hb.Div().
			Class("card-body").
			Child(hb.Canvas().ID("VisitorsChart").Style("width:100%;height:300px;"))).
		Style("margin-bottom: 30px;")

	dailyReport = dailyReport.Child(script)
	return dailyReport
}

func (c *homeController) tiles() []hb.TagInterface {
	// cmsTileOld := map[string]string{
	// 	"title": "Website Manager (Old)",
	// 	"icon":  "bi-globe",
	// 	"link":  links.Admin().CmsOld(),
	// }

	cmsTile := map[string]string{
		"title": "Website Manager",
		"icon":  "bi-globe",
		"link":  links.Admin().Cms(),
	}

	blogTile := map[string]string{
		"title": "Blog Manager",
		"icon":  "bi-newspaper",
		"link":  links.Admin().Blog(map[string]string{}),
	}

	userTile := map[string]string{
		"title": "User Manager",
		"icon":  "bi-people",
		"link":  links.Admin().Users(map[string]string{}),
	}

	shopTile := map[string]string{
		"title": "Shop Manager",
		"icon":  "bi-shop",
		"link":  links.Admin().Shop(map[string]string{}),
	}

	// faqTile := map[string]string{
	// 	"title": "FAQ Manager",
	// 	"icon":  "bi-question-circle",
	// 	"link":  links.NewAdminLinks().Faq(map[string]string{}),
	// }

	fileManagerTile := map[string]string{
		"title": "File Manager (New, DB)",
		"icon":  "bi-box",
		"link":  links.Admin().FileManager(map[string]string{}),
	}

	mediaManagerTile := map[string]string{
		"title": "Media Manager (Old, S3)",
		"icon":  "bi-box",
		"link":  links.Admin().MediaManager(map[string]string{}),
	}

	// cdnManagerTile := map[string]string{
	// 	"title": "CDN Manager",
	// 	"icon":  "bi-folder-symlink",
	// 	"link":  "https://gitlab.com/lesichkov/media",
	// }

	queueTile := map[string]string{
		"title": "Queue Manager",
		"icon":  "bi-heart-pulse",
		"link":  links.Admin().Tasks(map[string]string{}),
	}

	logsTile := map[string]string{
		"title": "Log Manager",
		"icon":  "bi-clipboard-data",
		"link":  links.Admin().Logs(map[string]string{}),
	}

	visitStatsTile := map[string]string{
		"title": "Visit Stats",
		"icon":  "bi-graph-up",
		"link":  links.Admin().Stats(map[string]string{}),
	}

	tiles := []map[string]string{}

	if c.app.GetConfig().GetCmsStoreUsed() {
		tiles = append(tiles, cmsTile)
	}

	if c.app.GetConfig().GetBlogStoreUsed() {
		tiles = append(tiles, blogTile)
	}

	if c.app.GetConfig().GetUserStoreUsed() {
		tiles = append(tiles, userTile)
	}

	if c.app.GetConfig().GetShopStoreUsed() {
		tiles = append(tiles, shopTile)
	}

	if c.app.GetConfig().GetSqlFileStoreUsed() {
		tiles = append(tiles, fileManagerTile)
	}

	if c.app.GetConfig().GetMediaDriver() != "" {
		tiles = append(tiles, mediaManagerTile)
	}

	if c.app.GetConfig().GetTaskStoreUsed() {
		tiles = append(tiles, queueTile)
	}

	if c.app.GetConfig().GetStatsStoreUsed() {
		tiles = append(tiles, visitStatsTile)
	}

	if c.app.GetConfig().GetLogStoreUsed() {
		tiles = append(tiles, logsTile)
	}

	cards := lo.Map(tiles, func(tile map[string]string, index int) hb.TagInterface {
		target := lo.ValueOr(tile, "target", "")
		card := bs.Card().
			Class("bg-transparent border round-10 shadow-lg h-100 pt-4").
			OnMouseOver(`
			this.style.setProperty('background-color', 'beige', 'important');
			this.style.setProperty('scale', 1.1);
			this.style.setProperty('border', '4px solid moccasin', 'important');
			`).
			OnMouseOut(`
			this.style.setProperty('background-color', 'transparent', 'important');
			this.style.setProperty('scale', 1);
			this.style.setProperty('border', '0px solid moccasin', 'important');
			`).
			Style("margin:0px 0px 20px 0px;").
			Children([]hb.TagInterface{
				bs.CardBody().
					Class("d-flex flex-column justify-content-evenly").
					Children([]hb.TagInterface{
						hb.Div().
							Child(hb.I().Class("bi " + tile["icon"]).Style("font-size: 36px; color: red")).
							Style("text-align:center;padding:10px;"),
						hb.Heading5().
							HTML(tile["title"]).
							Style("text-align:center;padding:10px;"),
					}),
			})

		link := hb.Hyperlink().
			Href(tile["link"]).
			AttrIf(target != "", "target", target).
			Child(card)

		column := bs.Column(3).
			Class("col-sm-6 col-md-4 col-lg-3").
			Child(link)

		return column
	})

	return cards
}

func (c *homeController) datesInRange(timeStart, timeEnd *carbon.Carbon) []string {
	rangeDates := []string{}

	if timeStart.Lte(timeEnd) {
		rangeDates = append(rangeDates, timeStart.ToDateString())
		for timeStart.Lt(timeEnd) {
			timeStart = timeStart.AddDays(1) // += 86400 // add 24 hours
			rangeDates = append(rangeDates, timeStart.ToDateString())
		}
	}

	return rangeDates
}

func (c *homeController) visitorsData() (dates []string, visits []int64, err error) {
	if c.app.GetStatsStore() == nil {
		return nil, nil, errors.New("statsstore is nil")
	}

	datesInRange := c.datesInRange(carbon.Now().SubDays(31), carbon.Now())

	dates = []string{}
	visits = []int64{}

	ctx := context.Background()
	for _, date := range datesInRange {
		visitorCount, err := c.app.GetStatsStore().VisitorCount(ctx, statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
			Distinct:     statsstore.COLUMN_IP_ADDRESS,
		})

		if err != nil {
			return nil, nil, err
		}

		dates = append(dates, date)
		visits = append(visits, visitorCount)
	}

	return dates, visits, nil
}
