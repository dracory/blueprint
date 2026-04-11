package social

// ShareLinksParams contains all parameters for generating share links
type ShareLinksParams struct {
	URL          string
	Title        string
	Description  string
	ImageURL     string
	Via          string
	HashTags     string
	EmailAddress string
	PhoneNumber  string
	UserID       string
	CCEmail      string
	BCCEmail     string
}

// ShareLinks provides lazy-evaluated social media share links
type ShareLinks struct {
	params ShareLinksParams
}

// New creates share links with full customization
func New(params ShareLinksParams) *ShareLinks {
	return &ShareLinks{
		params: params,
	}
}

// NewQuick creates share links for a URL with basic parameters
func NewQuick(url, title, imageURL string) *ShareLinks {
	params := ShareLinksParams{
		URL:         url,
		Title:       title,
		ImageURL:    imageURL,
		Description: "",
	}

	return New(params)
}

// GetBloggerShareUrl returns the Blogger share link
func (s *ShareLinks) GetBloggerShareUrl() string {
	return to(SubmitURLBlogger, map[string]string{
		"u": s.params.URL,
		"n": s.params.Title,
		"t": s.params.Description,
	})
}

// GetBlueskyShareUrl returns the Bluesky share link
func (s *ShareLinks) GetBlueskyShareUrl() string {
	return to(SubmitURLBluesky, map[string]string{
		"text": s.params.Title + " " + s.params.URL,
	})
}

// GetDiasporaShareUrl returns the Diaspora share link
func (s *ShareLinks) GetDiasporaShareUrl() string {
	return to(SubmitURLDiaspora, map[string]string{
		"title": s.params.Title,
		"url":   s.params.URL,
	})
}

// GetDiscordShareUrl returns the Discord share link (links to Discord or invite)
func (s *ShareLinks) GetDiscordShareUrl() string {
	return to(SubmitURLDiscord, map[string]string{})
}

// GetDoubanShareUrl returns the Douban share link
func (s *ShareLinks) GetDoubanShareUrl() string {
	return to(SubmitURLDouban, map[string]string{
		"url":   s.params.URL,
		"title": s.params.Title,
	})
}

// GetEmailShareUrl returns the email share link
func (s *ShareLinks) GetEmailShareUrl() string {
	return to("mailto:"+s.params.EmailAddress, map[string]string{
		"subject": s.params.Title,
		"body":    s.params.Description,
	})
}

// GetEvernoteShareUrl returns the Evernote share link
func (s *ShareLinks) GetEvernoteShareUrl() string {
	return to(SubmitURLEvernote, map[string]string{
		"url":   s.params.URL,
		"title": s.params.Title,
	})
}

// GetFacebookShareUrl returns the Facebook share link
func (s *ShareLinks) GetFacebookShareUrl() string {
	return to(SubmitURLFacebook, map[string]string{
		"u": s.params.URL,
	})
}

// GetFlipboardShareUrl returns the Flipboard share link
func (s *ShareLinks) GetFlipboardShareUrl() string {
	return to(SubmitURLFlipboard, map[string]string{
		"v":     "2",
		"title": s.params.Title,
		"url":   s.params.URL,
	})
}

// GetGitHubShareUrl returns the GitHub share link (links to profile or repository)
func (s *ShareLinks) GetGitHubShareUrl() string {
	return to(SubmitURLGitHub, map[string]string{})
}

// GetGmailShareUrl returns the Gmail share link
func (s *ShareLinks) GetGmailShareUrl() string {
	return to(SubmitURLGmail, map[string]string{
		"view": "cm",
		"to":   s.params.EmailAddress,
		"su":   s.params.Title,
		"body": s.params.URL,
		"cc":   s.params.CCEmail,
		"bcc":  s.params.BCCEmail,
	})
}

// GetGoogleBookmarksShareUrl returns the Google Bookmarks share link
func (s *ShareLinks) GetGoogleBookmarksShareUrl() string {
	return to(SubmitURLGoogleBookmarks, map[string]string{
		"op":         "edit",
		"bkmk":       s.params.URL,
		"title":      s.params.Title,
		"annotation": s.params.Title,
		"labels":     s.params.HashTags,
	})
}

// GetInstagramShareUrl returns the Instagram share link (links to Instagram, no direct share)
func (s *ShareLinks) GetInstagramShareUrl() string {
	return to(SubmitURLInstagram, map[string]string{})
}

// GetInstapaperShareUrl returns the Instapaper share link
func (s *ShareLinks) GetInstapaperShareUrl() string {
	return to(SubmitURLInstapaper, map[string]string{
		"url":         s.params.URL,
		"title":       s.params.Title,
		"description": s.params.Description,
	})
}

// GetLineMeShareUrl returns the Line.me share link
func (s *ShareLinks) GetLineMeShareUrl() string {
	return to(SubmitURLLineMe, map[string]string{
		"url":  s.params.URL,
		"text": s.params.Title,
	})
}

