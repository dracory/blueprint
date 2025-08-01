package testutils

import (
	"net/http"
	"project/internal/config"
	"project/internal/types"

	"github.com/dracory/base/str"
	"github.com/spf13/cast"
)

func FlashMessageFind(messageID string) (msg *types.FlashMessage, err error) {
	msgData, err := config.CacheStore.GetJSON(messageID+"_flash_message", "")
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

func FlashMessageFindFromBody(body string) (msg *types.FlashMessage, err error) {
	flashMessageID := str.LeftFrom(str.RightFrom(body, `/flash?message_id=`), `"`)
	return FlashMessageFind(flashMessageID)
}

func FlashMessageFindFromResponse(r *http.Response) (msg *types.FlashMessage, err error) {
	location := r.Header.Get("Location")
	flashMessageID := str.RightFrom(location, `/flash?message_id=`)
	return FlashMessageFind(flashMessageID)
}
