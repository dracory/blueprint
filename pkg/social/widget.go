package social

import (
	"html/template"
	"sort"
	"strings"
)

type platformInfo struct {
	getURL        func(*ShareLinks) string
	fontAwesome   string
	bootstrapIcon string
}

var platformRegistry = map[string]platformInfo{
	PlatformBlogger: {
		getURL:        (*ShareLinks).GetBloggerShareUrl,
		fontAwesome:   FontAwesomeBlogger,
		bootstrapIcon: BootstrapIconBlogger,
	},
	PlatformBluesky: {
		getURL:        (*ShareLinks).GetBlueskyShareUrl,
		fontAwesome:   FontAwesomeBluesky,
		bootstrapIcon: BootstrapIconBluesky,
	},
	PlatformDiaspora: {
		getURL:        (*ShareLinks).GetDiasporaShareUrl,
		fontAwesome:   FontAwesomeDiaspora,
		bootstrapIcon: BootstrapIconDiaspora,
	},
	PlatformDiscord: {
		getURL:        (*ShareLinks).GetDiscordShareUrl,
		fontAwesome:   FontAwesomeDiscord,
		bootstrapIcon: BootstrapIconDiscord,
	},
	PlatformDouban: {
		getURL:        (*ShareLinks).GetDoubanShareUrl,
		fontAwesome:   FontAwesomeDouban,
		bootstrapIcon: BootstrapIconDouban,
	},
	PlatformEmail: {
		getURL:        (*ShareLinks).GetEmailShareUrl,
		fontAwesome:   FontAwesomeEmail,
		bootstrapIcon: BootstrapIconEmail,
	},
	PlatformEvernote: {
		getURL:        (*ShareLinks).GetEvernoteShareUrl,
		fontAwesome:   FontAwesomeEvernote,
		bootstrapIcon: BootstrapIconEvernote,
	},
	PlatformFacebook: {
		getURL:        (*ShareLinks).GetFacebookShareUrl,
		fontAwesome:   FontAwesomeFacebook,
		bootstrapIcon: BootstrapIconFacebook,
	},
	PlatformFlipboard: {
		getURL:        (*ShareLinks).GetFlipboardShareUrl,
		fontAwesome:   FontAwesomeFlipboard,
		bootstrapIcon: BootstrapIconFlipboard,
	},
	PlatformGitHub: {
		getURL:        (*ShareLinks).GetGitHubShareUrl,
		fontAwesome:   FontAwesomeGitHub,
		bootstrapIcon: BootstrapIconGitHub,
	},
	PlatformGmail: {
		getURL:        (*ShareLinks).GetGmailShareUrl,
		fontAwesome:   FontAwesomeGmail,
		bootstrapIcon: BootstrapIconGmail,
	},
	PlatformGoogleBookmarks: {
		getURL:        (*ShareLinks).GetGoogleBookmarksShareUrl,
		fontAwesome:   FontAwesomeGoogleBookmarks,
		bootstrapIcon: BootstrapIconGoogleBookmarks,
	},
	PlatformInstagram: {
		getURL:        (*ShareLinks).GetInstagramShareUrl,
		fontAwesome:   FontAwesomeInstagram,
		bootstrapIcon: BootstrapIconInstagram,
	},
	PlatformInstapaper: {
		getURL:        (*ShareLinks).GetInstapaperShareUrl,
		fontAwesome:   FontAwesomeInstapaper,
		bootstrapIcon: BootstrapIconInstapaper,
	},
	PlatformLineMe: {
		getURL:        (*ShareLinks).GetLineMeShareUrl,
		fontAwesome:   FontAwesomeLineMe,
		bootstrapIcon: BootstrapIconLineMe,
	},
	PlatformLinkedIn: {
		getURL:        (*ShareLinks).GetLinkedInShareUrl,
		fontAwesome:   FontAwesomeLinkedIn,
		bootstrapIcon: BootstrapIconLinkedIn,
	},
	PlatformLiveJournal: {
		getURL:        (*ShareLinks).GetLiveJournalShareUrl,
		fontAwesome:   FontAwesomeLiveJournal,
		bootstrapIcon: BootstrapIconLiveJournal,
	},
	PlatformHackerNews: {
		getURL:        (*ShareLinks).GetHackerNewsShareUrl,
		fontAwesome:   FontAwesomeHackerNews,
		bootstrapIcon: BootstrapIconHackerNews,
	},
	PlatformMastodon: {
		getURL:        (*ShareLinks).GetMastodonShareUrl,
		fontAwesome:   FontAwesomeMastodon,
		bootstrapIcon: BootstrapIconMastodon,
	},
	PlatformMedium: {
		getURL:        (*ShareLinks).GetMediumShareUrl,
		fontAwesome:   FontAwesomeMedium,
		bootstrapIcon: BootstrapIconMedium,
	},
	PlatformOkRu: {
		getURL:        (*ShareLinks).GetOkRuShareUrl,
		fontAwesome:   FontAwesomeOkRu,
		bootstrapIcon: BootstrapIconOkRu,
	},
	PlatformPinterest: {
		getURL:        (*ShareLinks).GetPinterestShareUrl,
		fontAwesome:   FontAwesomePinterest,
		bootstrapIcon: BootstrapIconPinterest,
	},
	PlatformPocket: {
		getURL:        (*ShareLinks).GetPocketShareUrl,
		fontAwesome:   FontAwesomePocket,
		bootstrapIcon: BootstrapIconPocket,
	},
	PlatformPrint: {
		getURL:        (*ShareLinks).GetPrintShareUrl,
		fontAwesome:   FontAwesomePrint,
		bootstrapIcon: BootstrapIconPrint,
	},
	PlatformCopyLink: {
		getURL:        (*ShareLinks).GetCopyLinkShareUrl,
		fontAwesome:   FontAwesomeCopyLink,
		bootstrapIcon: BootstrapIconCopyLink,
	},
	PlatformNativeShare: {
		getURL:        (*ShareLinks).GetNativeShareShareUrl,
		fontAwesome:   FontAwesomeNativeShare,
		bootstrapIcon: BootstrapIconNativeShare,
	},
	PlatformQZone: {
		getURL:        (*ShareLinks).GetQZoneShareUrl,
		fontAwesome:   FontAwesomeQZone,
		bootstrapIcon: BootstrapIconQZone,
	},
	PlatformReddit: {
		getURL:        (*ShareLinks).GetRedditShareUrl,
		fontAwesome:   FontAwesomeReddit,
		bootstrapIcon: BootstrapIconReddit,
	},
	PlatformRenren: {
		getURL:        (*ShareLinks).GetRenrenShareUrl,
		fontAwesome:   FontAwesomeRenren,
		bootstrapIcon: BootstrapIconRenren,
	},
	PlatformSkype: {
		getURL:        (*ShareLinks).GetSkypeShareUrl,
		fontAwesome:   FontAwesomeSkype,
		bootstrapIcon: BootstrapIconSkype,
	},
	PlatformSMS: {
		getURL:        (*ShareLinks).GetSMSShareUrl,
		fontAwesome:   FontAwesomeSMS,
		bootstrapIcon: BootstrapIconSMS,
	},
	PlatformSnapchat: {
		getURL:        (*ShareLinks).GetSnapchatShareUrl,
		fontAwesome:   FontAwesomeSnapchat,
		bootstrapIcon: BootstrapIconSnapchat,
	},
	PlatformTelegramMe: {
		getURL:        (*ShareLinks).GetTelegramShareUrl,
		fontAwesome:   FontAwesomeTelegramMe,
		bootstrapIcon: BootstrapIconTelegramMe,
	},
	PlatformThreema: {
		getURL:        (*ShareLinks).GetThreemaShareUrl,
		fontAwesome:   FontAwesomeThreema,
		bootstrapIcon: BootstrapIconThreema,
	},
	PlatformThreads: {
		getURL:        (*ShareLinks).GetThreadsShareUrl,
		fontAwesome:   FontAwesomeThreads,
		bootstrapIcon: BootstrapIconThreads,
	},
	PlatformTikTok: {
		getURL:        (*ShareLinks).GetTikTokShareUrl,
		fontAwesome:   FontAwesomeTikTok,
		bootstrapIcon: BootstrapIconTikTok,
	},
	PlatformTumblr: {
		getURL:        (*ShareLinks).GetTumblrShareUrl,
		fontAwesome:   FontAwesomeTumblr,
		bootstrapIcon: BootstrapIconTumblr,
	},
	PlatformTwitter: {
		getURL:        (*ShareLinks).GetTwitterShareUrl,
		fontAwesome:   FontAwesomeTwitter,
		bootstrapIcon: BootstrapIconTwitter,
	},
	PlatformVK: {
		getURL:        (*ShareLinks).GetVKShareUrl,
		fontAwesome:   FontAwesomeVK,
		bootstrapIcon: BootstrapIconVK,
	},
	PlatformWeibo: {
		getURL:        (*ShareLinks).GetWeiboShareUrl,
		fontAwesome:   FontAwesomeWeibo,
		bootstrapIcon: BootstrapIconWeibo,
	},
	PlatformWhatsApp: {
		getURL:        (*ShareLinks).GetWhatsAppShareUrl,
		fontAwesome:   FontAwesomeWhatsApp,
		bootstrapIcon: BootstrapIconWhatsApp,
	},
	PlatformXing: {
		getURL:        (*ShareLinks).GetXingShareUrl,
		fontAwesome:   FontAwesomeXing,
		bootstrapIcon: BootstrapIconXing,
	},
	PlatformYahoo: {
		getURL:        (*ShareLinks).GetYahooShareUrl,
		fontAwesome:   FontAwesomeYahoo,
		bootstrapIcon: BootstrapIconYahoo,
	},
	PlatformYouTube: {
		getURL:        (*ShareLinks).GetYouTubeShareUrl,
		fontAwesome:   FontAwesomeYouTube,
		bootstrapIcon: BootstrapIconYouTube,
	},
}

