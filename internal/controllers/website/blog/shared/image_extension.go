package shared

import "strings"

func ImageExtension(postImageURL string) string {
	lowerURL := strings.ToLower(postImageURL)
	if strings.HasSuffix(lowerURL, ".jpg") {
		return "jpg"
	} else if strings.HasSuffix(lowerURL, ".jpeg") {
		return "jpg"
	} else if strings.HasSuffix(lowerURL, ".webp") {
		return "webp"
	}

	return "png"
}

// func placeholderThumbnailURL(post blogstore.Post) string {
// 	return sizedThumbnailURL(post, "30", "20", "50")
// }
