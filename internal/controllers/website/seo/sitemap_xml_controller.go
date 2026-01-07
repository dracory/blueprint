package seo

import (
	"context"
	"log/slog"
	"net/http"
	"project/internal/links"
	"project/internal/registry"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

type sitemapXmlController struct {
	registry registry.RegistryInterface
}

// NewSitemapXmlController creates a new instance of the sitemapXmlController struct.
func NewSitemapXmlController(registry registry.RegistryInterface) *sitemapXmlController {
	return &sitemapXmlController{registry: registry}
}

func (c sitemapXmlController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/xml")
	body := c.buildSitemapXML(w, r)
	w.Write([]byte(body))
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
	if c.registry == nil {
		slog.Warn("At sitemapXmlController > blogPostLocations", slog.String("reason", "registry is not configured"))
		return []string{}
	}

	if !c.registry.GetConfig().GetBlogStoreUsed() {
		return []string{}
	}

	if c.registry.GetBlogStore() == nil {
		slog.Warn("At sitemapXmlController > blogPostLocations", slog.String("reason", "blog store is not configured"))
		return []string{}
	}

	postList, err := c.registry.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{
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