type linkData struct {
	Href      template.URL
	HrefAttr  template.HTMLAttr // For javascript: URLs (entire href="..." attribute)
	IconClass string
	Platform  string
	Title     string
}

type widgetData struct {
	ShareText string
	Links     []linkData
}

// getPopupScript returns the JavaScript for popup functionality
func getPopupScript() string {
	script := `(function() {
var popupSize = {width: 780, height: 550};

function openSharePopup(e) {
	var link = e.currentTarget;
	var href = link.getAttribute('href');
	// Don't use popup for mailto:, sms:, tel:, viber:, skype:, and javascript: links
	if (href.indexOf('mailto:') === 0 || href.indexOf('sms:') === 0 || href.indexOf('tel:') === 0 || href.indexOf('viber:') === 0 || href.indexOf('skype:') === 0 || href.indexOf('javascript:') === 0) {
		return true;
	}
	
	var left = Math.floor((window.innerWidth - popupSize.width) / 2);
	var top = Math.floor((window.innerHeight - popupSize.height) / 2);
	
	var popup = window.open(
		href,
		'social-share',
		'width=' + popupSize.width + 
		',height=' + popupSize.height +
		',left=' + left + 
		',top=' + top +
		',location=0,menubar=0,toolbar=0,status=0,scrollbars=1,resizable=1'
	);
	
	if (popup) {
		popup.focus();
		e.preventDefault();
		return false;
	}
	
	return true;
}

function initSocialButtons() {
	var buttons = document.querySelectorAll('.social-button');
	for (var i = 0; i < buttons.length; i++) {
		buttons[i].addEventListener('click', openSharePopup);
	}
}

if (document.readyState === 'loading') {
	document.addEventListener('DOMContentLoaded', initSocialButtons);
} else {
	initSocialButtons();
}

})();`
	return `<script>` + script + `</script>`
}

