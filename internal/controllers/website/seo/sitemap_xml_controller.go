package seo

import (
	"context"
	"log/slog"
	"net/http"
	"project/internal/app"
	"project/internal/links"
	"strings"
	"time"

	"github.com/dracory/blogstore"
	"github.com/dracory/neat"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

const (
	// sitemapMemCacheKey/TTL back the in-memory copy of the built XML so
	// repeat crawler hits avoid the blog store SELECTs entirely.
	sitemapMemCacheKey = "website:sitemap.xml:mem"
	sitemapMemCacheTTL = time.Hour
)

type sitemapXmlController struct {
	app app.AppInterface
}

// NewSitemapXmlController creates a new instance of the sitemapXmlController struct.
func NewSitemapXmlController(app app.AppInterface) *sitemapXmlController {
	return &sitemapXmlController{app: app}
}

func (c sitemapXmlController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	body := c.cachedSitemapXML(w, r)
	_, _ = w.Write([]byte(body))
	return ""
}

// cachedSitemapXML serves the sitemap from the memory cache when possible.
// Crawlers hit the sitemap constantly and the document only changes when a
// post is published, so a 1h TTL is safe; on expiry one request rebuilds
// and repopulates the cache.
func (c sitemapXmlController) cachedSitemapXML(w http.ResponseWriter, r *http.Request) string {
	if c.app != nil {
		if mc := c.app.GetMemoryCache(); mc != nil {
			if item := mc.Get(sitemapMemCacheKey); item != nil {
				if cached, ok := item.Value().(string); ok && cached != "" {
					return cached
				}
			}
		}
	}

	body := c.buildSitemapXML(w, r)

	if c.app != nil {
		if mc := c.app.GetMemoryCache(); mc != nil {
			mc.Set(sitemapMemCacheKey, body, sitemapMemCacheTTL)
		}
	}

	return body
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

	if c.app.IsDisabledBlogStore() {
		slog.Warn("At sitemapXmlController > blogPostLocations", slog.String("reason", "blog store is not configured"))
		return []string{}
	}

	postList, err := c.app.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		OrderBy:   "title",
		SortOrder: neat.SortDesc,
		Limit:     1000,
	})

	if err != nil {
		slog.Error("At sitemapXmlController > blogPostLocations", slog.String("error", err.Error()))
		return nil
	}

	postLocations := make([]string, 0, len(postList))
	lo.ForEach(postList, func(post blogstore.PostInterface, index int) {
		postLocations = append(postLocations, links.Website().BlogPost(post.GetID(), post.GetSlug()))
	})

	return postLocations
}
