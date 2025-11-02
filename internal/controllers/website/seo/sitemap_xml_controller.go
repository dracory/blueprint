package seo

import (
	"log/slog"
	"net/http"
	"project/internal/links"
	"project/internal/types"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/responses"
	"github.com/samber/lo"
)

type sitemapXmlController struct {
	app types.AppInterface
}

// NewSitemapXmlController creates a new instance of the sitemapXmlController struct.
func NewSitemapXmlController(app types.AppInterface) *sitemapXmlController {
	return &sitemapXmlController{app: app}
}

func (c sitemapXmlController) Handler(w http.ResponseWriter, r *http.Request) string {
	responses.XMLResponseF(w, r, c.buildSitemapXML)
	return ""
}

func (c sitemapXmlController) buildSitemapXML(w http.ResponseWriter, r *http.Request) string {
	locations := []string{
		"/",
		// "/about",
		// "/contact",
		// "/faq",
		// "/marketplace",
		"/robots.txt",
		// "/privacy-policy",
		// "/terms-of-use",
	}

	locations = append(locations, c.blogPostLocations()...)

	timeNow := carbon.Now().ToIso8601String()

	xml := `<?xml version="1.0" encoding="UTF-8"?>`
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	lo.ForEach(locations, func(location string, index int) {
		url := lo.Ternary(strings.HasPrefix(location, "http"), location, links.RootURL()+location)

		priority := "0.80"
		if index == 0 {
			priority = "1.00"
		}

		xml += "<url>"
		xml += "<loc>" + url + "</loc>"
		xml += "<lastmod>" + timeNow + "</lastmod>"
		xml += "<priority>" + priority + "</priority>"
		xml += "</url>"
	})
	xml += "</urlset>"

	return xml
}

func (c sitemapXmlController) blogPostLocations() []string {
	if c.app == nil {
		slog.Warn("At sitemapXmlController > blogPostLocations", slog.String("reason", "app is not configured"))
		return []string{}
	}

	if !c.app.GetConfig().GetBlogStoreUsed() {
		return []string{}
	}

	if c.app.GetBlogStore() == nil {
		slog.Warn("At sitemapXmlController > blogPostLocations", slog.String("reason", "blog store is not configured"))
		return []string{}
	}

	postList, err := c.app.GetBlogStore().PostList(blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		OrderBy:   "title",
		SortOrder: sb.DESC,
		Limit:     1000,
	})

	if err != nil {
		slog.Error("At sitemapXmlController > blogPostLocations", slog.String("error", err.Error()))
		return nil
	}

	postLocations := make([]string, 0, len(postList))
	lo.ForEach(postList, func(post blogstore.Post, index int) {
		postLocations = append(postLocations, links.Website().BlogPost(post.ID(), post.Slug()))
	})

	return postLocations
}
