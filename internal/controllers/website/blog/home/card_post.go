package home

import (
	"project/internal/controllers/website/blog/shared"
	"project/internal/links"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func cardPost(post blogstore.Post) hb.TagInterface {

	publishedAt := lo.Ternary(post.PublishedAt() == "", "", post.PublishedAtCarbon().Format("d M, Y"))

	postURL := links.Website().BlogPost(post.ID(), post.Slug())

	postTitle := hb.Heading5().
		Class("card-title").
		Style("font-size: 16px; color: #224b8e; margin-bottom: 10px; text-align: left; font-weight: 800;").
		Text(post.Title())

	postPublished := hb.Paragraph().
		Style("font-size: 12px;	color: #6c757d;	margin-bottom: 20px; text-align: right;").
		Text(publishedAt)

	// postPublished := hb.Span().
	// 	Class(`small`).
	// 	Style(`font-size:12px;color:#666;display:inline-block;padding-right:10px;padding-top:10px;`).
	// 	HTML(publishedAt)

	// postImage := hb.Div().Class(`overflow-hidden rounded-3`).Children([]hb.TagInterface{
	// 	hb.Image().
	// 		Class(`card-img`).
	// 		Style(`object-fit:cover;max-height:180px;`).
	// 		Src(postImageURL).
	// 		Alt("course image").
	// 		Attr("loading", "lazy"),
	// 	hb.Div().
	// 		Class(`bg-overlay bg-dark opacity-4`),
	// 	// Badge
	// 	// bs.CardImageTop().Class(`d-flex align-items-start`).Children([]hb.TagInterface{
	// 	// 	hb.Div().Class(`badge text-bg-danger`).Style(`position:absolute;top:10px;left:10px;`).HTML("Student life"),
	// 	// }),
	// })

	postSummary := hb.Paragraph().
		Class("card-text").
		Text(post.Summary()).
		Style(`text-align: left;`).
		Style(`font-size: 14px;`).
		Style(`font-weight: 400;`).
		Style(`overflow: hidden;`).
		Style(`text-overflow: ellipsis;`).
		Style(`display: -webkit-box;`).
		Style(`-webkit-line-clamp: 2;`).
		Style(`-webkit-box-orient: vertical;`)

	separator := hb.HR().
		Style(`width: 80%`).
		Style(`margin: 0 auto`).
		Style(`border: 0`).
		Style(`height: 2px`).
		Style(`background-image: linear-gradient(to right, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0))`).
		Style(`opacity: 0.25`).
		Style(`margin-bottom: 20px`)

	card := hb.Div().
		Class("card").
		Style("border: none; margin-bottom: 20px;").
		Child(postImage(post)).
		Child(hb.Div().
			Class("card-body").
			Style(`padding: 20px 10px;`).
			Child(postTitle).
			Child(postSummary)).
		Child(hb.Div().
			Class("card-footer").
			Style(`background: none;border: none;padding: 0px;`).
			Child(postPublished).
			Child(separator))

	link := hb.Hyperlink().
		Href(postURL).
		Target("_blank").
		Style("text-decoration: none; color: inherit;").
		Style("display: flex; height: 100%;").
		Child(card)

	return hb.Div().
		Class("col-md-3 col-sm-6 d-flex align-items-stretch").
		Child(link)

}

func postImage(post blogstore.Post) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(nil, post, "300", "200", "80")

	postImage := hb.Image(``).
		Class("card-img-top rounded-3").
		Class(`slazy-placeholder`).
		Class(`slazy-resize`).
		Style("object-fit: cover;").
		Style("max-height: 180px;").
		Style("aspect-ratio: 9/6;").
		Style("border-radius: 0.5rem").
		Attr("data-slazy-src", thumbnailURL).
		Attr("loading", "lazy").
		Alt("")

	return postImage
}

// func placeholderThumbnailURL(post blogstore.Post) string {
// 	return sizedThumbnailURL(post, "30", "20", "50")
// }
