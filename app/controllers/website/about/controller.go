package about

import (
	"net/http"

	"project/app/config"
	"project/app/controllers/base"
	"project/internal/platform/database"

	"github.com/gouniverse/hb"
)

// Controller handles requests for the about page
type Controller struct {
	base.Controller
}

// NewController creates a new Controller instance
func NewController(cfg *config.Config, db *database.Database) *Controller {
	return &Controller{
		Controller: *base.NewController(cfg, db),
	}
}

// Index handles the about page request
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	// Create a view using the HB library
	content := hb.NewDiv().Class("row").AddChildren([]hb.TagInterface{
		hb.NewDiv().Class("col-md-12").AddChildren([]hb.TagInterface{
			hb.NewH1().Text("About Dracory"),
			hb.NewP().Text("Dracory is a lightweight and efficient Golang web framework designed to simplify the development of web applications. It emphasizes clarity, speed, and ease of use, providing a solid foundation for building scalable and maintainable applications."),
			hb.NewHR(),
			hb.NewH2().Text("Architecture"),
			hb.NewP().Text("Dracory follows a well-defined project structure that promotes maintainability and scalability:"),
			hb.NewPRE().Text(`
dracory-project/
├── cmd/
│   └── web/            # Main application entry point
├── internal/
│   ├── app/
│   │   ├── controllers/ # HTTP handlers and business logic
│   │   ├── models/      # Data models and database interactions
│   │   ├── config/      # Configuration handling
│   │   └── routes/      # Routing definitions
│   └── platform/
│       └── database/   # Database connection and utilities
├── web/                # Static assets (HTML, CSS, JS)
│   ├── static/
│   └── templates/
├── test/               # Integration and end-to-end tests
├── .env                # Environment variables
├── go.mod              # Go module definition
└── README.md           # Project documentation
			`),
			hb.NewHR(),
			hb.NewDiv().Class("d-grid gap-2 d-md-flex justify-content-md-start").AddChildren([]hb.TagInterface{
				hb.NewA().Class("btn btn-primary me-md-2").Attr("href", "/").Text("Home"),
				hb.NewA().Class("btn btn-outline-secondary").Attr("href", "/login").Text("Login"),
			}),
		}),
	})

	// Render the view using the layout
	c.Render(w, c.Layout("Dracory - About", content))
}
