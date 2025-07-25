package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"project/app/layouts"
	"project/app/links"
	"project/config"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/statsstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONTROLLER ===============================================================

type homeController struct{}

// == CONSTRUCTOR ==============================================================

func NewHomeController() *homeController {
	return &homeController{}
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(r, layouts.Options{
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

	return hb.Wrap().
		Child(header).
		Child(sectionTiles).
		Child(sectionDailyVisitors)
}

func (c *homeController) cardDailyVisitors() hb.TagInterface {
	if !config.StatsStoreUsed {
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
					Href(links.NewAdminLinks().Stats(map[string]string{})).
					Text("View all").
					Style("float: right;")))).
		Child(hb.Div().
			Class("card-body").
			Child(hb.Canvas().ID("VisitorsChart").Style("width:100%;height:300px;"))).
		Style("margin-bottom: 30px;")

	dailyReport = dailyReport.Child(script)
	return dailyReport
}

func (*homeController) tiles() []hb.TagInterface {
	cmsTileOld := map[string]string{
		"title": "Website Manager (Old)",
		"icon":  "bi-globe",
		"link":  links.Admin().Cms(map[string]string{}),
	}

	cmsTileNew := map[string]string{
		"title": "Website Manager (New)",
		"icon":  "bi-globe",
		"link":  links.Admin().CmsNew(),
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

	visitStatsTile := map[string]string{
		"title": "Visit Stats",
		"icon":  "bi-graph-up",
		"link":  links.Admin().Stats(map[string]string{}),
	}

	tiles := []map[string]string{}

	if config.CmsUsed {
		tiles = append(tiles, cmsTileOld)
	}

	if config.CmsStoreUsed {
		tiles = append(tiles, cmsTileNew)
	}

	tiles = append(tiles,
		blogTile,
		userTile,
		shopTile,
		fileManagerTile,
		mediaManagerTile,
		queueTile,
		visitStatsTile,
	)

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
							Child(icons.Icon(tile["icon"], 36, 36, "red")).
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
	if config.StatsStore == nil {
		return nil, nil, errors.New("statsstore is nil")
	}

	datesInRange := c.datesInRange(carbon.Now().SubDays(31), carbon.Now())

	dates = []string{}
	visits = []int64{}

	ctx := context.Background()
	for _, date := range datesInRange {
		visitorCount, err := config.StatsStore.VisitorCount(ctx, statsstore.VisitorQueryOptions{
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
