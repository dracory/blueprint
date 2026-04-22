package search

import (
	"context"
	"fmt"
	"net/http"
	"project/internal/links"
	"strconv"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
)

// SearchBlockType represents a block that provides search functionality
type SearchBlockType struct {
	cmsStore  cmsstore.StoreInterface
	blogStore blogstore.StoreInterface
}

// NewSearchBlockType creates a new Search block type
func NewSearchBlockType(cmsStore cmsstore.StoreInterface, blogStore blogstore.StoreInterface) *SearchBlockType {
	return &SearchBlockType{
		cmsStore:  cmsStore,
		blogStore: blogStore,
	}
}

// TypeKey returns the unique identifier for this block type
func (t *SearchBlockType) TypeKey() string {
	return "search"
}

// TypeLabel returns the human-readable display name
func (t *SearchBlockType) TypeLabel() string {
	return "Search"
}

// Render renders the search block for frontend display
func (t *SearchBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, options ...cmsstore.RenderOption) (string, error) {
	// Get configuration
	resultsPerPage := 10
	if val := block.Meta("results_per_page"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			resultsPerPage = n
		}
	}

	placeholder := block.Meta("placeholder")
	if placeholder == "" {
		placeholder = "Search..."
	}

	showPages := block.Meta("show_pages") != "false"
	showPosts := block.Meta("show_posts") != "false"

	// Get search query from request
	searchQuery := ""
	pageNum := 0
	if r := cmsstore.RequestFromContext(ctx); r != nil {
		searchQuery = strings.TrimSpace(r.URL.Query().Get("q"))
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				pageNum = p - 1 // Convert to 0-based
			}
		}
	}

	// Perform search if query exists
	var results []SearchResult
	totalResults := 0

	if searchQuery != "" && (showPages || showPosts) {
		results, totalResults = t.performSearch(ctx, searchQuery, showPages, showPosts, pageNum, resultsPerPage)
	}

	// Render
	return renderSearchHTML(searchQuery, placeholder, results, totalResults, pageNum, resultsPerPage)
}

// SearchResult represents a single search result
type SearchResult struct {
	Title   string
	Summary string
	URL     string
	Type    string // "page" or "post"
	Date    string
}

// performSearch searches across pages and blog posts
func (t *SearchBlockType) performSearch(ctx context.Context, query string, showPages, showPosts bool, pageNum, resultsPerPage int) ([]SearchResult, int) {
	var allResults []SearchResult
	searchLower := strings.ToLower(query)

	// Search pages if enabled and cmsStore is available
	if showPages && t.cmsStore != nil {
		pages, err := t.cmsStore.PageList(ctx, cmsstore.PageQuery().
			SetStatus(cmsstore.PAGE_STATUS_ACTIVE))

		if err == nil {
			for _, page := range pages {
				if pageMatchesSearch(page, searchLower) {
					summary := stripHTML(page.Content())
					if len(summary) > 200 {
						summary = summary[:200] + "..."
					}

					allResults = append(allResults, SearchResult{
						Title:   page.Title(),
						Summary: summary,
						URL:     "/" + strings.TrimPrefix(page.Alias(), "/"),
						Type:    "page",
					})
				}
			}
		}
	}

	// Search blog posts if enabled and blogStore is available
	if showPosts && t.blogStore != nil {
		postOpts := blogstore.PostQueryOptions{
			Status:    blogstore.POST_STATUS_PUBLISHED,
			SortOrder: "DESC",
			OrderBy:   "published_at",
		}

		posts, err := t.blogStore.PostList(ctx, postOpts)
		if err == nil {
			for _, post := range posts {
				if postMatchesSearch(post, searchLower) {
					summary := post.GetSummary()
					if summary == "" && post.GetContent() != "" {
						summary = stripHTML(post.GetContent())
						if len(summary) > 200 {
							summary = summary[:200] + "..."
						}
					}

					dateStr := ""
					if post.GetPublishedAt() != "" {
						dateStr = post.GetPublishedAtCarbon().Format("d M, Y")
					}

					allResults = append(allResults, SearchResult{
						Title:   post.GetTitle(),
						Summary: summary,
						URL:     links.Website().BlogPost(post.GetID(), post.GetSlug()),
						Type:    "post",
						Date:    dateStr,
					})
				}
			}
		}
	}

	// Calculate total
	total := len(allResults)

	// Paginate results
	offset := pageNum * resultsPerPage
	if offset >= total {
		return []SearchResult{}, total
	}

	end := offset + resultsPerPage
	if end > total {
		end = total
	}

	return allResults[offset:end], total
}

