package notfound

import (
	"net/http"

	"project/app/config"
	"project/app/controllers/base"
	"project/internal/platform/database"

	"github.com/gouniverse/hb"
)

// Controller handles 404 not found requests
type Controller struct {
	base.Controller
}

// NewController creates a new Controller instance
func NewController(cfg *config.Config, db *database.Database) *Controller {
	return &Controller{
		Controller: *base.NewController(cfg, db),
	}
}

// Index handles the 404 not found request
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	// Set the status code
	w.WriteHeader(http.StatusNotFound)

	// Get the layout and add the custom style
	layout := c.Layout("Dracory - Not Found", c.view())

	// Render the view
	c.Render(w, layout)
}

func (c *Controller) view() hb.TagInterface {

	// CSS styles for the 404 page
	customStyle := `
		.not-found-container {
			display: flex;
			min-height: 100vh;
		}
		.not-found-content {
			flex: 1;
			padding: 3rem;
			display: flex;
			flex-direction: column;
			justify-content: center;
		}
		.not-found-image {
			flex: 1;
			background: linear-gradient(135deg, #e25c67 0%, #f8b195 100%);
			display: flex;
			align-items: center;
			justify-content: center;
			position: relative;
			overflow: hidden;
		}
		.desert-scene {
			position: absolute;
			bottom: 0;
			width: 100%;
			height: 60%;
			background: #f9e0ae;
			border-radius: 50% 50% 0 0 / 100% 100% 0 0;
		}
		.camel {
			position: absolute;
			bottom: 30%;
			right: 15%;
			color: #8B4513;
			font-size: 2rem;
		}
		.error-code {
			font-size: 8rem;
			font-weight: bold;
			line-height: 1;
			margin-bottom: 0.5rem;
			color: #333;
		}
		.error-divider {
			width: 80px;
			height: 4px;
			background-color: #9c27b0;
			margin: 1.5rem 0;
		}
		.error-message {
			font-size: 1.5rem;
			color: #666;
			margin-bottom: 2rem;
		}
		.home-button {
			display: inline-block;
			padding: 0.75rem 1.5rem;
			background-color: #fff;
			color: #333;
			border: 2px solid #333;
			text-decoration: none;
			font-weight: bold;
			transition: all 0.3s ease;
		}
		.home-button:hover {
			background-color: #333;
			color: #fff;
		}
	`

	// Create a view using the HB library
	content := hb.NewDiv().Class("not-found-container").AddChildren([]hb.TagInterface{
		// Left content side
		hb.NewDiv().Class("not-found-content").AddChildren([]hb.TagInterface{
			hb.NewDiv().Class("error-code").Text("404"),
			hb.NewDiv().Class("error-divider"),
			hb.NewDiv().Class("error-message").Text("Sorry, the page you are looking for could not be found."),
			hb.NewA().Class("home-button").Attr("href", "/").Text("GO HOME"),
		}),
		// Right image side
		hb.NewDiv().Class("not-found-image").AddChildren([]hb.TagInterface{
			hb.NewDiv().Class("desert-scene"),
			hb.NewDiv().Class("camel").HTML("&#x1F42A;"), // Camel emoji
		}),
	})

	return hb.Wrap().
		Child(hb.Style(customStyle)).
		Child(content)
}
