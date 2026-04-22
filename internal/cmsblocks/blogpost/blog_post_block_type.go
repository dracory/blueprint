package blogpost

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"project/internal/links"

	"github.com/dracory/blogstore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
)

// Custom variable names set by the blog post block
const (
	VarBlogTitle         = "blog_title"
	VarBlogSummary       = "blog_summary"
	VarBlogSlug          = "blog_slug"
	VarBlogID            = "blog_id"
	VarBlogDate          = "blog_date"
	VarBlogDateFormatted = "blog_date_formatted"
	VarBlogMetaTitle     = "blog_meta_title"
	VarBlogMetaDesc      = "blog_meta_description"
	VarBlogCanonicalURL  = "blog_canonical_url"

	// Commented out variables (currently disabled)
	VarBlogFeaturedImage = "blog_featured_image"
	VarBlogAuthor        = "blog_author"
	VarBlogAuthorID      = "blog_author_id"
)

// BlogPostBlockType represents a block that renders a single blog post
type BlogPostBlockType struct {
	blogStore blogstore.StoreInterface
}

// NewBlogPostBlockType creates a new BlogPost block type
func NewBlogPostBlockType(blogStore blogstore.StoreInterface) *BlogPostBlockType {
	return &BlogPostBlockType{
		blogStore: blogStore,
	}
}

// TypeKey returns the unique identifier for this block type
func (t *BlogPostBlockType) TypeKey() string {
	return "blog_post"
}

// TypeLabel returns the human-readable display name
func (t *BlogPostBlockType) TypeLabel() string {
	return "Blog Post"
}

// Render renders the single blog post block for frontend display
func (t *BlogPostBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, options ...cmsstore.RenderOption) (string, error) {
	// Get post ID from block meta
	postID := block.Meta("post_id")

	// If no specific post ID, try to get from context/URL
	if postID == "" {
		// Try to get from request context using cmsstore helper
		r := cmsstore.RequestFromContext(ctx)

		if r != nil {
			uriParts := strings.Split(r.RequestURI, "/")

			// Handle both old format (/blog/{id}/{slug}) and new format (/blog/post/{id}/{slug})
			if len(uriParts) >= 3 {
				// Check if this is the new format with /blog/post/{id}/{slug}
				if len(uriParts) >= 4 && uriParts[2] == "post" {
					postID = uriParts[3] // New format: /blog/post/{id}/{slug}
				} else {
					postID = uriParts[2] // Old format: /blog/{id}/{slug}
				}
			}
		}
	}

	if postID == "" {
		return "", fmt.Errorf("no post ID specified for blog post block")
	}

	// Fetch the post
	post, err := t.blogStore.PostFindByID(ctx, postID)
	if err != nil {
		return "", fmt.Errorf("failed to find post: %w", err)
	}

	if post == nil {
		return "", fmt.Errorf("post not found: %s", postID)
	}

	// Check if post is published
	if post.IsUnpublished() {
		return "", fmt.Errorf("post is not published")
	}

	// Set custom variables for page title, meta tags, etc.
	if vars := cmsstore.VarsFromContext(ctx); vars != nil {
		vars.Set(VarBlogTitle, post.GetTitle())
		vars.Set(VarBlogSummary, post.GetSummary())
		vars.Set(VarBlogSlug, post.GetSlug())
		vars.Set(VarBlogID, post.GetID())

		// if post.GetFeaturedImage() != "" {
		// 	vars.Set(VarBlogFeaturedImage, post.GetFeaturedImage())
		// }

		// if post.GetAuthor() != nil {
		// 	vars.Set(VarBlogAuthor, post.GetAuthor().Name())
		// 	vars.Set(VarBlogAuthorID, post.GetAuthor().ID())
		// }

		if post.GetPublishedAt() != "" {
			vars.Set(VarBlogDate, post.GetPublishedAt())
			vars.Set(VarBlogDateFormatted, post.GetPublishedAtCarbon().Format("d M, Y"))
		}

		// SEO/meta variables
		vars.Set(VarBlogMetaTitle, post.GetTitle())
		if post.GetSummary() != "" {
			vars.Set(VarBlogMetaDesc, post.GetSummary())
		}

		// Canonical URL
		if r := cmsstore.RequestFromContext(ctx); r != nil {
			vars.Set(VarBlogCanonicalURL, links.Website().BlogPost(post.GetID(), post.GetSlug()))
		}
	}

	// Get display options
	showImage := block.Meta("show_image") != "false"
	showTitle := block.Meta("show_title") != "false"
	showDate := block.Meta("show_date") != "false"
	showAuthor := block.Meta("show_author") != "false"
	showSummary := block.Meta("show_summary") != "false"
	showPrevNext := block.Meta("show_prev_next") == "true"

	// Fetch previous and next posts if enabled
	var prevPost, nextPost blogstore.PostInterface
	if showPrevNext && t.blogStore != nil {
		prevPost, nextPost = t.getPrevNextPosts(ctx, post)
	}

	// Fetch tags for this post
	var postTags []blogstore.TermInterface
	if t.blogStore != nil {
		tags, err := t.blogStore.TermListByPostID(ctx, post.GetID(), blogstore.TAXONOMY_TAG)
		if err == nil {
			postTags = tags
		}
	}

	// Render
	return renderBlogPostHTML(
		post,
		showImage,
		showTitle,
		showDate,
		showAuthor,
		showSummary,
		showPrevNext,
		prevPost,
		nextPost,
		postTags,
	)
}

