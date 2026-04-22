package blogpostlist

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/spf13/cast"
)

// Custom variable names set by the blog post list block
const (
	VarTagName        = "tag_name"
	VarTagSlug        = "tag_slug"
	VarTagDescription = "tag_description"
	VarTagID          = "tag_id"
)

// BlogPostListBlockType represents a block that renders a list of blog posts
type BlogPostListBlockType struct {
	blogStore blogstore.StoreInterface
}

// NewBlogPostListBlockType creates a new BlogPostList block type
func NewBlogPostListBlockType(blogStore blogstore.StoreInterface) *BlogPostListBlockType {
	return &BlogPostListBlockType{
		blogStore: blogStore,
	}
}

// TypeKey returns the unique identifier for this block type
func (t *BlogPostListBlockType) TypeKey() string {
	return "blog_post_list"
}

// TypeLabel returns the human-readable display name
func (t *BlogPostListBlockType) TypeLabel() string {
	return "Blog Post List"
}

// Render renders the blog post list block for frontend display
func (t *BlogPostListBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, options ...cmsstore.RenderOption) (string, error) {
	// Get configuration
	postsPerPage := 12
	if val := block.Meta("posts_per_page"); val != "" {
		if n := cast.ToInt(val); n > 0 {
			postsPerPage = n
		}
	}

	showPagination := block.Meta("show_pagination") != "false"
	showImages := block.Meta("show_images") != "false"
	showSummary := block.Meta("show_summary") != "false"
	showDate := block.Meta("show_date") != "false"

	excerptLength := 150
	if val := block.Meta("excerpt_length"); val != "" {
		if n := cast.ToInt(val); n > 0 {
			excerptLength = n
		}
	}

	columns := 4
	if val := block.Meta("columns"); val != "" {
		if n := cast.ToInt(val); n > 0 && n <= 6 {
			columns = n
		}
	}

	// Check URL for tag filtering
	tagSlug := t.extractTagSlugFromURL(ctx)
	currentPage := t.extractCurrentPageFromURL(ctx)

	var postList []blogstore.PostInterface
	var postCount int64
	var err error

	if tagSlug != "" {
		// Find the tag by slug
		tag, err := t.blogStore.TermFindBySlug(ctx, blogstore.TAXONOMY_TAG, tagSlug)
		if err != nil || tag == nil {
			// Tag not found, return empty list
			postList = []blogstore.PostInterface{}
			postCount = 0
		} else {
			// Get posts with this tag
			offset := (currentPage - 1) * postsPerPage
			opts := blogstore.PostQueryOptions{
				Status:    blogstore.POST_STATUS_PUBLISHED,
				SortOrder: "DESC",
				OrderBy:   "published_at",
				Limit:     postsPerPage,
				Offset:    offset,
			}
			postList, err = t.blogStore.PostListByTermID(ctx, tag.GetID(), opts)
			if err != nil {
				return "", fmt.Errorf("failed to get posts for tag: %w", err)
			}
			// Get total count for pagination
			countOpts := blogstore.PostQueryOptions{
				Status:    blogstore.POST_STATUS_PUBLISHED,
				SortOrder: "DESC",
				OrderBy:   "published_at",
			}
			allPosts, err := t.blogStore.PostListByTermID(ctx, tag.GetID(), countOpts)
			if err != nil {
				return "", fmt.Errorf("failed to get post count for tag: %w", err)
			}
			postCount = int64(len(allPosts))

			// Set custom variables for tag information
			if vars := cmsstore.VarsFromContext(ctx); vars != nil {
				vars.Set(VarTagName, tag.GetName())
				vars.Set(VarTagSlug, tag.GetSlug())
				vars.Set(VarTagDescription, tag.GetDescription())
				vars.Set(VarTagID, tag.GetID())
			}
		}
	} else {
		// Get all posts
		offset := (currentPage - 1) * postsPerPage
		opts := blogstore.PostQueryOptions{
			Status:    blogstore.POST_STATUS_PUBLISHED,
			SortOrder: "DESC",
			OrderBy:   "published_at",
			Limit:     postsPerPage,
			Offset:    offset,
		}

		postList, err = t.blogStore.PostList(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to get posts: %w", err)
		}

		postCount, err = t.blogStore.PostCount(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to get post count: %w", err)
		}
	}

	// Render
	return renderBlogPostListHTML(postList, postCount, postsPerPage, currentPage, showPagination, showImages, showSummary, showDate, columns, excerptLength)
}