// GetLinkedInShareUrl returns the LinkedIn share link
func (s *ShareLinks) GetLinkedInShareUrl() string {
	return to(SubmitURLLinkedIn, map[string]string{
		"url": s.params.URL,
	})
}

// GetLiveJournalShareUrl returns the LiveJournal share link
func (s *ShareLinks) GetLiveJournalShareUrl() string {
	return to(SubmitURLLiveJournal, map[string]string{
		"subject": s.params.Title,
		"event":   s.params.URL,
	})
}

// GetHackerNewsShareUrl returns the Hacker News share link
func (s *ShareLinks) GetHackerNewsShareUrl() string {
	return to(SubmitURLHackerNews, map[string]string{
		"u": s.params.URL,
		"t": s.params.Title,
	})
}

// GetMastodonShareUrl returns the Mastodon share link
func (s *ShareLinks) GetMastodonShareUrl() string {
	return to(SubmitURLMastodon, map[string]string{
		"title": s.params.Title,
		"url":   s.params.URL,
	})
}

// GetMediumShareUrl returns the Medium share link (links to Medium)
func (s *ShareLinks) GetMediumShareUrl() string {
	return to(SubmitURLMedium, map[string]string{})
}

// GetOkRuShareUrl returns the OK.ru share link
func (s *ShareLinks) GetOkRuShareUrl() string {
	return to(SubmitURLOkRu, map[string]string{
		"st.cmd":      "WidgetSharePreview",
		"st.shareUrl": s.params.URL,
	})
}

// GetPinterestShareUrl returns the Pinterest share link
func (s *ShareLinks) GetPinterestShareUrl() string {
	return to(SubmitURLPinterest, map[string]string{
		"url":         s.params.URL,
		"description": s.params.Title,
		"media":       s.params.ImageURL,
	})
}

// GetPocketShareUrl returns the Pocket share link
func (s *ShareLinks) GetPocketShareUrl() string {
	return to(SubmitURLPocket, map[string]string{
		"url": s.params.URL,
	})
}

// GetQZoneShareUrl returns the QZone share link
func (s *ShareLinks) GetQZoneShareUrl() string {
	return to(SubmitURLQZone, map[string]string{
		"url": s.params.URL,
	})
}

// GetRedditShareUrl returns the Reddit share link
func (s *ShareLinks) GetRedditShareUrl() string {
	return to(SubmitURLReddit, map[string]string{
		"url":   s.params.URL,
		"title": s.params.Title,
	})
}

// GetRenrenShareUrl returns the Renren share link
func (s *ShareLinks) GetRenrenShareUrl() string {
	return to(SubmitURLRenren, map[string]string{
		"resourceUrl": s.params.URL,
		"srcUrl":      s.params.URL,
		"title":       s.params.Title,
		"description": s.params.Description,
	})
}

// GetSkypeShareUrl returns the Skype share link
func (s *ShareLinks) GetSkypeShareUrl() string {
	return to(SubmitURLSkype, map[string]string{
		"url":  s.params.URL,
		"text": s.params.Title,
	})
}

// GetSkypeCallShareUrl returns the Skype call link
func (s *ShareLinks) GetSkypeCallShareUrl() string {
	return "skype:" + s.params.UserID + "?call"
}

// GetSkypeChatShareUrl returns the Skype chat link
func (s *ShareLinks) GetSkypeChatShareUrl() string {
	return "skype:" + s.params.UserID + "?chat"
}

// GetSMSShareUrl returns the SMS share link
func (s *ShareLinks) GetSMSShareUrl() string {
	return to("sms:"+s.params.PhoneNumber, map[string]string{
		"body": s.params.Title,
	})
}

// GetSnapchatShareUrl returns the Snapchat share link (links to Snapchat)
func (s *ShareLinks) GetSnapchatShareUrl() string {
	return to(SubmitURLSnapchat, map[string]string{})
}

// GetTelegramShareUrl returns the Telegram share link
func (s *ShareLinks) GetTelegramShareUrl() string {
	return to(SubmitURLTelegramMe, map[string]string{
		"url":  s.params.URL,
		"text": s.params.Title,
		"to":   s.params.PhoneNumber,
	})
}

// GetTelephoneShareUrl returns the telephone/call link (tel: protocol)
func (s *ShareLinks) GetTelephoneShareUrl() string {
	return "tel:" + s.params.PhoneNumber
}

// GetThreemaShareUrl returns the Threema share link
func (s *ShareLinks) GetThreemaShareUrl() string {
	return to(SubmitURLThreema, map[string]string{
		"text": s.params.Title,
		"id":   s.params.UserID,
	})
}

// GetThreadsShareUrl returns the Threads share link
func (s *ShareLinks) GetThreadsShareUrl() string {
	return to(SubmitURLThreads, map[string]string{
		"text": s.params.Title + " " + s.params.URL,
	})
}

// GetTikTokShareUrl returns the TikTok share link (links to TikTok home, no direct share)
func (s *ShareLinks) GetTikTokShareUrl() string {
	return to(SubmitURLTikTok, map[string]string{})
}

