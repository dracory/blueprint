package helpers

import (
	"net/http"
	"project/internal/links"
	"strings"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/uid"
)

const FLASH_ERROR = "error"
const FLASH_SUCCESS = "success"
const FLASH_INFO = "info"
const FLASH_WARNING = "warning"

func IsFlashRoute(r *http.Request) bool {
	return strings.Contains(r.URL.Path, "/flash")
}

// ToFlashURL return a flash message URL
func ToFlashURL(cacheStore cachestore.StoreInterface, messageType string, message string, url string, seconds int) string {
	id := uid.HumanUid()
	cacheStore.SetJSON(id+"_flash_message", map[string]any{
		"type":    messageType,
		"message": message,
		"url":     url,
		"time":    seconds,
	}, int64(seconds)+10)

	return links.NewWebsiteLinks().Flash(map[string]string{
		"message_id": id,
	})
}

// ToFlash redirects the user to a flash page
func ToFlash(cachestore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, messageType string, message string, url string, seconds int) string {
	flashUrl := ToFlashURL(cachestore, messageType, message, url, seconds)

	http.Redirect(w, r, flashUrl, http.StatusSeeOther)
	return `<a href="` + flashUrl + `">See Other</a>`
}

func ToFlashError(cacheStore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, message string, url string, seconds int) string {
	return ToFlash(cacheStore, w, r, FLASH_ERROR, message, url, seconds)
}

func ToFlashInfo(cacheStore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, message string, url string, seconds int) string {
	return ToFlash(cacheStore, w, r, FLASH_INFO, message, url, seconds)
}

func ToFlashSuccess(cacheStore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, message string, url string, seconds int) string {
	return ToFlash(cacheStore, w, r, FLASH_SUCCESS, message, url, seconds)
}

func ToFlashWarning(cacheStore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, message string, url string, seconds int) string {
	return ToFlash(cacheStore, w, r, FLASH_WARNING, message, url, seconds)
}

func ToFlashErrorURL(cacheStore cachestore.StoreInterface, message string, url string, seconds int) string {
	return ToFlashURL(cacheStore, FLASH_ERROR, message, url, seconds)
}

func ToFlashInfoURL(cacheStore cachestore.StoreInterface, message string, url string, seconds int) string {
	return ToFlashURL(cacheStore, FLASH_INFO, message, url, seconds)
}

func ToFlashSuccessURL(cacheStore cachestore.StoreInterface, message string, url string, seconds int) string {
	return ToFlashURL(cacheStore, FLASH_SUCCESS, message, url, seconds)
}

func ToFlashWarningURL(cacheStore cachestore.StoreInterface, message string, url string, seconds int) string {
	return ToFlashURL(cacheStore, FLASH_WARNING, message, url, seconds)
}
