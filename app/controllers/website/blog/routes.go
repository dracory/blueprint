package blog

import (
	"project/app/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Guest > Articles").
			SetPath("/articles").
			SetHTMLHandler(NewBlogController().Handler),

		rtr.NewRoute().
			SetName("Guest > Articles > Post with ID > Index").
			SetPath("/article/{id:[0-9]+}").
			SetHTMLHandler(NewBlogPostController().Handler),

		rtr.NewRoute().
			SetName("Guest > Articles > Post with ID && Title > Index").
			SetPath("/article/{id:[0-9]+}/{title}").
			SetHTMLHandler(NewBlogPostController().Handler),

		rtr.NewRoute().
			SetName("Guest > Blog").
			SetPath(links.BLOG).
			SetHTMLHandler(NewBlogController().Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID > Index").
			SetPath(links.BLOG_POST_WITH_REGEX).
			SetHTMLHandler(NewBlogPostController().Handler),

		rtr.NewRoute().
			SetName("Guest > Blog > Post with ID && Title > Index").
			SetPath(links.BLOG_POST_WITH_REGEX2).
			SetHTMLHandler(NewBlogPostController().Handler),
	}
}
