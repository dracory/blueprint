package home

import (
	"context"
	"project/internal/app"
	"project/internal/controllers/website/blog/shared"
	"project/internal/links"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func cardPost(ctx context.Context, app app.AppInterface, post blogstore.PostInterface) hb.TagInterface {

	publishedAt := lo.Ternary(post.GetPublishedAt() == "", "", post.GetPublishedAtCarbon().Format("d M, Y"))

	postURL := links.Website().BlogPost(post.GetID(), post.GetSlug())

	postTitle := hb.Heading5().
		Class("card-title fw-black text-uppercase mb-2").
		Style("font-size: 1rem; text-align: left; letter-spacing: -0.5px;").
		Text(post.GetTitle())

	postPublished := hb.Span().
		Class("badge bg-dark fw-black small").
		Text(publishedAt)

	postSummary := hb.Paragraph().
		Class("card-text text-secondary fw-bold small mb-0").
		Text(post.GetSummary()).
		Style(`text-align: left;`).
		Style(`overflow: hidden;`).
		Style(`text-overflow: ellipsis;`).
		Style(`display: -webkit-box;`).
		Style(`-webkit-line-clamp: 3;`).
		Style(`-webkit-box-orient: vertical;`)

	card := hb.Div().
		Class("card rounded-4 h-100 overflow-hidden").
		Child(postImage(ctx, app, post)).
		Child(hb.Div().
			Class("card-body p-4").
			Child(postTitle).
			Child(postSummary)).
		Child(hb.Div().
			Class("card-footer border-0 bg-transparent p-4 pt-0").
			Child(postPublished))

	link := hb.Hyperlink().
		Href(postURL).
		Style("text-decoration: none; color: inherit;").
		Style("display: flex; height: 100%;").
		Child(card)

	return hb.Div().
		Class("col-md-6 col-lg-4 d-flex align-items-stretch").
		Child(link)

}

func postImage(ctx context.Context, app app.AppInterface, post blogstore.PostInterface) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(ctx, app, post, "300", "200", "80")

	postImage := hb.Image(``).
		Class("card-img-top").
		Class(`slazy-placeholder`).
		Class(`slazy-resize`).
		Style("object-fit: cover;").
		Style("height: 200px;").
		Style("aspect-ratio: 9/6;").
		Attr("data-slazy-src", thumbnailURL).
		Attr("loading", "lazy").
		Alt("")

	return postImage
}

// func placeholderThumbnailURL(post blogstore.Post) string {
// 	return sizedThumbnailURL(post, "30", "20", "50")
// }
