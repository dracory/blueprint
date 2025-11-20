package post

import (
	"context"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"project/internal/controllers/website/blog/shared"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"

	livefluxctl "project/internal/controllers/liveflux"
)

const (
	recommendationsQueryLimit   = 18
	recommendationsDisplayLimit = 3
)

type postRecommendationsComponent struct {
	liveflux.Base
	App           types.AppInterface
	CurrentPostID string
	Posts         []blogstore.Post
	errorMessage  string
}

func NewPostRecommendationsComponent(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&postRecommendationsComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if component, ok := inst.(*postRecommendationsComponent); ok {
		component.App = app
	}

	return inst
}

func (c *postRecommendationsComponent) GetKind() string {
	return "blog_post_recommendations"
}

func (c *postRecommendationsComponent) Mount(ctx context.Context, params map[string]string) error {
	c.CurrentPostID = strings.TrimSpace(params["post_id"])

	if c.App == nil {
		if app, ok := ctx.Value(livefluxctl.AppContextKey).(types.AppInterface); ok {
			c.App = app
		}
	}

	if c.App == nil {
		c.errorMessage = "Application not initialized"
		return nil
	}

	store := c.App.GetBlogStore()
	if store == nil {
		c.errorMessage = "Blog store is not configured"
		return nil
	}

	options := blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		SortOrder: "DESC",
		OrderBy:   "published_at",
		Limit:     recommendationsQueryLimit,
	}

	postList, err := store.PostList(options)
	if err != nil {
		if logger := c.App.GetLogger(); logger != nil {
			logger.Error("Failed to load recommended posts", "error", err.Error())
		}
		c.errorMessage = "Unable to load more posts right now"
		return nil
	}

	filtered := make([]blogstore.Post, 0, len(postList))
	for _, post := range postList {
		if post.ID() == c.CurrentPostID {
			continue
		}

		if !post.IsPublished() {
			continue
		}

		filtered = append(filtered, post)
	}

	if len(filtered) == 0 {
		c.Posts = nil
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	if len(filtered) > recommendationsDisplayLimit {
		filtered = filtered[:recommendationsDisplayLimit]
	}

	c.Posts = filtered
	return nil
}

func (c *postRecommendationsComponent) Handle(ctx context.Context, action string, data url.Values) error {
	return nil
}

func (c *postRecommendationsComponent) Render(ctx context.Context) hb.TagInterface {
	if c.errorMessage != "" || len(c.Posts) == 0 {
		return c.Root(hb.Div().Text(c.errorMessage))
	}

	cards := make([]hb.TagInterface, 0, len(c.Posts))
	for _, post := range c.Posts {
		cards = append(cards, bs.Column(12).
			Class("col-md-4 col-sm-6 d-flex").
			Child(c.postCard(post)))
	}

	heading := hb.Heading3().
		Class("text-center mb-2 text-primary fw-bold").
		Text("Keep Reading")

	summary := hb.Paragraph().
		Class("text-center text-muted mb-3 mx-auto").
		Style("max-width: 660px;").
		Text("Explore more practical guides and insights to stay ahead of contract review best practices.")

	row := bs.Row().
		Class("g-4 justify-content-center mt-2").
		Children(cards)

	viewAll := hb.Div().
		Class("text-center mt-4").
		Child(hb.Hyperlink().
			Class("btn btn-primary px-4 py-2").
			Href(links.Website().Blog()).
			Text("View All Blog Posts"))

	section := hb.Section().
		Style("background:#f6f9fc;padding:48px 0 56px;").
		Children([]hb.TagInterface{
			bs.Container().Children([]hb.TagInterface{
				heading,
				summary,
				row,
				viewAll,
			}),
		})

	return c.Root(section)
}

func (c *postRecommendationsComponent) postCard(post blogstore.Post) hb.TagInterface {
	postURL := links.Website().BlogPost(post.ID(), post.Slug())

	title := hb.Heading5().
		Class("card-title fw-semibold").
		Text(post.Title())

	summary := hb.Paragraph().
		Class("card-text text-muted mb-4").
		Text(c.truncatedSummary(post.Summary()))

	button := hb.Hyperlink().
		Class("btn btn-outline-primary mt-auto w-100").
		Href(postURL).
		Text("Read This Next")

	body := hb.Div().
		Class("card-body d-flex flex-column").
		Children([]hb.TagInterface{
			title,
			summary,
			button,
		})

	card := hb.Div().
		Class("card shadow-sm border-0 h-100 overflow-hidden")

	if post.ImageUrlOrDefault() != "" {
		card = card.Child(c.postImage(post))
	}

	card = card.Child(body)

	link := hb.Hyperlink().
		Class("text-decoration-none text-reset d-flex flex-column h-100 flex-fill").
		Href(postURL).
		Child(card)

	return hb.Div().
		Class("w-100 d-flex").
		Child(link)
}

func (c *postRecommendationsComponent) truncatedSummary(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	const maxLength = 160
	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}

	return strings.TrimSpace(string(runes[:maxLength])) + "..."
}

// func (c *postRecommendationsComponent) postImage(post blogstore.Post) *hb.Tag {
// 	thumbnailURL := thumbnailURL(post)

// 	return hb.Image(thumbnailURL).
// 		Class("card-img-top rounded-3").
// 		Style("object-fit: cover;").
// 		Style("max-height: 180px;").
// 		Style("aspect-ratio: 9/6;").
// 		Style("border-radius: 0.5rem").
// 		Alt("")
// }

func (c *postRecommendationsComponent) postImage(post blogstore.Post) *hb.Tag {
	thumbnailURL := shared.SizedThumbnailURL(post, "300", "200", "80")

	return hb.Image(``).
		Class("card-img-top").
		Class(`slazy-placeholder`).
		Class(`slazy-resize`).
		Style("object-fit: cover; height: 180px; border-radius: 0.75rem 0.75rem 0 0;").
		Data("slazy-src", thumbnailURL)
}

func init() {
	if err := liveflux.Register(&postRecommendationsComponent{}); err != nil {
		log.Printf("Failed to register postRecommendationsComponent: %v", err)
	}
}
