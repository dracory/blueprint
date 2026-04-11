// Package social provides comprehensive social media share link generation for web content.
// It supports 30+ social media platforms with proper URL encoding and customization options.
//
// Basic Usage:
//
//	shareLinks := social.NewQuick("https://example.com/page", "My Page", "image.jpg")
//	facebookURL := shareLinks.GetFacebookShareUrl()
//	twitterURL := shareLinks.GetTwitterShareUrl()
//
// Widget (with popup enabled by default):
//
//	html := shareLinks.Widget(social.WidgetOptions{
//		Platforms: []string{social.PlatformFacebook, social.PlatformTwitter},
//		ShareText: "Share",
//	})
//
// For full customization:
//
//	params := social.ShareLinksParams{
//		URL:          "https://example.com/page",
//		Title:        "My Page",
//		Description:  "Description",
//		ImageURL:     "https://example.com/image.jpg",
//		Via:          "handle",
//		HashTags:     "golang,web",
//	}
//	shareLinks := social.New(params)
package social