// extractTagSlugFromURL extracts tag slug from URL if present
// URL format: /tag/{tag-slug}/
func (t *BlogPostListBlockType) extractTagSlugFromURL(ctx context.Context) string {
	r := cmsstore.RequestFromContext(ctx)
	if r == nil {
		return ""
	}

	path := r.URL.Path
	// Look for /tag/ in the URL
	tagIndex := strings.Index(path, "/tag/")
	if tagIndex == -1 {
		return ""
	}

	// Extract everything after /tag/
	afterTag := path[tagIndex+5:]

	// Remove trailing slash if present
	afterTag = strings.TrimSuffix(afterTag, "/")

	// Handle case where there might be additional path segments
	// Take only the first segment as the tag slug
	if idx := strings.Index(afterTag, "/"); idx != -1 {
		afterTag = afterTag[:idx]
	}

	return afterTag
}

// extractCurrentPageFromURL extracts current page from URL query parameter
func (t *BlogPostListBlockType) extractCurrentPageFromURL(ctx context.Context) int {
	r := cmsstore.RequestFromContext(ctx)
	if r == nil {
		return 1
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		return 1
	}

	page := cast.ToInt(pageStr)
	if page < 1 {
		return 1
	}
	return page
}

// GetAdminFields returns form fields for editing block configuration
func (t *BlogPostListBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Posts Per Page",
			Name:  "posts_per_page",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta("posts_per_page"),
			Help:  "Number of posts to display per page (default: 12)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Columns",
			Name:  "columns",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("columns"),
			Help:  "Number of columns to display posts in",
			Options: []form.FieldOption{
				{Value: "1 column", Key: "1"},
				{Value: "2 columns", Key: "2"},
				{Value: "3 columns", Key: "3"},
				{Value: "4 columns", Key: "4"},
				{Value: "6 columns", Key: "6"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Pagination",
			Name:  "show_pagination",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_pagination"),
			Help:  "Show pagination controls",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Images",
			Name:  "show_images",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_images"),
			Help:  "Display post featured images",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Summary",
			Name:  "show_summary",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_summary"),
			Help:  "Display post summary/excerpt",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Date",
			Name:  "show_date",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_date"),
			Help:  "Display publication date",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Excerpt Length",
			Name:  "excerpt_length",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta("excerpt_length"),
			Help:  "Maximum characters for auto-generated excerpts (default: 150)",
		}),
	}

	return fields
}

// SaveAdminFields processes form submission and updates the block
func (t *BlogPostListBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	block.SetMeta("posts_per_page", r.FormValue("posts_per_page"))
	block.SetMeta("columns", r.FormValue("columns"))
	block.SetMeta("show_pagination", r.FormValue("show_pagination"))
	block.SetMeta("show_images", r.FormValue("show_images"))
	block.SetMeta("show_summary", r.FormValue("show_summary"))
	block.SetMeta("show_date", r.FormValue("show_date"))
	block.SetMeta("excerpt_length", r.FormValue("excerpt_length"))

	return nil
}

// Validate validates the block configuration
func (t *BlogPostListBlockType) Validate(block cmsstore.BlockInterface) error {
	return nil
}

// GetPreview returns a preview of the block
func (t *BlogPostListBlockType) GetPreview(block cmsstore.BlockInterface) string {
	return "Blog Post List"
}

// GetCustomVariables returns a list of custom variables that this block type sets.
// These variables can be used in page templates when viewing tag-filtered lists.
func (t *BlogPostListBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{
		{Name: VarTagName, Description: "The name of the tag when viewing a tag-filtered list."},
		{Name: VarTagSlug, Description: "The slug of the tag when viewing a tag-filtered list."},
		{Name: VarTagDescription, Description: "The description of the tag when viewing a tag-filtered list."},
		{Name: VarTagID, Description: "The unique ID of the tag when viewing a tag-filtered list."},
	}
}
