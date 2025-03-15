package base

import (
	"encoding/json"
	"net/http"

	"project/app/config"
	"project/internal/platform/database"

	"github.com/dracory/base/bbcode"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
)

// Controller represents a base controller with common functionality
type Controller struct {
	Config   *config.Config
	Database *database.Database
}

// NewController creates a new Controller instance
func NewController(cfg *config.Config, db *database.Database) *Controller {
	return &Controller{
		Config:   cfg,
		Database: db,
	}
}

// Render renders an HTML view using the HB library
func (c *Controller) Render(w http.ResponseWriter, view hb.TagInterface) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(view.ToHTML()))
}

// RenderBBCode renders BBCode content as HTML
func (c *Controller) RenderBBCode(w http.ResponseWriter, content string) {
	html := bbcode.BbcodeToHtml(content)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// JSON sends a JSON response
func (c *Controller) JSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Error sends an error response
func (c *Controller) Error(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// Redirect redirects to the specified URL
func (c *Controller) Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	http.Redirect(w, r, url, statusCode)
}

// GetParam gets a query parameter from the request
func (c *Controller) GetParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetFormValue gets a form value from the request
func (c *Controller) GetFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

// Layout creates a basic HTML layout with the given title and content
func (c *Controller) Layout(title string, content hb.TagInterface) hb.TagInterface {
	// 2. Webpage Favicon
	favicon := "data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAABNTU0AVKH/AOPj4wDExMQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACIiAREQEREAIiIBERAREQAiIgIiICIiACIiAiIgIiIAMzMDMzAzMwAzMwMzMDMzACIiAiIgIiIAIiICIiAiIgAzMwMzMDMzADMzAzMwMzMAIiICIiAiIgAiIgIiICIiAAAAAAAAAAAAIiICIiAiIgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	webpage := hb.Webpage().
		SetTitle(title).
		SetFavicon(favicon).
		AddStyleURLs([]string{
			cdn.BootstrapCss_5_3_3(),
		}).
		AddScriptURLs([]string{
			cdn.BootstrapJs_5_3_3(),
		}).
		AddScripts([]string{}).
		AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`).
		AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}`)
	webpage.AddChild(content)
	return webpage
}