// sortPlatformsAlphabetically returns platforms sorted alphabetically by their nice names
func sortPlatformsAlphabetically(platforms []string) []string {
	niceNames := SocialMediaNiceNames()
	sorted := make([]string, len(platforms))
	copy(sorted, platforms)
	sort.Slice(sorted, func(i, j int) bool {
		nameI := niceNames[sorted[i]]
		nameJ := niceNames[sorted[j]]
		if nameI == "" {
			nameI = sorted[i]
		}
		if nameJ == "" {
			nameJ = sorted[j]
		}
		return nameI < nameJ
	})
	return sorted
}

// sortPlatformsByPopularity returns platforms sorted by popularity (most popular first)
func sortPlatformsByPopularity(platforms []string) []string {
	popular := SocialMediaSitesByPopularity()
	popularityRank := make(map[string]int)
	for i, site := range popular {
		popularityRank[site] = i
	}

	sorted := make([]string, len(platforms))
	copy(sorted, platforms)
	sort.Slice(sorted, func(i, j int) bool {
		rankI, okI := popularityRank[sorted[i]]
		rankJ, okJ := popularityRank[sorted[j]]
		if !okI && !okJ {
			return sorted[i] < sorted[j]
		}
		if !okI {
			return false
		}
		if !okJ {
			return true
		}
		return rankI < rankJ
	})
	return sorted
}

