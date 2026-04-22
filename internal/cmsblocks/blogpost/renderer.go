package blogpost

import (
	"fmt"
	"project/internal/links"
	"project/pkg/social"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// renderBlogPostHTML renders a single blog post
func renderBlogPostHTML(post blogstore.PostInterface, showImage, showTitle, showDate, showAuthor, showSummary, showPrevNext bool, prevPost, nextPost blogstore.PostInterface, postTags []blogstore.TermInterface) (string, error) {
	postImageURL := post.GetImageUrlOrDefault()
	publishedAt := ""
	if post.GetPublishedAt() != "" {
		publishedAt = post.GetPublishedAtCarbon().Format("d M, Y")
	}

	var content []hb.TagInterface

	// Header section with image and social share
	if showImage && postImageURL != "" {
		postURL := links.Website().BlogPost(post.GetID(), post.GetSlug())
		shareLinks := social.NewQuick(postURL, post.GetTitle(), "")
		shareWidgetHTML := shareLinks.Widget(social.WidgetOptions{
			Platforms: []string{
				social.PlatformFacebook,
				social.PlatformTwitter,
				social.PlatformLinkedIn,
				social.PlatformPinterest,
				social.PlatformCopyLink,
				social.PlatformPrint,
			},
			ShareText: "Follow us",
		})

		headerSection := bs.Row().Style("margin-bottom: 30px;").Children([]hb.TagInterface{
			bs.Column(12).Class("col-md-8").Child(
				hb.Image("").
					Class("img-fluid img-thumbnail").
					Class("slazy-placeholder").
					Class("slazy-resize-zero").
					Style("width: 100%; height: auto; object-fit: cover;").
					Data("slazy-src", postImageURL),
			),
			bs.Column(12).Class("col-md-4").Child(
				hb.Div().
					Class("d-flex align-items-center justify-content-center").
					Style("height: 100%; padding: 20px; min-height: 200px;").
					HTML(shareWidgetHTML),
			),
		})
		content = append(content, headerSection)
	}

	// Title
	if showTitle {
		title := hb.Heading1().
			Class("blog-post-title").
			Style("padding:30px 0px 10px 0px; text-align:center;").
			HTML(post.GetTitle())
		content = append(content, title)
	}

	// Get author name from meta if available
	authorName := post.GetMeta("author_name")

	// Meta info (date, author)
	if showDate || (showAuthor && authorName != "") {
		metaParts := []string{}
		if showDate && publishedAt != "" {
			metaParts = append(metaParts, publishedAt)
		}
		if showAuthor && authorName != "" {
			metaParts = append(metaParts, "by "+authorName)
		}

		if len(metaParts) > 0 {
			metaText := hb.Paragraph().
				Class("blog-post-meta").
				Style("color:#6c757d; margin-bottom:30px; text-align:center;")
			metaText.Text(fmt.Sprintf("%s", joinStrings(metaParts, " | ")))
			content = append(content, metaText)
		}
	}

	// Summary
	if showSummary && post.GetSummary() != "" {
		summary := hb.Paragraph().
			Style("padding:10px 0px 20px 0px; text-align:left; font-family:'Times New Roman', Times, serif; font-size:20px; font-style:italic; color:#555;")
		summary.Text(post.GetSummary())
		content = append(content, summary)
	}

	// Main content
	if post.GetContent() != "" {
		contentDiv := hb.Div().
			Class("BlogContent").
			Style("padding:10px 0px 20px 0px; text-align:left; font-family:'Times New Roman', Times, serif; font-size:20px; line-height:1.6;")
		contentDiv.HTML(post.GetContent())
		content = append(content, contentDiv)
	}

	// Tags section (without share, as it's now in header)
	if len(postTags) > 0 {
		tagsSection := renderTagsOnly(postTags)
		content = append(content, tagsSection)
	}

	// Previous/Next navigation
	if showPrevNext && (prevPost != nil || nextPost != nil) {
		navContent := renderPrevNextNav(prevPost, nextPost)
		content = append(content, navContent)
	}

	// Back button (moved after navigation to be less prominent)
	backBtn := hb.Div().Style("margin-top:30px; text-align:center;").Child(
		hb.Hyperlink().
			Class("view-all-posts-link text-decoration-none").
			Style("color:#6c757d; font-size:16px; transition:color 0.2s;").
			Children([]hb.TagInterface{
				hb.Span().HTML("View All Posts"),
			}).
			Attr("href", "/blog"),
	)
	content = append(content, backBtn)

	// Wrap in section
	section := hb.Section().
		ID("SectionNewsItem").
		Style("background:#fff; padding:50px 0px 80px 0px;").
		Child(bs.Container().Children([]hb.TagInterface{
			bs.Row().Children([]hb.TagInterface{
				bs.Column(12).Children(content),
			}),
		}))

	return hb.Wrap(
		// custom styling for social widget and responsive fixes
		hb.Style(`
			.blog-post-title {
				font-size: 35px;
				line-height: 1.2;
				font-weight: 700;
			}
			.blog-post-meta {
				font-size: 16px;
				letter-spacing: 0.5px;
			}
			.BlogContent p {
				margin-bottom: 1.5rem;
			}
			@media (max-width: 768px) {
				.blog-post-title {
					font-size: 28px;
				}
				.blog-post-meta {
					font-size: 14px;
				}
				#SectionNewsItem {
					padding: 30px 15px !important;
				}
				.prev-next-nav-image {
					display: none !important;
				}
				.prev-next-nav-title {
					display: none !important;
				}
				.prev-next-nav-label {
					font-size: 16px !important;
					font-weight: 600 !important;
				}
			}
			#social-links {
				display: flex;
				align-items: center;
				gap: 10px;
			}
			#social-links .share-text {
				color: #6c757d;
				font-size: 14px;
				margin-right: 8px;
			}
			#social-links ul {
				display: flex;
				list-style: none;
				margin: 0;
				padding: 0;
				gap: 8px;
			}
			#social-links li {
				margin: 0;
			}
			#social-links a.social-button {
				color: #666;
				font-size: 18px;
				text-decoration: none;
				transition: color 0.2s;
			}
			#social-links a.social-button:hover {
				color: #36C199;
			}
			.view-all-posts-link:hover {
				color: #36C199;
				text-decoration: underline;
			}
			@media print {
				/* Hide everything except the blog post content */
				header, footer, nav, aside, .navbar, .footer, .social-share-widget, #social-links, .view-all-posts-link, .prev-next-nav {
					display: none !important;
				}
				/* Show only the blog post section */
				#SectionNewsItem {
					padding: 0 !important;
				}
				.blog-post-title {
					color: #000 !important;
					text-align: left !important;
				}
				.blog-post-meta {
					color: #666 !important;
				}
				.BlogContent {
					font-size: 12pt !important;
					line-height: 1.5 !important;
					color: #000 !important;
				}
				/* Hide the header image and social share section */
				#SectionNewsItem > .container > .row > .col-12 > .row:first-child {
					display: none !important;
				}
			}
		`),
		section,
		// include lazy loading script for images
		hb.ScriptURL(cdn.Slazy_0_5_0()),
	).ToHTML(), nil
}

// renderTagsOnly renders post tags without social share links
func renderTagsOnly(postTags []blogstore.TermInterface) hb.TagInterface {
	// Build tags HTML
	var tagChildren []hb.TagInterface
	for _, tag := range postTags {
		tagLink := hb.Hyperlink().
			Class("badge text-white text-decoration-none me-2 mb-2").
			Style("background:#36C199; padding:8px 12px; font-size:14px;").
			Href("/blog/tag/" + tag.GetSlug()).
			Text(tag.GetName())
		tagChildren = append(tagChildren, tagLink)
	}

	// Container
	container := hb.Div().
		Style("margin-top:30px; padding:20px 0; border-top:1px solid #eee;").
		Child(
			hb.Div().
				Class("d-flex flex-wrap").
				Children(tagChildren),
		)

	return container
}

// renderTagsAndShare renders post tags and social share links
func renderTagsAndShare(postTags []blogstore.TermInterface, postTitle, postSlug string, postID string) hb.TagInterface {
	// Create share links using helper
	postURL := links.Website().BlogPost(postID, postSlug)
	postImageURL := "" // Could be extracted from post if needed
	shareLinks := social.NewQuick(postURL, postTitle, postImageURL)

	// Build tags HTML
	var tagChildren []hb.TagInterface
	for _, tag := range postTags {
		tagLink := hb.Hyperlink().
			Class("badge text-white text-decoration-none me-2 mb-2").
			Style("background:#36C199; padding:8px 12px; font-size:14px;").
			Href("/blog/tag/" + tag.GetSlug()).
			Text(tag.GetName())
		tagChildren = append(tagChildren, tagLink)
	}

	// Generate social share widget with popup enabled (default) and share text
	shareWidgetHTML := shareLinks.Widget(social.WidgetOptions{
		Platforms: []string{
			social.PlatformFacebook,
			social.PlatformTwitter,
			social.PlatformLinkedIn,
			social.PlatformPinterest,
			social.PlatformCopyLink,
			social.PlatformPrint,
		},
		ShareText: "Share",
	})

	// Container
	container := hb.Div().
		Style("margin-top:30px; padding:20px 0; border-top:1px solid #eee;").
		Child(
			hb.Div().
				Class("d-flex flex-wrap justify-content-between align-items-center").
				Children([]hb.TagInterface{
					hb.Div().
						Class("d-flex flex-wrap").
						Children(tagChildren),
					hb.Div().
						Class("d-flex align-items-center mt-2 mt-md-0").
						HTML(shareWidgetHTML),
				}),
		)

	return container
}

func joinStrings(parts []string, separator string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += separator
		}
		result += part
	}
	return result
}

