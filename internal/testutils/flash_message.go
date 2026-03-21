package testutils

import (
	"errors"
	"net/http"
	"strings"

	basetypes "github.com/dracory/base/types"

	"github.com/dracory/cachestore"
	"github.com/dracory/str"
	"github.com/spf13/cast"
)

func FlashMessageFind(cacheStore cachestore.StoreInterface, messageID string) (msg *basetypes.FlashMessage, err error) {
	if cacheStore == nil {
		return msg, errors.New("flash message find: cache store is nil")
	}

	msgData, err := cacheStore.GetJSON(messageID+"_flash_message", "")
	if err != nil {
		return msg, err
	}

	if msgData == "" {
		return msg, nil
	}

	msgDataAny := msgData.(map[string]interface{})
	dataMap := &basetypes.FlashMessage{
		Type:    cast.ToString(msgDataAny["type"]),
		Message: cast.ToString(msgDataAny["message"]),
		Url:     cast.ToString(msgDataAny["url"]),
		Time:    cast.ToString(msgDataAny["time"]),
	}

	return dataMap, nil
}

func FlashMessageFindFromBody(cacheStore cachestore.StoreInterface, body string) (msg *basetypes.FlashMessage, err error) {
	flashMessageID := str.LeftFrom(str.RightFrom(body, `/flash?message_id=`), `"`)
	return FlashMessageFind(cacheStore, flashMessageID)
}

func FlashMessageFindFromResponse(cacheStore cachestore.StoreInterface, r *http.Response) (msg *basetypes.FlashMessage, err error) {
	if r == nil {
		return msg, errors.New("flash message find from response: response is nil")
	}

	location := r.Header.Get("Location")
	if location == "" {
		return msg, errors.New("flash message find from response: no Location header found")
	}

	// Check if the location contains the flash message pattern
	if !strings.Contains(location, "/flash?message_id=") {
		return msg, errors.New("flash message find from response: Location header does not contain flash message pattern")
	}

	flashMessageID := str.RightFrom(location, `/flash?message_id=`)
	if flashMessageID == "" {
		return msg, errors.New("flash message find from response: could not extract message ID from Location header")
	}

	return FlashMessageFind(cacheStore, flashMessageID)
}
