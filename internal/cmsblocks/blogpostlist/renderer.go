package blogpostlist

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"project/internal/links"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// renderBlogPostListHTML renders the blog post list
func renderBlogPostListHTML(postList []blogstore.PostInterface, postCount int64, postsPerPage int, currentPage int, showPagination, showImages, showSummary, showDate bool, columns int, excerptLength int) (string, error) {
	// Calculate column class
	colClass := fmt.Sprintf("col-md-%d", 12/columns)

	switch columns {
	case 6:
		colClass = "col-md-2 col-sm-4"
	case 4:
		colClass = "col-md-3 col-sm-6"
	case 3:
		colClass = "col-md-4 col-sm-6"
	case 2:
		colClass = "col-md-6 col-sm-6"
	case 1:
		colClass = "col-12"
	}

	columnCards := lo.Map(postList, func(post blogstore.PostInterface, index int) hb.TagInterface {
		return renderPostCard(post, colClass, showImages, showSummary, showDate, excerptLength)
	})

	section := hb.Section().
		Style("background:#fff;padding-top:40px; padding-bottom: 40px;").
		Children([]hb.TagInterface{
			bs.Container().Children([]hb.TagInterface{
				bs.Row().Class("g-4").Children(columnCards),
			}),
		})

	// Add pagination if enabled
	if showPagination && postCount > int64(postsPerPage) {
		section.Child(bs.Container().Child(
			hb.Div().Class("d-flex justify-content-center mt-5").
				HTML(renderPagination(int(postCount), postsPerPage, currentPage)),
		))
	}

	return hb.Wrap(
		section,
		// include lazy loading script for images
		hb.ScriptURL(cdn.Slazy_0_5_0()),
	).ToHTML(), nil
}

func renderPostCard(post blogstore.PostInterface, colClass string, showImages, showSummary, showDate bool, excerptLength int) hb.TagInterface {
	postImageURL := post.GetImageUrlOrDefault()
	publishedAt := lo.Ternary(post.GetPublishedAt() == "", "", post.GetPublishedAtCarbon().Format("d M, Y"))
	postURL := links.Website().BlogPost(post.GetID(), post.GetSlug())
	postTitle := hb.Heading5().
		Class("card-title").
		Style("font-size: 16px; margin-bottom: 10px; text-align: left; font-weight: 800;").
		Text(post.GetTitle())

	cardContent := []hb.TagInterface{}

	// Post image
	if showImages && postImageURL != "" {
		postImage := hb.Image("").
			Class("card-img-top").
			Class("slazy-placeholder").
			Class("slazy-resize-zero").
			Style("width: 100%; height: 100%; object-fit: cover;").
			Data("slazy-src", postImageURL).
			Alt("")
		postImageWrapper := hb.Div().
			Style("height: 180px; overflow: hidden; border-radius: 0.375rem;").
			Child(postImage)
		cardContent = append(cardContent, postImageWrapper)
	}

	// Card body
	cardBody := hb.Div().
		Class("card-body").
		Style("padding: 20px 10px;").
		Child(postTitle)

	// Summary - use summary if available, otherwise generate excerpt from content
	summaryText := lo.TernaryF(
		post.GetSummary() != "",
		func() string {
			return post.GetSummary()
		},
		func() string {
			if post.GetContent() != "" {
				return generateExcerpt(post.GetContent(), excerptLength)
			}
			return ""
		},
	)

	if showSummary && summaryText != "" {
		postSummary := hb.Paragraph().
			Class("card-text").
			Text(summaryText).
			Style("text-align: left; font-size: 14px; font-weight: 400; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;")
		cardBody.Child(postSummary)
	}

	cardContent = append(cardContent, cardBody)

	// Card footer with date
	if showDate {
		cardFooter := hb.Div().
			Class("card-footer").
			Style("background: none; border: none; padding: 0px;")

		postPublished := hb.Paragraph().
			Style("font-size: 12px; color: #6c757d; margin-bottom: 20px; text-align: right;").
			Text(publishedAt)
		cardFooter.Child(postPublished)

		separator := hb.HR().
			Style("width: 80%; margin: 0 auto; border: 0; height: 2px; background-image: linear-gradient(to right, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0))")
		cardFooter.Child(separator)

		cardContent = append(cardContent, cardFooter)
	}

	// Card wrapper
	card := hb.Div().
		Class("card").
		Style("border: none; width: 100%;").
		Children(cardContent)

	// Link wrapper
	link := hb.Hyperlink().
		Href(postURL).
		Style("text-decoration: none; color: inherit; display: flex; height: 100%; width: 100%;").
		Child(card)

	// Column wrapper
	return hb.Div().
		Class(colClass + " d-flex align-items-stretch").
		Child(link)
}

