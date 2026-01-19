package indexnow

import (
	"net/http"
	// "project/internal/controllers/website/shared"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/hb"
)

type indexNowController struct {
	app registry.RegistryInterface
}

func NewIndexNowController(app registry.RegistryInterface) *indexNowController {
	return &indexNowController{
		app: app,
	}
}

func (controller *indexNowController) Handler(w http.ResponseWriter, r *http.Request) string {
	options := layouts.Options{
		WebsiteSection: "IndexNow",
		Title:          "IndexNow",
		Content:        hb.Wrap().HTML(controller.page()),
		ScriptURLs:     []string{},
	}

	if controller.app.GetConfig().GetCmsStoreUsed() {
		return layouts.NewCmsLayout(
			controller.app,
			r,
			options).ToHTML()
	} else {
		return layouts.NewBlankLayout(
			controller.app,
			r,
			options).ToHTML()
	}
}

func (controller *indexNowController) page() string {
	return hb.Wrap().Children([]hb.TagInterface{
		// shared.PageHeader(
		// 	"IndexNow",
		// 	[]layouts.Breadcrumb{
		// 		{
		// 			Name: "IndexNow",
		// 			URL:  links.Website().IndexNow(),
		// 		},
		// 	},
		// ),
		controller.sectionIntro(),
		controller.sectionKey(),
		controller.sectionUsage(),
	}).ToHTML()
}

func (controller *indexNowController) sectionIntro() hb.TagInterface {
	title := hb.H1().
		Class("h2 fw-bold mb-3 text-center").
		Text("Faster indexing with IndexNow")

	lead := hb.P().
		Class("lead mb-3").
		HTML("We ping search engines via IndexNow whenever key pages change, so updates are discovered quickly.")

	body := hb.P().
		Class("mb-0").
		HTML("If you manage uptime or SEO, keep this URL handy and ensure it returns HTTP 200.")

	container := hb.Div().
		Class("container").
		Child(title).
		// Child(shared.HeadingSeparator()).
		Child(lead).
		Child(body)

	section := hb.Section().
		Style("background:#ffffff;").
		Style("padding:40px 0 20px 0;").
		Style("text-align:left;").
		Child(container)

	return section
}

func (controller *indexNowController) sectionKey() hb.TagInterface {
	cardTitle := hb.H5().
		Class("mb-2").
		Text("Verification key")

	keyCode := hb.Code().
		Class("d-block p-3 bg-light rounded text-break").
		Text("YOUR_INDEXNOW_KEY")

	keyLink := hb.A().
		Href("/YOUR_INDEXNOW_KEY.txt").
		Text("View key file")

	cardBody := hb.Div().
		Class("card-body").
		Child(cardTitle).
		Child(keyCode).
		Child(hb.P().Class("mb-0 mt-2").Child(keyLink))

	card := hb.Div().
		Class("card shadow-sm border-0").
		Child(cardBody)

	section := hb.Section().
		Style("background:#f8f9fb;").
		Style("padding:20px 0 20px 0;").
		Child(
			hb.Div().
				Class("container").
				Child(
					hb.Div().
						Class("row justify-content-center").
						Child(
							hb.Div().
								Class("col-md-8").
								Child(card),
						),
				),
		)

	return section
}

func (controller *indexNowController) sectionUsage() hb.TagInterface {
	list := hb.Ul().
		Class("mb-0").
		Children([]hb.TagInterface{
			hb.Li().Text("We submit page URLs to IndexNow on updates to accelerate crawling."),
			hb.Li().Text("Keep the key file accessible (HTTP 200) for verification."),
			hb.Li().Text("If you rotate the key, update both this page and the key file."),
		})

	cardBody := hb.Div().
		Class("card-body").
		Child(hb.H5().Class("mb-2").Text("How we use IndexNow")).
		Child(list)

	card := hb.Div().
		Class("card shadow-sm border-0").
		Child(cardBody)

	section := hb.Section().
		Style("background:#ffffff;").
		Style("padding:20px 0 40px 0;").
		Child(
			hb.Div().
				Class("container").
				Child(
					hb.Div().
						Class("row justify-content-center").
						Child(
							hb.Div().
								Class("col-md-8").
								Child(card),
						),
				),
		)

	return section
}
