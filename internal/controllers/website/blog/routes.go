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
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Guest > Blog").
			SetPath(links.BLOG).
			SetHTMLHandler(home.NewBlogController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_WITH_REGEX).
			SetHTMLHandler(post.NewPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_WITH_REGEX2).
			SetHTMLHandler(post.NewPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_01).
			SetHTMLHandler(post.NewPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_02).
			SetHTMLHandler(post.NewPostController(app).Handler),
	}
}
