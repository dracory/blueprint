package post

import (
	"bytes"
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"project/internal/app"
	"project/internal/controllers/website/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/rules"
	"project/pkg/blogadmin/post_update"
	"strings"
	"time"

	"github.com/dracory/base/blogtheme"
	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/rtr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type postController struct {
	app app.AppInterface
}

func NewPostController(
	app app.AppInterface,
) *postController {
	return &postController{
		app: app,
	}
}

func (c *postController) Handler(w http.ResponseWriter, r *http.Request) string {
	postID, _ := rtr.GetParam(r, "id")
	postSlug, _ := rtr.GetParam(r, "title")
	blogsUrl := links.Website().Blog(map[string]string{})

	if postID == "" {
		c.app.GetLogger().Error("anyPost: post ID is missing", slog.String("uri", r.RequestURI))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	post, errPost := c.app.GetBlogStore().PostFindByID(r.Context(), postID)

	if errPost != nil {
		c.app.GetLogger().Error("Error. At BlogPostController.AnyIndex. Post not found", slog.String("error", errPost.Error()))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	if post == nil {
		c.app.GetLogger().Error("ERROR: anyPost: post with ID "+postID+" is missing", slog.String("postID", postID))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return ""
	}

	if !c.accessAllowed(r, post) {
		c.app.GetLogger().Error("WARNING: anyPost: post with ID "+postID+" is unpublished", slog.String("postID", postID))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for is no longer active. Redirecting to the blog location...", blogsUrl, 5)
		return ""
	}

	if postSlug == "" || postSlug != post.GetSlug() {
		url := links.Website().BlogPost(post.GetID(), post.GetSlug())
		c.app.GetLogger().Error("ERROR: anyPost: post Title is missing for ID "+postID, slog.String("postID", postID))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "success", "The post location has changed. Redirecting to the new address...", url, 5)
		return ""
	}

	canonicalURL := strings.TrimSpace(post.GetCanonicalURL())
	if canonicalURL == "" {
		canonicalURL = links.Website().BlogPost(post.GetID(), post.GetSlug())
	}

	options := layouts.Options{
		WebsiteSection: "Blog",
		Title:          post.GetTitle(),
		ScriptURLs: []string{
			cdn.Slazy_0_5_0(),
		},
		Content:      hb.Wrap().HTML(c.page(post, r.Context())),
		CanonicalURL: canonicalURL,
	}

	return layouts.NewPageLayout(
		c.app,
		r,
		options).ToHTML()
}

func (controller *postController) accessAllowed(r *http.Request, post blogstore.PostInterface) bool {
	if post.IsPublished() {
		return true // everyone can access published posts
	}

	authUser := helpers.GetAuthUser(r)

	// If the user is not logged in, they can't access unpublished posts
	if authUser == nil {
		return false
	}

	// If the user is an administrator, they can access unpublished posts
	if authUser.IsAdministrator() {
		return true
	}

	return false // default to false
}

func (c *postController) page(post blogstore.PostInterface, ctx context.Context) string {
	sectionBanner := shared.SectionBanner()
	return hb.Wrap().Children([]hb.TagInterface{
		// hb.Style(c.cssSectionIntro()),
		hb.Style(c.css()),
		sectionBanner,
		// c.sectionIntro(),
		c.sectionPost(post),
		c.recommendationsSection(post),
	}).ToHTML()
}

// css returns the CSS code for the blog post section.
//
// No parameters.
// Returns a string containing the CSS code.
func (c *postController) css() string {
	return `
#SectionNewsItem {
	padding:20px 0px 60px 0px;
}

#SectionNewsItem .BlogTitle {
	padding:10px 0px 20px 0px;
	font-size: clamp(1.75rem, 4vw, 2.5rem);
	font-weight: 900;
	text-transform: uppercase;
	letter-spacing: -1px;
	text-align:left;
}

#SectionNewsItem .BlogContent {
	padding:10px 0px 20px 0px;
	text-align:left;
	font-size:18px;
	line-height: 1.7;
}

#SectionNewsItem .BlogContent h2,
#SectionNewsItem .BlogContent h3 {
	font-weight: 900;
	text-transform: uppercase;
	letter-spacing: -0.5px;
	margin-top: 1.5em;
}

#SectionNewsItem .BlogContent img {
	max-width: 100%;
	height: auto;
}
	`
}

func (controller *postController) processContent(content string, editor string, contentType string) (html string, css string) {
	// Use content type as primary determinant, fallback to editor for compatibility

	// blocks content type
	if contentType == blogstore.POST_CONTENT_TYPE_BLOCKS {
		if editor == blogstore.POST_EDITOR_BLOCKAREA {
			return helpers.BlogPostBlocksToString(content), ""
		}
		if editor == blogstore.POST_EDITOR_BLOCKEDITOR {
			theme, err := blogtheme.New(content)
			if err != nil {
				return "Error parsing content. Please try again later.", ""
			}
			return theme.ToHtml(), theme.Style()
		}
	}

	// markdown content type
	if contentType == blogstore.POST_CONTENT_TYPE_MARKDOWN {
		return controller.markdownToHtml(content), ""
	}

	// html content type
	if contentType == blogstore.POST_CONTENT_TYPE_HTML {
		return content, ""
	}

	// plain text content type (default)
	if contentType == blogstore.POST_CONTENT_TYPE_PLAIN_TEXT || contentType == "" {
		return content, ""
	}

	// Fallback to editor-based processing for backward compatibility
	// blockarea
	if editor == blogstore.POST_EDITOR_BLOCKAREA {
		return helpers.BlogPostBlocksToString(content), ""
	}

	// blockeditor
	if editor == blogstore.POST_EDITOR_BLOCKEDITOR {
		theme, err := blogtheme.New(content)

		if err != nil {
			return "Error parsing content. Please try again later.", ""
		}

		return theme.ToHtml(), theme.Style()
	}

	// markdown
	if editor == blogstore.POST_EDITOR_MARKDOWN {
		return controller.markdownToHtml(content), ""
	}

	// markdown easy mde
	if editor == post_update.PostEditorMarkdownEasyMDE {
		return controller.markdownToHtml(content), ""
	}

	return content, ""
}

func (c *postController) sectionPost(post blogstore.PostInterface) *hb.Tag {
	postHtml, themeStyle := c.processContent(post.GetContent(), post.GetEditor(), post.GetContentType())

	rowTitle := hb.Div().
		Class("BlogTitle").
		Child(hb.Heading1().
			HTML(post.GetTitle()))

	rowContent := hb.Div().
		Child(hb.Div().
			Class("BlogImage float-end").
			Style("padding-top:10px; padding-left:30px; padding-bottom:30px; width:600px; max-width:100%;").
			Child(c.postImage(post))).
		Child(hb.Div().
			Class("BlogContent").
			HTML(postHtml))

	sectionPost := hb.Section().
		ID("SectionNewsItem").
		Child(hb.Style(themeStyle)).
		Child(hb.Div().
			Class("container").
			Child(rowTitle).
			Child(rowContent))

	return sectionPost
}

func (c *postController) recommendationsSection(post blogstore.PostInterface) hb.TagInterface {
	if c.app == nil || c.app.GetBlogStore() == nil {
		return hb.Div()
	}

	options := blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		SortOrder: "DESC",
		OrderBy:   "published_at",
		Limit:     18,
	}

	postList, err := c.app.GetBlogStore().PostList(context.Background(), options)
	if err != nil {
		return hb.Div()
	}

	filtered := make([]blogstore.PostInterface, 0, len(postList))
	for _, p := range postList {
		if p.GetID() == post.GetID() {
			continue
		}
		if p.IsPublished() {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return hb.Div()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	if len(filtered) > 3 {
		filtered = filtered[:3]
	}

	cards := make([]hb.TagInterface, 0, len(filtered))
	for _, p := range filtered {
		cards = append(cards, bs.Column(12).
			Class("col-md-4 col-sm-6 d-flex").
			Child(c.postCard(p)))
	}

	heading := hb.Heading3().
		Class("text-center mb-2 fw-black text-uppercase").
		Style("font-size: clamp(1.25rem, 3vw, 1.75rem); letter-spacing: -1px;").
		Text("Keep Reading")

	summary := hb.Paragraph().
		Class("text-center text-secondary fw-bold mb-4 mx-auto small").
		Style("max-width: 660px;").
		Text("Explore more practical guides and insights to stay ahead of contract review best practices.")

	row := bs.Row().
		Class("g-4 justify-content-center mt-2").
		Children(cards)

	viewAll := hb.Div().
		Class("text-center mt-4").
		Child(hb.Hyperlink().
			Class("btn btn-primary rounded-4 px-4 py-2 fw-black text-uppercase tracking-wider text-decoration-none").
			Href(links.Website().Blog()).
			Text("View All Blog Posts"))

	section := hb.Section().
		Style("padding:48px 0 56px;").
		Children([]hb.TagInterface{
			bs.Container().Children([]hb.TagInterface{
				heading,
				summary,
				row,
				viewAll,
			}),
		})

	return section
}

func (c *postController) postCard(post blogstore.PostInterface) hb.TagInterface {
	postURL := links.Website().BlogPost(post.GetID(), post.GetSlug())

	title := hb.Heading5().
		Class("card-title fw-black text-uppercase mb-2").
		Style("font-size: 1rem; letter-spacing: -0.5px;").
		Text(post.GetTitle())

	summary := hb.Paragraph().
		Class("card-text text-secondary fw-bold small mb-3").
		Text(c.truncatedSummary(post.GetSummary()))

	button := hb.Hyperlink().
		Class("btn btn-outline-dark rounded-4 mt-auto w-100 fw-black text-uppercase tracking-wider text-decoration-none").
		Href(postURL).
		Text("Read This Next")

	body := hb.Div().
		Class("card-body d-flex flex-column p-4").
		Children([]hb.TagInterface{
			title,
			summary,
			button,
		})

	card := hb.Div().
		Class("card rounded-4 h-100 overflow-hidden")

	if rules.PostImageURL(context.Background(), c.app.GetBlogStore(), post) != "" {
		card = card.Child(c.postImage(post))
	}

	card = card.Child(body)

	link := hb.Hyperlink().
		Class("text-decoration-none text-reset d-flex flex-column h-100 flex-fill").
		Href(postURL).
		Child(card)

	return hb.Div().
		Class("w-100 d-flex").
		Child(link)
}

func (c *postController) truncatedSummary(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	const maxLength = 160
	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}

	return strings.TrimSpace(string(runes[:maxLength])) + "..."
}

func (c *postController) postImage(post blogstore.PostInterface) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(context.Background(), c.app, post, "300", "200", "80")

	return hb.Image(``).
		Class("card-img-top slazy-placeholder slazy-resize").
		Style("object-fit: cover; height: 180px;").
		Data("slazy-src", thumbnailURL)
}

// markdownToHtml converts a markdown text to html
//
// 1. the text is trimmed of any white spaces
// 2. if the text is empty, it returns an empty string
// 3. the text is converted to html using the goldmark library
func (controller *postController) markdownToHtml(text string) string {
	text = strings.TrimSpace(text)

	if text == "" {
		return ""
	}

	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	if err := md.Convert([]byte(text), &buf); err != nil {
		// Return a plain-text fallback instead of crashing the server.
		return "<p>" + strings.ReplaceAll(text, "\n", "<br>") + "</p>"
	}

	return buf.String()
}