// renderPrevNextNav renders the previous/next post navigation
func renderPrevNextNav(prevPost, nextPost blogstore.PostInterface) hb.TagInterface {
	// Previous link (left side)
	var prevLink hb.TagInterface
	if prevPost != nil {
		prevImage := ""
		if prevPost.GetImageUrlOrDefault() != "" {
			prevImage = hb.Image("").
				Class("slazy-placeholder").
				Class("slazy-resize-zero").
				Class("prev-next-nav-image").
				Style("width: 60px; height: 60px; object-fit: cover; border-radius: 4px; margin-right: 12px;").
				Data("slazy-src", prevPost.GetImageUrlOrDefault()).
				ToHTML()
		}

		prevLink = hb.Hyperlink().
			Href(links.Website().BlogPost(prevPost.GetID(), prevPost.GetSlug())).
			Style("display: flex; align-items: center; text-decoration: none; color: inherit;").
			Children([]hb.TagInterface{
				hb.Div().HTML(prevImage),
				hb.Div().Children([]hb.TagInterface{
					hb.Div().
						Class("prev-next-nav-label").
						Style("font-size: 12px; color: #888; margin-bottom: 4px;").
						Children([]hb.TagInterface{
							hb.Span().Style("margin-right: 4px;").Text("←"),
							hb.Span().Text("Previous"),
						}),
					hb.Div().
						Class("prev-next-nav-title").
						Style("font-size: 14px; font-weight: 500; color: #333;").
						Text(prevPost.GetTitle()),
				}),
			})
	} else {
		prevLink = hb.Div() // Empty placeholder
	}

	// Next link (right side)
	var nextLink hb.TagInterface
	if nextPost != nil {
		nextImage := ""
		if nextPost.GetImageUrlOrDefault() != "" {
			nextImage = hb.Image("").
				Class("slazy-placeholder").
				Class("slazy-resize-zero").
				Class("prev-next-nav-image").
				Style("width: 60px; height: 60px; object-fit: cover; border-radius: 4px; margin-left: 12px;").
				Data("slazy-src", nextPost.GetImageUrlOrDefault()).
				ToHTML()
		}

		nextLink = hb.Hyperlink().
			Href(links.Website().BlogPost(nextPost.GetID(), nextPost.GetSlug())).
			Style("display: flex; align-items: center; text-decoration: none; color: inherit; justify-content: flex-end;").
			Children([]hb.TagInterface{
				hb.Div().Children([]hb.TagInterface{
					hb.Div().
						Class("prev-next-nav-label").
						Style("font-size: 12px; color: #888; margin-bottom: 4px; text-align: right;").
						Children([]hb.TagInterface{
							hb.Span().Text("Next"),
							hb.Span().Style("margin-left: 4px;").Text("→"),
						}),
					hb.Div().
						Class("prev-next-nav-title").
						Style("font-size: 14px; font-weight: 500; color: #333; text-align: right;").
						Text(nextPost.GetTitle()),
				}),
				hb.Div().HTML(nextImage),
			})
	} else {
		nextLink = hb.Div() // Empty placeholder
	}

	// Navigation container
	nav := hb.Div().
		Style("margin-top: 40px; padding-top: 30px; border-top: 1px solid #eee;").
		Child(
			hb.Div().
				Style("display: flex; justify-content: space-between; align-items: center; gap: 20px;").
				Children([]hb.TagInterface{
					hb.Div().Style("flex: 1;").Child(prevLink),
					hb.Div().Style("flex: 1;").Child(nextLink),
				}),
		)

	return nav
}
