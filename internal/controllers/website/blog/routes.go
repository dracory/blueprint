package blog

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(
	app types.AppInterface,
) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Guest > Blog").
			SetPath(links.BLOG).
			SetHTMLHandler(NewBlogController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_WITH_REGEX).
			SetHTMLHandler(NewBlogPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_WITH_REGEX2).
			SetHTMLHandler(NewBlogPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_01).
			SetHTMLHandler(NewBlogPostController(app).Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_02).
			SetHTMLHandler(NewBlogPostController(app).Handler),
	}
}
