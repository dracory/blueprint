package social

import (
	"net/url"
	"sort"
)

// to is a helper to build a URL with query parameters
func to(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Set(key, value)
		}
	}
	if len(values) == 0 {
		return baseURL
	}
	return baseURL + "?" + values.Encode()
}

// SocialMediaNiceNames returns the display names for all social media platforms
func SocialMediaNiceNames() map[string]string {
	return map[string]string{
		PlatformBlogger:         "Blogger",
		PlatformBluesky:         "Bluesky",
		PlatformDiaspora:        "Diaspora",
		PlatformDiscord:         "Discord",
		PlatformDouban:          "Douban",
		PlatformEmail:           "EMail",
		PlatformEvernote:        "EverNote",
		PlatformFacebook:        "FaceBook",
		PlatformFlipboard:       "FlipBoard",
		PlatformPocket:          "Pocket",
		PlatformGitHub:          "GitHub",
		PlatformGmail:           "GMail",
		PlatformGoogleBookmarks: "GoogleBookmarks",
		PlatformHackerNews:      "HackerNews",
		PlatformInstagram:       "Instagram",
		PlatformInstapaper:      "InstaPaper",
		PlatformLineMe:          "Line.me",
		PlatformLinkedIn:        "LinkedIn",
		PlatformLiveJournal:     "LiveJournal",
		PlatformMastodon:        "Mastodon",
		PlatformMedium:          "Medium",
		PlatformOkRu:            "OK.ru",
		PlatformPinterest:       "Pinterest",
		PlatformPrint:           "Print",
		PlatformCopyLink:        "Copy Link",
		PlatformNativeShare:     "Native Share",
		PlatformQZone:           "QZone",
		PlatformReddit:          "Reddit",
		PlatformRenren:          "RenRen",
		PlatformSkype:           "Skype",
		PlatformSMS:             "SMS",
		PlatformSnapchat:        "Snapchat",
		PlatformTelegramMe:      "Telegram.me",
		PlatformThreema:         "Threema",
		PlatformThreads:         "Threads",
		PlatformTikTok:          "TikTok",
		PlatformTumblr:          "Tumblr",
		PlatformTwitter:         "Twitter",
		PlatformVK:              "VK",
		PlatformWeibo:           "Weibo",
		PlatformWhatsApp:        "WhatsApp",
		PlatformXing:            "Xing",
		PlatformYahoo:           "Yahoo",
		PlatformYouTube:         "YouTube",
	}
}

// SocialMediaSitesByPopularity returns social media sites ordered by popularity
func SocialMediaSitesByPopularity() []string {
	return []string{
		PlatformGoogleBookmarks,
		PlatformFacebook,
		PlatformReddit,
		PlatformWhatsApp,
		PlatformEmail,
		PlatformGmail,
		PlatformYahoo,
		PlatformBluesky,
		PlatformThreads,
		PlatformTikTok,
		PlatformInstagram,
		PlatformYouTube,
		PlatformGitHub,
		PlatformDiscord,
		PlatformMastodon,
		PlatformSnapchat,
		PlatformMedium,
		PlatformPrint,
		PlatformCopyLink,
		PlatformNativeShare,
		PlatformLinkedIn,
		PlatformTumblr,
		PlatformPinterest,
		PlatformBlogger,
		PlatformLiveJournal,
		PlatformEvernote,
		PlatformPocket,
		PlatformHackerNews,
		PlatformFlipboard,
		PlatformInstapaper,
		PlatformDiaspora,
		PlatformQZone,
		PlatformVK,
		PlatformWeibo,
		PlatformOkRu,
		PlatformDouban,
		PlatformXing,
		PlatformRenren,
		PlatformThreema,
		PlatformSMS,
		PlatformLineMe,
		PlatformSkype,
		PlatformTelegramMe,
	}
}

// SocialMediaSitesByAlphabet returns platforms sorted alphabetically
func SocialMediaSitesByAlphabet() []string {
	sites := make([]string, 0, len(SocialMediaNiceNames()))
	for site := range SocialMediaNiceNames() {
		sites = append(sites, site)
	}
	sort.Strings(sites)
	return sites
}

// SocialMediaColors returns a map of platform names to their brand colors (hex values)
// Colors are sourced from WPZoom/WPZOOM Social Icons widget reference
func SocialMediaColors() map[string]string {
	return map[string]string{
		PlatformBluesky:     ColorBluesky,
		PlatformDiscord:     ColorDiscord,
		PlatformFacebook:    ColorFacebook,
		PlatformGitHub:      ColorGitHub,
		PlatformInstagram:   ColorInstagram,
		PlatformLinkedIn:    ColorLinkedIn,
		PlatformMastodon:    ColorMastodon,
		PlatformMedium:      ColorMedium,
		PlatformPinterest:   ColorPinterest,
		PlatformReddit:      ColorReddit,
		PlatformSkype:       ColorSkype,
		PlatformSnapchat:    ColorSnapchat,
		PlatformTelegram:    ColorTelegram,
		PlatformTelegramMe:  ColorTelegram,
		PlatformTikTok:      ColorTikTok,
		PlatformThreads:     ColorThreads,
		PlatformTwitter:     ColorTwitter,
		PlatformWhatsApp:    ColorWhatsApp,
		PlatformYouTube:     ColorYouTube,
		PlatformEmail:       ColorDefault,
		PlatformPrint:       ColorDefault,
		PlatformCopyLink:    ColorDefault,
		PlatformNativeShare: ColorDefault,
	}
}