// getPrevNextPosts fetches previous and next posts based on published date
func (t *BlogPostBlockType) getPrevNextPosts(ctx context.Context, currentPost blogstore.PostInterface) (blogstore.PostInterface, blogstore.PostInterface) {
	var prevPost, nextPost blogstore.PostInterface

	// Fetch all published posts sorted by published_at
	opts := blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		OrderBy:   "published_at",
		SortOrder: "DESC",
		Limit:     100,
	}
	posts, err := t.blogStore.PostList(ctx, opts)
	if err != nil || len(posts) == 0 {
		return nil, nil
	}

	// Find current post index and get prev/next
	for i, p := range posts {
		if p.GetID() == currentPost.GetID() {
			// Previous post is the next one in the list (older, published earlier)
			if i+1 < len(posts) {
				prevPost = posts[i+1]
			}
			// Next post is the previous one in the list (newer, published later)
			if i-1 >= 0 {
				nextPost = posts[i-1]
			}
			break
		}
	}

	return prevPost, nextPost
}

// GetAdminFields returns form fields for editing block configuration
func (t *BlogPostBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	postOptions := []form.FieldOption{
		{Value: "- Select Post -", Key: ""},
	}

	// Only fetch posts if blogStore is available
	if t.blogStore != nil && r != nil {
		opts := blogstore.PostQueryOptions{
			Status:    blogstore.POST_STATUS_PUBLISHED,
			SortOrder: "DESC",
			OrderBy:   "published_at",
			Limit:     100,
		}

		postList, err := t.blogStore.PostList(r.Context(), opts)
		if err == nil {
			for _, post := range postList {
				postOptions = append(postOptions, form.FieldOption{
					Value: post.GetTitle(),
					Key:   post.GetID(),
				})
			}
		}
	}

	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Post",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_SELECT,
			Value:    block.Meta("post_id"),
			Required: false,
			Help:     "Select a specific post to display (leave empty to auto-detect from URL)",
			Options:  postOptions,
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Image",
			Name:  "show_image",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_image"),
			Help:  "Display the post featured image",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Title",
			Name:  "show_title",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_title"),
			Help:  "Display the post title",
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
			Help:  "Display the publication date",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Author",
			Name:  "show_author",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_author"),
			Help:  "Display the author name",
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
			Help:  "Display the post summary/excerpt",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Show Previous/Next",
			Name:  "show_prev_next",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("show_prev_next"),
			Help:  "Display previous and next post navigation",
			Options: []form.FieldOption{
				{Value: "Yes", Key: "true"},
				{Value: "No", Key: "false"},
			},
		}),
	}

	return fields
}

// SaveAdminFields processes form submission and updates the block
func (t *BlogPostBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	block.SetMeta("post_id", r.FormValue("post_id"))
	block.SetMeta("show_image", r.FormValue("show_image"))
	block.SetMeta("show_title", r.FormValue("show_title"))
	block.SetMeta("show_date", r.FormValue("show_date"))
	block.SetMeta("show_author", r.FormValue("show_author"))
	block.SetMeta("show_summary", r.FormValue("show_summary"))
	block.SetMeta("show_prev_next", r.FormValue("show_prev_next"))

	return nil
}

// Validate validates the block configuration
func (t *BlogPostBlockType) Validate(block cmsstore.BlockInterface) error {
	return nil
}

// GetPreview returns a preview of the block
func (t *BlogPostBlockType) GetPreview(block cmsstore.BlockInterface) string {
	postID := block.Meta("post_id")
	if postID == "" {
		return "Blog Post (auto-detect from URL)"
	}

	// If blogStore is nil, just return the ID
	if t.blogStore == nil {
		return fmt.Sprintf("Blog Post: %s", postID)
	}

	// Try to get post title
	ctx := context.Background()
	post, err := t.blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Sprintf("Blog Post: %s", postID)
	}

	return fmt.Sprintf("Blog Post: %s", post.GetTitle())
}

// GetCustomVariables returns a list of custom variables that this block type sets.
// These variables can be used in page templates and content.
func (t *BlogPostBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{
		{Name: VarBlogTitle, Description: "The title of the blog post."},
		{Name: VarBlogSummary, Description: "A short summary of the blog post."},
		{Name: VarBlogSlug, Description: "The URL-friendly slug for the post."},
		{Name: VarBlogID, Description: "The unique ID of the blog post."},
		{Name: VarBlogDate, Description: "The publication date of the post."},
		{Name: VarBlogDateFormatted, Description: "The formatted publication date."},
		{Name: VarBlogMetaTitle, Description: "The SEO meta title for the post."},
		{Name: VarBlogMetaDesc, Description: "The SEO meta description for the post."},
		{Name: VarBlogCanonicalURL, Description: "The canonical URL for the post."},
	}
}
