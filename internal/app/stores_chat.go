package app

import (
	"database/sql"
	"errors"

	"project/internal/types"

	"github.com/dracory/chatstore"
)

// chatStoreInitialize initializes the chat store when enabled via configuration.
func chatStoreInitialize(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetChatStoreUsed() {
		return nil
	}

	store, err := newChatStore(app.GetDB())
	if err != nil {
		return err
	}

	app.SetChatStore(store)
	return nil
}

func chatStoreMigrate(app types.AppInterface) error {
	if app.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !app.GetConfig().GetChatStoreUsed() {
		return nil
	}

	if app.GetChatStore() == nil {
		return errors.New("chat store is not initialized")
	}

	if err := app.GetChatStore().AutoMigrate(); err != nil {
		return err
	}

	return nil
}

func newChatStore(db *sql.DB) (chatstore.StoreInterface, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	store, err := chatstore.NewStore(chatstore.NewStoreOptions{
		DB:               db,
		TableChatName:    "snv_chat_chats",
		TableMessageName: "snv_chat_messages",
	})
	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("chatstore.NewStore returned a nil store")
	}

	return store, nil
}
