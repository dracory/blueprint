package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreChatMigrate)(nil)

type StoreChatMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreChatMigrate) Signature() string {
	return "2026_03_21_0007_store_chat_migrate"
}

func (m *StoreChatMigrate) Description() string {
	return "Run chat store MigrateUp to create chat tables"
}

func (m *StoreChatMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetChatStore()
	if store == nil {
		return errors.New("chat store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreChatMigrate) Down() error {
	store := m.app.GetChatStore()
	if store == nil {
		return errors.New("chat store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

