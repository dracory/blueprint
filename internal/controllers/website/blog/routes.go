package blog

import (
	"project/internal/links"
	"project/internal/types"

	"project/internal/controllers/website/blog/home"
	"project/internal/controllers/website/blog/post"

	"github.com/dracory/rtr"
)

func Routes(
	app types.AppInterface,
) []rtr.RouteInterface {
	blogRoute := rtr.NewRoute().
		SetName("Guest > Blog").
		SetPath(links.BLOG).
		SetHTMLHandler(NewBlogController(app).Handler)
		
	rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_WITH_REGEX).
			SetHTMLHandler(post.NewPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_WITH_REGEX2).
			SetHTMLHandler(post.NewPostController(app).Handler),

	blogPost01Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID)").
		SetPath(links.BLOG_01).
		SetHTMLHandler(NewBlogPostController(app).Handler)

	blogPost02Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID && Title)").
		SetPath(links.BLOG_02).
		SetHTMLHandler(NewBlogPostController(app).Handler)

	blogPost03Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID)").
		SetPath(links.BLOG_POST_01).
		SetHTMLHandler(NewBlogPostController(app).Handler)

	blogPost04Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID && Title)").
		SetPath(links.BLOG_POST_02).
		SetHTMLHandler(NewBlogPostController(app).Handler)

	return []rtr.RouteInterface{
		blogRoute,
		blogPost01Route,
		blogPost02Route,
		blogPost03Route,
		blogPost04Route,
	}
}