func renderPagination(totalItems, itemsPerPage, currentPage int) string {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages <= 1 {
		return ""
	}

	var pageItems []hb.TagInterface

	// Previous button
	prevDisabledClass := ""
	if currentPage <= 1 {
		prevDisabledClass = "disabled"
	}
	prevPage := currentPage - 1
	prevLink := hb.Hyperlink().
		Class("page-link").
		Href(fmt.Sprintf("?page=%d", prevPage)).
		Aria("label", "Previous").
		Child(hb.Span().Aria("hidden", "true").HTML("&laquo;"))
	prevItem := hb.Li().
		Class(fmt.Sprintf("page-item %s", prevDisabledClass)).
		Child(prevLink)
	pageItems = append(pageItems, prevItem)

	// Page numbers
	for i := 1; i <= totalPages; i++ {
		activeClass := ""
		if i == currentPage {
			activeClass = "active"
		}
		pageLink := hb.Hyperlink().
			Class("page-link").
			Href(fmt.Sprintf("?page=%d", i)).
			Text(fmt.Sprintf("%d", i))
		pageItem := hb.Li().Class(fmt.Sprintf("page-item %s", activeClass)).Child(pageLink)
		pageItems = append(pageItems, pageItem)
	}

	// Next button
	nextDisabledClass := ""
	nextPage := currentPage + 1
	if currentPage >= totalPages {
		nextDisabledClass = "disabled"
	}
	nextLink := hb.Hyperlink().
		Class("page-link").
		Href(fmt.Sprintf("?page=%d", nextPage)).
		Attr("aria-label", "Next").
		Child(hb.Span().Attr("aria-hidden", "true").HTML("&raquo;"))
	nextItem := hb.Li().
		Class(fmt.Sprintf("page-item %s", nextDisabledClass)).
		Child(nextLink)
	pageItems = append(pageItems, nextItem)

	// Build pagination
	pagination := hb.Ul().
		Class("pagination").
		Children(pageItems)
	nav := hb.Nav().
		Aria("label", "Page navigation").
		Child(pagination)

	return nav.ToHTML()
}

// generateExcerpt extracts plain text from HTML content and truncates to maxLength
func generateExcerpt(content string, maxLength int) string {
	// Strip HTML tags
	excerpt := stripHTML(content)
	// Truncate if too long
	if len(excerpt) > maxLength {
		excerpt = excerpt[:maxLength]
		// Find last space to avoid cutting mid-word
		if idx := strings.LastIndex(excerpt, " "); idx > 0 {
			excerpt = excerpt[:idx]
		}
		excerpt += "..."
	}
	return strings.TrimSpace(excerpt)
}

// stripHTML removes HTML tags from a string and decodes HTML entities
func stripHTML(input string) string {
	// Remove script and style tags with content
	reScript := regexp.MustCompile(`(?s)<script.*?</script>`)
	reStyle := regexp.MustCompile(`(?s)<style.*?</style>`)
	text := reScript.ReplaceAllString(input, " ")
	text = reStyle.ReplaceAllString(text, " ")

	// Remove HTML tags
	reTag := regexp.MustCompile(`<[^>]+>`)
	text = reTag.ReplaceAllString(text, " ")

	// Replace multiple spaces with single space
	reSpace := regexp.MustCompile(`\s+`)
	text = reSpace.ReplaceAllString(text, " ")

	// Decode HTML entities (&nbsp; -> space, &amp; -> &, etc.)
	text = html.UnescapeString(text)

	return strings.TrimSpace(text)
}
