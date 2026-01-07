package post

import (
	"bytes"
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/blog/post_update"
	"project/internal/controllers/website/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/blogtheme"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/rtr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type postController struct {
	registry registry.RegistryInterface
}

func NewPostController(
	registry registry.RegistryInterface,
) *postController {
	return &postController{
		registry: registry,
	}
}

func (c *postController) Handler(w http.ResponseWriter, r *http.Request) string {
	postID, _ := rtr.GetParam(r, "id")
	postSlug, _ := rtr.GetParam(r, "title")
	blogsUrl := links.Website().Blog(map[string]string{})

	if postID == "" {
		c.registry.GetLogger().Error("anyPost: post ID is missing", slog.String("uri", r.RequestURI))
		helpers.ToFlash(c.registry.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	post, errPost := c.registry.GetBlogStore().PostFindByID(r.Context(), postID)

	if errPost != nil {
		c.registry.GetLogger().Error("Error. At BlogPostController.AnyIndex. Post not found", slog.String("error", errPost.Error()))
		helpers.ToFlash(c.registry.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	if post == nil {
		c.registry.GetLogger().Error("ERROR: anyPost: post with ID "+postID+" is missing", slog.String("postID", postID))
		helpers.ToFlash(c.registry.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return ""
	}

	if !c.accessAllowed(r, *post) {
		c.registry.GetLogger().Error("WARNING: anyPost: post with ID "+postID+" is unpublished", slog.String("postID", postID))
		helpers.ToFlash(c.registry.GetCacheStore(), w, r, "warning", "The post you are looking for is no longer active. Redirecting to the blog location...", blogsUrl, 5)
		return ""
	}

	if postSlug == "" || postSlug != post.Slug() {
		url := links.Website().BlogPost(post.ID(), post.Slug())
		c.registry.GetLogger().Error("ERROR: anyPost: post Title is missing for ID "+postID, slog.String("postID", postID))
		helpers.ToFlash(c.registry.GetCacheStore(), w, r, "success", "The post location has changed. Redirecting to the new address...", url, 5)
		return ""
	}

	options := layouts.Options{
		WebsiteSection: "Blog",
		Title:          post.Title(),
		StyleURLs: []string{
			"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
		},
		ScriptURLs: []string{
			"/liveflux",
			"https://cdn.jsdelivr.net/gh/lesichkovm/slazy@latest/dist/slazy.min.js",
		},
		Content: hb.Wrap().HTML(c.page(*post)),
	}

	if c.registry.GetConfig().GetCmsStoreUsed() {
		return layouts.NewCmsLayout(
			c.registry,
			r,
			options).ToHTML()
	} else {
		return layouts.NewBlankLayout(
			c.registry,
			r,
			options).ToHTML()
	}
}

func (controller *postController) accessAllowed(r *http.Request, post blogstore.Post) bool {
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

func (c *postController) page(post blogstore.Post) string {
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
	padding:50px 0px 80px 0px;
}

#SectionNewsItem .BlogTitle {
	padding:10px 0px 20px 0px;
	font-size:35px;
	font-family: 'Roboto', serif;
	color:#224b8e;
	text-align:centre;
}

#SectionNewsItem .BlogContent {
	padding:10px 0px 20px 0px;
	text-align:left;
	font-family: 'Roboto', Times, serif;
	font-size:20px;
}
	`
}

func (controller *postController) processContent(content string, editor string) (html string, css string) {
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

func (c *postController) sectionPost(post blogstore.Post) *hb.Tag {
	postHtml, themeStyle := c.processContent(post.Content(), post.Editor())

	rowTitle := hb.Div().
		Class("BlogTitle").
		Child(hb.Heading1().
			HTML(post.Title()))

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
		Style(`background:#fff;`).
		Child(hb.Style(themeStyle)).
		Child(hb.Div().
			Class("container").
			Child(rowTitle).
			Child(rowContent))

	return sectionPost
}

func (c *postController) recommendationsSection(post blogstore.Post) hb.TagInterface {
	component := NewPostRecommendationsComponent(c.registry)
	rendered := liveflux.Placeholder(component, map[string]string{
		"post_id": post.ID(),
	})

	// if rendered == nil {
	// 	if c.app != nil && c.app.GetLogger() != nil {
	// 		c.app.GetLogger().Warn("blogPostController: recommendations component render returned nil", "post_id", post.ID())
	// 	}
	// 	return hb.Div()
	// }

	return rendered
}

func (c *postController) postImage(post blogstore.Post) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(c.registry, post, "300", "200", "80")

	return hb.Image(``).
		Class("img img-responsive img-thumbnail w-100").
		Class(`slazy-placeholder`).
		Class(`slazy-resize-zero`).
		// Style("object-fit: cover; height: 180px; border-radius: 0.75rem 0.75rem 0 0;").
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
		panic(err)
	}

	return buf.String()
}
