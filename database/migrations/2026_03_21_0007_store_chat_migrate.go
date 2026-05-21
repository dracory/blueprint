package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreChatMigrate)(nil)

type StoreChatMigrate struct {
	registry RegistryInterface
}

func (m *StoreChatMigrate) ID() string {
	return "2026_03_21_0007_store_chat_migrate"
}

func (m *StoreChatMigrate) Description() string {
	return "Run chat store MigrateUp to create chat tables"
}

func (m *StoreChatMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetChatStore()
	if store == nil {
		return errors.New("chat store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreChatMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreChatMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:07:00", "UTC").StdTime()
}