// pageMatchesSearch checks if a page matches the search query
func pageMatchesSearch(page cmsstore.PageInterface, query string) bool {
	if page == nil {
		return false
	}

	// Check title
	if strings.Contains(strings.ToLower(page.Title()), query) {
		return true
	}

	// Check content (HTML stripped)
	contentLower := strings.ToLower(stripHTML(page.Content()))
	if strings.Contains(contentLower, query) {
		return true
	}

	// Check alias
	if strings.Contains(strings.ToLower(page.Alias()), query) {
		return true
	}

	// Check meta description
	if strings.Contains(strings.ToLower(page.MetaDescription()), query) {
		return true
	}

	return false
}

// postMatchesSearch checks if a blog post matches the search query
func postMatchesSearch(post blogstore.PostInterface, query string) bool {
	// Check title
	if strings.Contains(strings.ToLower(post.GetTitle()), query) {
		return true
	}

	// Check content (HTML stripped)
	contentLower := strings.ToLower(stripHTML(post.GetContent()))
	if strings.Contains(contentLower, query) {
		return true
	}

	// Check summary
	if strings.Contains(strings.ToLower(post.GetSummary()), query) {
		return true
	}

	// Check slug
	if strings.Contains(strings.ToLower(post.GetSlug()), query) {
		return true
	}

	return false
}

// stripHTML removes HTML tags from a string
func stripHTML(html string) string {
	// Simple HTML tag removal
	result := html
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return strings.TrimSpace(result)
}

// GetAdminFields returns form fields for editing block configuration
func (t *SearchBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Placeholder Text",
			Name:  "placeholder",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta("placeholder"),
			Help:  "Placeholder text for the search input (default: Search...)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Results Per Page",
			Name:  "results_per_page",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta("results_per_page"),
			Help:  "Number of results to display per page (default: 10)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Pages",
			Name:  "show_pages",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_pages"),
			Help:  "Include CMS pages in search results",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Blog Posts",
			Name:  "show_posts",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_posts"),
			Help:  "Include blog posts in search results",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
	}

	return fields
}

// SaveAdminFields processes form submission and updates the block
func (t *SearchBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	block.SetMeta("placeholder", r.FormValue("placeholder"))
	block.SetMeta("results_per_page", r.FormValue("results_per_page"))
	block.SetMeta("show_pages", r.FormValue("show_pages"))
	block.SetMeta("show_posts", r.FormValue("show_posts"))

	return nil
}

// Validate validates the block configuration
func (t *SearchBlockType) Validate(block cmsstore.BlockInterface) error {
	return nil
}

// GetPreview returns a preview of the block
func (t *SearchBlockType) GetPreview(block cmsstore.BlockInterface) string {
	showPages := block.Meta("show_pages") != "false"
	showPosts := block.Meta("show_posts") != "false"

	var types []string
	if showPages {
		types = append(types, "Pages")
	}
	if showPosts {
		types = append(types, "Blog Posts")
	}

	if len(types) == 0 {
		return "Search (no content types selected)"
	}

	return fmt.Sprintf("Search (%s)", strings.Join(types, ", "))
}

// GetCustomVariables returns a list of custom variables that this block type sets.
// This block type does not currently set any custom variables.
func (t *SearchBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{}
}
