package links

import (
	"project/internal/cache"
	"strings"
	"time"

	"github.com/dracory/str"
	"github.com/samber/lo"
)

type websiteLinks struct{}

func Website() *websiteLinks {
	return &websiteLinks{}
}

func (l *websiteLinks) Home() string {
	return URL(HOME, map[string]string{})
}

func (l *websiteLinks) Blog(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(BLOG, p)
}

func (l *websiteLinks) BlogPost(postID string, postSlug string) string {
	uri := BLOG_POST
	uri += "/" + postID
	uri += "/" + postSlug
	return URL(uri, map[string]string{})
}

func (l *websiteLinks) Chat(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(CHAT_HOME, p)
}

func (l *websiteLinks) Contact() string {
	return URL(CONTACT, map[string]string{})
}

func (l *websiteLinks) File(filePath string) string {
	path := strings.TrimSuffix(FILES, CATCHALL) + "/" + filePath
	return URL(path, map[string]string{})
}

func (l *websiteLinks) Flash(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(FLASH, p)
}

func (l *websiteLinks) Shop(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(SHOP, p)
}

func (l *websiteLinks) ShopProduct(productID string, productSlug string, params map[string]string) string {
	uri := SHOP_PRODUCT
	uri += "/" + productID
	if productSlug != "" {
		uri += "/" + productSlug
	}
	return URL(uri, params)
}

func (l *websiteLinks) PaymentCanceled(paymentKey string) string {
	params := map[string]string{}
	params["payment_key"] = paymentKey
	return URL(PAYMENT_CANCELED, params)
}

func (l *websiteLinks) PaymentSuccess(paymentKey string) string {
	params := map[string]string{}
	params["payment_key"] = paymentKey
	return URL(PAYMENT_SUCCESS, params)
}

func (l *websiteLinks) Resource(resourcePath string, params ...map[string]string) string {
	if resourcePath == "" {
		return ""
	}
	if resourcePath[0] != '/' {
		resourcePath = "/" + resourcePath
	}

	resourcePath = strings.TrimSuffix(RESOURCES, CATCHALL) + resourcePath

	return URL(resourcePath, lo.FirstOr(params, map[string]string{}))
}

func (l *websiteLinks) Theme(params map[string]string) string {
	return URL(THEME, params)
}

func (l *websiteLinks) Thumbnail(extension, width, height, quality, path string) string {
	if quality == "" {
		quality = "80"
	}
	if width == "" {
		width = "100"
	}
	if height == "" {
		height = "100"
	}
	if extension == "" {
		extension = "png"
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		url := strings.ReplaceAll(path, "https://", "https/")
		path = strings.ReplaceAll(url, "http://", "http/")
	}
	if strings.HasPrefix(path, "data") {
		hash := str.MD5(path)
		if cache.File != nil {
			if err := cache.File.Save(hash, path, 5*time.Minute); err != nil {
				return RootURL() + "/th/cache-error"
			}
			path = "cache-" + hash
		}
	}
	return RootURL() + "/th/" + extension + "/" + width + "x" + height + "/" + quality + "/" + path
}

func (l *websiteLinks) Widget(alias string, params map[string]string) string {
	params["alias"] = alias
	return URL(WIDGET, params)
}
