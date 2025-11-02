package blog

import (
	"bytes"
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/blogtheme"
	"strings"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/rtr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type blogPostController struct {
	app types.AppInterface
}

func NewBlogPostController(
	app types.AppInterface,
) *blogPostController {
	return &blogPostController{
		app: app,
	}
}

func (c blogPostController) Handler(w http.ResponseWriter, r *http.Request) string {
	postID, _ := rtr.GetParam(r, "id")
	postSlug, _ := rtr.GetParam(r, "title")

	if postID == "" || postSlug == "" {
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) >= 3 && segments[0] == "blog" && segments[1] == "post" {
			if postID == "" {
				postID = segments[2]
			}
			if postSlug == "" && len(segments) >= 4 {
				postSlug = segments[3]
			}
		}
	}

	blogsUrl := links.Website().Blog(map[string]string{})

	if postID == "" {
		c.app.GetLogger().Error("anyPost: post ID is missing", slog.String("uri", r.RequestURI))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	post, errPost := c.app.GetBlogStore().PostFindByID(postID)

	if errPost != nil {
		c.app.GetLogger().Error("Error. At BlogPostController.AnyIndex. Post not found", slog.String("error", errPost.Error()))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	if post == nil {
		c.app.GetLogger().Error("ERROR: anyPost: post with ID "+postID+" is missing", slog.String("postID", postID))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for no longer exists. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	if !c.accessAllowed(r, *post) {
		c.app.GetLogger().Error("WARNING: anyPost: post with ID "+postID+" is unpublished", slog.String("postID", postID))
		helpers.ToFlash(c.app.GetCacheStore(), w, r, "warning", "The post you are looking for is no longer active. Redirecting to the blog location...", blogsUrl, 5)
		return "post is missing"
	}

	if postSlug == "" || postSlug != post.Slug() {
		url := links.Website().BlogPost(post.ID(), post.Slug())
		c.app.GetLogger().Error("ERROR: anyPost: post Title is missing for ID "+postID, slog.String("postID", postID))
		return helpers.ToFlash(c.app.GetCacheStore(), w, r, "success", "The post location has changed. Redirecting to the new address...", url, 5)
	}

	return layouts.NewCmsLayout(
		c.app,
		r,
		layouts.Options{
			WebsiteSection: "Blog",
			Title:          post.Title(),
			StyleURLs: []string{
				"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
			},
			Content: hb.Wrap().HTML(c.page(*post)),
		}).ToHTML()
}

func (controller blogPostController) accessAllowed(r *http.Request, post blogstore.Post) bool {
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

func (c blogPostController) page(post blogstore.Post) string {
	sectionBanner := blogController{}.sectionBanner()
	return hb.Wrap().Children([]hb.TagInterface{
		// hb.Style(c.cssSectionIntro()),
		hb.Style(c.css()),
		sectionBanner,
		// c.sectionIntro(),
		c.sectionPost(post),
	}).ToHTML()
}

// css returns the CSS code for the blog post section.
//
// No parameters.
// Returns a string containing the CSS code.
func (c blogPostController) css() string {
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

func (controller *blogPostController) processContent(content string, editor string) (html string, css string) {
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
	if editor == blogstore.POST_EDITOR_MARKDOWN {
		return controller.markdownToHtml(content), ""
	}

	return content, ""
}

func (c *blogPostController) sectionPost(post blogstore.Post) *hb.Tag {
	postHtml, themeStyle := c.processContent(post.Content(), post.Editor())

	// nextPost, _ := models.NewBlogRepository().PostFindNext(post)
	// prevPost, _ := models.NewBlogRepository().PostFindPrevious(post)
	sectionPost := hb.Section().
		ID("SectionNewsItem").
		Style(`background:#fff;`).
		Child(hb.Style(themeStyle)).
		Children([]hb.TagInterface{
			bs.Container().
				Children([]hb.TagInterface{
					bs.Row().
						Child(bs.Column(12).
							Child(hb.Div().
								Class("BlogTitle").
								Child(hb.Heading1().
									HTML(post.Title())))),
					bs.Row().
						Children([]hb.TagInterface{
							bs.Column(12).Children([]hb.TagInterface{
								hb.Div().Class("BlogImage float-end").Style("padding-top:30px; padding-left:30px; padding-bottom:30px; width:600px; max-width:100%;").Children([]hb.TagInterface{
									hb.Image(post.ImageUrlOrDefault()).
										Class("img img-responsive img-thumbnail"),
								}),
								hb.Div().
									Class("BlogContent").
									HTML(postHtml),
							}),
							// bs.Column(8).Children([]hb.TagInterface{
							// 	hb.Div().Class("BlogContent").Children([]hb.TagInterface{
							// 		hb.HTML(post.Content()),
							// 	}),
							// }),
							// bs.Column(4).Children([]hb.TagInterface{
							// 	hb.Div().Class("BlogImage").Children([]hb.TagInterface{
							// 		hb.Image().Class("img img-responsive img-thumbnail").Attrs(map[string]string{
							// 			"src": post.ImageUrlOrDefault(),
							// 		}),
							// 	}),
							// }),
						}),
					// bs.Row().Children([]hb.TagInterface{
					// 	bs.Column(6).Children([]hb.TagInterface{
					// 		lo.IfF(prevPost != nil, func() *hb.Tag {
					// 			link := links.NewWebsiteLinks().BlogPost(prevPost.ID(), prevPost.Title(), map[string]string{})
					// 			return hb.Div().Children([]hb.TagInterface{
					// 				hb.Hyperlink().Children([]hb.TagInterface{
					// 					icons.Icon("bi-chevron-left", 20, 20, "#333").Style("margin-right:5px;"),
					// 					hb.Span().HTML("Previous"),
					// 				}).Attr("href", link).
					// 					Style("font-weight:bold; font-size:20px;"),
					// 				hb.Div().HTML(prevPost.Title()),
					// 			})
					// 		}).ElseF(func() *hb.Tag {
					// 			return hb.Span().HTML("")
					// 		}),
					// 	}),
					// 	bs.Column(6).Children([]hb.TagInterface{
					// 		lo.IfF(nextPost != nil, func() *hb.Tag {
					// 			link := links.NewWebsiteLinks().BlogPost(nextPost.ID(), nextPost.Title(), map[string]string{})
					// 			return hb.Div().Children([]hb.TagInterface{
					// 				hb.Hyperlink().Children([]hb.TagInterface{
					// 					hb.Span().HTML("Next"),
					// 					icons.Icon("bi-chevron-right", 20, 20, "#333").Style("margin-right:5px;"),
					// 				}).Attr("href", link).
					// 					Style("font-weight:bold; font-size:20px;"),
					// 				hb.Div().HTML(nextPost.Title()),
					// 			}).Style("text-align:right;")
					// 		}).ElseF(func() *hb.Tag {
					// 			return hb.Span().HTML("")
					// 		}),
					// 	}),
					// }),
					bs.Row().Style("margin-top:40px;").Children([]hb.TagInterface{
						bs.Column(12).Children([]hb.TagInterface{
							hb.Div().Children([]hb.TagInterface{
								hb.Hyperlink().Class("btn text-white text-center").Style(`background:#1ba1b6;color:#fff;width:600px;max-width:100%;`).Children([]hb.TagInterface{
									// icons.Icon("bi-arrow-left", 16, 16, "#333").Style("margin-right:5px;"),
									hb.Span().HTML("View All Posts"),
								}).Attr("href", links.Website().Blog()),
							}),
						}),
					}),
				}),
		})
	return sectionPost
}

// markdownToHtml converts a markdown text to html
//
// 1. the text is trimmed of any white spaces
// 2. if the text is empty, it returns an empty string
// 3. the text is converted to html using the goldmark library
func (controller *blogPostController) markdownToHtml(text string) string {
	text = strings.TrimSpace(text)

	if text == "" {
		return ""
	}

	var buf bytes.Buffer
	md := goldmark.New(
		// goldmark.WithExtensions(extension.GFM),
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
