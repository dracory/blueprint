package testutils

import (
	"net/http"
	"project/internal/types"

	"github.com/dracory/str"
	"github.com/gouniverse/cachestore"
	"github.com/spf13/cast"
)

func FlashMessageFind(cacheStore cachestore.StoreInterface, messageID string) (msg *types.FlashMessage, err error) {
	msgData, err := cacheStore.GetJSON(messageID+"_flash_message", "")
	if err != nil {
		return msg, err
	}

	if msgData == "" {
		return msg, nil
	}

	msgDataAny := msgData.(map[string]interface{})
	dataMap := &types.FlashMessage{
		Type:    cast.ToString(msgDataAny["type"]),
		Message: cast.ToString(msgDataAny["message"]),
		Url:     cast.ToString(msgDataAny["url"]),
		Time:    cast.ToString(msgDataAny["time"]),
	}

	return dataMap, nil
}

func FlashMessageFindFromBody(cacheStore cachestore.StoreInterface, body string) (msg *types.FlashMessage, err error) {
	flashMessageID := str.LeftFrom(str.RightFrom(body, `/flash?message_id=`), `"`)
	return FlashMessageFind(cacheStore, flashMessageID)
}

func FlashMessageFindFromResponse(cacheStore cachestore.StoreInterface, r *http.Response) (msg *types.FlashMessage, err error) {
	location := r.Header.Get("Location")
	flashMessageID := str.RightFrom(location, `/flash?message_id=`)
	return FlashMessageFind(cacheStore, flashMessageID)
}
