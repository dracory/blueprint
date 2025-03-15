package testutils

import (
	"errors"
	"net/http"
	"project/app/config"
	"project/internal/types"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/utils"
)

func FlashMessageFind(messageID string) (msg *types.FlashMessage, err error) {
	cfg, err := config.New()
	if err != nil {
		return nil, errors.New("failed to create config: " + err.Error())
	}
	
	cacheStore, ok := cfg.CacheStore.(*cachestore.Store)
	if !ok {
		return nil, errors.New("cache store is not initialized or is of wrong type")
	}
	
	msgData, err := cacheStore.GetJSON(messageID+"_flash_message", "")
	if err != nil {
		return msg, err
	}

	if msgData == "" {
		return msg, nil
	}

	msgDataAny := msgData.(map[string]interface{})
	dataMap := &types.FlashMessage{
		Type:    utils.ToString(msgDataAny["type"]),
		Message: utils.ToString(msgDataAny["message"]),
		Url:     utils.ToString(msgDataAny["url"]),
		Time:    utils.ToString(msgDataAny["time"]),
	}

	return dataMap, nil
}

func FlashMessageFindFromBody(body string) (msg *types.FlashMessage, err error) {
	flashMessageID := utils.StrLeftFrom(utils.StrRightFrom(body, `/flash?message_id=`), `"`)
	return FlashMessageFind(flashMessageID)
}

func FlashMessageFindFromResponse(r *http.Response) (msg *types.FlashMessage, err error) {
	location := r.Header.Get("Location")
	flashMessageID := utils.StrRightFrom(location, `/flash?message_id=`)
	return FlashMessageFind(flashMessageID)
}
