package home

import (
	"net/http"

	"project/app/config"
	"project/app/controllers/base"
	"project/internal/platform/database"

	"github.com/gouniverse/hb"
)

// Controller handles requests for the home page
type Controller struct {
	base.Controller
}

// NewController creates a new Controller instance
func NewController(cfg *config.Config, db *database.Database) *Controller {
	return &Controller{
		Controller: *base.NewController(cfg, db),
	}
}

// Index handles the home page request
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	// Create a view using the HB library
	content := hb.NewDiv().Class("row").AddChildren([]hb.TagInterface{
		hb.NewDiv().Class("col-md-12").AddChildren([]hb.TagInterface{
			hb.NewH1().Text("Welcome to Dracory"),
			hb.NewP().Text("A lightweight and efficient Golang web framework designed to simplify the development of web applications."),
			hb.NewHR(),
			hb.NewH2().Text("Key Features"),
			hb.NewUL().AddChildren([]hb.TagInterface{
				hb.NewLI().Text("Simple and Efficient Routing"),
				hb.NewLI().Text("Clean Architecture"),
				hb.NewLI().Text("Configuration Management"),
				hb.NewLI().Text("Database Integration"),
				hb.NewLI().Text("Middleware Support"),
				hb.NewLI().Text("Testing Focus"),
				hb.NewLI().Text("Web Server"),
				hb.NewLI().Text("Web Authentication"),
				hb.NewLI().Text("Minimalistic Templating"),
			}),
			hb.NewHR(),
			hb.NewDiv().Class("d-grid gap-2 d-md-flex justify-content-md-start").AddChildren([]hb.TagInterface{
				hb.NewA().Class("btn btn-primary me-md-2").Attr("href", "/about").Text("Learn More"),
				hb.NewA().Class("btn btn-outline-secondary").Attr("href", "/login").Text("Login"),
			}),
		}),
	})

	// Render the view using the layout
	c.Render(w, c.Layout("Dracory - Home", content))
}
