package registry

import (
	"database/sql"
	"errors"

	"github.com/dracory/chatstore"
)

// chatStoreInitialize initializes the chat store when enabled via configuration.
func chatStoreInitialize(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetChatStoreUsed() {
		return nil
	}

	store, err := newChatStore(registry.GetDatabase())
	if err != nil {
		return err
	}

	registry.SetChatStore(store)
	return nil
}

func chatStoreMigrate(registry RegistryInterface) error {
	if registry.GetConfig() == nil {
		return errors.New("config is not initialized")
	}

	if !registry.GetConfig().GetChatStoreUsed() {
		return nil
	}

	chatStore := registry.GetChatStore()
	if chatStore == nil {
		return errors.New("chat store is not initialized")
	}

	err := chatStore.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// newChatStore constructs the Chat store without running migrations
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
