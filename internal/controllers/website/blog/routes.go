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
		SetHTMLHandler(home.NewBlogController(app).Handler)

	blogPostRegex01Route := rtr.NewRoute().
		SetName("Guest > Blog > Post with ID > Index").
		SetPath(links.BLOG_POST_WITH_REGEX).
		SetHTMLHandler(post.NewPostController(app).Handler)

	blogPostRegex02Route := rtr.NewRoute().
		SetName("Guest > Blog > Post with ID && Title > Index").
		SetPath(links.BLOG_POST_WITH_REGEX2).
		SetHTMLHandler(post.NewPostController(app).Handler)

	// blogPost01Route := rtr.NewRoute().
	// 	SetName("Guest > Blog > Post (ID)").
	// 	SetPath(links.BLOG_01).
	// 	SetHTMLHandler(post.NewPostController(app).Handler)

	// blogPost02Route := rtr.NewRoute().
	// 	SetName("Guest > Blog > Post (ID && Title)").
	// 	SetPath(links.BLOG_02).
	// 	SetHTMLHandler(post.NewPostController(app).Handler)

	blogPost01Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID)").
		SetPath(links.BLOG_POST_01).
		SetHTMLHandler(post.NewPostController(app).Handler)

	blogPost02Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID && Title)").
		SetPath(links.BLOG_POST_02).
		SetHTMLHandler(post.NewPostController(app).Handler)

	return []rtr.RouteInterface{
		blogRoute,
		blogPostRegex01Route,
		blogPostRegex02Route,
		blogPost01Route,
		blogPost02Route,
	}
}