// widget generates HTML for social share links (private)
func widget(s *ShareLinks, opts WidgetOptions) string {
	var links []linkData

	// Determine platform order: apply SortStrategy if specified
	platforms := opts.Platforms
	if opts.SortStrategy != "" && opts.SortStrategy != SortStrategyManual {
		switch opts.SortStrategy {
		case SortStrategyAlphabetical:
			platforms = sortPlatformsAlphabetically(platforms)
		case SortStrategyPopularity:
			platforms = sortPlatformsByPopularity(platforms)
		}
	}

	for _, platform := range platforms {
		info, ok := platformRegistry[platform]
		if !ok {
			continue
		}

		href := info.getURL(s)
		iconClass := info.fontAwesome
		if opts.IconLibrary == IconLibraryBootstrap {
			iconClass = info.bootstrapIcon
		}

		// Get nice name for title attribute
		niceNames := SocialMediaNiceNames()
		title := niceNames[platform]
		if title == "" {
			title = platform
		}

		// For executable-scheme URLs, build the entire href attribute to avoid URL encoding
		var hrefAttr template.HTMLAttr
		if strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "data:") || strings.HasPrefix(href, "vbscript:") {
			hrefAttr = template.HTMLAttr(`href="` + href + `"`)
		}

		links = append(links, linkData{
			Href:      template.URL(href),
			HrefAttr:  hrefAttr,
			IconClass: iconClass,
			Platform:  platform,
			Title:     title,
		})
	}

	tmpl := template.Must(template.New("widget").Parse(`<div id="social-links">{{if .ShareText}}<span class="share-text">{{.ShareText}}</span>{{end}}<ul>{{range .Links}}<li><a {{if .HrefAttr}}{{.HrefAttr}}{{else}}href="{{.Href}}"{{end}} class="social-button" id="social-{{.Platform}}" title="{{.Title}}"><span class="{{.IconClass}}"></span></a></li>{{end}}</ul></div>`))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, widgetData{
		ShareText: opts.ShareText,
		Links:     links,
	}); err != nil {
		return ""
	}
	html := buf.String()

	// Enable popup by default
	enablePopup := true
	if opts.EnablePopup != nil {
		enablePopup = *opts.EnablePopup
	}

	if !enablePopup {
		return html
	}

	return html + getPopupScript()
}