// GetTumblrShareUrl returns the Tumblr share link
func (s *ShareLinks) GetTumblrShareUrl() string {
	return to(SubmitURLTumblr, map[string]string{
		"canonicalUrl": s.params.URL,
		"title":        s.params.Title,
		"caption":      s.params.Description,
		"tags":         s.params.HashTags,
	})
}

// GetTwitterShareUrl returns the Twitter share link
func (s *ShareLinks) GetTwitterShareUrl() string {
	return to(SubmitURLTwitter, map[string]string{
		"url":      s.params.URL,
		"text":     s.params.Title,
		"via":      s.params.Via,
		"hashtags": s.params.HashTags,
	})
}

// GetViberShareUrl returns the Viber share link
func (s *ShareLinks) GetViberShareUrl() string {
	return to("viber://forward", map[string]string{
		"text": s.params.Title + " " + s.params.URL,
	})
}

// GetVKShareUrl returns the VK share link
func (s *ShareLinks) GetVKShareUrl() string {
	return to(SubmitURLVK, map[string]string{
		"url":     s.params.URL,
		"title":   s.params.Title,
		"comment": s.params.Description,
	})
}

// GetWeiboShareUrl returns the Weibo share link
func (s *ShareLinks) GetWeiboShareUrl() string {
	return to(SubmitURLWeibo, map[string]string{
		"url":       s.params.URL,
		"appkey":    "",
		"title":     s.params.Title,
		"pic":       s.params.ImageURL,
		"relateUid": "",
	})
}

// GetWhatsAppShareUrl returns the WhatsApp share link
func (s *ShareLinks) GetWhatsAppShareUrl() string {
	return to(SubmitURLWhatsApp, map[string]string{
		"text": s.params.Title + " " + s.params.URL,
	})
}

// GetXingShareUrl returns the Xing share link
func (s *ShareLinks) GetXingShareUrl() string {
	return to(SubmitURLXing, map[string]string{
		"url": s.params.URL,
	})
}

// GetYahooShareUrl returns the Yahoo share link
func (s *ShareLinks) GetYahooShareUrl() string {
	return to(SubmitURLYahoo, map[string]string{
		"to":      s.params.EmailAddress,
		"subject": s.params.Title,
		"body":    s.params.Title + " " + s.params.URL,
	})
}

// GetYouTubeShareUrl returns the YouTube share link (for sharing videos/channels)
func (s *ShareLinks) GetYouTubeShareUrl() string {
	return to(SubmitURLYouTube, map[string]string{})
}

// GetPrintShareUrl returns the print link (javascript:window.print())
func (s *ShareLinks) GetPrintShareUrl() string {
	return "javascript:window.print();"
}

// GetCopyLinkShareUrl returns the copy link button (uses clipboard API)
func (s *ShareLinks) GetCopyLinkShareUrl() string {
	return "javascript:navigator.clipboard.writeText('" + s.params.URL + "');"
}

// GetNativeShareShareUrl returns the native share link (Web Share API for mobile)
func (s *ShareLinks) GetNativeShareShareUrl() string {
	return "javascript:if(navigator.share){navigator.share({title:'" + s.params.Title + "',url:'" + s.params.URL + "'});}"
}

// WidgetOptions configures HTML generation for social share links.
// Use this struct to customize which platforms appear, their ordering,
// icon library choice, popup behavior, and display text.
type WidgetOptions struct {
	// Platforms is a list of social media platform names to include in the widget.
	// Use Platform* constants (e.g., PlatformFacebook, PlatformTwitter).
	// If empty, defaults to Facebook, Twitter, LinkedIn, WhatsApp, and Email.
	Platforms []string

	// IconLibrary specifies which icon library to use for the widget.
	// Use IconLibraryFontAwesome or IconLibraryBootstrap.
	// If empty, defaults to Bootstrap icons.
	IconLibrary string

	// EnablePopup controls whether share links open in a popup window.
	// When true (default), clicking a social icon opens a centered popup.
	// When false, links navigate directly in the current window.
	// Set to nil to use the default behavior (enabled).
	EnablePopup *bool

	// ShareText is optional text displayed before the share icons.
	// For example: "Share this article:" or "Share:".
	// If empty, no text is displayed.
	ShareText string

	// SortStrategy determines how platforms are sorted.
	// Use SortStrategyManual (default) to use Platforms order as specified,
	// SortStrategyAlphabetical to sort alphabetically, or SortStrategyPopularity
	// to sort by popularity (most popular first).
	SortStrategy string
}

// Widget generates HTML for social share links
// By default, includes popup JavaScript for better UX. Set EnablePopup to false to disable.
func (s *ShareLinks) Widget(opts WidgetOptions) string {
	if len(opts.Platforms) == 0 {
		opts.Platforms = []string{
			PlatformFacebook,
			PlatformTwitter,
			PlatformLinkedIn,
			PlatformWhatsApp,
			PlatformEmail,
		}
	}
	if opts.IconLibrary == "" {
		opts.IconLibrary = IconLibraryBootstrap
	}

	return widget(s, opts)
}
